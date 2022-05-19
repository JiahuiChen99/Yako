package model

import (
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"google.golang.org/grpc"
)

// ServiceInfo struct that represents the data transferred back to
// the client side
type ServiceInfo struct {
	Socket  string  `json:"socket"`
	SysInfo SysInfo `json:"sys_info"`
	CpuList []Cpu   `json:"cpu_list"`
	GpuList []Gpu   `json:"gpu_list"`
	Memory  Memory  `json:"memory"`
}

// Response struct represents the final response struct transferred back to
// the client side
type Response struct {
	YakoMasters map[string]*ServiceInfo `json:"yako_masters"`
	YakoAgents  map[string]*Agent       `json:"yako_agents"`
}

// Agent struct is used in the Service Registry
type Agent struct {
	ServiceInfo *ServiceInfo           `json:"info"`
	GrpcClient  yako.NodeServiceClient `json:"-"`
	GrpcConn    *grpc.ClientConn       `json:"-"`
}
