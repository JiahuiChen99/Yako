package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

const (
	CPU          = "cpu"
	GPU          = "gpu"
	Memory       = "memory"
	SysInfo      = "sysinfo"
	TopicCpu     = "topic/+/" + CPU
	TopicGpu     = "topic/+/" + GPU
	TopicMemory  = "topic/+/" + Memory
	TopicSysInfo = "topic/+/" + SysInfo
)

var (
	topics = []string{TopicCpu, TopicGpu, TopicMemory, TopicSysInfo}
)

// ConnectMqttBroker connects to an MQTT Broker and returns the connection
// YakoMaster to listen for subscribed channels
func ConnectMqttBroker(mqttBrokerIp string, mqttBrokerPort string) {
	// Create clients options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%s", mqttBrokerIp, mqttBrokerPort))

	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetOnConnectHandler(connectionHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)

	// Connect to broker and obtain a new connection
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error while connecting to the MQTT broker. MQTT service unavailable, please restart the service")
	}

	for _, topic := range topics {
		subToTopic(client, topic)
	}
}

// messageHandler Callback handler for subscribed events. Processes the event
// according to the message topic
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	// Topic parsing topic/<agent_socket>/<topic_name>
	mqttTopic := strings.Split(msg.Topic(), "/")
	agentSocket := mqttTopic[1]
	switch mqttTopic[2] {
	case CPU:
		var cpu []model.Cpu
		if err := json.Unmarshal(msg.Payload(), &cpu); err != nil {
			log.Println("Err", err)
		}
		updateRegistry(agentSocket, cpu)
	case GPU:
		var gpu []model.Gpu
		if err := json.Unmarshal(msg.Payload(), &gpu); err != nil {
			log.Println("Err", err)
		}
		updateRegistry(agentSocket, gpu)
	case Memory:
		var memory model.Memory
		if err := json.Unmarshal(msg.Payload(), &memory); err != nil {
			log.Println("Err", err)
		}
		updateRegistry(agentSocket, memory)
	case SysInfo:
		var sysinfo model.SysInfo
		if err := json.Unmarshal(msg.Payload(), &sysinfo); err != nil {
			log.Println("Err", err)
		}
		updateRegistry(agentSocket, sysinfo)
	}
}

// connectionHandler on connection to MQTT broker handler
func connectionHandler(client mqtt.Client) {
	if client.IsConnected() {
		log.Println("Connected to MQTT broker")
	}
}

// connectionLostHandler on connection lost to MQTT broker handler
func connectionLostHandler(client mqtt.Client, err error) {
	if !client.IsConnected() {
		log.Printf("Connection lost: %v\n", err)
	}
}

// subToTopic subscribes to the specified MQTT topic
func subToTopic(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Println(fmt.Sprintf("Subscribed to topic: %s", topic))
}

// updateRegistry creates a new entry in the service registry if it did not exist previously
// otherwise it updates the existing information
func updateRegistry(agentSocket string, data interface{}) {
	// Add a new entry if it doesn't exist in the registry
	if zookeeper.ServicesRegistry[agentSocket] == nil {
		var info model.ServiceInfo
		switch data.(type) {
		case []model.Cpu:
			info.CpuList = data.([]model.Cpu)
		case []model.Gpu:
			info.GpuList = data.([]model.Gpu)
		case model.Memory:
			info.Memory = data.(model.Memory)
		case model.SysInfo:
			info.SysInfo = data.(model.SysInfo)
		}
		// Save the information
		zookeeper.ServicesRegistry[agentSocket] = &model.Agent{
			ServiceInfo: &info,
		}
	} else {
		// Update the service information
		agent := zookeeper.ServicesRegistry[agentSocket]
		switch data.(type) {
		case []model.Cpu:
			agent.ServiceInfo.CpuList = data.([]model.Cpu)
		case []model.Gpu:
			agent.ServiceInfo.GpuList = data.([]model.Gpu)
		case model.Memory:
			agent.ServiceInfo.Memory = data.(model.Memory)
		case model.SysInfo:
			agent.ServiceInfo.SysInfo = data.(model.SysInfo)
		}
	}
}
