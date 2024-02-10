package app

import (
	"context"
	"dalkak/config"
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/interfaces"
	"fmt"
	"log"
	"net/http"
)

type APP struct {
	Origin     string
	Domain     string
	StaticLink string
	Database   *DB
	Storage    *Storage
	KmsSet     *appsecurity.KmsSet
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
	app.StaticLink = appConfig.StaticLink

	kmsSet, err := appsecurity.GetKMSClient(ctx, mode, appConfig.KmsKeyId)
	if err != nil {
		return nil, err
	}
	app.KmsSet = kmsSet

	db, err := NewDB(ctx, mode)
	if err != nil {
		return nil, err
	}
	app.Database = db

	storage, err := NewStorage(ctx, mode, app.StaticLink)
	if err != nil {
		return nil, err
	}
	app.Storage = storage

	return &app, nil
}

func (app *APP) StartServer(port int, userService interfaces.UserService, boardService interfaces.BoardService) error {
	log.Printf("Starting server on port %d", port)

	router := app.NewRouter(userService, boardService)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
