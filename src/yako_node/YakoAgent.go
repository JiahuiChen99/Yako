package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"yako/src/grpc/yako"
	"yako/src/utils/zookeeper"
	"yako/src/yako_node/services"
)

var (
	// Connect to Zookeeper and get singleton
	zkp = zookeeper.NewZookeeper()
	// TODO: Unregister node when yakoagent is killed
	zn_uuid = ""
)

func main() {
	log.Println("Starting YakoAgent")

	lis, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalln(err)
	}

	// Attempt to create Service Registry
	zookeeper.CreateServiceRegistryZnode(zkp)
	// Add YakoAgent to Service Registry for service discovery
	zn_uuid = zookeeper.RegisterToCluster(zkp, fmt.Sprintf("http://%s", lis.Addr().String()))

	// Start gRPC server
	s := grpc.NewServer()
	yako.RegisterNodeServiceServer(s, &yako_node_service.YakoNodeServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
