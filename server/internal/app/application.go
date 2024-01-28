package app

import (
	"context"
	"dalkak/config"
	"dalkak/pkg/interfaces"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type APP struct {
	Origin    string
	Database  *DB
	KmsClient *kms.Client
}

func NewApplication(ctx context.Context, mode string) (*APP, error) {
	var app APP

	// Load config
	appConfig, err := config.LoadConfig[config.AppConfig](ctx, mode, "AppConfig")
	if err != nil {
		return nil, err
	}
	app.Origin = appConfig.Origin

	// Load kms client
	kmsClient, err := config.GetKMSClient(ctx)
	if err != nil {
		return nil, err
	}
	app.KmsClient = kmsClient

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
