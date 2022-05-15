package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/directory_util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
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

	// MQTT broker configuration
	opts := mqtt.NewClientOptions()
	opts.AddBroker(BrokerSocket)
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error while connecting to the MQTT broker.")
	}

	// Regular reports to the MQTT broker
	go timedReport(client)

	// IoT YakoAgent Server listens for application deployment
	serve()
}

// serve starts a socket in listening mode and awaits for new deployment
// connections.
func serve() {
	ln, err := net.Listen("tcp", AgentSocket)
	if err != nil {
		log.Fatalln("Could not create IoT YakoAgent Server ", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Await for connections error", err)
		}
		// TODO: Close listener on SIGINT

		appName := make([]byte, 256)
		//readName := false
		nameBytes := 0
		// Read app name
		nameBytes, err = conn.Read(appName[:])
		if err != nil {
			log.Println("Error while receiving app name", err)
		} else {
			fmt.Println("App name", string(appName[:nameBytes]))
		}

		// Wait for connections and process application deployment
		appFrame := make([]byte, 1024)
		appData := bytes.Buffer{}
		reader := bufio.NewReader(conn)
		for {
			// Start application receptiion
			nBytesRead, err := reader.Read(appFrame)
			if err != nil && err != io.EOF {
				log.Println("Frame dropped while transferring the application ", err)
			}
			if err == io.EOF {
				// Store the application and spin it up
				log.Println("Spinning the application up")
				directory_util.WorkDir("yakoagentiot")
				deployedApp, err := os.Create("/usr/yakoagentiot/" + string(appName[:nameBytes]))

				if err != nil {
					log.Println(fmt.Sprintf("Could not create application file: %s", err))
				}

				// Write the binary application to the file system
				_, err = appData.WriteTo(deployedApp)
				if err != nil {
					log.Println(fmt.Sprintf("Could not write application file: %s", err))
				}

				err = deployedApp.Chmod(0710)
				if err != nil {
					log.Println("Could not change application permissions")
				}

				// Close the application file descriptor after writing & chmoding
				err = deployedApp.Close()
				if err != nil {
					log.Println("Error while closing the application file descriptor")
				}

				// Spin up the application
				cmd := exec.Command("/usr/yakoagentiot/" + string(appName[:nameBytes]))
				err = cmd.Start()
				if err != nil {
					log.Println("Error: Could not start", err)
				} else {
					// TODO: Report the PID back to the YakoMaster
					log.Println("Application up - PID: ", cmd.Process.Pid)
				}
				break
			}
			appData.Write(appFrame[:nBytesRead])
		}
	}
}

// timedReport schedules a timed system resources information report to the MQTT broker
func timedReport(client mqtt.Client) {
	var topic string
	for {
		log.Println("Reporting data to the broker")
		for _, mqttTopic := range topics {
			var data interface{}
			var err error
			switch mqttTopic {
			case CPU:
				cpu := model.Cpu{}
				data, err = cpu.GetResources()
				if err != nil {
					log.Println("Error retrieving system's CPU data")
				}
				topic = fmt.Sprintf("topic/%s/%s", AgentSocket, CPU)
			case GPU:
				data, err = model.Gpu{}.GetResources()
				if err != nil {
					log.Println("Error retrieving system's GPU data")
				}
				topic = fmt.Sprintf("topic/%s/%s", AgentSocket, GPU)
			case Memory:
				data, err = model.Memory{}.GetResources()
				if err != nil {
					log.Println("Error retrieving system's memory data")
				}
				topic = fmt.Sprintf("topic/%s/%s", AgentSocket, Memory)
			case SysInfo:
				data, err = model.SysInfo{}.GetResources()
				if err != nil {
					log.Println("Error retrieving system's information data")
				}
				topic = fmt.Sprintf("topic/%s/%s", AgentSocket, SysInfo)
			}
			// Publish the data
			pubTopic(client, topic, data)
		}
		// Sleep for 10 seconds
		time.Sleep(10 * time.Second)
	}
}

// pubTopic encodes and publish a message with topic to the MQTT broker server
func pubTopic(client mqtt.Client, topic string, data interface{}) {
	// Encode data to json format
	message, err := json.Marshal(data)
	if err != nil {
		log.Print("Error while trying to encode the data", err)
	}
	token := client.Publish(topic, 0, true, message)
	token.Wait()
}
