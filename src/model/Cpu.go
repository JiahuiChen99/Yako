package model

import "os"

// Cpu
//Representation of a processor modeled after /proc/cpuinfo
type Cpu struct {
	Id     uint `json:"id"`
	Core   uint `json:"core"`
	Socket uint `json:"socket"`
}

// GetResources Retrieves information related to the system cpu
// currently only Linux is supported
func (c *Cpu) GetResources() {
	// Open the file
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		panic("Failed to open /proc/cpuinfo")
	}
	// Deferred file close
	defer f.Close()

}
