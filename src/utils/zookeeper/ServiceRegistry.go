package zookeeper

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"time"
	"yako/src/model"
)

const (
	RegistryZnode = "/service_registry"
)

var (
	Zookeeper        *zk.Conn                      // Zookeeper instance
	NewServiceChan   chan string                   // Handle Service Channel
	ServicesRegistry map[string]*model.ServiceInfo // Service list
)

// NewZookeeper will create a new singleton of Zookeeper client
func NewZookeeper() {
	// Connect to Zookeeper
	zookeeper, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		log.Fatalln("Error connecting to Apache Zookeeper instance")
	}

	Zookeeper = zookeeper
	ServicesRegistry = make(map[string]*model.ServiceInfo)
	MasterRegistry = make(map[string]*model.ServiceInfo)
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
func RegisterToCluster(yakoNodeAddress string) string {
	// Create YakoAgent ephemeral znode
	path, err := Zookeeper.Create(RegistryZnode+"/n_", []byte(yakoNodeAddress), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("Error while adding %s znode to Service Registry", path)
	}

	log.Printf("Registered to the Service Registry: %s", path)

	return path[len(RegistryZnode+"/"):]
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

	var exists bool
	var yakoagentWatch <-chan zk.Event

	for _, service := range yakoagents {
		yakoagentPath := fmt.Sprintf("%s/%s", RegistryZnode, service)
		exists, _, yakoagentWatch, err = Zookeeper.ExistsW(yakoagentPath)
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
		socketPath := string(socket[:])
		if ServicesRegistry[service] == nil {
			// Store socket path in the registry
			ServicesRegistry[service] = &model.ServiceInfo{
				Socket: socketPath,
			}
			// A new service has connected
			NewServiceChan <- service
		}
		go WatchServices(yakoagentWatch)
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

// Unregister receives a zookeeper node UUID, unregisters it from the service registry
// and closes the established connection
func Unregister(znUUID string) {
	log.Println(fmt.Sprintf("Unregistering %s from Service Registry", znUUID))

	// Unregister node from Service Registry
	if err := Zookeeper.Delete(RegistryZnode+"/"+znUUID, -1); err != nil {
		log.Fatalln(err)
	}

	// Close Zookeeper connection
	Zookeeper.Close()
}
