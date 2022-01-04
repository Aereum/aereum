package main

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

var ignorecheck chan consensus.ValidatedConnection = func() chan consensus.ValidatedConnection {
	ok := make(chan consensus.ValidatedConnection)
	go func() {
		check := <-ok
		check.Ok <- true
	}()
	return ok
}()

func Serve(port int, prvKey crypto.PrivateKey, db *DB) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			secureConnection, err := network.PerformServerHandShake(conn, prvKey, ignorecheck)
			if err != nil {
				conn.Close()
			} else {
				go ReceiveMsg(secureConnection, db)
			}
		}
	}
}

func ReceiveMsg(conn *network.SecureConnection, db *DB) {
	for {
		msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if len(msg) == crypto.PublicKeySize {
			var token crypto.Token
			copy(token[:], msg)
			if listeners, ok := db.listeners[token]; ok {
				db.listeners[token] = append(listeners, conn)
			} else {
				db.listeners[token] = []*network.SecureConnection{conn}
			}
			// send all data on another thread
			go func(db *DB) {
				if messages, ok := db.index[token]; ok {
					for _, msg := range messages {
						if data := db.ReadMessage(msg); data != nil {
							conn.WriteMessage(data)
						}
					}
				}
			}(db)
		} else {
			if instruction := instructions.ParseInstruction(msg); instruction != nil {
				db.broker <- msg
			}
		}
	}
}
