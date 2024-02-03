package board

import "dalkak/pkg/interfaces"

type BoardServiceImpl struct {
	mode   string
	domain string
	db     interfaces.BoardRepository
}

func NewBoardService(mode string, domain string, db interfaces.Database) *BoardServiceImpl {
	boardRepo := NewBoardRepository(db)

	return &BoardServiceImpl{
		mode:   mode,
		domain: domain,
		db:     boardRepo,
	}
}
