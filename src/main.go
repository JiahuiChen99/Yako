package main

import (
	"google.golang.org/grpc"
	"log"
	"yako/src/grpc/yako"
)

func main() {
	cc, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error al connectar")
	}
	defer cc.Close()

	c := yako.NewNodeServiceClient(cc)

}
