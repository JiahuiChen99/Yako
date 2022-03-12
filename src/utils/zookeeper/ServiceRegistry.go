package zookeeper

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"time"
)

const (
	RegistryZnode = "/service_registry"
)

var (
	ServicesRegistry []string // Service list
	Zookeeper        *zk.Conn // Zookeeper instance
)

// NewZookeeper will create a new singleton of Zookeeper client
func NewZookeeper() {
	// Connect to Zookeeper
	zookeeper, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		log.Fatalln("Error connecting to Apache Zookeeper instance")
	}

	Zookeeper = zookeeper
}

// CreateServiceRegistryZnode will only be ran once
// It creates a non-ephemeral znode in Zookeeper for
// Service Registry at RegistryZnode
func CreateServiceRegistryZnode() {
	// Create if the service registry znode doesn't exist
	if exists, _, _ := Zookeeper.Exists(RegistryZnode); !exists {
		log.Println("Creating Service Registry")
		path, err := Zookeeper.Create(RegistryZnode, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln("Error while creating Service Registry znode")
		}

		log.Printf("Service Registry successfully created: %s", path)
	}
}

// RegisterToCluster will register an ephemeral znode for the current YakoAgent
// Called by YakoAgents on start up for YakoMaster service discovery
func RegisterToCluster(zkp *zk.Conn, yakoNodeAddress string) string {
	// Create YakoAgent ephemeral znode
	path, err := zkp.Create(RegistryZnode+"/n_", []byte(yakoNodeAddress), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("Error while adding %s znode to Service Registry", path)
	}

	log.Printf("Registered to the Service Registry: %s", path)

	return path
}

// updateServices is called whenever an event happens in zookeeper
// it could be either a service disconnection or a new service registry
func updateServices() {
	log.Println("Updating cluster services list")
	GetAllServiceAddresses()
}

// GetAllServiceAddresses consults zookeeper service registry and
// watches for any change
func GetAllServiceAddresses() {
	// Retrieve all znodes from service registry
	yakoagents, _, registryWatch, err := Zookeeper.ChildrenW(RegistryZnode)
	if err != nil {
		log.Fatalln(err.Error())
	}
	go WatchServices(registryWatch)

	var addresses []string

	for _, service := range yakoagents {
		yakoagentPath := fmt.Sprintf("%s/%s", RegistryZnode, service)
		exists, _, err := zkp.Exists(yakoagentPath)
		if err != nil {
			log.Fatalln("Error while trying to check for " + yakoagentPath)
		}

		// Race condition, check if yakoagent exists
		if !exists {
			continue
		}

		// Get yakoagent socket
		socket, _, err := Zookeeper.Get(yakoagentPath)
		if err != nil {
			log.Fatalln("Error while trying to fetch data from " + yakoagentPath)
		}

		// Add the socket to the service registry list
		addresses = append(addresses, string(socket[:]))
	}

	ServicesRegistry = addresses

	// Log cluster available services
	for i, yakoagent := range ServicesRegistry {
		log.Println(fmt.Sprintf("[%d] %s", i, yakoagent))
	}

}

// WatchServices processes watched events
func WatchServices(watch <-chan zk.Event) {
	event := <-watch
	switch event.Type {
	case zk.EventNodeCreated:
		log.Println("A znode has been created for the new service")
	case zk.EventNodeDeleted:
		delete(ServicesRegistry, event.Path)
		fmt.Println("Service at " + event.Path + " znode, has been disconnected")
	case zk.EventNodeChildrenChanged:
		updateServices()
	default:
		fmt.Println(event)
	}
}
