package application

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
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

	type TransactionResult struct {
		newBoard *boardaggregate.BoardAggregate
		mediaNft *mediaaggregate.MediaNftAggregate
		newOrder *orderaggregate.OrderAggregate
	}

	txResult, err := ExecuteTransaction[*TransactionResult](app, func(txId string) (*TransactionResult, error) {
		// 보드 생성
		boardCreateDto := boarddto.NewCreateBoardDto(userInfo, payload.Name, payload.Description, payload.ExternalLink, payload.BackgroundColor, payload.Attributes)
		newBoard, err := app.BoardDomain.CreateBoard(boardCreateDto)
		if err != nil {
			return nil, err
		}

		// 미디어 변경
		mediaNftCreateDto := mediadto.NewCreateMediaNftDto(userInfo, newBoard.BoardEntity.Timestamp, "board", newBoard.BoardEntity.Id, payload.ImageId, payload.VideoId)
		mediaNft, err := app.MediaDomain.CreateMediaNft(mediaNftCreateDto)
		if err != nil {
			return nil, err
		}

		// 오더 생성
		orderCreateDto := orderdto.NewCreateOrderDto(userInfo, string(ordervalueobject.OrderCategoryTypeNft), newBoard.BoardEntity.Id, newBoard.BoardMetadata.Name, nil)
		newOrder, err := app.OrderDomain.CreateOrder(orderCreateDto)
		if err != nil {
			return nil, err
		}

		// 스토리지 이동

		// 트랜잭션 // 보드 저장 // 오더 저장	// 미디어 변경
		err = app.Database.CreateBoard(newBoard, nil, nil)
		if err != nil {
			return nil, err
		}

		return &TransactionResult{newBoard, mediaNft, newOrder}, nil
	})

	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴 // todo update
	result := boarddto.NewCreateBoardResponse(txResult.newBoard, txResult.newOrder)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeCreated), nil)
}
