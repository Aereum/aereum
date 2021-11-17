package store

import (
	"github.com/Aereum/aereum/core/crypto"
)

func GetOrSetAudience(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
	get := false
	if len(param) == 0 {
		get = true
	}
	if found {
		if get {
			keys := b.ReadItem(item)
			return OperationResult{
				result: QueryResult{ok: true, data: keys[size:]},
			}
		} else {
			updated := make([]byte, crypto.Size+3*crypto.PublicKeySize)
			copy(updated[0:size], hash[:])
			copy(updated[size:], param)
			b.WriteItem(item, updated)
			return OperationResult{
				result: QueryResult{ok: true},
			}

		}
	} else {
		if !get {
			newKeys := make([]byte, crypto.Size+3*crypto.PublicKeySize)
			copy(newKeys[0:size], hash[:])
			copy(newKeys[size:], param)
			b.WriteItem(item, newKeys)
			return OperationResult{
				added:  &Item{bucket: b, item: item},
				result: QueryResult{ok: false},
			}
		} else {
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	}
}

type Audience struct {
	hs *HashStore
}

func (w *Audience) GetKeys(hash crypto.Hash) (bool, []byte) {
	response := make(chan QueryResult)
	ok, keys := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	return ok, keys
}

func (w *Audience) Exists(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	return ok
}

func (w *Audience) SetKeys(hash crypto.Hash, keys []byte) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: keys, response: response})
	return ok
}

func (w *Audience) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewMemoryAudienceStore(epoch uint64, bitsForBucket int64) *Wallet {
	itemsize := int64(32 + 2*crypto.PublicKeySize)
	nbytes := 56 + int64(1<<bitsForBucket)*(itemsize*6+8)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(itemsize, 6, bytestore)
	w := &Wallet{
		hs: NewHashStore("audience", bucketstore, int(bitsForBucket), GetOrSetAudience),
	}
	w.hs.Start()
	return w
}
