package application

import (
	"dalkak/internal/infrastructure/eventbus"
	mediadto "dalkak/pkg/dto/media"
	responseutil "dalkak/pkg/utils/response"
)

func (app *ApplicationImpl) RegisterMediaEventListeners() {
	app.EventManager.Subscribe("post.media.presigned", app.handleCreateTempMedia)
}

func (app *ApplicationImpl) handleCreateTempMedia(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*mediadto.CreateTempMediaRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 미디어 생성
	dto := mediadto.NewCreateTempMediaDto(userInfo, payload.MediaType, payload.Ext, payload.Prefix)
	newMedia, err := app.MediaDomain.CreateMediaTemp(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 미디어 저장
	err = app.Database.CreateUserMediaTemp(userInfo.GetUserId(), newMedia)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := mediadto.NewUserCreateMediaResponse(newMedia.MediaEntity.Id, newMedia.MediaTempUrl.AccessUrl, *newMedia.MediaTempUrl.UploadUrl)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
