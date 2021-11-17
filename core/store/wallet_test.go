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

// Package wallet contains the implementation of the hashtable data structure
// used to store current state of the aereum system.
package store

import (
	"crypto/rand"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestWallet(t *testing.T) {
	var w = NewMemoryWalletStore(0, 6)
	hash := crypto.Hash{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	w.Credit(hash, 10)
	if ok, balance := w.Balance(hash); !ok || balance != 10 {
		t.Errorf("wrong balance %v, %v\n", ok, balance)
	}
	if ok := w.Debit(hash, 5); !ok {
		t.Errorf("cannot debit a valid account\n")
	}
	if ok, balance := w.Balance(hash); !ok || balance != 5 {
		t.Errorf("wrong balance after debit %v, %v\n", ok, balance)
	}
	if ok := w.Debit(hash, 6); ok {
		t.Errorf("could debit an account withou sufficient balance\n")
	}
}

func TestWalletDoubling(t *testing.T) {
	var w = NewMemoryWalletStore(0, 6)
	maptest := make(map[crypto.Hash]uint64)
	for n := 0; n < 4096; n++ {
		h := make([]byte, 32)
		rand.Read(h)
		var hash crypto.Hash
		for m := 0; m < size; m++ {
			hash[m] = h[m]
		}
		if b, ok := maptest[hash]; ok {
			maptest[hash] = b + 10
		} else {
			maptest[hash] = 10
		}
		w.Credit(hash, 10)
	}
	for hash, balance := range maptest {
		if ok, b := w.Balance(hash); !ok || b != balance {
			t.Fatalf("wrong balance after duplication: %v, %v, %v\n", ok, b, balance)
		}
	}
}
