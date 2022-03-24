package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"yako/src/grpc/yako"
	"yako/src/utils/zookeeper"
	"yako/src/yako_node/services"
)

var (
	// TODO: Unregister node when yakoagent is killed
	zn_uuid = ""
)

func main() {
	port := os.Args[1]
	addr := fmt.Sprintf("localhost:%s", port)
	log.Println("Starting YakoAgent at " + addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	zookeeper.NewZookeeper()
	// Attempt to create Service Registry
	zookeeper.CreateServiceRegistryZnode()
	// Add YakoAgent to Service Registry for service discovery
	zn_uuid = zookeeper.RegisterToCluster(fmt.Sprintf("%s", lis.Addr().String()))

	// Start gRPC server
	s := grpc.NewServer()
	yako.RegisterNodeServiceServer(s, &yako_node_service.YakoNodeServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
