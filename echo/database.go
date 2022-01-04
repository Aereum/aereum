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
	var tokens []crypto.Token
	switch instructions.InstructionKind(msg) {
	case instructions.IContent:
		if content := instructions.ParseContent(msg); content != nil {
			tokens = []crypto.Token{content.Author, content.Attorney, content.Wallet, content.Audience}
		}
	case instructions.ITransfer:
		if transfer := instructions.ParseTransfer(msg); transfer != nil {
			tokens = []crypto.Token{transfer.From}
			for _, reciepient := range transfer.To {
				tokens = append(tokens, reciepient.Token)
			}
		}
	case instructions.IDeposit:
		if deposit := instructions.ParseDeposit(msg); deposit != nil {
			tokens = []crypto.Token{deposit.Token}
		}
	case instructions.IWithdraw:
		if withdraw := instructions.ParseWithdraw(msg); withdraw != nil {
			tokens = []crypto.Token{withdraw.Token}
		}
	case instructions.IJoinNetwork:
		if join := instructions.ParseJoinNetwork(msg); join != nil {
			tokens = authoredTokens(join.Authored)
		}
	case instructions.IUpdateInfo:
		if update := instructions.ParseUpdateInfo(msg); update != nil {
			tokens = authoredTokens(update.Authored)
		}
	case instructions.ICreateAudience:
		if join := instructions.ParseCreateStage(msg); join != nil {
			tokens = append(authoredTokens(join.Authored), join.Audience)
		}
	case instructions.IJoinAudience:
		if join := instructions.ParseJoinStage(msg); join != nil {
			tokens = append(authoredTokens(join.Authored), join.Audience)
		}
	case instructions.IAcceptJoinRequest:
		if join := instructions.ParseAcceptJoinStage(msg); join != nil {
			tokens = append(authoredTokens(join.Authored), join.Stage)
		}
	case instructions.IUpdateAudience:
		// TODO
	case instructions.IGrantPowerOfAttorney:
		if grant := instructions.ParseGrantPowerOfAttorney(msg); grant != nil {
			tokens = append(authoredTokens(grant.Authored), grant.Attorney)
		}
	case instructions.IRevokePowerOfAttorney:
		if revoke := instructions.ParseRevokePowerOfAttorney(msg); revoke != nil {
			tokens = append(authoredTokens(revoke.Authored), revoke.Attorney)
		}
	case instructions.ISponsorshipOffer:
		if offer := instructions.ParseSponsorshipOffer(msg); offer != nil {
			tokens = append(authoredTokens(offer.Authored), offer.Stage)
		}
	case instructions.ISponsorshipAcceptance:
		if accept := instructions.ParseSponsorshipAcceptance(msg); accept != nil {
			tokens = append(authoredTokens(accept.Authored), accept.Stage, accept.Offer.Authored.Author)
		}
	case instructions.ICreateEphemeral:
		if ephemeral := instructions.ParseCreateEphemeral(msg); ephemeral != nil {
			tokens = append(authoredTokens(ephemeral.Authored), ephemeral.EphemeralToken)
		}
	case instructions.ISecureChannel:
		if secure := instructions.ParseSecureChannel(msg); secure != nil {
			// TODO: token range
			tokens = authoredTokens(secure.Authored)
		}
	case instructions.IReact:
		if react := instructions.ParseReact(msg); react != nil {
			// TODO: token range
			tokens = authoredTokens(react.Authored)
		}
		/*
			ISecureChannel
			IReact
		*/
	}
	if len(tokens) > 0 {
		db.IndexTokens(newMsg, msg, tokens)
	}
}

func authoredTokens(authored *instructions.AuthoredInstruction) []crypto.Token {
	return []crypto.Token{authored.Author, authored.Attorney, authored.Wallet}
}
