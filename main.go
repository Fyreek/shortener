package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/random"
	"github.com/fyreek/shortener/shorts"
)

var sMap map[string]*shorts.Shorts
var shortLength = 7

func main() {
	random.Seed()
	sMap = make(map[string]*shorts.Shorts)
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
		sMap[s.Short] = s
		c.JSON(201, gin.H{
			"short": s.Short,
		})
	})
	router.GET("/:url", func(c *gin.Context) {
		url := c.Param("url")
		fmt.Println("/:url route:", url)
		if url == "" || len(url) < shortLength {
			c.JSON(400, gin.H{
				"message": "no valid url was provided",
			})
			return
		}
		s, ok := sMap[url]
		if !ok {
			c.JSON(404, gin.H{
				"message": "url was not found",
			})
			return
		}
		s.Visit()
		fmt.Println("This url was used", s.Visits, "times")
		c.Redirect(301, s.URL)
	})
	log.Fatal(router.Run(":8080"))
}
