package model

// Memory
// Representation of the main memory, the unit is in kB
type Memory struct {
	Total     int `json:"total"`    // "MemTotal" system installed memory
	Free      int `json:"free"`     // "MemFree" system unused memory
	TotalSwap int `json:"swap"`     // "SwapTotal"
	FreeSwap  int `json:"freeSwap"` // "SwapFree"
}
