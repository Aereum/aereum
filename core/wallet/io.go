package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	BeginOfFile int = iota
	CurrentPosition
	EndOfFile
)

var errOverflow = errors.New("outside store")
var errWhence = errors.New("unrecognized whence")

type ByteStore interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Seek(int64, int) (int64, error)
	New(int64) ByteStore
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

func (f *FileStore) Seek(offset int64, whence int) (int64, error) {
	return f.data.Seek(offset, whence)
}

func (f *FileStore) Write(b []byte) (int, error) {
	return f.data.Write(b)
}

func (f *FileStore) Read(b []byte) (int, error) {
	return f.data.Read(b)
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
	length   int64
	data     []byte
	position int64
}

func (m *MemoryStore) Size() int64 {
	return m.length
}

func (m *MemoryStore) Seek(offset int64, whence int) (int64, error) {
	var newPosition int64
	if whence == BeginOfFile {
		newPosition = offset
	} else if whence == CurrentPosition {
		newPosition += offset
	} else if whence == EndOfFile {
		newPosition = m.length - offset
	} else {
		return -1, errWhence
	}
	if newPosition < 0 {
		return -1, errOverflow
	}
	if newPosition >= m.length {
		return -1, errOverflow
	}
	return newPosition, nil
}

func (m *MemoryStore) Write(b []byte) (int, error) {
	len64 := int64(len(b))
	overflow := m.position + len64 - m.length
	if overflow > 0 {
		copy(m.data[m.position:m.length], b[0:m.length-m.position])
		m.data = append(m.data, b[overflow:]...)
		m.length += overflow
	} else {
		copy(m.data[m.position:m.position+len64], b[0:m.length-m.position])
	}
	m.position += len64
	return len(b), nil
}

func (m *MemoryStore) Read(b []byte) (int, error) {
	len64 := int64(len(b))
	overflow := m.position + len64 - m.length
	if overflow > 0 {
		copy(b[0:overflow], m.data[m.position:])
		return int(overflow), errOverflow
	}
	copy(b, m.data[m.position:m.position+len64])
	return len(b), nil
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
	m.length = newStore.length
	m.position = 0
}

func NewMemoryStore(size int64) *MemoryStore {
	return &MemoryStore{
		data:     make([]byte, size),
		length:   size,
		position: 0,
	}
}
