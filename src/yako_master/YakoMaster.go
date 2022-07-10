package main

import (
	"context"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/directory_util"
	"github.com/JiahuiChen99/Yako/src/utils/mqtt"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/JiahuiChen99/Yako/src/yako_master/API"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	addr    = "" // Socket ip + port
	zn_uuid = ""
)

// signalHandler Traps UNIX SIGINT, SIGTERM signals and processes them
func signalHandler(signalChannel chan os.Signal) {
	for {
		// Receive the SIGNAL ID
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
			log.Println("Shutting down YakoMaster")

			// Close all connections to YakoAgents gRPC Server
			for _, agent := range zookeeper.ServicesRegistry {
				agent.GrpcConn.Close()
			}

			// Shutdown YakoAgent gracefully with no errors
			os.Exit(0)
		}
	}
}

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
			Socket:  addr,
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

	go updateAgentsInformation()

	go zookeeper.GetAllServiceAddresses()

	// UNIX signal channel for events
	signalChannel := make(chan os.Signal, 1)
	// Signals to trap
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL)

	// Goroutine for signal processing
	go signalHandler(signalChannel)

	// Connect to MQTT broker for IoT edge YakoAgents
	mqttBrokerIp := os.Args[5]
	mqttBrokerPort := os.Args[6]
	mqtt.ConnectMqttBroker(mqttBrokerIp, mqttBrokerPort)

	// YakoMaster working directory
	directory_util.WorkDir("yakomaster")

	// go routine for gin gonic rest API
	go APIServer()

	for {
		newServiceNodeUUID := <-newService
		newServiceSocket := zookeeper.ServicesRegistry[newServiceNodeUUID]
		log.Println("Call the new service " + newServiceSocket.ServiceInfo.Socket)

		cc, err := grpc.Dial(newServiceSocket.ServiceInfo.Socket, grpc.WithInsecure())
		if err != nil {
			log.Fatalln("Error while dialing the service" + newServiceSocket.ServiceInfo.Socket)
		}

		c := yako.NewNodeServiceClient(cc)

		var sysInfo *yako.SysInfo
		var cpuInfo *yako.CpuList
		var gpuInfo *yako.GpuList
		var memInfo *yako.Memory

		var gpuList []model.Gpu
		var cpuList []model.Cpu

		sysInfo, err = c.GetSystemInformation(context.Background(), &empty.Empty{})
		cpuInfo, err = c.GetSystemCpuInformation(context.Background(), &empty.Empty{})
		gpuInfo, err = c.GetSystemGpuInformation(context.Background(), &empty.Empty{})
		if err != nil {
			log.Println(err)
		} else {
			for _, gpu := range gpuInfo.GetGpuList() {
				gpuList = append(gpuList, model.UnmarshallGPU(gpu))
			}
		}
		memInfo, err = c.GetSystemMemoryInformation(context.Background(), &empty.Empty{})
		for _, cpu := range cpuInfo.GetCpuList() {
			cpuList = append(cpuList, model.UnmarshallCPU(cpu))
		}

		// Update service information to the cluster schema
		if zookeeper.ServicesRegistry[newServiceNodeUUID] != nil {
			newServiceSocket.ServiceInfo.CpuList = cpuList
			newServiceSocket.ServiceInfo.GpuList = gpuList
			newServiceSocket.ServiceInfo.SysInfo = model.UnmarshallSysInfo(sysInfo)
			newServiceSocket.ServiceInfo.Memory = model.UnmarshallMemory(memInfo)
			newServiceSocket.GrpcClient = c
			newServiceSocket.GrpcConn = cc
		}
	}
}

// updateAgentsInformation schedules a timed job. By default, every 10 seconds YakoMaster will
// ask all connected YakoAgents to report back
func updateAgentsInformation() {
	for {
		// Sleep for 10 seconds
		time.Sleep(10 * time.Second)
		for agentID, agent := range zookeeper.ServicesRegistry {
			// Skip if is YakoAgent (IoT)
			if !strings.HasPrefix(agentID, "n") {
				continue
			}
			var err error
			var sysInfo *yako.SysInfo
			var cpuInfo *yako.CpuList
			var gpuInfo *yako.GpuList
			var memInfo *yako.Memory

			var gpuList []model.Gpu
			var cpuList []model.Cpu

			sysInfo, err = agent.GrpcClient.GetSystemInformation(context.Background(), &empty.Empty{})
			cpuInfo, err = agent.GrpcClient.GetSystemCpuInformation(context.Background(), &empty.Empty{})
			gpuInfo, err = agent.GrpcClient.GetSystemGpuInformation(context.Background(), &empty.Empty{})
			if err != nil {
				log.Println(err)
			} else {
				for _, gpu := range gpuInfo.GetGpuList() {
					gpuList = append(gpuList, model.UnmarshallGPU(gpu))
				}
			}
			memInfo, err = agent.GrpcClient.GetSystemMemoryInformation(context.Background(), &empty.Empty{})
			for _, cpu := range cpuInfo.GetCpuList() {
				cpuList = append(cpuList, model.UnmarshallCPU(cpu))
			}

			// Update service information to the cluster schema
			if zookeeper.ServicesRegistry[agentID] != nil {
				agent.ServiceInfo.CpuList = cpuList
				agent.ServiceInfo.GpuList = gpuList
				agent.ServiceInfo.SysInfo = model.UnmarshallSysInfo(sysInfo)
				agent.ServiceInfo.Memory = model.UnmarshallMemory(memInfo)
			}
		}
	}
}
