package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ConnectMqttBroker connects to an MQTT Broker and returns the connection
// YakoMaster to listen for subscribed channels
func ConnectMqttBroker(mqttBrokerIp string, mqttBrokerPort string) {
	// Create clients options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%s", mqttBrokerIp, mqttBrokerPort))

	// Connect to broker and obtain a new connection
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error while connecting to the MQTT broker. MQTT service unavailable, please restart the service")
	}
}
