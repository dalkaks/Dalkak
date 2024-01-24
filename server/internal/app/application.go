package app

import (
	"context"
	"dalkak/config"
	"dalkak/domain/user"
	"fmt"
	"log"
	"net/http"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Application struct {
	Origin   string
	Database *DB
}

type DB struct {
	client *dynamodb.Client
}

func NewApplication(ctx context.Context, mode string) (*Application, error) {
	var app Application

	// Load config
	appConfig, err := config.LoadConfig[config.AppConfig](ctx, mode, "AppConfig")
	if err != nil {
		return nil, err
	}
	app.Origin = appConfig.Origin

	// Connect to database
	db, err := NewDB(ctx)
	if err != nil {
		return nil, err
	}
	app.Database = db

	return &app, nil
}

func (app *Application) StartServer(port int, userService *user.UserService) error {
	log.Printf("Starting server on port %d", port)

	router := app.NewRouter(userService)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func NewDB(ctx context.Context) (*DB, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	return &DB{client: dbClient}, nil
}
