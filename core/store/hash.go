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
// along with the aereum library. If not, see <http://www.gnu.org/licenses/>.
package store

import "github.com/Aereum/aereum/core/crypto"

const (
	delete byte = iota
	exists
	insert
)

func DeleteOrInsert(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
	if found {
		if param[0] == delete { //Delete
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true},
			}
		} else if param[0] == exists { // exists?
			return OperationResult{
				result: QueryResult{ok: true},
			}
		} else { // insert
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	} else {
		if param[0] == insert {
			b.WriteItem(item, hash[:])
			return OperationResult{
				added:  &Item{bucket: b, item: item},
				result: QueryResult{ok: true},
			}
		} else {
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	}
}

type HashVault struct {
	hs *HashStore
}

func (w *HashVault) Exists(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: []byte{1}, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashVault) Insert(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: []byte{2}, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashVault) Remove(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: []byte{0}, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashVault) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewHashVault(name string, epoch uint64, bitsForBucket int64) *Wallet {
	nbytes := 8 + int64(1<<bitsForBucket)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(40, 6, bytestore)
	return &Wallet{
		hs: NewHashStore(name, bucketstore, int(bitsForBucket), CreditOrDebit),
	}
}
