package application

import (
	"dalkak/internal/infrastructure/eventbus"
	boarddto "dalkak/pkg/dto/board"
	mediadto "dalkak/pkg/dto/media"
	orderdto "dalkak/pkg/dto/order"
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

	// 보드 생성
	boardCreateDto := boarddto.NewCreateBoardDto(userInfo, payload.Name, payload.Description, payload.ExternalLink, payload.BackgroundColor, payload.Attributes)
	newBoard, err := app.BoardDomain.CreateBoard(boardCreateDto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 미디어 변경
	mediaNftCreateDto := mediadto.NewCreateMediaNftDto(userInfo, "board", newBoard.BoardEntity.Id, &payload.ImageId, &payload.VideoId)
	mediaNft, err := app.MediaDomain.CreateMediaNft(mediaNftCreateDto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 오더 생성
	// todo pay price
	orderCreateDto := orderdto.NewCreateOrderDto(userInfo, "board", newBoard.BoardEntity.Id, newBoard.BoardMetadata.Name, nil, 0, 0, 0)
	newOrder, err := app.OrderDomain.CreateOrder(orderCreateDto)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 스토리지 이동

	// 트랜잭션 // 보드 저장 // 오더 저장	// 미디어 변경

	// 리턴 // todo update
	result := boarddto.NewCreateBoardResponse(mediaNft.MediaEntity.Id, newBoard.BoardEntity.Status, "", newOrder.OrderPrice.OriginPrice, newOrder.OrderPrice.DiscountPrice, newOrder.OrderPrice.PaymentPrice)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
