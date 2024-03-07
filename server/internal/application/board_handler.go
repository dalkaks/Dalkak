package application

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"
	"dalkak/internal/infrastructure/eventbus"
	boarddto "dalkak/pkg/dto/board"
	mediadto "dalkak/pkg/dto/media"
	orderdto "dalkak/pkg/dto/order"
	responseutil "dalkak/pkg/utils/response"
)

func (app *ApplicationImpl) RegisterBoardEventListeners() {
	app.EventManager.Subscribe("post.board", app.handleCreateBoard)
	app.EventManager.Subscribe("get.board.list.processing", app.handleGetBoardListProcessing)
	app.EventManager.Subscribe("delete.board", app.handleDeleteBoard)
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

	txResult, err := ExecuteOptimisticTransactionWithRetry(app, func(txId string) (*TransactionResult, error) {
		// 보드 생성
		boardCreateDto := boarddto.NewCreateBoardDto(userInfo, payload.Name, payload.Description, payload.CategoryType, payload.Network, payload.ExternalLink, payload.BackgroundColor, payload.Attributes)
		newBoard, err := app.BoardDomain.CreateBoard(boardCreateDto)
		if err != nil {
			return nil, err
		}

		// 미디어 변경
		mediaNftCreateDto := mediadto.NewCreateMediaNftDto(userInfo, newBoard.BoardEntity.Timestamp, "board", newBoard.BoardEntity.Id, payload.ImageId, payload.VideoId)
		mediaNft, tempImage, tempVideo, err := app.MediaDomain.CreateMediaNft(mediaNftCreateDto)
		if err != nil {
			return nil, err
		}

		// 오더 생성
		orderCreateDto := orderdto.NewCreateOrderDto(userInfo, string(ordervalueobject.OrderCategoryTypeNft), newBoard.BoardEntity.Id, newBoard.BoardMetadata.Name, nil)
		newOrder, err := app.OrderDomain.CreateOrder(orderCreateDto)
		if err != nil {
			return nil, err
		}

		// 트랜잭션 // 보드 저장 // 오더 저장	// 미디어 변경
		err = app.Database.CreateBoard(txId, newBoard, newOrder, mediaNft.MediaImageResource, mediaNft.MediaVideoResource)
		if err != nil {
			return nil, err
		}

		// 스토리지 이동
		if tempImage != nil {
			go func() {
				app.Storage.CopyObject(tempImage.MediaUrl.AccessUrl, mediaNft.MediaImageUrl.AccessUrl)
			}()
		}
		if tempVideo != nil {
			go func() {
				app.Storage.CopyObject(tempVideo.MediaUrl.AccessUrl, mediaNft.MediaVideoUrl.AccessUrl)
			}()
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

func (app *ApplicationImpl) handleGetBoardListProcessing(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*boarddto.GetBoardListProcessingRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	// 보드 리스트 조회
	getBoardFilter := app.BoardDomain.GetBoardListProcessingFilter(userInfo, payload)
	boardDaos, page, err := app.Database.FindBoardByUserId(getBoardFilter, &dao.RequestPageDao{Limit: payload.Limit, ExclusiveStartKey: &payload.ExclusiveStartKey})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}
	if len(boardDaos) == 0 {
		result := boarddto.NewGetBoardListProcessingResponse(nil, nil, page)
		app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeSuccess), nil)
		return
	}

	// 보드 리스트 변환
	boards, err := app.BoardDomain.ConvertBoardDaos(boardDaos)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 보드 리스트 미디어 변환
	medias, err := app.MediaDomain.ConvertBoardDaosToMediaNft(boardDaos)
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	result := boarddto.NewGetBoardListProcessingResponse(boards, medias, page)
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(result, responseutil.DataCodeSuccess), nil)
}

func (app *ApplicationImpl) handleDeleteBoard(event eventbus.Event) {
	userInfo := event.UserInfo
	if userInfo == nil {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth))
		return
	}
	payload, ok := event.Payload.(*boarddto.DeleteBoardRequest)
	if !ok {
		app.SendResponse(event.ResponseChan, nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid))
		return
	}

	_, err := ExecuteOptimisticTransactionWithRetry(app, func(txId string) (interface{}, error) {
		// 보드 조회
		boardDao, err := app.BoardDomain.GetBoardById(userInfo, payload.Id)
		if err != nil {
			return nil, err
		}

		// 보드 변환
		_, err = app.BoardDomain.ConvertBoardDao(boardDao)
		if err != nil {
			return nil, err
		}

		// 보드 상태 체크

		// // 트랜잭션 // 보드 삭제 // 오더 삭제
		// err = app.Database.DeleteBoard(txId, board)
		// if err != nil {
		// 	return nil, err
		// }

		// // 스토리지 삭제

		return nil, nil
	})
	if err != nil {
		app.SendResponse(event.ResponseChan, nil, err)
		return
	}

	// 리턴
	app.SendResponse(event.ResponseChan, responseutil.NewAppData(nil, responseutil.DataCodeSuccess), nil)
}
