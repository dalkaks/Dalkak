package board

import "dalkak/pkg/dtos"

type BoardImageData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id          string
	Extension   string
	ContentType string
	Url         string
}

const BoardImageDataType = "BoardImage"

func GenerateBoardDataPk(boardId string) string {
	return BoardImageDataType + `#` + boardId
}

func (b *BoardImageData) ToBoardImageDto() *dtos.BoardImageDto {
	return &dtos.BoardImageDto{
		Id:          b.Id,
		Extension:   b.Extension,
		ContentType: b.ContentType,
		Url:         b.Url,
	}
}
