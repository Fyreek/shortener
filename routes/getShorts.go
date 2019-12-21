package routes

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/security"
	"github.com/fyreek/shortener/shorts"
)

// GetShorts will get all shorts for a the provided manage id
func GetShorts(dBase db.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		s, err := shorts.GetShortsForManageID(sort, manageID, iLimit, dBase)
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
	}
	return gin.HandlerFunc(fn)
}
