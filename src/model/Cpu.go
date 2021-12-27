package model

// Cpu
//Representation of a processor modeled after /proc/cpuinfo
type Cpu struct {
	Id     uint `json:"id"`
	Core   uint `json:"core"`
	Socket uint `json:"socket"`
}
