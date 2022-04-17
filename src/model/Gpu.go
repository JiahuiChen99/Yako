package model

import (
	"bufio"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"os"
	"regexp"
	"strconv"
)

// Gpu
// Abstraction of a GPU
// TODO: Get CUDA Cores
type Gpu struct {
	GpuName string `json:"gpuName"` // GPU model name
	GpuID   string `json:"gpuID"`   // "GPU UUID"
	BusID   string `json:"busID"`   // "Bus Location" PCIe bus ID
	IRQ     uint64 `json:"IRQ"`     // "IRQ" GPU Interrupt lane
	Major   uint64 `json:"major"`   //
	Minor   uint64 `json:"minor"`   // "Device Minor" for /dev/nvidia<minor> character device
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

	// Create a GPU slice for multi-GPU devices
	gpuList := make([]Gpu, 1)

	// Get GPU directories names (multi-gpu support)
	for gpusCount, file := range gpusDirFiles {
		f, err := os.Open("/proc/driver/nvidia/gpus/" + file.Name() + "/information")
		if err != nil {
			panic("Failed to open /proc/driver/nvidia/gpus/" + file.Name() + "/information")
		}

		// Add new GPU if there's multiple
		if gpusCount >= 1 {
			gpuList = append(gpuList, Gpu{})
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if scannerLine := twoColRegexGPU.Split(scanner.Text(), 2); scannerLine != nil {
				switch scannerLine[0] {
				case "Model":
					gpuList[gpusCount].GpuName = scannerLine[1]
				case "GPU UUID":
					gpuList[gpusCount].GpuID = scannerLine[1]
				case "Bus Location":
					gpuList[gpusCount].BusID = scannerLine[1]
				case "Device Minor":
					minorNumber, err := strconv.Atoi(scannerLine[1])
					if err != nil {
						panic("Error while parsing GPU 'Device Minor'")
					}
					gpuList[gpusCount].Minor = uint64(minorNumber)
				case "IRQ":
					irqNumber, err := strconv.Atoi(scannerLine[1])
					if err != nil {
						panic("Error while parsing GPU 'IRQ'")
					}
					gpuList[gpusCount].IRQ = uint64(irqNumber)
				}
			}
		}

		f.Close()
	}

	return gpuList
}

// UnmarshallGPU converts protobuf gpu model into yako gpu model
func UnmarshallGPU(pbGPU *yako.Gpu) Gpu {
	return Gpu{
		GpuID:   pbGPU.GetGpuID(),
		GpuName: pbGPU.GetGpuName(),
		BusID:   pbGPU.GetBusID(),
		Major:   pbGPU.GetMajor(),
		Minor:   pbGPU.GetMinor(),
		IRQ:     pbGPU.GetIRQ(),
	}
}
