package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yako/src/model"
	"yako/src/utils/directory_util"
	"yako/src/utils/zookeeper"
	"yako/src/yako_master/API/utils"
)

// UploadApp handles the file that the user wants
// to deploy in the cluster
func UploadApp(c *gin.Context) {
	file, formErr := c.FormFile("app")
	if formErr != nil {
		err := utils.BadRequestError(formErr.Error())
		c.JSON(err.Status, err)
		return
	}

	// Check if YakoMaster's working directory is available
	directory_util.WorkDir("yakomaster")

	// Save the file on the server
	if saveErr := c.SaveUploadedFile(file, "/usr/yakomaster/"+file.Filename); saveErr != nil {
		err := utils.InternalServerError(saveErr.Error())
		c.JSON(err.Status, err)
		return
	}

	// Get the app's resources configuration
	var config model.Config
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	// Compute and find the best nodes to deploy the app
	findYakoAgents(config)

	// File uploaded and stored
	c.JSON(http.StatusOK, map[string]string{"status": "uploaded successfully"})
}

// findYakoAgents calculates and finds the top X best & suitable
// yakoagents where the app could be deployed according to the
// requested resources from the client
// Default X = 3
func findYakoAgents(config model.Config) {
	// Priority queue with max heap to rank higher the nodes
	// with more brownie points

	var browniePoints uint64
	// Loop through all the available yakoagents, computes the
	// brownie points and adds it to a priority queue
	for agentID, agentInfo := range zookeeper.ServicesRegistry {
		// Set brownie points to 0
		browniePoints = 0
		compliesWithCPUCores(agentInfo, config, &browniePoints)
		compliesWithMemory(agentInfo, config, &browniePoints)
		//pq.add(agentID, browniePoints)
	}

	// Select the top X ones to be recommended
	// X is the number of nodes specified by the user

}

// compliesWithCPUCores check if the CPU has the specified cores
// If it does, it adds a brownie point
func compliesWithCPUCores(agent *model.ServiceInfo, config model.Config, browniePoints *uint64) {
	for _, cpu := range agent.CpuList {
		if uint64(len(cpu.Cores)) >= config.CpuCores {
			*browniePoints++
		}
	}
}

// compliesWithMemory check if the specified amount of
// free memory from the system can be allocated for the app
// If it does, it adds a brownie point
func compliesWithMemory(agent *model.ServiceInfo, config model.Config, browniePoints *uint64) {
	if agent.Memory.Free >= config.Memory {
		*browniePoints++
	}
}
