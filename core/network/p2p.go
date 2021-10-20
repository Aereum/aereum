package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/logger"
	"github.com/Aereum/aereum/core/message"
)

const (
	validationNodePort             = 7080
	blockBroadcastPort             = 7801
	messageReceiveConnectionPort   = 7802
	messageBroadcastConnectionPort = 7803
	syncPort                       = 7804
)

type MessageQueueRequest struct {
	message  message.Message
	response chan bool
}

type MessageUnqueueRequest struct {
	response chan message.Message
}

type ValidMessageQueue struct {
	queue   chan MessageQueueRequest
	unqueue chan MessageUnqueueRequest
}

func NewValidMessageQueue() *ValidMessageQueue {
	queue := make(chan MessageQueueRequest)
	unqueue := make(chan MessageUnqueueRequest)
	//messages := make([]message.Message, 0)
	//hashes := make(map[hashdb.Hash]struct{})
	go func() {
		for {
			select {
			case <-queue:
				//
			case <-unqueue:
				//
			}
		}
	}()
	return &ValidMessageQueue{
		queue:   queue,
		unqueue: unqueue,
	}
}

func NewMessageServer() {
	//messages := NewValidMessageQueue()
	service := fmt.Sprintf(":%v", messageReceiveConnectionPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	logger.MustOrPanic(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	logger.MustOrPanic(err)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
		}
	}()
}

type MessageReceiveConnection struct {
	conn     *net.TCPConn
	msgqueue *ValidMessageQueue
}

// exchange keys for secure channel
func Handshake(conn *net.TCPConn) error {
	return nil
}
