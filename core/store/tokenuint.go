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
package store

import (
	"github.com/Aereum/aereum/core/crypto"
)

func ReadOrSetToken(found bool, hash crypto.Hash, b *Bucket, item int64, param []byte) OperationResult {
	if param[0] == 1 { // get
		if found {
			data := b.ReadItem(item)
			return OperationResult{
				result: QueryResult{ok: true, data: data[size:]},
			}
		} else {
			return OperationResult{
				result: QueryResult{ok: false, data: nil},
			}
		}
	} else { //set
		hashAndValue := make([]byte, size+8)
		copy(hashAndValue[0:size], hash[:])
		copy(hashAndValue[size:], param[1:])
		b.WriteItem(item, hashAndValue)
		return OperationResult{
			result: QueryResult{ok: true, data: nil},
		}
	}
}

type TokenByteArrayStore struct {
	hs *HashStore
}

func (w *TokenByteArrayStore) SetToken(token crypto.Token, value [8]byte) bool {
	response := make(chan QueryResult)
	param := make([]byte, 9)
	copy(param, value[:])
	hash := crypto.HashToken(token)
	ok, _ := w.hs.Query(Query{hash: hash, param: param, response: response})
	return ok
}

func (w *TokenByteArrayStore) GetToken(token crypto.Token) [8]byte {
	response := make(chan QueryResult)
	param := []byte{1}
	hash := crypto.HashToken(token)
	ok, data := w.hs.Query(Query{hash: hash, param: param, response: response})
	var output [8]byte
	if !ok {
		return output
	}
	copy(output[:], data)
	return output
}

func (w *TokenByteArrayStore) Close() bool {
	ok := make(chan bool)
	w.hs.stop <- ok
	return <-ok
}

func NewTokenByteArrayStore(storage string, bitsForBucket int64) *TokenByteArrayStore {
	nbytes := 56 + int64(1<<bitsForBucket)*(40*6+8)
	var bytestore ByteStore
	if storage == "RAM" {
		bytestore = NewMemoryStore(nbytes)
	} else {
		bytestore = NewFileStore(storage, nbytes)
	}
	bucketstore := NewBucketStore(40, 6, bytestore)
	w := &TokenByteArrayStore{
		hs: NewHashStore("TokenUint64", bucketstore, int(bitsForBucket), ReadOrSetToken),
	}
	w.hs.Start()
	return w
}
