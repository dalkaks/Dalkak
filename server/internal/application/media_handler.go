package application

import (
	mediadomain "dalkak/internal/domain/media"
	"dalkak/internal/infrastructure/eventbus"
	mediadto "dalkak/pkg/dto/media"
	responseutil "dalkak/pkg/utils/response"
)

func (app *ApplicationImpl) RegisterMediaEventListeners() {
	app.EventManager.Subscribe("post.media.presigned", app.handleCreateMediaTemp)
	app.EventManager.Subscribe("get.media", app.handleGetMediaTemp)
	app.EventManager.Subscribe("post.media.confirm", app.handleConfirmMediaTemp)
	app.EventManager.Subscribe("delete.media", app.handleDeleteMediaTemp)
}

func (app *ApplicationImpl) handleCreateMediaTemp(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*mediadto.CreateMediaTempRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 미디어 생성
	dto := mediadto.NewCreateMediaTempDto(userInfo, payload.MediaType, payload.Ext, payload.Prefix)
	newMedia, err := app.MediaDomain.CreateMediaTemp(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 미디어 저장
	err = app.Database.CreateMediaTemp(userInfo.GetUserId(), newMedia)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := mediadto.NewCreateMediaTempResponse(newMedia.MediaEntity.Id, newMedia.MediaUrl.AccessUrl, *newMedia.MediaUrl.UploadUrl)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}

func (app *ApplicationImpl) handleGetMediaTemp(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*mediadto.GetMediaTempRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 미디어 조회
	dto := mediadto.NewGetMediaTempDto(userInfo, payload.MediaType, payload.Prefix)
	media, err := app.MediaDomain.GetMediaTemp(dto, mediadomain.GetMediaTempOptions{CheckPublic: true})
	if media == nil || err != nil {
		app.SendResponse(event.ResponseChan, media, err)
		return
	}

	// 리턴
	result := mediadto.NewGetMediaTempResponse(media.MediaEntity.Id, media.MediaResource.ContentType.String(), media.MediaUrl.AccessUrl)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeSuccess), nil)
}

func (app *ApplicationImpl) handleConfirmMediaTemp(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*mediadto.ConfirmMediaTempRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 미디어 컨펌
	dto := mediadto.NewConfirmMediaTempDto(userInfo, payload.Id, payload.MediaType, payload.Prefix)
	mediaTempUpdate, err := app.MediaDomain.ConfirmMediaTemp(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 미디어 저장
	err = app.Database.UpdateMediaTempConfirm(userInfo.GetUserId(), mediaTempUpdate)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := mediadto.NewConfirmMediaTempResponse(mediaTempUpdate.MediaUrl.AccessUrl)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeSuccess), nil)
}

func (app *ApplicationImpl) handleDeleteMediaTemp(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*mediadto.DeleteMediaTempRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 미디어 조회
	dto := mediadto.NewGetMediaTempDto(userInfo, payload.MediaType, payload.Prefix)
	media, err := app.MediaDomain.GetMediaTemp(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}
	if media == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgDataNotFound))
		return
	}

	// 미디어 삭제
	err = app.Database.DeleteMediaTemp(userInfo.GetUserId(), media)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}
	err = app.Storage.DeleteObject(media.MediaUrl.GetUrlKey(app.AppConfig.StaticLink))
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(nil, responseutil.DataCodeSuccess), nil)
}
