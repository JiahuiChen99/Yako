package Controller

import (
	"container/heap"
	"context"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/directory_util"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/JiahuiChen99/Yako/src/yako_master/API/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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
	appPath := "/usr/yakomaster/" + file.Filename
	if saveErr := c.SaveUploadedFile(file, appPath); saveErr != nil {
		err := utils.InternalServerError(saveErr.Error())
		c.JSON(err.Status, err)
		return
	}

	// Get the app's resources configuration
	var config model.Config
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	var iot bool
	if c.Query("iot") == "true" {
		iot = true
	} else {
		iot = false
	}

	// Compute and find the best nodes to deploy the app
	recommendedNodes := findYakoAgents(config, iot)

	// Check if automation is enabled
	if autoDeploy := c.Query("autodeploy"); autoDeploy == "true" {
		// Auto-deploy the app to the best computed node
		yakoAgentID := recommendedNodes[0].ID
		log.Println("Autodeploying application to " + yakoAgentID)
		if !iot {
			// Deploy to a non-IoT agent
			deployStatus := deployApp(&zookeeper.ServicesRegistry[yakoAgentID].GrpcClient, appPath, file.Filename)
			log.Println(deployStatus.Message)
		} else {
			// Deploy to the best IoT agent
			deployAppIoT(yakoAgentID, appPath, file.Filename)
		}
	}

	// File uploaded and stored
	c.JSON(http.StatusOK, recommendedNodes)
}

// findYakoAgents calculates and finds the top X best & suitable
// yakoagents where the app could be deployed according to the
// requested resources from the client
// Default X = 3
func findYakoAgents(config model.Config, iot bool) []*model.YakoAgent {
	// Priority queue with max heap to rank higher the nodes
	// with more brownie points
	agentsCount := 0
	for agentID, _ := range zookeeper.ServicesRegistry {
		if iot && string(agentID[0]) != "n" {
			agentsCount++
		} else if !iot && string(agentID[0]) == "n" {
			agentsCount++
		}
	}
	pq := make(model.PQNodes, agentsCount)

	var browniePoints uint64
	counter := 0
	// Loop through all the available yakoagents, computes the
	// brownie points and adds it to a priority queue
	for agentID, agentInfo := range zookeeper.ServicesRegistry {
		// Filter out the devices depending on the IoT field
		if iot && string(agentID[0]) != "n" {
			// Set brownie points to 0
			browniePoints = 0
			compliesWithCPUCores(agentInfo.ServiceInfo, config, &browniePoints)
			compliesWithMemory(agentInfo.ServiceInfo, config, &browniePoints)
			pq[counter] = &model.YakoAgent{
				ID:            agentID,
				BrowniePoints: browniePoints,
			}
			counter++
		} else if !iot && string(agentID[0]) == "n" {
			// Set brownie points to 0
			browniePoints = 0
			compliesWithCPUCores(agentInfo.ServiceInfo, config, &browniePoints)
			compliesWithMemory(agentInfo.ServiceInfo, config, &browniePoints)
			pq[counter] = &model.YakoAgent{
				ID:            agentID,
				BrowniePoints: browniePoints,
			}
			counter++
		}
	}
	heap.Init(&pq)

	// Select the top X ones to be recommended
	// X is the number of nodes specified by the user
	x := pq.Len()
	if x > 3 {
		x = 3
	}
	recommendedYakoAgents := make([]*model.YakoAgent, x)
	for i := 0; i < x; i++ {
		recommendedYakoAgents[i] = heap.Pop(&pq).(*model.YakoAgent)
	}

	return recommendedYakoAgents
}

// compliesWithCPUCores check if the CPU has the specified cores
// If it does, it adds a brownie point
func compliesWithCPUCores(agent *model.ServiceInfo, config model.Config, browniePoints *uint64) {
	for _, cpu := range agent.CpuList {
		if uint64(len(cpu.Cores)) >= config.SysCpuCores {
			*browniePoints++
		}
	}
}

// compliesWithMemory check if the specified amount of
// free memory from the system can be allocated for the app
// If it does, it adds a brownie point
func compliesWithMemory(agent *model.ServiceInfo, config model.Config, browniePoints *uint64) {
	if agent.Memory.Free >= config.SysMemory {
		*browniePoints++
	}
}

// deployApp opens application binary file, splices it up into chunks of 1KB
// and sends it through gRPC streaming service
func deployApp(c *yako.NodeServiceClient, appPath string, appName string) *yako.DeployStatus {
	file, err := os.Open(appPath)
	if err != nil {
		log.Println("Could not open the file")
		return nil
	}

	stream, err := (*c).DeployAppToAgent(context.Background())

	// 1KB buffer
	buf := make([]byte, 1024)

	// Send application meta-data
	err = stream.Send(&yako.Chunk{
		Data: &yako.Chunk_Info{
			Info: &yako.AppInfo{
				AppName: appName,
			},
		},
	})
	log.Println("Sending application meta-data", appName)
	if err != nil {
		log.Println("Error while sending application meta-data", err)
	}

	// Start transmission
	transmitting := true
	nBytes := 0
	for transmitting {
		nBytes, err = file.Read(buf)

		// End of File
		if err == io.EOF {
			transmitting = false
		}

		err = stream.Send(&yako.Chunk{
			Data: &yako.Chunk_Content{
				Content: buf[:nBytes],
			},
		})
		if err != nil {
			log.Println("Error while deploying the application ", err)
		}
	}

	var deployStatus *yako.DeployStatus
	deployStatus, err = stream.CloseAndRecv()

	return deployStatus
}

// deployAppIoT opens application binary file, splices it up into chunks of 1KB
// and sends it to the IoT YakoAgent
func deployAppIoT(agentSocket string, appPath string, appName string) {
	file, err := os.Open(appPath)
	if err != nil {
		log.Println("Could not open the file")
	}

	conn, err := net.Dial("tcp", agentSocket)
	if err != nil {
		log.Println("Could not dial the agent ", err)
	}

	// Send name
	if _, err := conn.Write([]byte(appName)); err != nil {
		fmt.Println("Could not send the app name")
	}

	// 1KB buffer
	buf := make([]byte, 1024)
	for {
		nBytesRead, err := file.Read(buf)
		// End of File
		if err == io.EOF {
			conn.Close()
			break
		}

		nBytesSent, err := conn.Write(buf[:nBytesRead])
		if err != nil {
			log.Println("Error while deploying the application ", err)
		}

		if nBytesRead != nBytesSent {
			log.Println("Some bytes went missing during the application transmission")
		}
	}
}
