package app

import (
	"context"
	"dalkak/config"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/kmsutils"
	"fmt"
	"log"
	"net/http"
)

type APP struct {
	Origin   string
	Domain   string
	Database *DB
	KmsSet   *kmsutils.KmsSet
}

func NewApplication(ctx context.Context, mode string) (*APP, error) {
	var app APP

	// Load config
	appConfig, err := config.LoadConfig[config.AppConfig](ctx, mode, "AppConfig")
	if err != nil {
		return nil, err
	}
	app.Origin = appConfig.Origin
  app.Domain = appConfig.Domain

	// Load kms client
	kmsSet, err := kmsutils.GetKMSClient(ctx, appConfig.KmsKeyId)
	if err != nil {
		return nil, err
	}
	app.KmsSet = kmsSet

	// Connect to database
	db, err := NewDB(ctx, mode)
	if err != nil {
		return nil, err
	}
	app.Database = db

	return &app, nil
}

func (app *APP) StartServer(port int, userService interfaces.UserService) error {
	log.Printf("Starting server on port %d", port)

	router := app.NewRouter(userService)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
