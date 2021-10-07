package main

import (
	"errors"
	"fmt"
	"os"
)

var errOverflow = errors.New("outside store")
var errWhence = errors.New("unrecognized whence")

// ByteStore is a panic on error semantics to store and retrieve raw bytes
// on any medium.
// It panics on any IO error and if ReadAt or WriteAt cannot be executed
// at the current store size. Use append to enlarge the store.
// New creates a new bytestore. If filebased, it should be interpreted as a
// temporary file. On Merge, the temporary file is renamed to the current file.
// On memory, current data is released to garbage collector and temporary data
// is promoted to main data.
type ByteStore interface {
	ReadAt(int64, int64) []byte
	WriteAt(int64, []byte)
	Append([]byte)
	New(int64) ByteStore // create a new empty bytestore os size int64
	Merge(ByteStore)
	Size() int64
}

type FileStore struct {
	name string
	size int64
	data *os.File
}

func NewFileStore(name string, size int64) *FileStore {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	file.Close()
	if err := os.Truncate(name, size); err != nil {
		panic(err)
	}
	file, err = os.OpenFile(name, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		panic(err)
	}
	return &FileStore{
		name: name,
		size: size,
		data: file,
	}
}

func (f *FileStore) Size() int64 {
	return f.size
}

func (f *FileStore) WriteAt(offset int64, b []byte) {
	if offset+int64(len(b)) >= f.size || offset < 0 {
		panic("invalid offset")
	}
	if _, err := f.data.Seek(offset, 0); err != nil {
		panic(err)
	}
	if n, err := f.data.Write(b); n != len(b) {
		panic(err)
	}
}

func (f *FileStore) Append(b []byte) {
	if _, err := f.data.Seek(0, 2); err != nil {
		panic(err)
	}
	if n, err := f.data.Write(b); n != len(b) {
		panic(err)
	}
	f.size += int64(len(b))
}

func (f *FileStore) ReadAt(offset int64, nbytes int64) []byte {
	if offset+nbytes >= f.size || offset < 0 || nbytes < 1 {
		panic("invalid read parameters")
	}
	data := make([]byte, nbytes)
	if _, err := f.data.Seek(offset, 0); err != nil {
		panic(err)
	}
	if n, err := f.data.Read(data); int64(n) != nbytes {
		panic(err)
	}
	return data
}

func (f *FileStore) New(size int64) ByteStore {
	return NewFileStore(fmt.Sprintf("%v_temp", f.name), size)
}

func (f *FileStore) Merge(another ByteStore) {
	other, ok := another.(*FileStore)
	if !ok {
		panic("can only merge FileStore with FileStore")
	}
	f.data.Close()
	other.data.Close()
	os.Rename(other.name, f.name)
	file, err := os.OpenFile(f.name, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	f.data = file
	f.size = other.size
}

type MemoryStore struct {
	data []byte
}

func NewMemoryStore(size int64) *MemoryStore {
	return &MemoryStore{
		data: make([]byte, size),
	}
}

func (m *MemoryStore) Size() int64 {
	return int64(len(m.data))
}

func (m *MemoryStore) WriteAt(offset int64, b []byte) {
	if offset+int64(len(b)) >= int64(len(m.data)) || offset < 0 {
		panic("invalid offset")
	}
	copy(m.data[offset:offset+int64(len(b))], b)
}

func (m *MemoryStore) Append(b []byte) {
	m.data = append(m.data, b...)
}

func (m *MemoryStore) ReadAt(offset int64, ncount int64) []byte {
	if offset+ncount >= int64(len(m.data)) || offset < 0 {
		panic("invalid offset")
	}
	data := make([]byte, ncount)
	copy(data, m.data[offset:offset+ncount])
	return data
}

func (m *MemoryStore) New(size int64) ByteStore {
	return NewMemoryStore(size)
}

func (m *MemoryStore) Merge(another ByteStore) {
	newStore, ok := another.(*MemoryStore)
	if !ok {
		panic("MemoryStore can only be merged with memory store")
	}
	m.data = newStore.data
}
