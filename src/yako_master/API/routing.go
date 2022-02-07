package API

import (
	"github.com/gin-gonic/gin"
	"yako/src/yako_master/API/Controller"
)

// AddRoutes specifies yakoAPI routes and its handlers
func AddRoutes(router *gin.Engine) {
	// handler for system heartbeat
	router.GET("/alive", Controller.IsAlive)
	// handler for app deploying
	router.GET("/deploy", nil)
	// handler for cluster information
	router.GET("/cluster", nil)
}
