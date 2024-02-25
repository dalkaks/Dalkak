package boarddomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	boarddto "dalkak/pkg/dto/board"
)

type BoardDomainService interface {
	CreateBoard(dto *boarddto.CreateBoardDto) (string, error)
}

type BoardDomainServiceImpl struct {
	StaticLink string
	Database BoardRepository
	Storage core.StorageManager
	EventManager core.EventManager
}

func NewBoardDomainService(appConfig *config.AppConfig, database BoardRepository, storage core.StorageManager, eventManager core.EventManager) BoardDomainService {
	return &BoardDomainServiceImpl{
		StaticLink: appConfig.StaticLink,
		Database: database,
		Storage: storage,
		EventManager: eventManager,
	}
}

func (service *BoardDomainServiceImpl) CreateBoard(dto *boarddto.CreateBoardDto) (string, error) {
	return "", nil
}
