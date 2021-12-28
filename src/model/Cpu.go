package model

import (
	"bufio"
	"os"
	"regexp"
)

// Cpu
//Representation of a processor modeled after /proc/cpuinfo
type Cpu struct {
	Id     uint `json:"id"`
	Core   uint `json:"core"`
	Socket uint `json:"socket"`
}

var (
	twoColRegex = regexp.MustCompile("\t+: ")
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

		}
	}
}