package main

import (
	"fmt"
	"os"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

const maxFileLength = 1 << 32
const fileNameTemaplate = "aereum_rawdb_%v.dat"

type DBMessage struct {
	file     int
	position int
}

type DB struct {
	permanent []*os.File
	current   *os.File
	length    int
	index     map[crypto.Token][]DBMessage
	listeners map[crypto.Token][]*network.SecureConnection
}

func (db *DB) CreateNewFile() error {
	if db.current != nil {
		if err := db.current.Close(); err != nil {
			return err
		}
		filepath := db.current.Name()
		if file, err := os.Open(filepath); err != nil {
			db.permanent = append(db.permanent, file)
		} else {
			return nil
		}
	}
	if newFile, err := os.Create(fmt.Sprintf(fileNameTemaplate, len(db.permanent)+1)); err != nil {
		db.current = newFile
		db.length = 0
		return nil
	} else {
		db.current = nil
		db.length = 0
		return err
	}
}

func (db *DB) AppendMsg(msg []byte) error {
	dataLen, err := db.current.Write(msg)
	if err != nil || dataLen < len(msg) {
		return fmt.Errorf("could not persist message on the database: %v", err)
	}
	db.length += dataLen
	if db.length > maxFileLength {
		return db.CreateNewFile()
	}
	return nil
}

func (db *DB) Broadcast(token crypto.Token, msg []byte) {
	if listeners, ok := db.listeners[token]; ok {
		for _, conn := range listeners {
			conn.WriteMessage(msg)
		}
	}
}

func (db *DB) IndexContent(position, file int, content instructions.Content) {
	token := content.Audience
	if index, ok := db.index[token]; ok {
		db.index[token] = append(index, DBMessage{file: file, position: position})
	} else {
		db.index[token] = []DBMessage{{file: file, position: position}}
	}
	db.Broadcast(content.Audience, content.Serialize())
}
