package model

import (
	"syscall"
)

// SysInfo
// Wrapper struct of 'Utsname' returned by 'uname' system call
type SysInfo struct {
	SysName  string `json:"sys_name"`
	NodeName string `json:"node_name"`
	Release  string `json:"release"`
	Version  string `json:"version"`
	Machine  string `json:"machine"`
}

func (s SysInfo) GetResources() interface{} {
	var uname syscall.Utsname

	if err := syscall.Uname(&uname); err != nil {
		panic("Couldn't get system information" + err.Error())
	}

	return nil
}
