package main

import (
	"context"
	"dalkak/domain/user"
	"dalkak/internal/app"
	"log"
)

var Mode string

const port = 80

func main() {
	ctx := context.TODO()

	app, err := app.NewApplication(ctx, Mode)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	userService := user.NewUserService()

	err = app.StartServer(port, &userService)
	if err != nil {
		log.Fatal(err)
	}
}
