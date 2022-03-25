package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
	"yako/src/grpc/yako"
	"yako/src/model"
	"yako/src/utils/directory_util"
	"yako/src/utils/zookeeper"
	"yako/src/yako_master/API"
)

var (
	addr    = "" // Socket ip + port
	zn_uuid = ""
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
	err := router.Run(addr)
	if err != nil {
		// TODO: Use logger
		panic("API gin Server could not be started!")
	}
}

// registerMasterSystemInfo gets YakoMaster system information
// and saves it to the regsitry
func registerMasterSystemInfo() {
	// Get all the information
	sf := model.SysInfo{}
	sysInfo := sf.GetResources().(model.SysInfo)
	cpu := model.Cpu{}
	cpuInfo := cpu.GetResources().([]model.Cpu)
	gpu := model.Gpu{}
	gpuInfo := gpu.GetResources().([]model.Gpu)
	mem := model.Memory{}
	memInfo := mem.GetResources().(model.Memory)

	// Add data to the master registry object
	zookeeper.MasterRegistry[zn_uuid].SysInfo = sysInfo
	zookeeper.MasterRegistry[zn_uuid].CpuList = cpuInfo
	zookeeper.MasterRegistry[zn_uuid].GpuList = gpuInfo
	zookeeper.MasterRegistry[zn_uuid].Memory = memInfo
}

func main() {
	// YakoMaster socket address
	port := os.Args[1]
	addr = fmt.Sprintf("localhost:%s", port)

	zookeeper.NewZookeeper()
	// Attempt to create Master Registry
	zookeeper.CreateMasterRegistryZnode()
	// Add YakoMaster to Master Registry
	zn_uuid = zookeeper.RegisterToMasterCluster(addr)

	registerMasterSystemInfo()

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
		log.Println("Call the new service " + newServiceSocket.Socket)

		cc, err := grpc.Dial(newServiceSocket.Socket, grpc.WithInsecure())
		if err != nil {
			log.Fatalln("Error while dialing the service" + newServiceSocket.Socket)
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

		var cpuList []model.Cpu
		for _, cpu := range cpuInfo.GetCpuList() {
			cpuList = append(cpuList, model.UnmarshallCPU(cpu))
		}

		var gpuList []model.Gpu
		for _, gpu := range gpuInfo.GetGpuList() {
			gpuList = append(gpuList, model.UnmarshallGPU(gpu))
		}

		// Update service information to the cluster schema
		if zookeeper.ServicesRegistry[newServiceNodeUUID] != nil {
			zookeeper.ServicesRegistry[newServiceNodeUUID] = &model.ServiceInfo{
				CpuList: cpuList,
				GpuList: gpuList,
				Memory:  model.UnmarshallMemory(memInfo),
				SysInfo: model.UnmarshallSysInfo(sysInfo),
			}
		}
	}
}
