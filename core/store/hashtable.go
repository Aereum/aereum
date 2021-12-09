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

// Package message contains data types related to aereum network.
package store

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"time"

	"github.com/Aereum/aereum/core/crypto"
)

var cloneInterval time.Duration

func init() {
	var err error
	cloneInterval, err = time.ParseDuration("10ms")
	if err != nil {
		panic(err)
	}
}

const (
	size       = int(sha256.Size)
	size64     = int64(size)
	NBuckets   = int64(2048) // TODO: ajusta depois
	loadFactor = int64(2)    // number of overflow buckets that will trigger duplication
)

type Item struct {
	bucket *Bucket
	item   int64
}

type OperationResult struct {
	added   *Item
	deleted *Item
	result  QueryResult
}

type QueryResult struct {
	ok   bool
	data []byte
}

type Query struct {
	hash     crypto.Hash
	param    []byte
	response chan QueryResult
}

type QueryOperation func(found bool, hash crypto.Hash, b *Bucket, item int64, data []byte) OperationResult

type HashStore struct {
	name             string
	store            *BucketStore
	bitsForBucket    int
	mask             int64
	bitsCount        []int // number of items in the bucket
	freeOverflows    []int64
	isReady          bool
	operation        QueryOperation
	query            chan Query
	doubleJob        chan int64
	cloneJob         chan int64
	stop             chan chan bool
	clone            chan chan bool
	cloned           chan bool
	isDoubling       bool
	bitsTransferered int64
	newHashStore     *HashStore
}

func NewHashStore(name string, buckets *BucketStore, bitsForBucket int, operation QueryOperation) *HashStore {
	if bitsForBucket < 6 {
		panic("bitsForBucket too small")
	}
	return &HashStore{
		name:             name,
		store:            buckets,
		bitsForBucket:    bitsForBucket,
		mask:             int64(1<<bitsForBucket - 1),
		bitsCount:        make([]int, 1<<bitsForBucket),
		freeOverflows:    make([]int64, 0),
		isReady:          true,
		operation:        operation,
		query:            make(chan Query),
		doubleJob:        make(chan int64),
		stop:             make(chan chan bool),
		cloneJob:         make(chan int64),
		clone:            make(chan chan bool),
		cloned:           make(chan bool),
		isDoubling:       false,
		bitsTransferered: 0,
		newHashStore:     nil,
	}
}

func (hs *HashStore) Query(q Query) (bool, []byte) {
	hs.query <- q
	resp := <-q.response
	return resp.ok, resp.data
}

func (hs *HashStore) Start() {
	go func() {
		for {
			select {
			case q := <-hs.query:
				resp := hs.findAndOperate(q)
				q.response <- resp
			case bucket := <-hs.doubleJob:
				hs.continueDuplication(bucket)
			case hs.cloned = <-hs.clone:
				hs.StartCloning()
			case <-hs.cloneJob:
				hs.continueCloning()
			case ok := <-hs.stop:
				// wait until cloning and doubling is complete
				if hs.store.isCloning || hs.isDoubling {
					ok <- false
				}
				close(hs.query)
				close(hs.doubleJob)
				close(hs.cloneJob)
				close(hs.stop)
				ok <- true
				return
			}
		}
	}()
}

func (ws *HashStore) findAndOperate(q Query) QueryResult {
	hashMask := q.hash.ToInt64() & ws.mask
	wallet := ws
	if ws.isDoubling && hashMask < ws.bitsTransferered {
		hashMask = q.hash.ToInt64() & ws.newHashStore.mask
		wallet = ws.newHashStore
	}
	bucket := wallet.store.ReadBucket(hashMask)
	countAccounts, totalAccounts := 0, wallet.bitsCount[hashMask]
	for {
		for item := int64(0); item < ws.store.itemsPerBucket; item++ {
			countAccounts += 1
			if countAccounts > int(totalAccounts) {
				resp := ws.operation(false, q.hash, bucket, item, q.param)
				wallet.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts) // ws -> wallet
				return resp.result
			}
			data := bucket.ReadItem(item)
			if q.hash.Equals(data) {
				resp := ws.operation(true, q.hash, bucket, item, q.param)
				wallet.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts) // ws -> walltet
				return resp.result
			}
		}
		bucket = bucket.NextBucket()
		if bucket == nil {
			panic(fmt.Sprintf("could not get here: %v %v", countAccounts, totalAccounts))
		}
	}
}

