package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"time"
)

// hashdb implements a hash disk storage solution with a simple interface
// Exits(hash) bool
// Remove(hash) bool
// Add(hash) bool

// these configurations should be altered latter
const segmentsBytes = 2
const maxHashesPerSegment = 10

const size = sha256.Size

var errSegmentIsFull = errors.New("hash segment is full")

type Hash [size]byte

type MultipleFileHashDB struct {
	multipleFiles bool
	singleFile    *os.File
	files         []*os.File
	segmentsBytes int
	collisions    map[Hash]struct{}
}

func (m *MultipleFileHashDB) Exists(hash Hash) (bool, error) {
	if _, ok := m.collisions[hash]; ok {
		return true, nil
	}
	return m.FindHash(hash, false, false)
}

func (m *MultipleFileHashDB) RemoveIfExists(hash Hash) (exists bool, error error) {
	if _, ok := m.collisions[hash]; ok {
		delete(m.collisions, hash)
		return true, nil
	}
	return m.FindHash(hash, false, true)
}

func (m *MultipleFileHashDB) InsertIfExists(hash Hash) (exists bool, error error) {
	if _, ok := m.collisions[hash]; ok {
		return true, nil
	}
	exists, err := m.FindHash(hash, true, false)
	if err != nil {
		if err == errSegmentIsFull {
			m.collisions[hash] = struct{}{}
			return false, nil
		} else {
			return exists, err
		}
	}
	return exists, nil
}

func (m *MultipleFileHashDB) FindHash(hash Hash, insert bool, remove bool) (exists bool, error error) {

	var file *os.File
	var segment, empty, sizeOnSegment int
	var hashEnd []byte
	if m.multipleFiles {
		file = m.files[hash[0]]
		segment = int(hash[1])
		for n := 1; n < m.segmentsBytes; n++ {
			segment += int(hash[1+n]) << (8 * n)
		}
		hashEnd = hash[1+m.segmentsBytes+1:]
		sizeOnSegment = (size - m.segmentsBytes - 1)
	} else {
		file = m.singleFile
		segment = int(hash[0])
		for n := 1; n < m.segmentsBytes; n++ {
			segment += int(hash[n]) << (8 * n)
		}
		hashEnd = hash[m.segmentsBytes+1:]
		sizeOnSegment = (size - m.segmentsBytes)
	}
	segmentSize := maxHashesPerSegment * sizeOnSegment
	data := make([]byte, segmentSize)
	segmentStart := int64(segment * segmentSize)
	_, err := file.ReadAt(data, segmentStart)
	if err != nil {
		// TODO: EOF error
		return false, err
	}
	empty = -1
	zeros := make([]byte, sizeOnSegment)
	for n := 0; n < maxHashesPerSegment; n++ {
		if bytes.Equal(hashEnd, data[n*sizeOnSegment:(n+1)*sizeOnSegment]) {
			if remove {
				position := segmentStart + int64(n*sizeOnSegment)
				_, err := file.WriteAt(zeros, position)
				return true, err
			}
			return true, nil
		}
		if insert && (empty == -1) {
			if bytes.Equal(zeros, data[n*sizeOnSegment:(n+1)*sizeOnSegment]) {
				empty = n
			}
		}
	}
	if insert {
		if empty != -1 {
			position := int64(empty*sizeOnSegment) + segmentStart
			_, err := file.WriteAt(hashEnd, position)
			return false, err
		} else {
			return false, errSegmentIsFull
		}
	}
	return false, nil
}

func OpenOrCreateSingleFileHashDb(filePath string, segmentsBytes int) *MultipleFileHashDB {
	db := &MultipleFileHashDB{
		segmentsBytes: segmentsBytes,
		multipleFiles: false,
		collisions:    make(map[Hash]struct{}),
	}
	segments := 1 << (8 * segmentsBytes)
	segmentSize := maxHashesPerSegment * (size - segmentsBytes)
	fileSize := segmentSize * segments
	if stats, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err != nil {
			return nil
		} else {
			bytes := make([]byte, fileSize)
			n, err := file.Write(bytes)
			if err != nil || n != len(bytes) {
				return nil
			}
			db.singleFile = file
			return db
		}
	} else {
		if stats.Size() != int64(fileSize) {
			return nil
		}
		file, err := os.OpenFile(filePath, os.O_RDWR, os.ModeExclusive)
		if err != nil {
			return nil
		}
		db.singleFile = file
		return db
	}
}

func main() {
	var db = OpenOrCreateSingleFileHashDb("hash_teste.tmp", 3)
	if db == nil {
		fmt.Println("Could not access database")
		return
	}
	b := make([]byte, size)
	var h Hash
	start := time.Now()
	for r := 0; r < 500000; r++ {
		rand.Read(b)
		for n := 0; n < size; n++ {
			h[n] = b[n]
		}
		_, err := db.InsertIfExists(h)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}
	fmt.Println(time.Since(start))
	fmt.Println(len(db.collisions))
}
