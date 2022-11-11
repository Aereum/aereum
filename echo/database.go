package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

const maxFileLength = 1 << 32
const fileNameTemaplate = "aereum_rawdb_%v.dat"

type DBMessage struct {
	file     int
	position int
	size     int
}

type DB struct {
	files     []*os.File
	length    int
	current   int
	index     map[crypto.Token][]DBMessage
	listeners map[crypto.Token][]*network.SecureConnection
	broker    chan []byte
	mu        sync.Mutex
}

func (db *DB) ReadMessage(msg DBMessage) []byte {
	file := db.files[msg.file]
	data := make([]byte, msg.size)
	bytes, err := file.ReadAt(data, int64(msg.position))
	if err != nil || bytes != msg.size {
		return nil
	}
	return data
}

func NewDB(broker chan []byte) *DB {
	return &DB{
		files:     make([]*os.File, 0),
		index:     make(map[crypto.Token][]DBMessage),
		listeners: make(map[crypto.Token][]*network.SecureConnection),
	}
}

func (db *DB) CreateNewFile() {
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(db.files) > 0 {
		current := db.files[len(db.files)-1]
		if err := current.Close(); err != nil {
			panic(err)
		}
		// reopen as readonly
		filepath := current.Name()
		if file, err := os.Open(filepath); err != nil {
			panic(err)
		} else {
			db.files[len(db.files)-1] = file
		}
	}
	if newFile, err := os.Create(fmt.Sprintf(fileNameTemaplate, len(db.files)+1)); err != nil {
		panic(err)
	} else {
		db.files = append(db.files, newFile)
		db.length = 0
		db.current += 1
	}
}

func (db *DB) AppendMsg(msg []byte) (*DBMessage, error) {
	dataLen, err := db.files[db.current-1].Write(msg)
	if err != nil || dataLen < len(msg) {
		return nil, fmt.Errorf("could not persist message on the database: %v", err)
	}
	position := db.length
	fileNum := len(db.files) + 1
	db.length += dataLen
	if db.length > maxFileLength {
		db.CreateNewFile()
	}
	return &DBMessage{file: fileNum, position: position, size: len(msg)}, nil
}

func (db *DB) Broadcast(token crypto.Token, msg []byte) {
	if listeners, ok := db.listeners[token]; ok {
		for _, conn := range listeners {
			conn.WriteMessage(msg)
		}
	}
}

func (db *DB) IndexTokens(msg *DBMessage, data []byte, tokens []crypto.Token) {
	for _, token := range tokens {
		if index, ok := db.index[token]; ok {
			db.index[token] = append(index, *msg)
		} else {
			db.index[token] = []DBMessage{*msg}
		}
		db.Broadcast(token, data)
	}
}

func (db *DB) Incorporate(msg []byte) {
	newMsg, err := db.AppendMsg(msg)
	if err != nil {
		return
	}
	tokens := getTokens(msg)
	if len(tokens) > 0 {
		db.IndexTokens(newMsg, msg, tokens)
	}
}

func authoredTokens(authored *instructions.AuthoredInstruction) []crypto.Token {
	return []crypto.Token{authored.Author, authored.Attorney, authored.Wallet}
}
