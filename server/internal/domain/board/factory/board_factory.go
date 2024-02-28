package boardfactory

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	boarddto "dalkak/pkg/dto/board"
)

type BoardAggregateFactory interface {
	CreateBoardAggregate() (*boardaggregate.BoardAggregate, error)
}

type CreateBoardDtoFactory struct {
	dto       *boarddto.CreateBoardDto
	boardType boardentity.BoardType
}

func NewCreateBoardDtoFactory(dto *boarddto.CreateBoardDto, boardType boardentity.BoardType) *CreateBoardDtoFactory {
	return &CreateBoardDtoFactory{
		dto:       dto,
		boardType: boardType,
	}
}

func (factory *CreateBoardDtoFactory) CreateBoardAggregate() (*boardaggregate.BoardAggregate, error) {
	board := boardentity.NewBoardEntity(factory.boardType, factory.dto.UserInfo.GetUserId())
	boardMetadata, err := boardvalueobject.NewNftMetadata(factory.dto.Name, factory.dto.Description, factory.dto.ExternalLink, factory.dto.BackgroundColor, factory.dto.Attributes)
	if err != nil {
		return nil, err
	}

	boardAggregate, err := boardaggregate.NewBoardAggregate(board,	boardMetadata)
	if err != nil {
		return nil, err
	}
	return boardAggregate, nil
}
