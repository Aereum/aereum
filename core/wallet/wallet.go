package wallet

import "encoding/binary"

func CreditOrDebit(found bool, hash Hash, b *Bucket, item int64, param int64) OperationResult {
	if found {
		acc := b.ReadItem(item)
		balance := int64(binary.LittleEndian.Uint64(acc[size:]))
		newbalance := balance + param
		if newbalance > 0 {
			return OperationResult{
				result: QueryResult{ok: true, data: acc},
			}
		} else if newbalance == 0 {
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true, data: acc},
			}
		} else {
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	} else {
		if param >= 0 {
			acc := make([]byte, size+8)
			binary.LittleEndian.PutUint64(acc[size:], uint64(param))
			copy(acc[0:size], hash[:])
			b.WriteItem(item, acc)
			return OperationResult{
				added:  &Item{bucket: b, item: item},
				result: QueryResult{ok: false, data: acc},
			}
		} else {
			return OperationResult{
				result: QueryResult{
					ok: false,
				},
			}
		}
	}
}

type Wallet struct {
	hs *HashStore
}

func (w *Wallet) Credit(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: int64(value), response: response})
	resp := <-response
	return resp.ok
}

func (w *Wallet) Debit(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: -int64(value), response: response})
	resp := <-response
	return resp.ok
}

func (w *Wallet) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewMemoryWalletStore(epoch uint64, bitsForBucket int64) *Wallet {
	nbytes := 8 + int64(1<<bitsForBucket)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(40, 6, 8, bytestore)
	return &Wallet{
		hs: NewHashStore("wallet", bucketstore, int(bitsForBucket), CreditOrDebit),
	}
}
