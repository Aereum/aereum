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
	hash      Hash
	operation func(found bool, b *Bucket, item int64) OperationResult
	response  chan QueryResult
}

type walletStore struct {
	store         BucketStore
	bitsForBucket int
	mask          int64
	bitsCount     []int // number of items in the bucket
	freeOverflows []int64
	isReady       bool
	// for user when the map is doubling
	isDoubling       bool
	bitsTransferered int64
	newWallet        *walletStore
	query            chan Query
}

func (ws *walletStore) Get(q Query) QueryResult {
	hashMask := q.hash.ToInt64() & ws.mask
	wallet := ws
	if ws.isDoubling && hashMask < ws.bitsTransferered {
		hashMask = q.hash.ToInt64() & ws.newWallet.mask
		wallet = ws.newWallet
	}
	bucket := wallet.store.ReadBucket(hashMask)
	countAccounts, totalAccounts := 0, wallet.bitsCount[hashMask]
	for {
		for item := int64(0); item < bucketSize64; item++ {
			countAccounts += 1
			if countAccounts > int(totalAccounts) {
				resp := q.operation(false, bucket, item)
				ws.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts)
				return resp.result
			}
			data := bucket.Read(item)
			if q.hash.Equals(data) {
				resp := q.operation(true, bucket, item)
				ws.ProcessMutation(hashMask, resp.added, resp.deleted, countAccounts)
				return resp.result
			}
		}
		bucket = bucket.GetOverflow()
	}
}

func (ws *walletStore) ProcessMutation(hashMask int64, added *Item, deleted *Item, count int) {
	if added != nil {
		ws.bitsCount[hashMask] += 1
		if added.item == bucketSize64-1 {
			if len(ws.freeOverflows) > 0 {
				added.bucket.SetOverflow(ws.freeOverflows[0])
				ws.freeOverflows = ws.freeOverflows[1:]
			} else {
				added.bucket.SetOverflow(ws.freeOverflows[0])
			}
		}
	}
	if deleted != nil {
		lastItem := ws.bitsCount[hashMask] - 1
		ws.bitsCount[hashMask] -= 1
		if count == lastItem {
			deleted.bucket.Write(deleted.item, make([]byte, ws.store.itemBytes))
			return
		}
		var previousBucket *Bucket
		lastBucket := deleted.bucket
		for {
			if nextBucket := lastBucket.GetOverflow(); nextBucket != nil {
				previousBucket = lastBucket
				lastBucket = nextBucket
			} else {
				item := lastItem % bucketSize
				lastBucket.Write(int64(item), make([]byte, ws.store.itemBytes))
				if item == 0 && previousBucket != nil {
					ws.freeOverflows = append(ws.freeOverflows, lastBucket.n)
					previousBucket.ZeroOverflow()
				}
			}
		}
	}
}

func (w *walletStore) transferBuckets(starting, N int64) {
	// mask to test new bit
	highBit := uint64(1 << (w.newWallet.bitsForBucket - 1))
	for bucket := starting; bucket < starting+N; bucket++ {
		lBitBucket := make([]byte, 0, w.store.itemBytes*bucketSize64)
		hBitBucket := make([]byte, 0, w.store.itemBytes*bucketSize64)
		// loop over buckets and overflowbuckets
		count := 0
		itemsCount := w.bitsCount[starting]
		bucketchain := bucket
		for {
			oldBucket := w.store.ReadBucket(bucketchain)
			for n := int64(0); n < bucketSize64; n++ {
				count++
				if count > itemsCount {
					break
				}
				data := oldBucket.Read(n)
				hashBit := (uint64(data[0]) +
					(uint64(data[1]) << 8) +
					(uint64(data[2]) << 16) +
					(uint64(data[3]) << 24)) & highBit
				if hashBit == highBit {
					hBitBucket = append(hBitBucket, data...)
				} else {
					lBitBucket = append(lBitBucket, data...)
				}
			}
			oldBucket = oldBucket.GetOverflow()
			if oldBucket == nil {
				break
			}
		}
		w.newWallet.store.insertBulk(lBitBucket, bucket)
		w.newWallet.store.insertBulk(hBitBucket, bucket+int64(highBit))
	}
	w.bitsTransferered = starting + N - 1
}

func (w *walletStore) continueDuplication(bucket int64) {
	//for bucket := int64(0); bucket < 1<<w.bitsForBucket; bucket += NBuckets {
	w.transferBuckets(bucket, NBuckets)
	if bucket+NBuckets < 1<<w.bitsForBucket {
		go func() {
			sleep, _ := time.ParseDuration("10ms")
			time.Sleep(sleep)
			w.doubleJob <- bucket + NBuckets
		}()
	} else {
		w.store.store.Merge(w.newWallet.store.store)
		w.bitsForBucket = w.newWallet.bitsForBucket
		w.mask = w.newWallet.mask
		w.bitsCount = w.newWallet.bitsCount
		w.freeOverflows = w.newWallet.freeOverflows
		w.store.bucketCount = w.newWallet.store.bucketCount
		w.isDoubling = false
		w.bitsTransferered = 0
		w.newWallet = nil
		w.isReady = true
	}
}

func (w *walletStore) startDuplication() {
	w.isDoubling = true
	w.newWallet = CreateWalletStore("wallet_duplicate.tmp", byte(w.bitsForBucket)+1, 0)
	w.newWallet.isReady = false
	w.bitsTransferered = 0
	w.continueDuplication(0)
}
