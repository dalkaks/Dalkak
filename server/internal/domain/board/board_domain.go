package boarddomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	boardfactory "dalkak/internal/domain/board/factory"
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boarddto "dalkak/pkg/dto/board"
)

type BoardDomainService interface {
	CreateBoard(dto *boarddto.CreateBoardDto) (*boardaggregate.BoardAggregate, error)
}

type BoardDomainServiceImpl struct {
	StaticLink   string
	Database     BoardRepository
	Storage      core.StorageManager
	EventManager core.EventManager
}

func NewBoardDomainService(appConfig *config.AppConfig, database BoardRepository, storage core.StorageManager, eventManager core.EventManager) BoardDomainService {
	return &BoardDomainServiceImpl{
		StaticLink:   appConfig.StaticLink,
		Database:     database,
		Storage:      storage,
		EventManager: eventManager,
	}
}

func (service *BoardDomainServiceImpl) CreateBoard(dto *boarddto.CreateBoardDto) (*boardaggregate.BoardAggregate, error) {
	// 보드 조회 (중복 확인 및 결제 상태에 있는 보드 있는지 확인)
	factory := boardfactory.NewCreateBoardDtoFactory(dto, boardentity.BoardDefaultNft)
	board, err := factory.CreateBoardAggregate()
	if err != nil {
		return nil, err
	}

	return board, nil
}
