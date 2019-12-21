package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/security"
	"github.com/fyreek/shortener/shorts"
)

// DeleteShort will delete a short
func DeleteShort(dBase db.Database) gin.HandlerFunc {
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

		url := c.Param("url")
		if url == "" {
			c.JSON(400, gin.H{
				"success": "false",
				"error":   "No short was provided",
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

		if s.ManageID != manageID {
			c.JSON(404, gin.H{
				"success": "false",
				"error":   "Provided url does not match manage id",
				"data":    "",
			})
			return
		}

		err = s.Delete(dBase)
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
	}
	return gin.HandlerFunc(fn)
}
