package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/routes"
	"github.com/fyreek/shortener/security"
)

func main() {
	logging.SetLogLevel(0)

	mDB := &db.MongoDB{}
	err := mDB.Connect("localhost", 27017, 10)
	if err != nil {
		logging.Log(logging.Failure, "Could not connect to database:", err)
		return
	}
	b := mDB.IsConnected()
	if !b {
		logging.Log(logging.Failure, "Database connection is not open")
		return
	}
	mDB.SetDatabase("shortener")

	security.RandomSeed()
	router := gin.Default()
	routes.SetupRoutes(router, mDB)
	log.Fatal(router.Run(":8080"))
}
