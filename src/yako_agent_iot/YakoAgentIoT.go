package main

import (
	"fmt"
	"os"
)

const (
	CPU     = "cpu"
	GPU     = "gpu"
	Memory  = "memory"
	SysInfo = "sysinfo"
)

var (
	AgentSocket  = "" // IoT YakoAgent IP + Port
	BrokerSocket = "" // MQTT Broker IP + Port
	topics       = []string{CPU, GPU, Memory, SysInfo}
)

func main() {
	// Get IoT YakoAgent socket
	agentIp := os.Args[1]
	agentPort := os.Args[2]
	AgentSocket = fmt.Sprintf("%s:%s", agentIp, agentPort)

	// Get MQTT Broker socket
	brokerIp := os.Args[3]
	brokerPort := os.Args[4]
	BrokerSocket = fmt.Sprintf("%s:%s", brokerIp, brokerPort)

}
