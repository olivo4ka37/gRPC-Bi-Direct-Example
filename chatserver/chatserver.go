package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {

	clientUNiqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// receive messages - init a go routine
	go receiveFromStream(csi, clientUNiqueCode, errch)

	// send messages - init a go routine
	go sendToStream(csi, clientUNiqueCode, errch)

	return <-errch
}

// receive messages
func receiveFromStream(csi_ Services_ChatServiceServer, clientUNiqueCode int, errch chan error) {
	for {
		mssg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch <- err
		} else {
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        mssg.Name,
				MessageBody:       mssg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUNiqueCode,
			})

			messageHandleObject.mu.Unlock()

			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}
	}
}

// send message
func sendToStream(csi_ Services_ChatServiceServer, clientUNiqueCode int, errch chan error) {

	// implement a loop
	for {

		// loop through messages in MQue
		for {

			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObject.MQue[0].ClientName
			message4Client := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			// send message to designated client
			if senderUniqueCode != clientUNiqueCode {

				// Работа с сгенерированным кодом
				err := csi_.Send(&FromServer{
					Name: senderName4Client,
					Body: message4Client,
				})

				if err != nil {
					errch <- err
				}

				messageHandleObject.mu.Lock()

				if len(messageHandleObject.MQue) > 1 {
					messageHandleObject.MQue = messageHandleObject.MQue[1:] // delete the message at index 0 after sending it
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}

				messageHandleObject.mu.Unlock()

			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}
