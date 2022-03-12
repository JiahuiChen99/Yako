package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"yako/src/grpc/yako"
	"yako/src/utils/directory_util"
	"yako/src/utils/zookeeper"
	"yako/src/yako_master/API"
)

func APIServer() {
	// Default gin router with default middleware:
	// logger & recovery
	router := gin.Default()

	// Add CORS support
	router.Use(cors.Default())

	// Attach routes to gin router
	API.AddRoutes(router)

	// TODO: Use environment variables or secrets managers like Hashicorp Vault
	err := router.Run(":8001")
	if err != nil {
		// TODO: Use logger
		panic("API gin Server could not be started!")
	}
}

func main() {
	zookeeper.NewZookeeper()

	// Channel for services registration events
	newService := make(chan string)
	zookeeper.NewServiceChan = newService

	go zookeeper.GetAllServiceAddresses()

	// YakoMaster working directory
	directory_util.WorkDir("yakomaster")

	// go routine for gin gonic rest API
	go APIServer()

	for {
		newServiceNodeUUID := <-newService
		newServiceSocket := zookeeper.ServicesRegistry[newServiceNodeUUID]
		log.Println("Call the new service " + newServiceSocket)

		cc, err := grpc.Dial(newServiceSocket, grpc.WithInsecure())
		if err != nil {
			log.Fatalln("Error while dialing the service" + newServiceSocket)
		}
		defer cc.Close()

		c := yako.NewNodeServiceClient(cc)

		var sysInfo *yako.SysInfo
		var cpuInfo *yako.CpuList
		var gpuInfo *yako.GpuList
		var memInfo *yako.Memory

		sysInfo, err = c.GetSystemInformation(context.Background(), &empty.Empty{})
		cpuInfo, err = c.GetSystemCpuInformation(context.Background(), &empty.Empty{})
		gpuInfo, err = c.GetSystemGpuInformation(context.Background(), &empty.Empty{})
		memInfo, err = c.GetSystemMemoryInformation(context.Background(), &empty.Empty{})

		if err != nil {
			fmt.Println("Error")
		}

		fmt.Println(sysInfo)
		fmt.Println(cpuInfo)
		fmt.Println(gpuInfo)
		fmt.Println(memInfo)
	}
}
