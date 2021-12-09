package network

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

const (
	ISyncRequest byte = iota
	IResumeSyncRequest
	ISyncResponse
	IBlockListenerRequest
	IBlockBroadcast
	ISendInstruction
	IInstructionReceived
	IBroadcastInstruction
	IDenounceInstruction
	IPing
	IPong
	INewBlock
	IBlockValidation
	IDenounceCheckpoint
	IChecksumReceive
	IChecksumBrodcast
	IDenounceChecksum
	IDropFromPool
)

type Serializer interface {
	Serialize() []byte
	Kind() byte
}

type NetworkMessageTemplate struct {
	Version      byte
	MessageType  byte
	Timestamp    time.Time
	Nonce        []byte
	Data         Serializer
	Confirmation bool
	Signature    []byte
}

func NewNetworkMessage(msg Serializer, token crypto.PrivateKey, confirm bool) *NetworkMessageTemplate {
	netMsg := NetworkMessageTemplate{
		Version:      0,
		MessageType:  msg.Kind(),
		Timestamp:    time.Now(),
		Nonce:        crypto.Nonce(),
		Data:         msg,
		Confirmation: confirm,
	}
	hash := crypto.Hasher(netMsg.serializeWithoutSignatute())
	signature, err := token.Sign(hash[:])
	if err != nil {
		return nil
	}
	netMsg.Signature = signature
	return &netMsg
}

func (msg *NetworkMessageTemplate) serializeWithoutSignatute() []byte {
	output := []byte{0, msg.MessageType}
	util.PutUint64(uint64(msg.Timestamp.Unix()), &output)
	util.PutByteArray(msg.Nonce, &output)
	util.PutByteArray(msg.Data.Serialize(), &output)
	if msg.Confirmation {
		output = append(output, 1)
	} else {
		output = append(output, 0)
	}
	return output
}

func (msg *NetworkMessageTemplate) Serialize() []byte {
	output := msg.serializeWithoutSignatute()
	util.PutByteArray(msg.Signature, &output)
	return output
}

type SyncRequest struct{}

func (s *SyncRequest) Serialize() []byte {
	return []byte{}
}

func (s *SyncRequest) Kind() byte {
	return ISyncRequest
}

type ResumeSync struct{}

func (s *ResumeSync) Serialize() []byte {
	return []byte{}
}

func (s *ResumeSync) Kind() byte {
	return IResumeSyncRequest
}

type SyncResponse struct{}

func (s *SyncResponse) Serialize() []byte {
	return []byte{}
}

func (s *SyncResponse) Kind() byte {
	return ISyncResponse
}

type BlockListenerRequest struct{}

func (s *BlockListenerRequest) Serialize() []byte {
	return []byte{}
}

func (s *BlockListenerRequest) Kind() byte {
	return IBlockListenerRequest
}

type BlockBroadcast struct{}

func (s *BlockBroadcast) Serialize() []byte {
	return []byte{}
}

func (s *BlockBroadcast) Kind() byte {
	return IBlockBroadcast
}

type SendInstruction struct {
	MessageType byte
	Instruction []byte
}

func (s *SendInstruction) Serialize() []byte {
	return append([]byte{IBroadcastInstruction}, s.Instruction...)
}

func (s *SendInstruction) Kind() byte {
	return ISendInstruction
}

type InstructionReceived struct{}

func (s *InstructionReceived) Serialize() []byte {
	return []byte{}
}

func (s *InstructionReceived) Kind() byte {
	return IInstructionReceived
}

type BroadcastInstruction []byte

func (s BroadcastInstruction) Serialize() []byte {
	return []byte(s)
}

func (s BroadcastInstruction) Kind() byte {
	return IBroadcastInstruction
}

type DenounceInstruction struct{}

func (s *DenounceInstruction) Serialize() []byte {
	return []byte{}
}

func (s *DenounceInstruction) Kind() byte {
	return IDenounceInstruction
}

type Ping struct{}

func (s *Ping) Serialize() []byte {
	return []byte{}
}

func (s *Ping) Kind() byte {
	return IPing
}

type Pong struct{}

func (s *Pong) Serialize() []byte {
	return []byte{}
}

func (s *Pong) Kind() byte {
	return IPong
}

type NewBlock struct{}

func (s *NewBlock) Serialize() []byte {
	return []byte{}
}

func (s *NewBlock) Kind() byte {
	return INewBlock
}

type BlockValidation struct{}

func (s *BlockValidation) Serialize() []byte {
	return []byte{}
}

func (s *BlockValidation) Kind() byte {
	return IBlockValidation
}

type DenounceCheckpoint struct{}

func (s *DenounceCheckpoint) Serialize() []byte {
	return []byte{}
}

func (s *DenounceCheckpoint) Kind() byte {
	return IDenounceCheckpoint
}

type ChecksumReceive struct{}

func (s *ChecksumReceive) Serialize() []byte {
	return []byte{}
}

func (s *ChecksumReceive) Kind() byte {
	return IChecksumReceive
}

type ChecksumBrodcast struct{}

func (s *ChecksumBrodcast) Serialize() []byte {
	return []byte{}
}

func (s *ChecksumBrodcast) Kind() byte {
	return IChecksumBrodcast
}

type DenounceChecksum struct{}

func (s *DenounceChecksum) Serialize() []byte {
	return []byte{}
}

func (s *DenounceChecksum) Kind() byte {
	return IDenounceChecksum
}

type DropFromPool struct{}

func (s *DropFromPool) Serialize() []byte {
	return []byte{}
}

func (s *DropFromPool) Kind() byte {
	return IDropFromPool
}
