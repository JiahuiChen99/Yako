package model

import (
	"bufio"
	"os"
	"regexp"
)

// Gpu
// Abstraction of a GPU
type Gpu struct {
	Major uint `json:"major"`
	Minor uint `json:"minor"`
}

var (
	twoColRegexGPU = regexp.MustCompile(":( +)?(\t+)?( +)?")
)

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
		f, err := os.Open("/proc/driver/nvidia/gpus/" + file.Name() + "/information")
		if err != nil {
			panic("Failed to open /proc/driver/nvidia/gpus/" + file.Name() + "/information")
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if scannerLine := twoColRegexGPU.Split(scanner.Text(), 2); scannerLine != nil {
				switch scannerLine[0] {
				case "Model":

				case "GPU UUID":

				case "Bus Location":

				case "Device Minor":

				case "IRQ":

				}
			}
		}

		f.Close()
	}

	return nil
}
