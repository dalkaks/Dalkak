package boardaggregate

import (
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	responseutil "dalkak/pkg/utils/response"
)

type BoardAggregate struct {
	BoardEntity   boardentity.BoardEntity
	BoardCategory boardvalueobject.BoardCategory
	BoardMetadata boardvalueobject.NftMetadata
}

func NewBoardAggregate(board *boardentity.BoardEntity, category *boardvalueobject.BoardCategory, metadata *boardvalueobject.NftMetadata) (*BoardAggregate, error) {
	if board == nil || metadata == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	aggregate := &BoardAggregate{
		BoardEntity:   *board,
		BoardCategory: *category,
		BoardMetadata: *metadata,
	}

	return aggregate, nil
}

func (ag *BoardAggregate) CheckBoardDeleteAble() bool {
	return ag.BoardEntity.GetStatus() == string(boardentity.BoardCreated)
}

func (ag *BoardAggregate) UpdateBoardCancel() {
	ag.BoardEntity.SetStatus(boardentity.BoardCancelled)
}
