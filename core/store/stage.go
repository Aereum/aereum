package store

import (
	"github.com/Aereum/aereum/core/crypto"
)

type StageKeys struct {
	Moderate crypto.Token
	Submit   crypto.Token
	Stage    crypto.Token
	Flag     byte
}

func GetOrSetStage(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
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

type Stage struct {
	hs *HashStore
}

func (w *Stage) GetKeys(hash crypto.Hash) *StageKeys {
	response := make(chan QueryResult)
	ok, keys := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	if !ok {
		return nil
	}
	stage := StageKeys{}
	copy(stage.Moderate[:], keys[0:crypto.TokenSize])
	copy(stage.Submit[:], keys[crypto.TokenSize:2*crypto.TokenSize])
	copy(stage.Stage[:], keys[2*crypto.TokenSize:3*crypto.TokenSize])
	stage.Flag = keys[3*crypto.TokenSize]
	return &stage
}

func (w *Stage) Exists(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	return ok
}

func (w *Stage) SetKeys(hash crypto.Hash, stage *StageKeys) bool {
	keys := make([]byte, 2*crypto.TokenSize+1)
	copy(keys[0:crypto.TokenSize], stage.Moderate[:])
	copy(keys[crypto.TokenSize:2*crypto.TokenSize], stage.Submit[:])
	copy(keys[2*crypto.TokenSize:3*crypto.TokenSize], stage.Stage[:])
	keys[3*crypto.TokenSize] = stage.Flag
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: keys, response: response})
	return ok
}

func (w *Stage) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewMemoryAudienceStore(epoch uint64, bitsForBucket int64) *Stage {
	itemsize := int64(crypto.Size + 3*crypto.TokenSize + 1)
	nbytes := 56 + int64(1<<bitsForBucket)*(itemsize*6+8)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(itemsize, 6, bytestore)
	w := &Stage{
		hs: NewHashStore("audience", bucketstore, int(bitsForBucket), GetOrSetStage),
	}
	w.hs.Start()
	return w
}
