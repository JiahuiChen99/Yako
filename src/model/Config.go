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

// Len returns the number of elements in the priority queue
func (h PQNodes) Len() int {
	return len(h)
}

// Less is the comparator used by the heap data structure
// to order my the brownie points property
func (h PQNodes) Less(i, j int) bool {
	return h[i].BrowniePoints > h[j].BrowniePoints
}

// Swap elements
func (h PQNodes) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push adds a new yakoagent to the max heap
func (h *PQNodes) Push(agent interface{}) {
	*h = append(*h, agent.(*YakoAgent))
}

// Pop retrieves the root element which is the
// yakoagent with the most brownie points
func (h *PQNodes) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	*h = old[0 : n-1]
	return item
}
