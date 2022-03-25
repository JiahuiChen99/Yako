package model

// ServiceInfo struct that represents the data transferred back to
// the client side
type ServiceInfo struct {
	Socket  string
	SysInfo SysInfo
	CpuList []Cpu
	GpuList []Gpu
	Memory  Memory
}

// Response struct represents the final response struct transferred back to
// the client side
type Response struct {
	YakoMasters map[string]*ServiceInfo `json:"yako_masters"`
	YakoAgents  map[string]*ServiceInfo `json:"yako_agents"`
}
