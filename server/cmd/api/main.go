package main

import (
	"context"
	"dalkak/config"
	"dalkak/internal/application"
	"dalkak/internal/core"
	"dalkak/internal/infrastructure/database"
	"dalkak/internal/infrastructure/eventbus"
	"dalkak/internal/infrastructure/key"
	"dalkak/internal/infrastructure/storage"
	"dalkak/internal/infrastructure/web"
	"log"
)

var Mode string

func main() {
	ctx := context.TODO()

	appConfig, err := config.LoadConfig[config.AppConfig](ctx, Mode, "AppConfig")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	infra, err := initInfra(ctx, Mode, appConfig)
	if err != nil {
		log.Fatalf("Error initializing infrastructure: %v", err)
	}

	application.NewApplication(appConfig, infra)

	router := web.NewRouter(Mode, appConfig.Origin, infra)
	err = router.Listen(":" + config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func initInfra(ctx context.Context, mode string, appConfig *config.AppConfig) (*core.Infra, error) {
	keymanager, err := key.NewKeyManager(ctx, mode, appConfig.KmsKeyId, appConfig.Domain)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDB(ctx, mode)
	if err != nil {
		return nil, err
	}

	storage, err := storage.NewStorage(ctx, mode, appConfig.StaticLink)
	if err != nil {
		return nil, err
	}

	eventmanager := eventbus.NewEventBus()

	return &core.Infra{
		Database:     db,
		Storage:      storage,
		Keymanager:   keymanager,
		EventManager: eventmanager,
	}, nil
}
