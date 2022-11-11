package index

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
)

type InstructionIndex struct {
	bck *bucket
	idx *store.TokenByteArrayStore
}

type InstructionPosition struct {
	file     uint16
	flag     uint8
	position uint64
}

func (e InstructionPosition) Bytes() [8]byte {
	var bytes [8]byte
	bytes[0] = byte(e.file)
	bytes[1] = byte(e.file >> 8)
	bytes[2] = e.flag
	bytes[3] = byte(e.position)
	bytes[4] = byte(e.position >> 8)
	bytes[5] = byte(e.position >> 16)
	bytes[6] = byte(e.position >> 24)
	bytes[7] = byte(e.position >> 32)
	return bytes
}

func ParseEntryPoisition(bytes []byte) InstructionPosition {
	var pos InstructionPosition
	pos.file = uint16(bytes[0]) | uint16(bytes[1])<<8
	pos.flag = bytes[2]
	pos.position = uint64(bytes[3]) | uint64(bytes[4])<<8 |
		uint64(bytes[5])<<16 | uint64(bytes[6])<<24 |
		uint64(bytes[7])<<32
	return pos
}

func OpenInstructionIndex(identity string) *InstructionIndex {
	bucketFileName := fmt.Sprintf("%v_buckets.dat", identity)
	indexFileName := fmt.Sprintf("%v_buckets_index.dat", identity)
	return &InstructionIndex{
		bck: OpenBucket(bucketFileName),
		idx: store.NewTokenByteArrayStore(indexFileName, 6),
	}
}

func (i *InstructionIndex) Append(token crypto.Token, position InstructionPosition) {
	value := position.Bytes()
	if bucket := i.idx.GetToken(token); bucket[0] == 0 {
		bucket = i.bck.new(position)
		i.idx.SetToken(token, bucket)
	} else {
		bucket = i.bck.append(bucket, position)
		if bucket != 0 {
			i.idx.SetToken(token, bucket)
		}
	}
}

func (i *InstructionIndex) Retrieve(token crypto.Token) []uint64 {
	bucket := i.idx.GetToken(token)
	if bucket == 0 {
		return nil
	}
	return i.bck.readAll(bucket)
}
