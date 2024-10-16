package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"time"
)

func main() {
	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Создаем клиента Health Check
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Бесконечный цикл с интервалом в 2 секунды
	for {
		// Выполняем Health Check запрос
		req := &grpc_health_v1.HealthCheckRequest{Service: "chatserver"}
		res, err := healthClient.Check(context.Background(), req)

		// Проверяем ответ
		if err != nil {
			log.Printf("Health Check failed: %v", err)
		} else {
			status := res.Status.String()
			log.Printf("Health Check status: %s", status)
			break
		}

		// Ждем 2 секунды перед следующим запросом
		time.Sleep(2 * time.Second)
	}
}
