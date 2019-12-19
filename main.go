package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/random"
	"github.com/fyreek/shortener/shorts"
)

// var sMap map[string]*shorts.Shorts
var shortLength = 7

func main() {
	// dbClient, err := db.GetClient("localhost:6379", "", 0)
	// if err != nil {
	// 	fmt.Println("Could not connect to database " + err.Error())
	// 	return
	// }

	// logging.SetLogLevel(config.Configuration.Logging.LogLevel)
	logging.SetLogLevel(0)

	mDB := &db.MongoDB{}
	// err := mDB.Connect(config.Configuration.Database.IP, config.Configuration.Database.Port, config.Configuration.Database.Timeout)
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
	// mDB.SetDatabase(config.Configuration.Database.Database)
	mDB.SetDatabase("shortener")

	random.Seed()
	// sMap = make(map[string]*shorts.Shorts)
	router := gin.Default()
	// TODO: Move routes from main into own file
	router.POST("/url", func(c *gin.Context) {
		var sInput = shorts.Input{}
		err := c.ShouldBind(&sInput)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "could not parse provided data",
			})
			return
		}
		if sInput.URL == "" {
			c.JSON(400, gin.H{
				"message": "no url to shorten was provided",
			})
			return
		}

		// TODO: Implement way to know if actual url

		s := shorts.New(sInput.URL, shortLength)
		// err = dbClient.SetValueStruct(s.Short, s)
		err = s.Save(mDB)
		if err != nil {
			c.JSON(500, gin.H{
				"success": "false",
				"error":   "Could not insert into db: " + err.Error(),
				"data":    "",
			})
			return
		}

		c.JSON(201, gin.H{
			"success": "true",
			"error":   "",
			"data": gin.H{
				"short": s.Short,
			},
		})
	})
	router.GET("/:url", func(c *gin.Context) {
		url := c.Param("url")
		fmt.Println("/:url route:", url)
		if url == "" || len(url) < shortLength {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "No valid url was provided",
				"data":    "",
			})
			return
		}
		// s, ok := sMap[url]
		// val, err := dbClient.GetValue(url)
		s, err := shorts.GetShort(url, mDB)
		if err != nil {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "Url not found",
				"data":    "",
			})
			return
		}

		_ = s.Visit(mDB)
		fmt.Println("This url was used", s.Visits, "times")
		fmt.Println("s got " + string(s.Visits) + " visits")
		// _ = dbClient.SetValueStruct(s.Short, s)
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, s.URL)
		c.Abort()
	})
	log.Fatal(router.Run(":8080"))
}
