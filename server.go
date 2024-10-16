package main

import (
	"log"
	"net"
	"os"

	"gRPC-Bi-Direct-Example/chatserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	// Создаем и регистрируем сервис Health Check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// Регистрируем ваш сервис
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcServer, &cs)

	// Отмечаем сервис как "SERVING" для Health Check
	healthServer.SetServingStatus("chatserver", grpc_health_v1.HealthCheckResponse_SERVING)

	// Запускаем сервер
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Could not serve grpc server @ %v :: %v", Port, err)
	}
}
