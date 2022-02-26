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
