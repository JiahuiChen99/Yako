package model

import (
	"bufio"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"os"
	"regexp"
	"strconv"
)

// Core
// Representation of a CPU core
type Core struct {
	Processor uint64 `json:"processor"` // "processor" field
	CoreID    uint64 `json:"coreID"`    // "core id" field
}

// Cpu
// Representation of a processor modeled after /proc/cpuinfo
type Cpu struct {
	CpuName  string `json:"cpuName"`  // CPU model name
	CpuCores uint64 `json:"cpuCores"` // Number of CPU cores
	Socket   uint64 `json:"socket"`   // "physical id" for multiprocessor systems
	Cores    []Core `json:"cores"`
}

var (
	twoColRegexCPU = regexp.MustCompile("(\t+)?: ?")
)

// GetResources Retrieves information related to the system cpu
// currently only Linux is supported
func (c *Cpu) GetResources() (interface{}, error) {
	// Open the file
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		panic("Failed to open /proc/cpuinfo")
	}
	// Deferred file close
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Create a CPU slice for multiprocessor devices
	cpusList := make([]Cpu, 1)
	cpusCount := 0
	var tmpCore Core
	cpuModel := ""
	cpuCores, socket := 0, 0

	// Scan the file line by line
	for scanner.Scan() {
		// Split each row and get the column name
		if scannerLine := twoColRegexCPU.Split(scanner.Text(), 2); scannerLine != nil {
			switch scannerLine[0] {
			case "model name":
				// Add new CPU if there's a new CPU
				if cpuModel != scannerLine[1] && cpuModel != "" {
					cpusCount++
					cpusList = append(cpusList, Cpu{})
				}

				cpuModel = scannerLine[1]
			case "cpu cores":
				cpuCores, err = strconv.Atoi(scannerLine[1])
				if err != nil {
					panic("Error while parsing CPU 'cpu cores' field ")
				}
			case "physical id":
				socket, err = strconv.Atoi(scannerLine[1])
				if err != nil {
					panic("Error while parsing CPU 'cpu cores' field ")
				}
			case "processor":
				processor, err := strconv.Atoi(scannerLine[1])
				if err != nil {
					panic("Error while parsing CPU 'processor' field ")
				}

				tmpCore.Processor = uint64(processor)

			case "core id":
				coreID, err := strconv.Atoi(scannerLine[1])
				if err != nil {
					panic("Error while parsing CPU 'core id' field ")
				}

				tmpCore.CoreID = uint64(coreID)

			case "power management":
				// Store data after parsing the last property
				saveCPU(cpusList, cpusCount, cpuModel, cpuCores, socket, tmpCore)
			}
		}
	}

	return cpusList, err
}

// saveCPU Stores a CPU data if everything is parsed
// It takes the current cpu listing 'cpuList' and stores general information { 'cpuModel', 'cpuCores', 'socket' }
// as well as per-core information from 'tmpCores'
func saveCPU(cpuList []Cpu, cpusCount int, cpuModel string, cpuCores int, socket int, tmpCores Core) {
	// Store general CPU general information
	cpuList[cpusCount].CpuName = cpuModel
	cpuList[cpusCount].CpuCores = uint64(cpuCores)
	cpuList[cpusCount].Socket = uint64(socket)

	// Store CPU cores information
	cpuList[cpusCount].Cores = append(cpuList[cpusCount].Cores, tmpCores)
}

// UnmarshallCPU converts protobuf cpu model into yako cpu model
func UnmarshallCPU(pbCPU *yako.Cpu) Cpu {
	coresList := make([]Core, len(pbCPU.Cores))
	for _, core := range pbCPU.GetCores() {
		coresList = append(coresList, Core{
			Processor: core.GetProcessor(),
			CoreID:    core.GetCoreID(),
		})
	}

	return Cpu{
		Socket:   pbCPU.GetSocket(),
		CpuCores: pbCPU.GetCpuCores(),
		CpuName:  pbCPU.GetCpuName(),
		Cores:    coresList,
	}
}
