package model

import (
	"bufio"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Memory
// Representation of the main memory, the unit is in kB
type Memory struct {
	Total     uint64 `json:"total"`    // "MemTotal" system installed memory
	Free      uint64 `json:"free"`     // "MemFree" system unused memory
	TotalSwap uint64 `json:"swap"`     // "SwapTotal"
	FreeSwap  uint64 `json:"freeSwap"` // "SwapFree"
}

var (
	twoColRegexMemory = regexp.MustCompile(":( +)?")
)

func (m Memory) GetResources() (interface{}, error) {
	// Open the file
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		panic("Failed to open /proc/meminfo")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var memory Memory

	// Scan the file line by line
	for scanner.Scan() {
		// Split each row and get the column name
		if scannerLine := twoColRegexMemory.Split(scanner.Text(), 2); scannerLine != nil {
			// Remove the memory units
			noUnits := strings.Split(scannerLine[1], " ")[0]
			memoryQuantity, err := strconv.Atoi(noUnits)
			if err != nil {
				panic("Couldn't parse " + scannerLine[0] + "with" + scannerLine[1] + "format.")
			}

			switch scannerLine[0] {
			case "MemTotal":
				memory.Total = uint64(memoryQuantity)
			case "MemAvailable":
				memory.Free = uint64(memoryQuantity)
			case "SwapTotal":
				memory.TotalSwap = uint64(memoryQuantity)
			case "SwapFree":
				memory.FreeSwap = uint64(memoryQuantity)
			}
		}
	}

	return memory, nil
}

// UnmarshallMemory converts protobuf memory model into yako memory model
func UnmarshallMemory(memory *yako.Memory) Memory {
	return Memory{
		Total:     memory.GetTotal(),
		Free:      memory.GetFree(),
		TotalSwap: memory.GetTotalSwap(),
		FreeSwap:  memory.GetFreeSwap(),
	}
}
