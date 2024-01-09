package main

import (
	"context"
	"dalkak/config"
	"fmt"
	"log"
	"net/http"
)

var Mode string
const port = 8080

type application struct {
	Origin string
}

func main() {
	var app application

	ctx := context.TODO()

	appConfig, err := config.LoadConfig[config.AppConfig](ctx, Mode, "AppConfig")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app.Origin = appConfig.Origin

	log.Printf("Starting server on port %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
