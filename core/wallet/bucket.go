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
package wallet

import (
	"bytes"
	"encoding/binary"
)

const maxCloningBlockSize = 1 << 20

// BucketStore is a sequential appendable collection of buckets of equal size.
// Each bucket consists of a fixed number of items + a link to a next bucket.
// if the link is zero, the bucket is the final bucket in a chain.
// In order to allow bucket instant cloning, BucketStore can be set into
// journaling mode were any item writes or overflow link writes will be
// registered to journal in order to retrieve the original state after the
// mutating store is cloned into memory or into file.
type BucketStore struct {
	bytes          ByteStore // ByteStore size = bucketCount * bucketBytes + headerBytes
	bucketCount    int64     // number of buckets
	itemBytes      int64     // bytes per item
	bucketBytes    int64     // bytes per bucket = items per bucket * bytes per item + 8
	itemsPerBucket int64     // items per bucket
	headerBytes    int64     // bytes alocated for header
	isCloning      bool      // for cloning state
	journal        JournalStore
	cloning        JournalStore
	bucketsCloned  int64 // number of buckets already cloned
	bucketToClone  int64 // number of buckets at the start of cloning
}

type Bucket struct {
	n       int64  // n = sequential position in the bucket store
	data    []byte // memory cache of data (size = bucketBytes)
	buckets *BucketStore
}

func (b *BucketStore) toJournal(bucket int64, item byte, oldData, newData []byte) {
	data := make([]byte, 2*len(oldData)+9)
	binary.BigEndian.PutUint64(data[0:8], uint64(bucket))
	data[8] = item
	copy(data[9:9+b.bucketBytes], oldData)
	copy(data[9+b.bucketBytes:9+2*b.bucketBytes], newData)
	b.journal.Append(data)
}

// create a new bucket store of size itemBytes
func NewBucketStore(itemBytes, itemsPerBucket int64, bytes ByteStore) *BucketStore {
	headerBytes := int64(56) // epoch + hash + itemBytes + itemsPerBucket
	if (bytes.Size()-headerBytes)%(itemsPerBucket*itemBytes+8) != 0 {
		panic("ByteStore size incompatible with bucket store")
	}
	header := make([]byte, headerBytes)
	binary.LittleEndian.PutUint64(header[40:48], uint64(itemBytes))
	binary.LittleEndian.PutUint64(header[48:56], uint64(itemsPerBucket))
	return &BucketStore{
		bytes:          bytes,
		bucketCount:    (bytes.Size() - headerBytes) / (itemsPerBucket*itemBytes + 8),
		itemBytes:      itemBytes,
		bucketBytes:    (itemsPerBucket*itemBytes + 8),
		itemsPerBucket: itemsPerBucket,
		headerBytes:    headerBytes,
	}
}

// Read the n-th sequential (begining at zero up to bucketCount - 1)
func (b *BucketStore) ReadBucket(n int64) *Bucket {
	return &Bucket{
		n:       n,
		data:    b.bytes.ReadAt(b.headerBytes+n*b.bucketBytes, b.bucketBytes),
		buckets: b,
	}
}

// Append a bucket to the store and associate it as the next bucket of the
// current bucket.
func (b *BucketStore) Append() *Bucket {
	data := make([]byte, b.bucketBytes)
	b.bytes.Append(data)
	n := b.bucketCount
	b.bucketCount++
	return &Bucket{
		n:       n,
		data:    data,
		buckets: b,
	}
}

// Inserts a sequential collection of items exposed as a bytearray of size
// # of items * itemBytes into a chain of buckets starting at the calling bucket
// and returns the number of buckets appended.
// Writebulk does not preseve journaling semantics
func (b *Bucket) WriteBulk(data []byte) int64 {
	store := b.buckets.bytes
	itemBytes := b.buckets.itemBytes
	remaning := int64(len(data)) / b.buckets.itemBytes
	if len(data)%int(b.buckets.itemBytes) != 0 {
		panic("incongruent data")
	}
	processed := int64(0)
	bucket := b
	appended := 0
	for {
		if remaning == 0 {
			return int64(appended)
		} else if remaning <= b.buckets.itemsPerBucket {
			offset := bucket.n*b.buckets.bucketBytes + b.buckets.headerBytes
			bucketData := data[processed*itemBytes : (processed+remaning)*itemBytes]
			store.WriteAt(offset, bucketData)
			return int64(appended)
		} else {
			offset := bucket.n*b.buckets.bucketBytes + b.buckets.headerBytes
			bucketData := data[processed*itemBytes : (processed+remaning)*itemBytes]
			store.WriteAt(offset, bucketData)
		}
		bucket = b.buckets.Append()
		appended++
	}
}

