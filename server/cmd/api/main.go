package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct{
  Origin string
}

func main() {
	var app application
  flag.StringVar(&app.Origin, "origin", "http://dev-api.dalkak.com", "the origin url")

	log.Printf("Starting server on port %d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
