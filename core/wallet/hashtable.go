package main

import (
	"crypto/sha256"
	"time"
)

type Hash [size]byte

func (hash Hash) ToInt64() int64 {
	return int64(hash[0]) + (int64(hash[1]) << 8) + (int64(hash[2]) << 16) + (int64(hash[3]) << 24)
}

func (hash Hash) Equals(b []byte) bool {
	for n := 0; n < size; n++ {
		if hash[n] != b[n] {
			return false
		}
	}
	return true
}

const (
	size       = int(sha256.Size)
	size64     = int64(size)
	NBuckets   = int64(2024)
	loadFactor = int64(2)
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
	hash     Hash
	param    int64
	response chan QueryResult
}

type QueryOperation func(found bool, hash Hash, b *Bucket, item int64, param int64) OperationResult

type HashStore struct {
	store         *BucketStore
	bitsForBucket int
	mask          int64
	bitsCount     []int // number of items in the bucket
	freeOverflows []int64
	isReady       bool
	operation     QueryOperation
	//
	query     chan Query
	doubleJob chan int64
	stop      chan bool
	// for user when the map is doubling
	isDoubling       bool
	bitsTransferered int64
	newHashStore     *HashStore
}

func NewHashStore(buckets *BucketStore, bitsForBucket int, operation QueryOperation) *HashStore {
	if bitsForBucket < 12 {
		panic("bitsForBucket too small")
	}
	return &HashStore{
		store:            buckets,
		bitsForBucket:    bitsForBucket,
		mask:             int64(1<<bitsForBucket - 1),
		bitsCount:        make([]int, 1<<bitsForBucket),
		freeOverflows:    make([]int64, 0),
		isReady:          true,
		operation:        operation,
		query:            make(chan Query),
		doubleJob:        make(chan int64),
		stop:             make(chan bool),
		isDoubling:       false,
		bitsTransferered: 0,
		newHashStore:     nil,
	}
}

func (ws *HashStore) Get(q Query) QueryResult {
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
				ws.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts)
				return resp.result
			}
			data := bucket.ReadItem(item)
			if q.hash.Equals(data) {
				resp := ws.operation(true, q.hash, bucket, item, q.param)
				ws.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts)
				return resp.result
			}
		}
		bucket = bucket.NextBucket()
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
			}
		}
	}
}

func (w *HashStore) transferBuckets(starting, N int64) {
	// mask to test the newer bit
	highBit := uint64(1 << (w.newHashStore.bitsForBucket - 1))

	for bucket := starting; bucket < starting+N; bucket++ {
		// read items for the bucket
		itemsCount := int64(w.bitsCount[bucket])
		items := w.store.ReadBucket(bucket).ReadBulk(itemsCount)
		// divide items by lBit (newer bit = 0) and hBit (newer bit = 1)
		lBitBucket := make([]byte, 0, len(items)/2)
		hBitBucket := make([]byte, 0, len(items)/2)
		for _, item := range items {
			hashBit := (uint64(item[0]) + (uint64(item[1]) << 8) + (uint64(item[2]) << 16) +
				(uint64(item[3]) << 24)) & highBit
			if hashBit == highBit {
				hBitBucket = append(hBitBucket, item...)
			} else {
				lBitBucket = append(lBitBucket, item...)
			}
		}
		// put lBit and hBit items in new wallter
		w.newHashStore.store.ReadBucket(bucket).WriteBulk(lBitBucket)
		w.newHashStore.store.ReadBucket(bucket + int64(highBit)).WriteBulk(hBitBucket)
	}
	w.bitsTransferered = starting + N - 1
}

func (w *HashStore) continueDuplication(bucket int64) {
	//for bucket := int64(0); bucket < 1<<w.bitsForBucket; bucket += NBuckets {
	w.transferBuckets(bucket, NBuckets)
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
	newBucketStore := NewBucketStore(w.store.itemBytes, w.store.itemsPerBucket, w.store.headerBytes, newByteStore)
	w.newHashStore = NewHashStore(newBucketStore, int(newStoreBitsForBucket), w.operation)
	w.newHashStore.isReady = false
	w.bitsTransferered = 0
	w.continueDuplication(0)
}