// Read #count items starting at current bucket and following linked buckets.
// Panics if there are not enough items in the bucket chain to match count.
func (b *Bucket) ReadBulk(count int64) [][]byte {
	data := make([][]byte, count)
	bucket := b
	bItem := int64(0) // item on a
	for n := int64(0); n < count; n++ {
		data[n] = bucket.data[bItem*b.buckets.itemBytes : (bItem+1)*b.buckets.itemBytes]
		bItem++
		if bItem%b.buckets.itemsPerBucket == 0 {
			bItem = 0
			bucket = bucket.NextBucket()
			if bucket == nil {
				panic("could not read enough items")
			}
		}
	}
	return data
}

// Saves the item with data (of size itemBytes) into the ByteStore and the
// cached bucket data.
// Panics if either item or data size are outside bucketstore specification.
// Applies to journaling
func (b *Bucket) WriteItem(item int64, data []byte) {
	if item > b.buckets.itemsPerBucket || item < 0 {
		panic("invalid bucket read")
	}
	if len(data) != int(b.buckets.itemBytes) {
		panic("bucket only read entire items")
	}
	if (item+1)*b.buckets.itemBytes > b.buckets.bucketBytes-8 {
		panic(errOverflow)
	}
	if b.buckets.isCloning {
		b.buckets.toJournal(b.n, byte(item),
			b.data[item*b.buckets.itemBytes:(item+1)*b.buckets.itemBytes], data)
	}
	copy(b.data[item*b.buckets.itemBytes:(item+1)*b.buckets.itemBytes], data)
	offset := b.buckets.headerBytes + b.n*b.buckets.bucketBytes + item*b.buckets.itemBytes
	b.buckets.bytes.WriteAt(offset, data)
}

// Read content of the item in the bucket.
// Panics if item is outside bucketstore specification
func (b *Bucket) ReadItem(item int64) []byte {
	if item > b.buckets.itemsPerBucket || item < 0 {
		panic("invalid bucket read")
	}
	return b.data[item*b.buckets.itemBytes : (item+1)*b.buckets.itemBytes]
}

// Get next bucket in the bucket chain (return nil if it is a final bucket)
func (b *Bucket) NextBucket() *Bucket {
	overflow := b.ReadOverflow()
	if overflow == 0 {
		return nil
	}
	return b.buckets.ReadBucket(overflow)
}

// Read the sequential numbering of the next bucket in the bucket chain.
// Returns zero if it is a final bucket.
func (b *Bucket) ReadOverflow() int64 {
	overflow := int64(binary.LittleEndian.Uint64(b.data[b.buckets.itemsPerBucket*b.buckets.itemBytes:]))
	return overflow
}

// Write overflow as next bucket in the bucket chain.
// Applies to journaling
func (b *Bucket) WriteOverflow(overflow int64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(overflow))
	if b.buckets.isCloning {
		b.buckets.toJournal(b.n, 255,
			data[b.buckets.itemsPerBucket*b.buckets.itemBytes:], data)
	}
	copy(b.data[b.buckets.itemsPerBucket*b.buckets.itemBytes:], data)
	offset := b.n*b.buckets.bucketBytes + b.buckets.itemsPerBucket*b.buckets.itemBytes + b.buckets.headerBytes
	b.buckets.bytes.WriteAt(offset, data)
}

// Append a new Bucket into the store, marks it as the next bucket of the
// the current bucket and returns the new bucket.
func (b *Bucket) AppendOverflow() *Bucket {
	overflow := b.buckets.Append()
	b.WriteOverflow(overflow.n)
	return overflow
}

// Recreate state of the bucket at clone request by undoing all the alterations
// processed in the journal.
// We follow the journal from end to start, checkin the state and undoing all
// the modifications
// The clone ByteStore is modificated in the process.
func RecreateBucket(clone ByteStore, journal ByteStore) *BucketStore {
	// read header
	itemBytes := int64(binary.LittleEndian.Uint64(clone.ReadAt(40, 8)))
	itemsPerBucket := int64(binary.LittleEndian.Uint64(clone.ReadAt(48, 8)))
	bs := NewBucketStore(itemBytes, itemsPerBucket, clone)
	eof := journal.Size()
	journalEntry := 2*itemBytes + 9
	for position := eof - journalEntry; position >= 0; position -= journalEntry {
		entry := journal.ReadAt(position, journalEntry)
		bucketPosition := int64(binary.LittleEndian.Uint64(entry[0:8]))
		bucket := bs.ReadBucket(bucketPosition)
		item := int64(entry[8])
		if item == 255 {
			oldOverflow := int64(binary.LittleEndian.Uint64(entry[9:17]))
			newOverflow := int64(binary.LittleEndian.Uint64(entry[17:25]))
			if bucket.ReadOverflow() != newOverflow {
				panic("clone and journal are incompatible")
			}
			bs.Append().WriteOverflow(oldOverflow)
		} else {
			itemBytes := bucket.ReadItem(item)
			if !bytes.Equal(itemBytes, entry[9:9+bs.itemBytes]) {
				panic("clone and journal are incompatible")
			}
			bucket.WriteItem(item, entry[9:9+bs.itemBytes])
		}
	}
	return bs
}
