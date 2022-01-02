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

// GetResources retrives uname information
func (s SysInfo) GetResources() interface{} {
	var uname syscall.Utsname

	if err := syscall.Uname(&uname); err != nil {
		panic("Couldn't get system information" + err.Error())
	}

	var sysinfo SysInfo
	sysinfo.SysName = parseUname(uname.Sysname)
	sysinfo.NodeName = parseUname(uname.Nodename)
	sysinfo.Release = parseUname(uname.Release)
	sysinfo.Version = parseUname(uname.Version)
	sysinfo.Machine = parseUname(uname.Machine)

	return sysinfo
}

// parseUname
// Parse Utsname fields and returns a stringified version
func parseUname(unameBuff [65]int8) string {
	var byteString [65]byte
	index := 0

	for ; unameBuff[index] != 0; index++ {
		byteString[index] = uint8(unameBuff[index])
	}

	return string(byteString[:index])
}
