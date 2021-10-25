package network

import (
	"net"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

type ciphernonce struct {
	msg   []byte
	nonce []byte
}

var validator chan ValidatedConnection = func() chan ValidatedConnection {
	validator := make(chan ValidatedConnection)
	go func() {
		validate := <-validator
		validate.ok <- true
	}()
	return validator
}()

func TestSecureConnection(t *testing.T) {

	pubSv, prvSv := crypto.RandomAsymetricKey()
	listener, _ := net.Listen("tcp", ":7780")
	cipher := make(chan ciphernonce)
	go func() {
		conn, _ := listener.Accept()
		sec, err := PerformServerHandShake(conn, prvSv, validator)
		if err != nil {
			t.Error(err)
		}
		var msg ciphernonce
		msg.msg, msg.nonce = sec.cipher.Seal([]byte("thats correct"))
		cipher <- msg
		msg.msg, msg.nonce = sec.cipherRemote.Seal([]byte("thats also correct"))
		cipher <- msg
	}()

	_, prvCl := crypto.RandomAsymetricKey()
	client, _ := net.Dial("tcp", ":7780")
	sec, err := PerformClientHandShake(client, prvCl, pubSv)
	if err != nil {
		t.Error(err)
	}
	msg := <-cipher
	msgData, err := sec.cipherRemote.Open(msg.msg, msg.nonce)
	if err != nil {
		t.Fatal(err)
	}
	if string(msgData) != "thats correct" {
		t.Fatalf("wrong message:%v", string(msgData))
	}
	msg = <-cipher
	msgData, err = sec.cipher.Open(msg.msg, msg.nonce)
	if err != nil {
		t.Fatal(err)
	}
	if string(msgData) != "thats also correct" {
		t.Fatalf("wrong message:%v", string(msgData))
	}

}
