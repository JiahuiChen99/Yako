package yako_node_service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"yako/src/grpc/yako"
	"yako/src/model"
)

// YakoNodeServer implements NodeServiceServer interface
type YakoNodeServer struct {
}

func (ns *YakoNodeServer) GetSystemInformation(ctx context.Context, empty *empty.Empty) (*yako.SysInfo, error) {
	sysinfo := model.SysInfo{}
	i := sysinfo.GetResources().(model.SysInfo)

	info := &yako.SysInfo{
		SysName:  i.SysName,
		Machine:  i.Machine,
		Version:  i.Version,
		Release:  i.Release,
		NodeName: i.NodeName,
	}

	return info, nil
}

func (ns *YakoNodeServer) GetSystemCpuInformation(ctx context.Context, empty *empty.Empty) (*yako.CpuList, error) {
	// Get system CPU information
	cpu := model.Cpu{}
	cpuInfo := cpu.GetResources().([]model.Cpu)

	var cpuList []*yako.Cpu
	var coresList []*yako.Core

	// Build CPU list with gRPC structs
	for _, cpu := range cpuInfo {
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
	gpuInfo := gpu.GetResources().([]model.Gpu)

	var gpuList []*yako.Gpu

	// Build GPU list with gRPC structs
	for _, gpu := range gpuInfo {
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
	i := meminfo.GetResources().(model.Memory)

	info := &yako.Memory{
		Total:     i.Total,
		Free:      i.Free,
		FreeSwap:  i.FreeSwap,
		TotalSwap: i.TotalSwap,
	}

	return info, nil
}
