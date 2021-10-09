// Copyright 2021 The aereum Authors
// This file is part of the aereum library.
//
// The aereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The aereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package message contains data types related to aereum network.
package wallet

import (
	"encoding/binary"
)

func CreditOrDebit(found bool, hash Hash, b *Bucket, item int64, param int64) OperationResult {
	if found {
		acc := b.ReadItem(item)
		balance := int64(binary.LittleEndian.Uint64(acc[size:]))
		if param == 0 {
			return OperationResult{
				result: QueryResult{ok: true, data: acc},
			}
		}
		newbalance := balance + param
		if newbalance > 0 {
			// update balance
			acc := make([]byte, size+8)
			binary.LittleEndian.PutUint64(acc[size:], uint64(newbalance))
			copy(acc[0:size], hash[:])
			b.WriteItem(item, acc)
			return OperationResult{
				result: QueryResult{ok: true, data: acc},
			}
		} else if newbalance == 0 {
			// account is market to be deleted
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
	ok, _ := w.hs.Query(Query{hash: hash, param: int64(value), response: response})
	return ok
}

func (w *Wallet) Balance(hash Hash) (bool, uint64) {
	response := make(chan QueryResult)
	ok, data := w.hs.Query(Query{hash: hash, param: 0, response: response})
	if ok {
		return true, binary.LittleEndian.Uint64(data[32:])
	}
	return false, 0
}

func (w *Wallet) Debit(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	ok, _ := w.hs.Query(Query{hash: hash, param: -int64(value), response: response})
	return ok
}

func (w *Wallet) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewMemoryWalletStore(epoch uint64, bitsForBucket int64) *Wallet {
	nbytes := 56 + int64(1<<bitsForBucket)*(40*6+8)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(40, 6, bytestore)
	w := &Wallet{
		hs: NewHashStore("wallet", bucketstore, int(bitsForBucket), CreditOrDebit),
	}
	w.hs.Start()
	return w
}
