package model

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Memory
// Representation of the main memory, the unit is in kB
type Memory struct {
	Total     int `json:"total"`    // "MemTotal" system installed memory
	Free      int `json:"free"`     // "MemFree" system unused memory
	TotalSwap int `json:"swap"`     // "SwapTotal"
	FreeSwap  int `json:"freeSwap"` // "SwapFree"
}

var (
	twoColRegexMemory = regexp.MustCompile(":( +)?")
)

func (m Memory) GetResources() interface{} {
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
				memory.Total = memoryQuantity
			case "MemFree":
				memory.Free = memoryQuantity
			case "SwapTotal":
				memory.TotalSwap = memoryQuantity
			case "SwapFree":
				memory.FreeSwap = memoryQuantity
			}
		}
	}

	return memory
}
