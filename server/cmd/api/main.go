package main

import (
	"context"
	"dalkak/domain/user"
	"dalkak/internal/app"
	"dalkak/pkg/interfaces"
	"log"
)

var Mode string

const port = 80

func main() {
	ctx := context.TODO()

	appInstance, err := app.NewApplication(ctx, Mode)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	var db interfaces.Database = appInstance.Database

	userService := user.NewUserService(Mode, appInstance.Domain, db, appInstance.KmsSet)

	err = appInstance.StartServer(port, userService)
	if err != nil {
		log.Fatal(err)
	}
}
