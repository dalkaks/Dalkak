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
	dto *boarddto.CreateBoardDto
}

func NewCreateBoardDtoFactory(dto *boarddto.CreateBoardDto) *CreateBoardDtoFactory {
	return &CreateBoardDtoFactory{
		dto: dto,
	}
}

func (factory *CreateBoardDtoFactory) CreateBoardAggregate() (*boardaggregate.BoardAggregate, error) {
	board := boardentity.NewBoardEntity(factory.dto.UserInfo.GetUserId())

	boardCategory, err := boardvalueobject.NewBoardCategory(factory.dto.CategoryType, factory.dto.Network)
	if err != nil {
		return nil, err
	}

	boardMetadata, err := boardvalueobject.NewNftMetadata(factory.dto.Name, factory.dto.Description, factory.dto.ExternalLink, factory.dto.BackgroundColor, factory.dto.Attributes)
	if err != nil {
		return nil, err
	}

	boardAggregate, err := boardaggregate.NewBoardAggregate(board, boardCategory, boardMetadata)
	if err != nil {
		return nil, err
	}
	return boardAggregate, nil
}
