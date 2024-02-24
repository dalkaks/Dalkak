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

	// 보드 생성	// 보드 조회 (중복 확인 및 결제 상태에 있는 보드 있는지 확인)

	// 오더 생성 

	// 미디어 변경

	// 스토리지 이동

	// 트랜잭션	// 보드 저장 // 오더 저장	// 미디어 변경

	// 리턴 // redirect 고려
	result := payload
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
