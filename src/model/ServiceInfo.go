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
