package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"log"
	"net"
	"os"
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
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Await for connections error", err)
	}
	// TODO: Close listener on SIGINT
	defer ln.Close()

	// Wait for connections and process application deployment
	var app []byte
	var appFrame []byte
	reader := bufio.NewReader(conn)
	for {
		log.Println("Deploy application start...")
		// Start application tranmission
		_, err := reader.Read(appFrame)
		if err != nil && err != io.EOF {
			log.Println("Frame dropped while transferring the application ", err)
		}
		app = append(app, appFrame...)
		if err == io.EOF {
			// Store the application and spin it up
			log.Println("Spinning the application up")
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
