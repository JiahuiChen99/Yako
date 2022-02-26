package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"log"
	"time"
)

var (
	RegistryZnode = "/service_registry"
)

// CreateServiceRegistryZnode will only be ran once
// It creates a non-ephemeral znode in Zookeeper for
// Service Registry at RegistryZnode
func CreateServiceRegistryZnode() {
	// Connect to Zookeeper
	zookeeper, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		log.Fatalln("Error connecting to Apache Zookeeper instance")
	}

	// Create if the service registry znode doesn't exist
	if exists, _, _ := zookeeper.Exists(RegistryZnode); !exists {
		log.Println("Creating Service Registry")
		path, err := zookeeper.Create(RegistryZnode, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln("Error while creating Service Registry znode")
		}

		log.Printf("Service Registry successfully created: %s", path)
	}
}

// RegisterToCluster will register an ephemeral znode for the current YakoAgent
// Called by YakoAgents on start up for YakoMaster service discovery
func RegisterToCluster() {
	// Connect to Zookeeper
	zookeeper, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		log.Fatalln("Error connecting to Apache Zookeeper instance")
	}

	// Create YakoAgent ephemeral znode
	path, err := zookeeper.Create(RegistryZnode+"/n_", []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("Error while adding %s znode to Service Registry", path)
	}

	log.Printf("Registered to the Service Registry: %s", path)
}
