package board

type BoardImageTable struct {
	Id          string
	BoardId     *string
	Extension   string
	ContentType string
	Url         string
	UserId      string
	Timestamp   int64
}

const BoardImageTableName = "board_image"
