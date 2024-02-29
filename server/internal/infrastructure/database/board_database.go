package database

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	Type   string

	NftMetaName   string
	NftMetaDesc   string
	NftMetaExtUrl *string
	NftMetaBgCol  *string
	NftMetaAttrib []*boardvalueobject.NftAttribute

	NftImageExt *string
	NftVideoExt *string
}

func (repo *Database) CreateBoard(board *boardaggregate.BoardAggregate, nftImageExt, nftVideoExt *string) error {
	pk := GenerateBoardDataPk(board.BoardEntity.Id)
	newBoard := &BoardData{
		Pk:         pk,
		Sk:         pk,
		EntityType: BoardDataType,
		Timestamp:  board.BoardEntity.Timestamp,

		Id:     board.BoardEntity.Id,
		Status: board.BoardEntity.Status.String(),
		UserId: board.BoardEntity.UserId,
		Type:   board.BoardEntity.Type.String(),

		NftMetaName:   board.BoardMetadata.Name,
		NftMetaDesc:   board.BoardMetadata.Description,
		NftMetaExtUrl: board.BoardMetadata.ExternalUrl,
		NftMetaBgCol:  board.BoardMetadata.BackgroundColor,
		NftMetaAttrib: board.BoardMetadata.Attributes,

		NftImageExt: nftImageExt,
		NftVideoExt: nftVideoExt,
	}

	err := repo.PutDynamoDBItem(newBoard)
	if err != nil {
		return err
	}
	return nil
}

// get board by id(all) and sometimes filter

func (repo *Database) FindBoardByUserId(userId string, status *boardentity.BoardStatus, pageDao *dao.RequestPageDao) ([]*dao.BoardDao, *dao.ResponsePageDao, error) {
	index := UserIdEntityTypeIndex
	pk := userId
	sk := BoardDataType
	var boardToFind []*BoardData

	keyCond := expression.Key("UserId").Equal(expression.Value(pk)).
		And(expression.Key("EntityType").Equal(expression.Value(sk)))

	var expr expression.Expression
	var err error
	if status != nil {
		statusStr := status.String()
		filt := expression.Name("Status").Equal(expression.Value(statusStr))
		expr, err = GenerateQueryExpression(keyCond, &filt)
	} else {
		expr, err = GenerateQueryExpression(keyCond, nil)
	}
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
			Type:      board.Type,
			Timestamp: board.Timestamp,

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
