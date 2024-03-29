package application

import (
	"dalkak/config"
	"dalkak/internal/core"
	boarddomain "dalkak/internal/domain/board"
	mediadomain "dalkak/internal/domain/media"
	orderdomain "dalkak/internal/domain/order"
	userdomain "dalkak/internal/domain/user"
	appdto "dalkak/pkg/dto/app"
)

type ApplicationImpl struct {
	AppConfig    *config.AppConfig
	Database     core.DatabaseManager
	Storage      core.StorageManager
	Keymanager   core.KeyManager
	EventManager core.EventManager
	UserDomain   userdomain.UserDomainService
	MediaDomain  mediadomain.MediaDomainService
	BoardDomain  boarddomain.BoardDomainService
	OrderDomain  orderdomain.OrderDomainService
}

func NewApplication(appConfig *config.AppConfig, infra *core.Infra) {
	app := &ApplicationImpl{
		AppConfig:    appConfig,
		Database:     infra.Database,
		Storage:      infra.Storage,
		Keymanager:   infra.Keymanager,
		EventManager: infra.EventManager,
		UserDomain:   userdomain.NewUserDomainService(infra.Database, infra.Keymanager, infra.EventManager),
		MediaDomain:  mediadomain.NewMediaDomainService(appConfig, infra.Database, infra.Storage, infra.EventManager),
		BoardDomain:  boarddomain.NewBoardDomainService(appConfig, infra.Database, infra.Storage, infra.EventManager),
		OrderDomain:  orderdomain.NewOrderDomainService(infra.Database, infra.EventManager),
	}

	app.RegisterUserEventListeners()
	app.RegisterMediaEventListeners()
	app.RegisterBoardEventListeners()
}

func (app *ApplicationImpl) SendResponse(responseChan chan<- appdto.Response, data interface{}, err error) {
	responseChan <- appdto.Response{
		Data:  data,
		Error: err,
	}
}
