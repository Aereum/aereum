package instructionsnew

import (
	"encoding/json"
)

// Join network
type JoinNetwork struct {
	authored *authoredInstruction
	caption  string
	details  string
}

func (join *JoinNetwork) Kind() byte {
	return iJoinNetwork
}

func (join *JoinNetwork) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutString(join.caption, &bytes)
	PutString(join.details, &bytes)
	return bytes
}

func (join *JoinNetwork) Serialize() []byte {
	return join.authored.serialize(iJoinNetwork, join.serializeBulk())
}

func ParseJoinNetwork(data []byte) *JoinNetwork {
	if data[0] != 0 || data[1] != iJoinNetwork {
		return nil
	}
	join := JoinNetwork{
		authored: &authoredInstruction{},
	}
	position := join.authored.parseHead(data)
	join.caption, position = ParseString(data, position)
	join.details, position = ParseString(data, position)
	if !json.Valid([]byte(join.details)) {
		return nil
	}
	if join.authored.parseTail(data, position) {
		return &join
	}
	return nil
}
