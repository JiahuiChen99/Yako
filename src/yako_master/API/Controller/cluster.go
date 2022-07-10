package Controller

import (
	"encoding/json"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

// GetClusterApps returns a list of applications that have
// been uploaded to the system
func GetClusterApps(c *gin.Context) {
	// Read from working directory
	apps, err := ioutil.ReadDir("/usr/yakomaster/")
	if err != nil {
		log.Println("Could not read from working directory /usr/yakomaster")
	}
	var appsNames []string
	for _, app := range apps {
		appsNames = append(appsNames, app.Name())
	}

	c.JSON(http.StatusOK, appsNames)
}
