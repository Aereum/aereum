package main

import (
	"encoding/binary"
	"errors"
)

const (
	//header = int(9)
	//header64     = int64(header)
	bucketSize   = int(6)
	bucketSize64 = int64(bucketSize)
	//bucketBytes   = int(bucketSize*(size+8) + 8)
	//bucketBytes64 = int64(bucketBytes)
)

var errBucketOverflow = errors.New("bucket overflow")

// BucketStore (BS) is a sequential indefinite size collection of buckets, which
// are collections of BS.size items each of size BS.itemBytes and a link (the
// sequential number of the chained bucket) to another bucket.
//
// bucketbytes = size * itemBytes + sizeof(int64)
// bucketCount * bucketBytes + header = size of bytes
type BucketStore struct {
	bytes       ByteStore
	bucketCount int64 // number of buckets
	itemBytes   int64
	bucketBytes int64
	size        int64
	headerBytes int64
}

// n = sequential position in the bucket store
// size of data = itemBytes
type Bucket struct {
	n       int64
	data    []byte
	buckets *BucketStore
}

// create a new bucket store of size itemBytes
func NewBucketStore(itemBytes, headerBytes int64, bytes ByteStore) *BucketStore {
	if (bytes.Size()-headerBytes)%(bucketSize64*itemBytes+8) != 0 {
		panic("ByteStore size incompatible with bucket store")
	}
	return &BucketStore{
		bytes:       bytes,
		bucketCount: (bytes.Size() - headerBytes) / (bucketSize64*itemBytes + 8),
		itemBytes:   itemBytes,
		size:        bucketSize64,
		bucketBytes: (bucketSize64*itemBytes + 8),
	}
}

func (b *BucketStore) ReadBucket(n int64) *Bucket {
	return &Bucket{
		n:       n,
		data:    b.bytes.ReadAt(b.headerBytes+n*b.bucketBytes, b.bucketBytes),
		buckets: b,
	}
}

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

// Writebulk inserts a sequential collection of items exposed as a bytearray
// into a chain of buckets starting at the calling bucket and returns the number
// of buckets created.
func (b *Bucket) WriteBulk(data []byte) {
	store := b.buckets.bytes
	itemBytes := b.buckets.itemBytes
	remaning := int64(len(data)) / b.buckets.itemBytes
	if len(data)%int(b.buckets.itemBytes) != 0 {
		panic("incongruent data")
	}
	processed := int64(0)
	bucket := b
	for {
		if remaning == 0 {
			return
		} else if remaning <= b.buckets.size {
			offset := bucket.n*b.buckets.bucketBytes + b.buckets.headerBytes
			bucketData := data[processed*itemBytes : (processed+remaning)*itemBytes]
			store.WriteAt(offset, bucketData)
			return
		} else {
			offset := bucket.n*b.buckets.bucketBytes + b.buckets.headerBytes
			bucketData := data[processed*itemBytes : (processed+remaning)*itemBytes]
			store.WriteAt(offset, bucketData)
		}
		bucket = b.buckets.Append()
	}
}

// Read count items starting at current bucket and and the subsequent linked
// buckets. Panics if there are not enough buckets to read itemCount items.
func (b *Bucket) ReadBulk(count int64) []byte {
	data := make([]byte, 0, 2*b.buckets.itemBytes*b.buckets.size)
	processed := int64(0)
	bucket := b
	for {
		remaning := count - processed
		if remaning == 0 {
			return data
		} else if remaning < b.buckets.size {
			data = append(data, bucket.data[0:remaning*b.buckets.itemBytes]...)
			return data
		} else {
			data = append(data, bucket.data[0:b.buckets.size*b.buckets.itemBytes]...)
		}
		bucket = bucket.NextBucket()
		if bucket == nil {
			panic("could not read enough items")
		}
	}
}

// Write saves the item with data of size itemBytes into the ByteStore and the
// bucket data.
func (b *Bucket) WriteItem(item int64, data []byte) {
	if item > b.buckets.size || item < 0 {
		panic("invalid bucket read")
	}
	if len(data) != int(b.buckets.itemBytes) {
		panic("bucket only read entire items")
	}
	if (item+1)*b.buckets.itemBytes > b.buckets.bucketBytes-8 {
		panic(errOverflow)
	}
	copy(b.data[item*b.buckets.itemBytes:(item+1)*b.buckets.itemBytes], data)
	offset := b.buckets.headerBytes + b.n*b.buckets.bucketBytes + item*b.buckets.itemBytes
	b.buckets.bytes.WriteAt(offset, data)
}

// Read content of the item in the bucket
func (b *Bucket) ReadItem(item int64) []byte {
	if item > b.buckets.size || item < 0 {
		panic("invalid bucket read")
	}
	return b.data[item*b.buckets.itemBytes : (item+1)*b.buckets.itemBytes]
}

// Get the overflow link of the bucket (zero if a final bucket)
func (b *Bucket) NextBucket() *Bucket {
	overflow := int64(binary.LittleEndian.Uint64(b.data[b.buckets.size*b.buckets.itemBytes:]))
	if overflow > b.buckets.bucketCount {
		panic("bucket overflow")
	}
	if overflow == 0 {
		return nil
	}
	return b.buckets.ReadBucket(overflow)
}

func (b *Bucket) ReadOverflow() int64 {
	overflow := int64(binary.LittleEndian.Uint64(b.data[b.buckets.size*b.buckets.itemBytes:]))
	return overflow
}

func (b *Bucket) WriteOverflow(overflow int64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(overflow))
	copy(b.data[b.buckets.size*b.buckets.itemBytes:], data)
	offset := b.n*b.buckets.bucketBytes + b.buckets.size*b.buckets.itemBytes + b.buckets.headerBytes
	b.buckets.bytes.WriteAt(offset, data)
}

func (b *Bucket) AppendOverflow() *Bucket {
	overflow := b.buckets.Append()
	b.WriteOverflow(overflow.n)
	return overflow
}