func (ws *HashStore) ProcessMutation(hashMask int64, added *Item, deleted *Item, count int) {
	if added != nil {
		ws.bitsCount[hashMask] += 1
		if added.item == ws.store.itemsPerBucket-1 {
			if len(ws.freeOverflows) > 0 {
				added.bucket.WriteOverflow(ws.freeOverflows[0])
				ws.freeOverflows = ws.freeOverflows[1:]
			} else {
				added.bucket.AppendOverflow()
			}
		}
		if (ws.store.bucketCount > 2*int64(1<<ws.bitsForBucket)) && !ws.store.isCloning && !ws.isDoubling {
			ws.startDuplication()
		}
	}
	if deleted != nil {
		lastItem := ws.bitsCount[hashMask] - 1
		ws.bitsCount[hashMask] -= 1
		if count == lastItem {
			deleted.bucket.WriteItem(deleted.item, make([]byte, ws.store.itemBytes))
			return
		}
		var previousBucket *Bucket
		lastBucket := deleted.bucket
		for {
			if nextBucket := lastBucket.NextBucket(); nextBucket != nil {
				previousBucket = lastBucket
				lastBucket = nextBucket
			} else {
				item := lastItem % int(ws.store.itemsPerBucket)
				lastBucket.WriteItem(int64(item), make([]byte, ws.store.itemBytes))
				if item == 0 && previousBucket != nil {
					ws.freeOverflows = append(ws.freeOverflows, lastBucket.n)
					previousBucket.WriteOverflow(0)
				}
				break
			}
		}
	}
}

func (w *HashStore) transferBuckets(starting, N int64) {
	// mask to test the newer bit
	highBit := uint64(1 << w.bitsForBucket)
	for bucket := starting; bucket < starting+N; bucket++ {
		// read items for the bucket
		itemsCount := int64(w.bitsCount[bucket])
		items := w.store.ReadBucket(bucket).ReadBulk(itemsCount)
		// divide items by lBit (newer bit = 0) and hBit (newer bit = 1)
		lBitBucket := make([]byte, 0, len(items)/2)
		hBitBucket := make([]byte, 0, len(items)/2)
		for _, item := range items {
			hashBit := (uint64(item[0]) + (uint64(item[1]) << 8) + (uint64(item[2]) << 16) +
				(uint64(item[3]) << 24))
			hashBit = hashBit & highBit
			if hashBit > 0 {
				hBitBucket = append(hBitBucket, item...)
			} else {
				lBitBucket = append(lBitBucket, item...)
			}
		}
		// put lBit and hBit items in new wallter
		w.newHashStore.bitsCount[bucket] = len(lBitBucket) / int(w.store.itemBytes)
		w.newHashStore.store.ReadBucket(bucket).WriteBulk(lBitBucket)
		w.newHashStore.bitsCount[bucket+int64(highBit)] = len(hBitBucket) / int(w.store.itemBytes)
		w.newHashStore.store.ReadBucket(bucket + int64(highBit)).WriteBulk(hBitBucket)
	}
	w.bitsTransferered = starting + N - 1
}

func (w *HashStore) continueDuplication(bucket int64) {
	//for bucket := int64(0); bucket < 1<<w.bitsForBucket; bucket += NBuckets {
	if bucket+NBuckets > 1<<w.bitsForBucket {
		w.transferBuckets(bucket, 1<<w.bitsForBucket-bucket)
	} else {
		w.transferBuckets(bucket, NBuckets)
	}
	if bucket+NBuckets < 1<<w.bitsForBucket {
		go func() {
			sleep, _ := time.ParseDuration("10ms")
			time.Sleep(sleep)
			w.doubleJob <- bucket + NBuckets
		}()
	} else {
		// task completed merge stores
		w.store.bytes.Merge(w.newHashStore.store.bytes)
		w.bitsForBucket = w.newHashStore.bitsForBucket
		w.mask = w.newHashStore.mask
		w.bitsCount = w.newHashStore.bitsCount
		w.freeOverflows = w.newHashStore.freeOverflows
		w.store.bucketCount = w.newHashStore.store.bucketCount
		w.isDoubling = false
		w.bitsTransferered = 0
		w.newHashStore = nil
		w.isReady = true
	}
}

