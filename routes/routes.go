package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fyreek/shortener/db"
)

// SetupRoutes sets up every route for this service
func SetupRoutes(router *gin.Engine, dBase db.Database) {
	go router.POST("/url", PostShort(dBase))
	go router.GET("/m/:manageId", GetShorts(dBase))
	go router.GET("/l/:url", GetShort(dBase))
	go router.DELETE("/m/:manageId/:url", DeleteShort(dBase))
}
