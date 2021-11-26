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

import (
	"encoding/binary"

	"github.com/Aereum/aereum/core/crypto"
)

func DeleteOrInsertExpire(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
	if found {
		if param[0] == delete { //Delete
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true},
			}
		} else if param[0] == exists { // exists?
			acc := b.ReadItem(item)
			return OperationResult{
				result: QueryResult{ok: true, data: acc[size:]},
			}
		} else { // insert
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	} else {
		if param[0] == insert {
			value := binary.LittleEndian.Uint64(param[1:])
			acc := make([]byte, size+8)
			binary.LittleEndian.PutUint64(acc[size:], uint64(value))
			copy(acc[0:size], hash[:])
			b.WriteItem(item, acc)
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

type HashExpireVault struct {
	hs *HashStore
}

func (w *HashExpireVault) Exists(hash crypto.Hash) uint64 {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: []byte{1}, response: response})
	resp := <-response
	if !resp.ok {
		return 0
	}
	return binary.LittleEndian.Uint64(resp.data)
}

func (w *HashExpireVault) Insert(hash crypto.Hash, expire uint64) bool {
	response := make(chan QueryResult)
	param := make([]byte, 8+1)
	param[0] = 2
	binary.LittleEndian.PutUint64(param[1:], expire)
	w.hs.Query(Query{hash: hash, param: param, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashExpireVault) Remove(hash crypto.Hash) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: []byte{0}, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashExpireVault) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewExpireHashVault(name string, epoch uint64, bitsForBucket int64) *HashExpireVault {
	nbytes := 8 + int64(1<<bitsForBucket)
	bytestore := NewMemoryStore(nbytes)
	bucketstore := NewBucketStore(40, 6, bytestore)
	return &HashExpireVault{
		hs: NewHashStore(name, bucketstore, int(bitsForBucket), CreditOrDebit),
	}
}
