package index

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
)

type InstructionIndex struct {
	bck *bucket
	idx *persistentMap
}

func OpenInstructionIndex(identity string) *InstructionIndex {
	bucketFileName := fmt.Sprintf("%v_buckets.dat", identity)
	indexFileName := fmt.Sprintf("%v_buckets_index.dat", identity)
	return &InstructionIndex{
		bck: OpenBucket(bucketFileName),
		idx: openPersistentMap(indexFileName),
	}
}

func (i *InstructionIndex) Append(token crypto.Token, position uint64) {
	if bucket := i.idx.Get(token); bucket == 0 {
		bucket = i.bck.new(position)
		i.idx.Set(token, bucket)
	} else {
		bucket = i.bck.append(bucket, position)
		if bucket != 0 {
			i.idx.Set(token, bucket)
		}
	}
}

func (i *InstructionIndex) Retrieve(token crypto.Token) []uint64 {
	bucket := i.idx.Get(token)
	if bucket == 0 {
		return nil
	}
	return i.bck.readAll(bucket)
}
