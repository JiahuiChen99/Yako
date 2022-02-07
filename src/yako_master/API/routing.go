package API

import "github.com/gin-gonic/gin"

// AddRoutes specifies yakoAPI routes and its handlers
func AddRoutes(router *gin.Engine) {
	// handler for system heartbeat
	router.GET("/alive", nil)
	// handler for app deploying
	router.GET("/deploy", nil)
	// handler for cluster information
	router.GET("/cluster", nil)
}
