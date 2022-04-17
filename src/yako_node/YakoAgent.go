package main

import (
	"fmt"
	"github.com/JiahuiChen99/Yako/src/grpc/yako"
	"github.com/JiahuiChen99/Yako/src/utils/zookeeper"
	"github.com/JiahuiChen99/Yako/src/yako_node/services"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	znUUID = ""
	server *grpc.Server
)

// signalHandler Traps UNIX SIGINT, SIGTERM signals and processes them
func signalHandler(signalChannel chan os.Signal) {
	for {
		// Receive the SIGNAL ID
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
			// Shutdown gRPC server
			server.GracefulStop()
			// Unregister from service registry
			zookeeper.Unregister(znUUID)
			// Shutdown YakoAgent gracefully with no errors
			os.Exit(0)
		}
	}
}

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	addr := fmt.Sprintf("%s:%s", ip, port)
	log.Println("Starting YakoAgent at " + addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	zookeeper.NewZookeeper()
	// Attempt to create Service Registry
	zookeeper.CreateServiceRegistryZnode()
	// Add YakoAgent to Service Registry for service discovery
	znUUID = zookeeper.RegisterToCluster(fmt.Sprintf("%s", lis.Addr().String()))

	// UNIX signal channel for events
	signalChannel := make(chan os.Signal, 1)
	// Signals to trap
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL)

	// Goroutine for signal processing
	go signalHandler(signalChannel)

	// Start gRPC server
	server = grpc.NewServer()
	yako.RegisterNodeServiceServer(server, &yako_node_service.YakoNodeServer{})

	if err := server.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
