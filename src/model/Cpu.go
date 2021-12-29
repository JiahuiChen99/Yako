package model

import (
	"bufio"
	"os"
	"regexp"
)

// Core
// Representation of a CPU core
type Core struct {
	Processor uint `json:"processor"` // "processor" field
	CoreID    uint `json:"coreID"`    // "core id" field
}

// Cpu
// Representation of a processor modeled after /proc/cpuinfo
type Cpu struct {
	CpuName  string `json:"cpuName"`  // CPU model name
	CpuCores string `json:"cpuCores"` // Number of CPU cores
	Socket   uint   `json:"socket"`   // "physical id" for multiprocessor systems
	Cores    []Core `json:"cores"`
}

var (
	twoColRegex = regexp.MustCompile("(\t+)?: ")
)

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

	scanner := bufio.NewScanner(f)

	// Scan the file line by line
	for scanner.Scan() {
		// Split each row and get the column name
		if scannerLine := twoColRegex.Split(scanner.Text(), 2); scannerLine != nil {
			switch scannerLine[0] {
			case "model name":

			case "cpu cores":

			case "physical id":

			case "processor":

			case "core id":

			}
		}
	}
}
