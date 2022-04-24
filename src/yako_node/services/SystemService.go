package yako_node_service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/directory_util"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
	"log"
	"os"
	"os/exec"
)

// YakoNodeServer implements NodeServiceServer interface
type YakoNodeServer struct {
}

func (ns *YakoNodeServer) GetSystemInformation(ctx context.Context, empty *empty.Empty) (*yako.SysInfo, error) {
	sysinfo := model.SysInfo{}
	i, err := sysinfo.GetResources()

	if err != nil {
		return nil, err
	}

	sInfo := i.(model.SysInfo)
	info := &yako.SysInfo{
		SysName:  sInfo.SysName,
		Machine:  sInfo.Machine,
		Version:  sInfo.Version,
		Release:  sInfo.Release,
		NodeName: sInfo.NodeName,
	}

	return info, nil
}

func (ns *YakoNodeServer) GetSystemCpuInformation(ctx context.Context, empty *empty.Empty) (*yako.CpuList, error) {
	// Get system CPU information
	cpu := model.Cpu{}
	cpuInfo, err := cpu.GetResources()

	if err != nil {
		return nil, err
	}

	var cpuList []*yako.Cpu
	var coresList []*yako.Core

	// Build CPU list with gRPC structs
	for _, cpu := range cpuInfo.([]model.Cpu) {
		// Build CPU cores list with gRPC structs
		for i := range cpu.Cores {
			coresList = append(coresList, &yako.Core{
				CoreID:    cpu.Cores[i].CoreID,
				Processor: cpu.Cores[i].Processor,
			})
		}

		cpuList = append(cpuList, &yako.Cpu{
			CpuName:  cpu.CpuName,
			CpuCores: cpu.CpuCores,
			Socket:   cpu.Socket,
			Cores:    coresList,
		})
	}

	// Build CpuList gRPC struct
	info := &yako.CpuList{
		CpuList: cpuList,
	}

	return info, nil
}

func (ns *YakoNodeServer) GetSystemGpuInformation(ctx context.Context, empty *empty.Empty) (*yako.GpuList, error) {
	gpu := model.Gpu{}
	gpuInfo, err := gpu.GetResources()

	if err != nil {
		return nil, err
	}

	gpuInfoData := gpuInfo.([]model.Gpu)
	var gpuList []*yako.Gpu

	// Build GPU list with gRPC structs
	for _, gpu := range gpuInfoData {
		gpuList = append(gpuList, &yako.Gpu{
			GpuName: gpu.GpuName,
			GpuID:   gpu.GpuID,
			BusID:   gpu.BusID,
			IRQ:     gpu.IRQ,
			Major:   gpu.Major,
			Minor:   gpu.Minor,
		})
	}

	info := &yako.GpuList{
		GpuList: gpuList,
	}

	return info, nil
}

func (ns *YakoNodeServer) GetSystemMemoryInformation(ctx context.Context, empty *empty.Empty) (*yako.Memory, error) {
	meminfo := model.Memory{}
	i, err := meminfo.GetResources()

	if err != nil {
		return nil, err
	}

	mInfo := i.(model.Memory)

	info := &yako.Memory{
		Total:     mInfo.Total,
		Free:      mInfo.Free,
		FreeSwap:  mInfo.FreeSwap,
		TotalSwap: mInfo.TotalSwap,
	}

	return info, nil
}

// DeployAppToAgent implements the Upload method of the YakoNodeServer
// interface which is responsible for receiving a stream of
// chunks that form a complete application to spin up.
func (ns *YakoNodeServer) DeployAppToAgent(stream yako.NodeService_DeployAppToAgentServer) error {
	// Get application meta-data
	meta, metaErr := stream.Recv()
	if metaErr != nil {
		return metaErr
	}
	appName := meta.GetInfo().GetAppName()
	log.Println("Received application meta-data", appName)

	appData := bytes.Buffer{}
	// While there are app's chunks coming
	for {
		// Receive stream
		req, err := stream.Recv()
		if err != nil {
			// Finish receiving application byte stream
			if err == io.EOF {
				break
			}
			return err
		}

		// Get byte data and compose app's binary file from the stream
		chunk := req.GetContent()
		appData.Write(chunk)
	}

	// Check for YakoAgent working directory, create if it doesn't exist
	directory_util.WorkDir("yakoagent")

	// Save the application
	deployedApp, err := os.Create("/usr/yakoagent/" + appName)

	if err != nil {
		log.Println(fmt.Sprintf("Could not create application file: %s", err))
	}

	// Write the binary application to the file system
	_, err = appData.WriteTo(deployedApp)
	if err != nil {
		log.Println(fmt.Sprintf("Could not write application file: %s", err))
	}

	// Spin up the application
	cmd := exec.Command("/usr/yakoagent/" + appName)
	err = cmd.Start()
	if err != nil {
		log.Println("Error: Could not start", err)
	} else {
		log.Println("Application up - PID: ", cmd.Process.Pid)
	}

	// Transmission finished successfully with no errors
	err = stream.SendAndClose(&yako.DeployStatus{
		Message: "Successfully deployed",
		Code:    yako.DeployStatusCode_Ok,
	})

	return err
}
