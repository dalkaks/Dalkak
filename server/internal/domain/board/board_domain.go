package boarddomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	boardfactory "dalkak/internal/domain/board/factory"
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	"dalkak/internal/infrastructure/database/dao"
	appdto "dalkak/pkg/dto/app"
	boarddto "dalkak/pkg/dto/board"
	responseutil "dalkak/pkg/utils/response"
)

type BoardDomainService interface {
	CreateBoard(dto *boarddto.CreateBoardDto) (*boardaggregate.BoardAggregate, error)
	ConvertBoardDaos(daos []*dao.BoardDao) ([]*boardaggregate.BoardAggregate, error)
	ConvertBoardDao(dao *dao.BoardDao) (*boardaggregate.BoardAggregate, error)
	GetBoardListProcessingFilter(userInfo *appdto.UserInfo, payload *boarddto.GetBoardListProcessingRequest) *dao.BoardFindFilter
	GetBoardById(userInfo *appdto.UserInfo, id string) (*dao.BoardDao, error)
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
	StatusIncluded := string(boardentity.BoardCreated)
	existCreatedBoard, _, err := service.Database.FindBoardByUserId(
		&dao.BoardFindFilter{
			UserId:         dto.UserInfo.GetUserId(),
			StatusIncluded: &StatusIncluded,
		}, nil)
	if err != nil {
		return nil, err
	}
	if len(existCreatedBoard) > 0 {
		return nil, responseutil.NewAppError(responseutil.ErrCodeConflict, responseutil.ErrMsgBoardExistCreatedStatusBoard)
	}

	factory := boardfactory.NewCreateBoardFactory()
	board, err := factory.CreateBoardAggregateFromDto(dto)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (service *BoardDomainServiceImpl) ConvertBoardDaos(daos []*dao.BoardDao) ([]*boardaggregate.BoardAggregate, error) {
	factory := boardfactory.NewCreateBoardFactory()
	boards, err := factory.CreateBoardAggregatesFromDaos(daos)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (service *BoardDomainServiceImpl) ConvertBoardDao(dao *dao.BoardDao) (*boardaggregate.BoardAggregate, error) {
	factory := boardfactory.NewCreateBoardFactory()
	board, err := factory.CreateBoardAggregateFromDao(dao)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (service *BoardDomainServiceImpl) GetBoardListProcessingFilter(userInfo *appdto.UserInfo, payload *boarddto.GetBoardListProcessingRequest) *dao.BoardFindFilter {
	categoryId := "default"
	statusExcluded := string(boardentity.BoardPosted)

	return &dao.BoardFindFilter{
		UserId:         userInfo.GetUserId(),
		StatusExcluded: &statusExcluded,
		CategoryType:   &payload.CategoryType,
		CategoryId:     &categoryId,
	}
}

func (service *BoardDomainServiceImpl) GetBoardById(userInfo *appdto.UserInfo, id string) (*dao.BoardDao, error) {
	board, err := service.Database.FindBoardById(id)
	if err != nil {
		return nil, err
	}
	if board == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgDataNotFound)
	}

	return board, nil
}