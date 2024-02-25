package boardaggregate

import (
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
)

type BoardAggregate struct {
	BoardEntity   *boardentity.BoardEntity
	BoardMetadata *boardvalueobject.NftMetadata
}

type BoardAggregateOption func(*BoardAggregate)

func WithBoardEntity(boardEntity *boardentity.BoardEntity) BoardAggregateOption {
	return func(aggregate *BoardAggregate) {
		aggregate.BoardEntity = boardEntity
	}
}

func WithBoardMetadata(boardMetadata *boardvalueobject.NftMetadata) BoardAggregateOption {
	return func(aggregate *BoardAggregate) {
		aggregate.BoardMetadata = boardMetadata
	}
}

func NewBoardAggregate(options ...BoardAggregateOption) *BoardAggregate {
	aggregate := &BoardAggregate{}
	for _, option := range options {
		option(aggregate)
	}
	return aggregate
}
