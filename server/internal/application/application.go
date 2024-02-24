package application

import (
	"dalkak/config"
	"dalkak/internal/core"
	mediadomain "dalkak/internal/domain/media"
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
}

func NewApplication(appConfig *config.AppConfig, infra *core.Infra) {
	app := &ApplicationImpl{
		AppConfig:    appConfig,
		Database:     infra.Database,
		Storage:      infra.Storage,
		Keymanager:   infra.Keymanager,
		EventManager: infra.EventManager,
		UserDomain:   userdomain.NewUserDomainService(infra.Database, infra.Keymanager, infra.EventManager),
		MediaDomain:	mediadomain.NewMediaDomainService(appConfig, infra.Database, infra.Storage, infra.EventManager),
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
