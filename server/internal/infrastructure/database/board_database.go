package database

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	"dalkak/internal/infrastructure/database/dao"
	responseutil "dalkak/pkg/utils/response"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const BoardDataType = "Board"

func GenerateBoardDataPk(boardId string) string {
	return BoardDataType + `#` + boardId
}

type BoardData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id     string
	Status string
	UserId string

	Type    string
	TypeId  string
	Network string

	NftMetaName   string
	NftMetaDesc   string
	NftMetaExtUrl *string
	NftMetaBgCol  *string
	NftMetaAttrib []*boardvalueobject.NftAttribute

	NftImageExt *string
	NftVideoExt *string
}

func (repo *Database) CreateBoard(txId string, board *boardaggregate.BoardAggregate, order *orderaggregate.OrderAggregate, imageResource, videoResource *mediavalueobject.MediaResource) error {
	builder := NewTransactionBuilder(repo.table, txId)
	pk := GenerateBoardDataPk(board.BoardEntity.Id)

	var nftImageExt, nftVideoExt string
	if imageResource != nil {
		nftImageExt = imageResource.GetExtension()
		deleteMediaData := CreateDeleteMediaData(board.BoardEntity.UserId, imageResource)
		builder.AddDeleteItem(deleteMediaData)
	}
	if videoResource != nil {
		nftVideoExt = videoResource.GetExtension()
		deleteMediaData := CreateDeleteMediaData(board.BoardEntity.UserId, videoResource)
		builder.AddDeleteItem(deleteMediaData)
	}

	newBoard := &BoardData{
		Pk:         pk,
		Sk:         pk,
		EntityType: BoardDataType,
		Timestamp:  board.BoardEntity.Timestamp,

		Id:     board.BoardEntity.Id,
		Status: board.BoardEntity.GetStatus(),
		UserId: board.BoardEntity.UserId,

		Type:    board.BoardCategory.GetCategoryType(),
		TypeId:  board.BoardCategory.GetCategoryId(),
		Network: board.BoardCategory.GetNetwork(),

		NftMetaName:   board.BoardMetadata.Name,
		NftMetaDesc:   board.BoardMetadata.Description,
		NftMetaExtUrl: board.BoardMetadata.ExternalUrl,
		NftMetaBgCol:  board.BoardMetadata.BackgroundColor,
		NftMetaAttrib: board.BoardMetadata.Attributes,

		NftImageExt: &nftImageExt,
		NftVideoExt: &nftVideoExt,
	}
	builder.AddPutItem(newBoard)

	orderData := CreateOrderData(order)
	builder.AddPutItem(orderData)

	err := repo.WriteTransaction(builder)
	if err != nil {
		return err
	}
	return nil
}

// get board by id(all) and sometimes filter

func (repo *Database) FindBoardByUserId(filter *dao.BoardFindFilter, pageDao *dao.RequestPageDao) ([]*dao.BoardDao, *dao.ResponsePageDao, error) {
	index := EntityTypeTimestampIndex
	pk := BoardDataType
	var boardToFind []*BoardData

	keyCond := expression.Key("EntityType").Equal(expression.Value(pk))

	var builder expression.Builder
	builder = builder.WithKeyCondition(keyCond)

	if filter != nil {
		if filter.UserId != "" {
			builder = builder.WithFilter(expression.Name("UserId").Equal(expression.Value(filter.UserId)))
		}
		if filter.StatusIncluded != nil {
			builder = builder.WithFilter(expression.Name("Status").Equal(expression.Value(filter.StatusIncluded)))
		}
		if filter.StatusExcluded != nil {
			builder = builder.WithFilter(expression.Name("Status").NotEqual(expression.Value(filter.StatusExcluded)))
		}

		if filter.CategoryType != nil {
			builder = builder.WithFilter(expression.Name("Type").Equal(expression.Value(*filter.CategoryType)))
		}
		if filter.CategoryId != nil {
			builder = builder.WithFilter(expression.Name("TypeId").Equal(expression.Value(*filter.CategoryId)))
		}
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, nil, err
	}

	page, err := repo.QueryItems(expr, &index, pageDao, &boardToFind)
	if err != nil {
		return nil, nil, err
	}
	if len(boardToFind) == 0 {
		return nil, page, nil
	}

	var boardDaos []*dao.BoardDao
	for _, board := range boardToFind {
		boardDaos = append(boardDaos, &dao.BoardDao{
			Id:        board.Id,
			Status:    board.Status,
			UserId:    board.UserId,
			Timestamp: board.Timestamp,

			Type:    board.Type,
			TypeId:  board.TypeId,
			Network: board.Network,

			NftMetaName:   board.NftMetaName,
			NftMetaDesc:   board.NftMetaDesc,
			NftMetaExtUrl: board.NftMetaExtUrl,
			NftMetaBgCol:  board.NftMetaBgCol,
			NftMetaAttrib: board.NftMetaAttrib,

			NftImageExt: board.NftImageExt,
			NftVideoExt: board.NftVideoExt,
		})
	}
	return boardDaos, page, nil
}

