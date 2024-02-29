package database

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
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

// get board by userId and status
