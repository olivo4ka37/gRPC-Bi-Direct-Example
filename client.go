package main

import (
	"bufio"
	"context"
	"fmt"
	"gRPC-Bi-Direct-Example/chatserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
)

func main() {

	fmt.Println("Enter Serve IP:PORT ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Printf("\nConnecting to server :%s", serverID)

	// connect to grpc server
	conn, err := grpc.NewClient(serverID, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server :: %v", err)
	}
	defer conn.Close()

	// call ChatService to create a stream
	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService: %v", err)
	}

	// implement communications with gRPC server
	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	// blocker
	bl := make(chan bool)
	<-bl

}

// clientHandle

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter client name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
}

// send message
func (ch *clientHandle) sendMessage() {
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)
		if err != nil {
			log.Printf("Failed to send message to server: %v", err)
		}

	}
}

// receive message
func (ch *clientHandle) receiveMessage() {

	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message from server: %v", err)
		}

		//print message to console
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)
	}
}
