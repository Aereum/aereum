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

func (s *JoinNetwork) Kind() byte {
	return iJoinNetwork
}

func (s *JoinNetwork) Serialize() []byte {
	bytes := make([]byte, 0)
	PutString(s.caption, &bytes)
	PutString(s.details, &bytes)
	return s.authored.serialize(iJoinNetwork, bytes)
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
