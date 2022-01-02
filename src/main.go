package main

import (
	"yako/src/model"
)

func main() {
	cpu := model.Cpu{}
	cpuList := cpu.GetResources()
	gpu := model.Gpu{}
	gpuList := gpu.GetResources()
}
