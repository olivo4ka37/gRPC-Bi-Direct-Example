package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}

	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Printf("\nListening @ %v", Port)

	grpcServer := grpc.NewServer()

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Could not serve grpc server @ %v :: %v", Port, err)
	}
}
