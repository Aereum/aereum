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
package hashdb

/*
import (
	"bytes"
	"errors"
	"os"
)

// these configurations should be altered latter
const maxHashesPerSegment = 10

var errSegmentIsFull = errors.New("hash segment is full")
var errFileOfWrongSize = errors.New("data file size incompatible with segments bytes")

type HashStore interface {
	Exists(hash Hash) (bool, error)
	RemoveIfExists(hash Hash) (bool, error)
	InsertIfNotExists(hash Hash) (bool, error)
}

type ReaderWriterAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
	WriteAt(p []byte, off int64) (n int, err error)
}

type MemoryHashStore struct {
	data map[Hash]struct{}
}

func (m *MemoryHashStore) Exists(hash Hash) (bool, error) {
	_, ok := m.data[hash]
	return ok, nil
}

func (m *MemoryHashStore) RemoveIfExists(hash Hash) (bool, error) {
	_, ok := m.data[hash]
	if ok {
		delete(m.data, hash)
	}
	return ok, nil
}

func (m *MemoryHashStore) InsertIfNotExists(hash Hash) (bool, error) {
	_, ok := m.data[hash]
	if !ok {
		m.data[hash] = struct{}{}
	}
	return ok, nil
}

type MultipleFileHashDB struct {
	multipleFiles bool
	singleFile    ReaderWriterAt
	files         []ReaderWriterAt
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

func (m *MultipleFileHashDB) InsertIfNotExists(hash Hash) (exists bool, error error) {
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

	var file ReaderWriterAt
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

func OpenOrCreateFile(filePath string, segmentsBytes int) (*os.File, error) {
	segments := 1 << (8 * segmentsBytes)
	segmentSize := maxHashesPerSegment * (size - segmentsBytes)
	fileSize := segmentSize * segments
	if stats, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err != nil {
			return nil, err
		} else {
			bytes := make([]byte, fileSize)
			n, err := file.Write(bytes)
			if err != nil || n != len(bytes) {
				return nil, err
			}
			return file, nil
		}
	} else {
		if stats.Size() != int64(fileSize) {
			return nil, errFileOfWrongSize
		}
		file, err := os.OpenFile(filePath, os.O_RDWR, os.ModeExclusive)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func SingleFileHashDb(file ReaderWriterAt, segmentsBytes int) *MultipleFileHashDB {
	return &MultipleFileHashDB{
		segmentsBytes: segmentsBytes,
		multipleFiles: false,
		collisions:    make(map[Hash]struct{}),
		singleFile:    file,
	}
}

func MultipleFileHashDb(files []ReaderWriterAt, segmentsBytes int) *MultipleFileHashDB {
	return &MultipleFileHashDB{
		segmentsBytes: segmentsBytes,
		multipleFiles: true,
		collisions:    make(map[Hash]struct{}),
		files:         files,
	}
}
*/
