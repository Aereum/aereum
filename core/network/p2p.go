package network

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"fmt"
	"hash"
	"net"

	"github.com/Aereum/aereum/core/hashdb"
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
	messages := make([]message.Message, 0)
	hashes := make(map[hashdb.Hash]struct{})
	go func() {
		for {
			select {
			case req := <-queue:
				//
			case req := <-unqueue:
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
	messages := NewValidMessageQueue()
	service := fmt.Sprintf(":%v", messageReceiveConnectionPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	logger.MustOrPanic(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	logger.MustOrPanic(err)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {

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

}

type Conn struct {
	dialDest *ecdsa.PublicKey
	conn     net.Conn
	session  *sessionState

	// These are the buffers for snappy compression.
	// Compression is enabled if they are non-nil.
	snappyReadBuffer  []byte
	snappyWriteBuffer []byte
}

// sessionState contains the session keys.
type sessionState struct {
	enc cipher.Stream
	dec cipher.Stream

	egressMAC  hashMAC
	ingressMAC hashMAC
	rbuf       readBuffer
	wbuf       writeBuffer
}

// hashMAC holds the state of the RLPx v4 MAC contraption.
type hashMAC struct {
	cipher     cipher.Block
	hash       hash.Hash
	aesBuffer  [16]byte
	hashBuffer [32]byte
	seedBuffer [32]byte
}
