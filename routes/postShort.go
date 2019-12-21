package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/security"
	"github.com/fyreek/shortener/shorts"
)

// PostShort will create a new short for a new or existing manage id
func PostShort(dBase db.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		err = s.Save(dBase)
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
	}
	return gin.HandlerFunc(fn)
}
