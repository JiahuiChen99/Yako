package Controller

import (
	"encoding/json"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Cluster returns the cluster schema
func Cluster(c *gin.Context) {
	var response model.Response
	// Get master & service registries and compose response
	response.YakoMasters = zookeeper.MasterRegistry
	response.YakoAgents = zookeeper.ServicesRegistry
	clusterSchema, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, string(clusterSchema))
}
