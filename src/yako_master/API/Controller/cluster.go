package Controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yako/src/utils/zookeeper"
)

// Cluster returns the cluster schema
func Cluster(c *gin.Context) {
	clusterSchema, err := json.Marshal(zookeeper.ServicesRegistry)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, string(clusterSchema))
}
