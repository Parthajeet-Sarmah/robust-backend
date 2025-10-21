package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatalf("Failed to listen to port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9000, %v", err)
	}
}
