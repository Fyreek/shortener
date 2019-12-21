package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/security"
	"github.com/fyreek/shortener/shorts"
)

var shortLength = 7

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

		manageID, err := security.ParseUUIDString(sInput.ManageID)
		if err != nil && err != security.ErrEmptyUUID {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "Provided manage id was not valid",
				"data":    "",
			})
			return
		}

		s := shorts.New(sInput.URL, manageID, shortLength)
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
				"short":    s.Short,
				"manageId": s.ManageID,
			},
		})
	})
	router.GET("/m/:manageId", func(c *gin.Context) {
		manageID := c.Param("manageId")
		if manageID == "" {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "No manage id was provided",
				"data":    "",
			})
			return
		}
		manageID, err := security.ParseUUIDString(manageID)
		if err != nil && err != security.ErrEmptyUUID {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "Provided manage id was not valid",
				"data":    "",
			})
			return
		}

		limit := c.Query("limit")
		sort := c.Query("sort")
		iLimit := 0
		if limit != "" {
			i, err := strconv.Atoi(limit)
			if err != nil {
				message := "limit has to be number"
				logging.Log(logging.Failure, message, err)
				c.JSON(400, gin.H{
					"success": "false",
					"error":   "Limit has to be number",
					"data":    "",
				})
				return
			}
			iLimit = i
		}
		if iLimit > 100 {
			iLimit = 100
		} else if iLimit <= 0 {
			iLimit = 10
		}

		s, err := shorts.GetShortsForManageID(sort, manageID, iLimit, mDB)
		if err != nil {
			message := "Unknown error"
			logging.Log(logging.Failure, message, err)
			c.JSON(500, gin.H{
				"success": "false",
				"error":   "Could not get shorts",
				"data":    "",
			})
			return
		}

		if len(*s) == 0 {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "No shorts found for this manage id",
				"data":    "",
			})
			return
		}

		response := gin.H{
			"success": "true",
			"error":   "",
			"data":    s,
		}

		byteArray, err := json.Marshal(response)
		if err != nil {
			logging.Log(logging.Failure, "Error on marshalling documents", err)
			c.JSON(500, gin.H{
				"success": "false",
				"error":   "Could not parse loaded data into json",
				"data":    "",
			})
			return
		}

		c.String(201, string(byteArray))
	})
	router.GET("/l/:url", func(c *gin.Context) {
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
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, s.URL)
		c.Abort()
	})
	router.DELETE("/m/:manageId/:url", func(c *gin.Context) {
		manageID := c.Param("manageId")
		if manageID == "" {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "No manage id was provided",
				"data":    "",
			})
			return
		}
		manageID, err := security.ParseUUIDString(manageID)
		if err != nil && err != security.ErrEmptyUUID {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "Provided manage id was not valid",
				"data":    "",
			})
			return
		}

		url := c.Param("url")
		if url == "" {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "No short was provided",
				"data":    "",
			})
			return
		}

		s, err := shorts.GetShort(url, mDB)
		if err != nil {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "Url not found",
				"data":    "",
			})
			return
		}

		if s.ManageID != manageID {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "Provided url does not match manage id",
				"data":    "",
			})
			return
		}

		err = s.Delete(mDB)
		if err != nil {
			logging.Log(logging.Failure, "Error on deleting short from db", err)
			c.JSON(500, gin.H{
				"success": "false",
				"error":   "Could not delete short from database",
				"data":    "",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": "true",
			"error":   "",
			"data":    "",
		})
	})
	log.Fatal(router.Run(":8080"))
}
