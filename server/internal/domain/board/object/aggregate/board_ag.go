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

type BoardAggregateOption func(*BoardAggregate)

func NewBoardAggregate(board *boardentity.BoardEntity, category *boardvalueobject.BoardCategory, metadata *boardvalueobject.NftMetadata, options ...BoardAggregateOption) (*BoardAggregate, error) {
	if board == nil || metadata == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	aggregate := &BoardAggregate{
		BoardEntity:   *board,
		BoardCategory: *category,
		BoardMetadata: *metadata,
	}

	for _, option := range options {
		option(aggregate)
	}
	return aggregate, nil
}
