package boarddomain

import (
	boardentity "dalkak/internal/domain/board/object/entity"
	"dalkak/internal/infrastructure/database/dao"
)

type BoardRepository interface {
	FindBoardByUserId(userId string, status *boardentity.BoardStatus, pageDao *dao.RequestPageDao) ([]*dao.BoardDao, *dao.ResponsePageDao, error)
}
