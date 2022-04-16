package API

import (
	"github.com/JiahuiChen99/Yako/src/yako_master/API/Controller"
	"github.com/gin-gonic/gin"
)

// AddRoutes specifies yakoAPI routes and its handlers
func AddRoutes(router *gin.Engine) {
	// handler for system heartbeat
	router.GET("/alive", Controller.IsAlive)
	// handler for app deploying
	router.POST("/deploy", Controller.UploadApp)
	// handler for cluster information
	router.GET("/cluster", Controller.Cluster)
}
