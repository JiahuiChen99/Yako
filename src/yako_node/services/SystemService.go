package yako_node_service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"yako/src/grpc/yako"
)

// YakoNodeServer implements NodeServiceServer interface
type YakoNodeServer struct {
}

func (ns *YakoNodeServer) GetSystemInformation(ctx context.Context, empty *empty.Empty) (*yako.SysInfo, error) {
}

func (ns *YakoNodeServer) GetSystemCpuInformation(ctx context.Context, empty *empty.Empty) (*yako.CpuList, error) {
}

func (ns *YakoNodeServer) GetSystemGpuInformation(ctx context.Context, empty *empty.Empty) (*yako.GpuList, error) {
}

func (ns *YakoNodeServer) GetSystemMemoryInformation(ctx context.Context, empty *empty.Empty) (*yako.Memory, error) {
}
