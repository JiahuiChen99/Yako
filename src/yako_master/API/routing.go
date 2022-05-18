package API

import (
	"github.com/JiahuiChen99/Yako/src/yako_master/API/Controller"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

// AddRoutes specifies yakoAPI routes and its handlers
func AddRoutes(router *gin.Engine) {
	// handler for system heartbeat
	router.GET("/alive", Controller.IsAlive)
	// handler for app deploying
	router.POST("/deploy", Controller.UploadApp)
	// handler for cluster information
	router.GET("/cluster", Controller.Cluster)
	// handler for uploaded apps
	router.GET("/cluster/apps", Controller.GetClusterApps)

	opts := middleware.RedocOpts{
		SpecURL: "./src/yako_master/API/swagger.yaml",
		Title:   "YakoAPI",
	}

	swg := middleware.Redoc(opts, nil)
	router.GET("/docs", gin.WrapH(swg))
	router.GET("./src/yako_master/API/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./"))))
}
