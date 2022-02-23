package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"sync"
	"yako/src/grpc/yako"
	"yako/src/utils/directory_util"
	"yako/src/yako_master/API"
)

func APIServer(wg *sync.WaitGroup) {
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

	// Tell wait group once the go routine is ended
	wg.Done()
}

func main() {
	cc, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error al connectar")
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

	// YakoMaster working directory
	directory_util.WorkDir("yakomaster")

	// create new wait group
	wg := new(sync.WaitGroup)

	// add 1 go routines to 'wg' wait group
	wg.Add(1)

	// go routine for gin gonic rest API
	go APIServer(wg)

	// Wait for all go routines to finish
	wg.Wait()
}
