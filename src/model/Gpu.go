package model

import (
	"os"
)

// Gpu
// Abstraction of a GPU
type Gpu struct {
	Major uint `json:"major"`
	Minor uint `json:"minor"`
}

func (g Gpu) GetResources() interface{} {
	// Check if Nvidia drivers are installed
	if _, err := os.Stat("/proc/driver/nvidia"); os.IsNotExist(err) {
		panic("Please install nvidia drivers")
	}

	// Check for GPU information
	gpusDir, err := os.Open("/proc/driver/nvidia/gpus")
	if err != nil {
		panic("Error while parsing GPUs directory " + err.Error())
	}

	gpusDirFiles, err := gpusDir.Readdir(-1)
	if err != nil {
		panic("Error while parsing GPUs directory " + err.Error())
	}
	gpusDir.Close()

	// Get GPU directories names (multi-gpu support)
	for _, file := range gpusDirFiles {

	}

	return nil
}
