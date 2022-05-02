package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

const (
	TopicCpu     = "topic/cpu"
	TopicGpu     = "topic/gpu"
	TopicMemory  = "topic/memory"
	TopicSysInfo = "topic/sysinfo"
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
	switch msg.Topic() {
	case TopicCpu:
		fmt.Println("CPU topic")
	case TopicGpu:
		fmt.Println("GPU topic")
	case TopicMemory:
		fmt.Println("Memory topic")
	case TopicSysInfo:
		fmt.Println("Sysinfo topic")
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
