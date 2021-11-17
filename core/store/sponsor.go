package store

import (
	"github.com/Aereum/aereum/core/crypto"
)

func GetOrSetSponsor(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
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

type Sponsor struct {
	hs *HashStore
}

func (w *Sponsor) GetContentHash(hash crypto.Hash) (bool, []byte) {
	response := make(chan QueryResult)
	ok, keys := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	return ok, keys
}

func (w *Sponsor) Exists(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: []byte{}, response: response})
	return ok
}

func (w *Sponsor) SetContentHash(hash crypto.Hash, keys []byte) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: keys, response: response})
	return ok
}

func (w *Sponsor) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewSponsorShipOfferStore(epoch uint64, bitsForBucket int64) *Wallet {
	itemsize := int64(crypto.Size)
	nbytes := 56 + int64(1<<bitsForBucket)*(itemsize*6+8)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(itemsize, 6, bytestore)
	w := &Wallet{
		hs: NewHashStore("sponsor", bucketstore, int(bitsForBucket), GetOrSetSponsor),
	}
	w.hs.Start()
	return w
}
