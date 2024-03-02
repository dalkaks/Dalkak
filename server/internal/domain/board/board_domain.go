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
	GetBoardListProcessing(userInfo *appdto.UserInfo, payload *boarddto.GetBoardListProcessingRequest) ([]*boardaggregate.BoardAggregate, *dao.ResponsePageDao, error)
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

func (service *BoardDomainServiceImpl) GetBoardListProcessing(userInfo *appdto.UserInfo, payload *boarddto.GetBoardListProcessingRequest) ([]*boardaggregate.BoardAggregate, *dao.ResponsePageDao, error) {
	// todo remove hard coding
	categoryId := "default"
	statusExcluded := string(boardentity.BoardPosted)

	boardDaos, page, err := service.Database.FindBoardByUserId(
		&dao.BoardFindFilter{
			UserId:         userInfo.GetUserId(),
			StatusExcluded: &statusExcluded,
			CategoryType:   &payload.CategoryType,
			CategoryId:     &categoryId,
		}, &dao.RequestPageDao{
			Limit:             payload.Limit,
			ExclusiveStartKey: &payload.ExclusiveStartKey,
		})
	if err != nil {
		return nil, nil, err
	}
	if len(boardDaos) == 0 {
		return []*boardaggregate.BoardAggregate{}, page, nil
	}

	factory := boardfactory.NewCreateBoardFactory()
	boards, err := factory.CreateBoardAggregatesFromDaos(boardDaos)
	if err != nil {
		return nil, nil, err
	}

	return boards, page, nil
}
