package network

import (
	"github.com/Aereum/aereum/core/message"
	"gopkg.in/zeromq/goczmq.v4"
)

const (
	TCPMessageVoid byte = iota
	TCPMessageAuthor
	TCPTransfer
	TCPReceiveBlock
	TCPRequestBlock
)

type SyncSocket struct {
}

type TCPMessage []byte

func (t TCPMessage) Kind() byte {
	if t == nil || len(TCPMessage{}) == 0 {
		return 0
	}
	return t[0]
}

func (t TCPMessage) AsMessage() (*message.Message, error) {
	return message.ParseMessage(t[1:])
}

func (t TCPMessage) AsTransaction() (*message.Transfer, error) {
	return message.ParseTranfer(t[1:])
}

func NewRouter() {
	goczmq.NewRouter("tcp://*:5555")

}
