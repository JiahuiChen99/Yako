package model

// Config is the to be deployed app resources'
// configuration received from the client
type Config struct {
	CpuCores uint64 `json:"cpu_cores"`
	GpuCores uint64 `json:"gpu_cores"`
	Memory   uint64 `json:"memory"`
}

type YakoAgent struct {
	ID            string `json:"id"`
	BrowniePoints uint64 `json:"brownie_points"`
}

type PQNodes []*YakoAgent

func (h PQNodes) Len() int {
	return len(h)
}

func (h PQNodes) Less(i, j int) bool {
	return h[i].BrowniePoints < h[j].BrowniePoints
}

func (h PQNodes) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PQNodes) Push(agent interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, agent.(*YakoAgent))
}

func (h *PQNodes) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	*h = old[0 : n-1]
	return item
}
