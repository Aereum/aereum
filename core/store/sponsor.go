package store

import (
	"github.com/Aereum/aereum/core/crypto"
)

func GetOrSetSponsor(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
	if found {
		if param[0] == 0 { // get
			keys := b.ReadItem(item)
			return OperationResult{
				result: QueryResult{ok: true, data: keys},
			}
		} else if param[0] == 1 { // set
			return OperationResult{
				result: QueryResult{ok: false},
			}
		} else { // remove
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true},
			}
		}
	} else {
		if param[0] == 0 { //get
			return OperationResult{
				result: QueryResult{ok: false},
			}
		} else if param[0] == 1 { // set
			contentHash := make([]byte, crypto.Size)
			copy(contentHash[0:crypto.Size], param[1:])
			b.WriteItem(item, contentHash)
			return OperationResult{
				added:  &Item{bucket: b, item: item},
				result: QueryResult{ok: true},
			}
		} else { // remove
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
	ok, keys := w.hs.Query(Query{hash: hash, param: []byte{0}, response: response})
	if ok {
		return ok, keys
	}
	return false, nil
}

func (w *Sponsor) Exists(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: []byte{0}, response: response})
	return ok
}

func (w *Sponsor) SetContentHash(hash crypto.Hash, keys []byte) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: append([]byte{1}, hash[:]...), response: response})
	return ok
}

func (w *Sponsor) RemoveContentHash(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: []byte{2}, response: response})
	return ok
}

func (w *Sponsor) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewSponsorShipOfferStore(epoch uint64, bitsForBucket int64) *Sponsor {
	itemsize := int64(crypto.Size)
	nbytes := 56 + int64(1<<bitsForBucket)*(itemsize*6+8)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(itemsize, 6, bytestore)
	w := &Sponsor{
		hs: NewHashStore("sponsor", bucketstore, int(bitsForBucket), GetOrSetSponsor),
	}
	w.hs.Start()
	return w
}
