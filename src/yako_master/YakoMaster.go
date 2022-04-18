package main

import (
	"context"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/directory_util"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/JiahuiChen99/Yako/src/yako_master/API"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
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
	sysInfo, err := sf.GetResources()
	if err != nil {
		log.Println(err)
	}
	cpu := model.Cpu{}
	cpuInfo, err := cpu.GetResources()
	if err != nil {
		log.Println(err)
	}
	gpu := model.Gpu{}
	gpuInfo, err := gpu.GetResources()
	if err != nil {
		log.Println(err)
	}
	mem := model.Memory{}
	memInfo, err := mem.GetResources()
	if err != nil {
		log.Println(err)
	}

	// Add data to the master registry object
	if zookeeper.MasterRegistry[zn_uuid] == nil {
		// Try to type cast
		gInfo, _ := gpuInfo.([]model.Gpu)
		zookeeper.MasterRegistry[zn_uuid] = &model.ServiceInfo{
			CpuList: cpuInfo.([]model.Cpu),
			GpuList: gInfo,
			Memory:  memInfo.(model.Memory),
			SysInfo: sysInfo.(model.SysInfo),
		}
	}
}

func main() {
	// YakoMaster socket address
	ip := os.Args[1]
	port := os.Args[2]
	addr = fmt.Sprintf("%s:%s", ip, port)

	zkIp := os.Args[3]
	zkPort := os.Args[4]
	zookeeper.NewZookeeper(zkIp, zkPort)
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
			log.Println(err)
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
			newServiceSocket.CpuList = cpuList
			newServiceSocket.GpuList = gpuList
			newServiceSocket.SysInfo = model.UnmarshallSysInfo(sysInfo)
			newServiceSocket.Memory = model.UnmarshallMemory(memInfo)
		}
	}
}
