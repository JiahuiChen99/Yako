package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"yako/src/grpc/yako"
)

func APIServer(wg *sync.WaitGroup) {
	r := gin.Default()

	// TODO: Use environment variables or secrets managers like Hashicorp Vault
	err := r.Run(":8001")
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
}
