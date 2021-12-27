package model

type Resource struct {
	Cpus int `json:"cpus"`
	Gpus int `json:"gpus"`
	Mem  int `json:"mem"`
	Disk int `json:"disk"`
}
