package yako_node_service

import (
	"context"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/golang/protobuf/ptypes/empty"
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
