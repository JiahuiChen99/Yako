package zookeeper

import (
	"github.com/JiahuiChen99/Yako/src/model"
	"github.com/go-zookeeper/zk"
	"log"
)

const (
	MasterRegistryZnode = "/master_registry"
)

var (
	MasterRegistry map[string]*model.ServiceInfo // Masters list
)

// CreateMasterRegistryZnode will only be run once
// It creates a non-ephemeral znode in Zookeeper for
// master nodes tracking
func CreateMasterRegistryZnode() {
	// Create if the mater registry znode doesn't exist
	if exists, _, _ := Zookeeper.Exists(MasterRegistryZnode); !exists {
		log.Println("Creating Master Registry")
		path, err := Zookeeper.Create(MasterRegistryZnode, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln("Error while creating Master Registry znode")
		}
		log.Printf("Master Registry successfully created at %s", path)
	}
}

// RegisterToMasterCluster registers an ephemeral znode for the current YakoMaster
// Called on YakoMaster start up
func RegisterToMasterCluster(yakoMasterAddress string) string {
	// Create YakoMaster ephemeral znode
	path, err := Zookeeper.Create(MasterRegistryZnode+"/m_", []byte(yakoMasterAddress), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("Error while adding %s znode to Master Registry", path)
	}

	log.Printf("Registered to the Master Registry: %s", path)

	return path[len(MasterRegistryZnode+"/"):]
}
