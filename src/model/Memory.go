package model

import "os"

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

}