func (repo *Database) FindBoardById(boardId string) (*dao.BoardDao, error) {
	pk := GenerateBoardDataPk(boardId)
	var boardToFind *BoardData

	keyCond := expression.Key("Pk").Equal(expression.Value(pk)).
		And(expression.Key("Sk").Equal(expression.Value(pk)))
	expr, err := GenerateQueryExpression(keyCond, nil)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	err = repo.QuerySingleItem(expr, &boardToFind)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}
	if boardToFind == nil {
		return nil, nil
	}

	boardDao := &dao.BoardDao{
		Id:        boardToFind.Id,
		Status:    boardToFind.Status,
		UserId:    boardToFind.UserId,
		Timestamp: boardToFind.Timestamp,

		Type:    boardToFind.Type,
		TypeId:  boardToFind.TypeId,
		Network: boardToFind.Network,

		NftMetaName:   boardToFind.NftMetaName,
		NftMetaDesc:   boardToFind.NftMetaDesc,
		NftMetaExtUrl: boardToFind.NftMetaExtUrl,
		NftMetaBgCol:  boardToFind.NftMetaBgCol,
		NftMetaAttrib: boardToFind.NftMetaAttrib,

		NftImageExt: boardToFind.NftImageExt,
		NftVideoExt: boardToFind.NftVideoExt,
	}
	return boardDao, nil
}

func (repo *Database) UpdateBoardCancel(txId string, board *boardaggregate.BoardAggregate) error {
	builder := NewTransactionBuilder(repo.table, txId)

	key := CreateBoardKey(board.BoardEntity.Id)

	update := expression.Set(expression.Name("Status"), expression.Value(board.BoardEntity.Status)).
		Set(expression.Name("Timestamp"), expression.Value(board.BoardEntity.Timestamp))
	expr, err := GenerateUpdateExpression(update)
	if err != nil {
		return err
	}
	
	builder.AddUpdateItem(key, expr)

	err = repo.WriteTransaction(builder)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Database) DeleteBoard(txId string, board *boardaggregate.BoardAggregate, order *orderaggregate.OrderAggregate) error {
	builder := NewTransactionBuilder(repo.table, txId)
	
	deleteBoardKey := CreateBoardKey(board.BoardEntity.Id)
	builder.AddDeleteItem(deleteBoardKey)

	deleteOrderKey := CreateOrderKey(order.OrderEntity.Id)
	builder.AddDeleteItem(deleteOrderKey)

	err := repo.WriteTransaction(builder)
	if err != nil {
		return err
	}
	return nil
}

func CreateBoardKey(boardId string) map[string]types.AttributeValue {
	pk := GenerateBoardDataPk(boardId)
	return map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk},
		"Sk": &types.AttributeValueMemberS{Value: pk},
	}
}