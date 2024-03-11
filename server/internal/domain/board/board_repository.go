package boarddomain

import (
	"dalkak/internal/infrastructure/database/dao"
)

type BoardRepository interface {
	FindBoardByUserId(dao *dao.BoardFindFilter, pageDao *dao.RequestPageDao) ([]*dao.BoardDao, *dao.ResponsePageDao, error)
	FindBoardById(boardId string) (*dao.BoardDao, error)
}
