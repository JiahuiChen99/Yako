package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"log"
)

const (
	MasterRegistryZnode = "/master_registry"
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
