package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/shorts"
)

var shortLength = 7

// GetShort will get a short and redirect you to its underlying url
func GetShort(dBase db.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		s, err := shorts.GetShort(url, dBase)
		if err != nil {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "Url not found",
				"data":    "",
			})
			return
		}

		_ = s.Visit(dBase)
		fmt.Println("This url was used", s.Visits, "times")
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, s.URL)
		c.Abort()
	}
	return gin.HandlerFunc(fn)
}
