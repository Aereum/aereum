package network

import (
	"errors"
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/crypto"
)

var errMessageTooLarge = errors.New("message size cannot be larger than 65.536 bytes")

type SecureConnection struct {
	conn         net.Conn
	cipher       crypto.Cipher
	cipherRemote crypto.Cipher
}

func (s *SecureConnection) WriteMessage(msg []byte) error {
	msgSealed, nonce := s.cipher.Seal(msg)
	if len(msgSealed) > 1<<16-1 {
		return errMessageTooLarge
	}
	msgToSend := append(nonce, byte(len(msgSealed)), byte(len(msgSealed)>>8))
	msgToSend = append(msgToSend, msgSealed...)
	if n, err := s.conn.Write(msgToSend); n != len(msgToSend) {
		return err
	}
	return nil
}

func (s *SecureConnection) ReadMessage() ([]byte, error) {
	nonce := make([]byte, crypto.NonceSize)
	if n, err := s.conn.Read(nonce); n != crypto.NonceSize {
		return nil, err
	}
	lengthBytes := make([]byte, 2)
	lenght := lengthBytes[0] + (lengthBytes[1] << 8)
	sealedMsg := make([]byte, lenght)
	if n, err := s.conn.Read(sealedMsg); n != int(lenght) {
		return nil, err
	}
	if msg, err := s.cipherRemote.Open(sealedMsg, nonce); err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

type ReceiveMessage struct {
	Sender []byte
	Msg    []byte
}

type SendMessage struct {
	Msg []byte
	ok  bool
}

func ListenTCP(port int, receive chan ReceiveMessage, send chan SendMessage, prvKey crypto.PrivateKey) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			PerformServerHandShake(conn, prvKey)

		}
	}
}
