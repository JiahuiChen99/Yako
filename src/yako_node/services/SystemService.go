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

func (ns *YakoNodeServer) GetSystemCpuInformation(ctx context.Context, empty *empty.Empty) (*yako.Cpu, error) {
}

func (ns *YakoNodeServer) GetSystemGpuInformation(ctx context.Context, empty *empty.Empty) (*yako.Gpu, error) {
}

func (ns *YakoNodeServer) GetSystemMemoryInformation(ctx context.Context, empty *empty.Empty) (*yako.Memory, error) {
}

func (ns *YakoNodeServer) mustEmbedUnimplementedNodeServiceServer() {
}
