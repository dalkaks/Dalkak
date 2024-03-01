package boarddomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	boardfactory "dalkak/internal/domain/board/factory"
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boarddto "dalkak/pkg/dto/board"
	responseutil "dalkak/pkg/utils/response"
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
	createStatus := boardentity.BoardCreated
	existCreatedBoard, _, err := service.Database.FindBoardByUserId(dto.UserInfo.GetUserId(), &createStatus, nil)
	if err != nil {
		return nil, err
	}
	if len(existCreatedBoard) > 0 {
		return nil, responseutil.NewAppError(responseutil.ErrCodeConflict, responseutil.ErrMsgBoardExistCreatedStatusBoard)
	}

	factory := boardfactory.NewCreateBoardDtoFactory(dto)
	board, err := factory.CreateBoardAggregate()
	if err != nil {
		return nil, err
	}

	return board, nil
}
