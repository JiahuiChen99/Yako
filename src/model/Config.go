package model

// Config is the to be deployed app resources'
// configuration received from the client
type Config struct {
	CpuCores uint64 `json:"cpu_cores"`
	GpuCores uint64 `json:"gpu_cores"`
	Memory   uint64 `json:"memory"`
}
