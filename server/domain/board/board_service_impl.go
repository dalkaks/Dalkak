package board

import (
	"dalkak/pkg/interfaces"
)

type BoardServiceImpl struct {
	mode    string
	domain  string
	db      interfaces.BoardRepository
	storage interfaces.Storage
}

const boardStoragePath = "board"

func NewBoardService(mode string, domain string, db interfaces.Database, storage interfaces.Storage) *BoardServiceImpl {
	boardRepo := NewBoardRepository(db)

	return &BoardServiceImpl{
		mode:    mode,
		domain:  domain,
		db:      boardRepo,
		storage: storage,
	}
}