func (w *HashStore) startDuplication() {
	w.isDoubling = true
	newStoreBitsForBucket := int64(w.bitsForBucket + 1)
	newStoreInitialBuckets := int64(1 << newStoreBitsForBucket)
	newStoreSize := newStoreInitialBuckets*w.store.bucketBytes + w.store.headerBytes
	newByteStore := w.store.bytes.New(newStoreSize)
	header := w.store.bytes.ReadAt(0, w.store.headerBytes)
	newByteStore.WriteAt(0, header)
	newBucketStore := NewBucketStore(w.store.itemBytes, w.store.itemsPerBucket, newByteStore)
	w.newHashStore = NewHashStore(w.name, newBucketStore, int(newStoreBitsForBucket), w.operation)
	w.newHashStore.isReady = false
	w.bitsTransferered = 0
	w.continueDuplication(0)
}

func (hs *HashStore) StartCloning() {
	hs.store.isCloning = true
	timeStamp := time.Now().Format("2006_01_02_15_04_05")
	hs.store.journal = NewJournalStore(fmt.Sprintf("%v_journal_%v.dat", hs.name, timeStamp))
	hs.store.cloning = NewJournalStore(fmt.Sprintf("%v_clone_%v.dat", hs.name, timeStamp))
}

func (hs *HashStore) continueCloning() {
	bucketsToClone := maxCloningBlockSize / hs.store.bucketBytes
	if hs.store.bucketsCloned+bucketsToClone > hs.store.bucketToClone {
		bucketsToClone = hs.store.bucketToClone - hs.store.bucketsCloned
	}
	bytesCount := hs.store.bucketBytes * bucketsToClone
	offset := hs.store.headerBytes + hs.store.bucketsCloned*hs.store.bucketBytes
	data := hs.store.bytes.ReadAt(offset, bytesCount)
	hs.store.bucketsCloned += bucketsToClone
	go func() {
		hs.store.cloning.Append(data)
		time.Sleep(cloneInterval)
		if hs.store.bucketsCloned < hs.store.bucketToClone {
			hs.cloneJob <- hs.store.bucketsCloned
		} else {
			hs.store.journal.Close()
			hs.store.cloning.Close()
			hs.store.isCloning = false
			hs.store.journal = nil
			hs.store.cloning = nil
			hs.cloned <- true
		}
	}()
}

type itemsArray [][]byte

func (ia itemsArray) Len() int {
	return len(ia)
}

func (ia itemsArray) Less(i, j int) bool {
	for n := 0; n < size; n++ {
		if ia[i][n] < ia[j][n] {
			return true
		}
		if ia[i][n] > ia[j][n] {
			return false
		}
	}
	return false
}

func (ia itemsArray) Swap(i, j int) {
	ia[i], ia[j] = ia[j], ia[i]
}

func (hs *HashStore) Hash() crypto.Hash {
	hasharray := make([]byte, 0)
	hashBlock := 256 * 256 * 16 * hs.store.itemBytes
	bucketCollection := make([]byte, 0, hashBlock)
	for n := int64(0); n < 1<<hs.bitsForBucket; n++ {
		buckets := itemsArray(hs.store.ReadBucket(n).ReadBulk(int64(hs.bitsCount[n])))
		sort.Sort(buckets)
		for _, b := range buckets {
			bucketCollection = append(bucketCollection, b...)
			if len(bucketCollection) >= int(hashBlock) {
				hash := sha256.Sum256(bucketCollection)
				hasharray = append(hasharray, hash[:]...)
				bucketCollection = make([]byte, 0, hashBlock)
			}
		}
	}
	if len(bucketCollection) > 0 {
		hash := sha256.Sum256(bucketCollection)
		hasharray = append(hasharray, hash[:]...)
	}
	return crypto.Hasher(hasharray)
}
