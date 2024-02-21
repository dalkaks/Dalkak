package application

import (
	"dalkak/config"
	"dalkak/internal/core"
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
}

func NewApplication(appConfig *config.AppConfig, infra *core.Infra) {
	app := &ApplicationImpl{
		AppConfig:    appConfig,
		Database:     infra.Database,
		Storage:      infra.Storage,
		Keymanager:   infra.Keymanager,
		EventManager: infra.EventManager,
		UserDomain:   userdomain.NewUserDomainService(infra.Database, infra.Storage, infra.Keymanager, infra.EventManager),
	}

	app.RegisterUserEventListeners()
}

func (app *ApplicationImpl) SendResponse(responseChan chan<- appdto.Response, data interface{}, err error) {
	responseChan <- appdto.Response{
		Data:  data,
		Error: err,
	}
}
