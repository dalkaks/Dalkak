package boardfactory

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"
	boarddto "dalkak/pkg/dto/board"
)

type BoardAggregateFactory interface {
	CreateBoardAggregateFromDto(dto *boarddto.CreateBoardDto) (*boardaggregate.BoardAggregate, error)
	CreateBoardAggregateFromDad(dao *dao.BoardDao) (*boardaggregate.BoardAggregate, error)
	CreateBoardAggregatesFromDaos(daos []*dao.BoardDao) ([]*boardaggregate.BoardAggregate, error)
}

type CreateBoardFactory struct {
}

func NewCreateBoardFactory() *CreateBoardFactory {
	return &CreateBoardFactory{}
}

func (factory *CreateBoardFactory) CreateBoardAggregateFromDto(dto *boarddto.CreateBoardDto) (*boardaggregate.BoardAggregate, error) {
	board := boardentity.NewBoardEntity(dto.UserInfo.GetUserId())

	boardCategory, err := boardvalueobject.NewBoardCategory(dto.CategoryType, dto.Network)
	if err != nil {
		return nil, err
	}

	boardMetadata, err := boardvalueobject.NewNftMetadata(dto.Name, dto.Description, dto.ExternalLink, dto.BackgroundColor, dto.Attributes)
	if err != nil {
		return nil, err
	}

	boardAggregate, err := boardaggregate.NewBoardAggregate(board, boardCategory, boardMetadata)
	if err != nil {
		return nil, err
	}
	return boardAggregate, nil
}

func (factory *CreateBoardFactory) CreateBoardAggregateFromDao(dao *dao.BoardDao) (*boardaggregate.BoardAggregate, error) {
	board, err := boardentity.ConvertBoardEntity(dao.Id, dao.UserId, dao.Timestamp, dao.Status)
	if err != nil {
		return nil, err
	}

	boardCategory, err := boardvalueobject.NewBoardCategory(dao.Type, dao.Network)
	if err != nil {
		return nil, err
	}

	boardMetadata, err := boardvalueobject.NewNftMetadata(dao.NftMetaName, dao.NftMetaDesc, dao.NftMetaExtUrl, dao.NftMetaBgCol, dao.NftMetaAttrib)
	if err != nil {
		return nil, err
	}

	boardAggregate, err := boardaggregate.NewBoardAggregate(board, boardCategory, boardMetadata)
	if err != nil {
		return nil, err
	}
	return boardAggregate, nil
}

func (factory *CreateBoardFactory) CreateBoardAggregatesFromDaos(daos []*dao.BoardDao) ([]*boardaggregate.BoardAggregate, error) {
	aggregates := make([]*boardaggregate.BoardAggregate, 0, len(daos))
	for _, dao := range daos {
		aggregate, err := factory.CreateBoardAggregateFromDao(dao)
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, aggregate)
	}
	return aggregates, nil
}
