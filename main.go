package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/marianodsr/nura-api/router"
	"github.com/marianodsr/nura-api/storage"
)

var db = storage.GetDbConnection()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	storage.InitDB()
	//scripts.FillDb()
	router.HandleRoutes()
}
