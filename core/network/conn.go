package network

import (
	"errors"
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
)

var errMessageTooLarge = errors.New("message size cannot be larger than 65.536 bytes")

type SecureConnection struct {
	hash         crypto.Hash
	conn         net.Conn
	cipher       crypto.CipherNonce
	cipherRemote crypto.CipherNonce
}

func (s *SecureConnection) WriteMessage(msg []byte) error {
	sealed, nonce := s.cipher.SealWithNewNonce(msg)
	if len(sealed) > 1<<16-1 {
		return errMessageTooLarge
	}
	msgToSend := append(nonce, byte(len(sealed)), byte(len(sealed)>>8))
	msgToSend = append(msgToSend, sealed...)
	if n, err := s.conn.Write(msgToSend); n != len(sealed) {
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
	if n, err := s.conn.Read(lengthBytes); n != 2 {
		return nil, err
	}
	lenght := int(lengthBytes[0]) + (int(lengthBytes[1]) << 8)
	sealedMsg := make([]byte, lenght)
	if n, err := s.conn.Read(sealedMsg); n != int(lenght) {
		return nil, err
	}
	if msg, err := s.cipherRemote.OpenNewNonce(sealedMsg, nonce); err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

type handlePort func(conn *SecureConnection)

func ListenTCP(port int, handler handlePort, prvKey crypto.PrivateKey, validator chan consensus.ValidatedConnection) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			secureConnection, err := PerformServerHandShake(conn, prvKey, validator)
			if err != nil {
				conn.Close()
			} else {
				handler(secureConnection)
			}
		}
	}
}

func ConnectTCP(address string, prvKey crypto.PrivateKey, pubKey crypto.PublicKey) *SecureConnection {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	secureConnection, err := PerformClientHandShake(conn, prvKey, pubKey)
	if err != nil {
		conn.Close()
		return nil
	}
	return secureConnection
}

type connResult struct {
	hash crypto.Hash
	conn *SecureConnection
}

func ConnectTCPPool(trusted map[crypto.PublicKey]string, prvKey crypto.PrivateKey) map[crypto.Hash]*SecureConnection {
	remaining := len(trusted)
	resp := make(chan connResult)
	connections := make(map[crypto.Hash]*SecureConnection)
	for pubKey, addr := range trusted {
		go func(pubKey crypto.PublicKey, addr string) {
			conn := ConnectTCP(addr, prvKey, pubKey)
			resp <- connResult{
				hash: crypto.Hasher(pubKey.ToBytes()),
				conn: conn,
			}
		}(pubKey, addr)
	}
	go func() {
		for {
			newConn := <-resp
			remaining -= 1
			if newConn.conn != nil {
				connections[newConn.hash] = newConn.conn
			}
			if remaining == 0 {
				break
			}
		}
	}()
	return connections
}
