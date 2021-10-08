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
package wallet

func DeleteOrInsert(found bool, hash Hash, b *Bucket, item int64, param int64) OperationResult {
	if found {
		if param == -1 { //Delete
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true},
			}
		} else if param == 0 { // exists?
			return OperationResult{
				result: QueryResult{ok: true},
			}
		} else { // insert
			return OperationResult{
				result: QueryResult{ok: false},
			}
		}
	} else {
		if param == 1 {
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

func (w *HashVault) Exists(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: 0, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashVault) Insert(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: 1, response: response})
	resp := <-response
	return resp.ok
}

func (w *HashVault) Remove(hash Hash, value uint64) bool {
	response := make(chan QueryResult)
	w.hs.Query(Query{hash: hash, param: -1, response: response})
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
