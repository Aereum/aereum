package index

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/Aereum/aereum/core/crypto"
)

const indexSize = crypto.TokenSize + 8

type persistentMap struct {
	tokenOrder map[crypto.Token]uint64
	buckets    []uint64
	storage    IO
}

func (m *persistentMap) Close() {
	m.storage.Close()
}

func (m *persistentMap) Get(token crypto.Token) uint64 {
	pos, ok := m.tokenOrder[token]
	if !ok {
		return 0
	}
	return m.buckets[pos]
}

func (m *persistentMap) Set(token crypto.Token, bucket uint64) {
	pos, ok := m.tokenOrder[token]
	if !ok {
		bytes := make([]byte, 8, indexSize)
		binary.LittleEndian.PutUint64(bytes, bucket)
		bytes = append(bytes, token[:]...)
		WriteOrPanic(m.storage, bytes, int64(len(m.buckets)*indexSize))
		m.tokenOrder[token] = uint64(len(m.buckets))
		m.buckets = append(m.buckets, bucket)
		return
	}
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bucket)
	WriteOrPanic(m.storage, bytes, int64(pos)*indexSize+8)
	m.buckets[pos] = bucket
}

func openPersistentMap(FileName string) *persistentMap {
	file, err := os.OpenFile(FileName, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes := make([]byte, indexSize)
	newMap := persistentMap{
		tokenOrder: make(map[crypto.Token]uint64),
		buckets:    make([]uint64, 0),
		storage:    file,
	}
	for {
		if n, err := file.Read(bytes); n != indexSize {
			if err == io.EOF {
				return &newMap
			}
			return nil
		}
		var token crypto.Token
		copy(token[:], bytes[8:])
		if _, ok := newMap.tokenOrder[token]; ok {
			return nil
		}
		bucket := binary.LittleEndian.Uint64(bytes[0:8])
		newMap.tokenOrder[token] = uint64(len(newMap.buckets))
		newMap.buckets = append(newMap.buckets, bucket)
	}
}
