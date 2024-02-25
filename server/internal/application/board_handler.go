package application

import (
	"dalkak/internal/infrastructure/eventbus"
	boarddto "dalkak/pkg/dto/board"
	responseutil "dalkak/pkg/utils/response"
)

func (app *ApplicationImpl) RegisterBoardEventListeners() {
	app.EventManager.Subscribe("post.board", app.handleCreateBoard)
}

func (app *ApplicationImpl) handleCreateBoard(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*boarddto.CreateBoardRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 트랜잭션 시작

	// 미디어 변경

	// 보드 생성
	dto := boarddto.NewCreateBoardDto(userInfo, payload.Name, payload.Description, payload.ExternalLink, payload.BackgroundColor, payload.Attributes)
	newBoard, err := app.BoardDomain.CreateBoard(dto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 오더 생성

	// 스토리지 이동

	// 트랜잭션 // 보드 저장 // 오더 저장	// 미디어 변경

	// 리턴 // todo update
	result := boarddto.NewCreateBoardResponse("", newBoard.BoardEntity.Status, "", 0, 0, 0)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
