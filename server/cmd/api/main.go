package main

import (
	"context"
	"dalkak/config"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Mode string

const port = 8080

type application struct {
	Origin   string
	dbClient *dynamodb.Client
}

func main() {
	var app application

  // Load config
	ctx := context.TODO()
	appConfig, err := config.LoadConfig[config.AppConfig](ctx, Mode, "AppConfig")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app.Origin = appConfig.Origin

  // Connect to database
	dbClient, err := app.connectToDB(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
  app.dbClient = dbClient

	log.Printf("Starting server on port %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
