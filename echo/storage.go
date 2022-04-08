package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Aereum/aereum/core/instructions"
)

type Uint64Storage interface {
	ReadAt(int64, int64) []uint64
	WriteAt(int64, []uint64) error
	Append([]uint64)
	Size() int64
	Close()
}

type Uint64RawFile struct {
	file *os.File
}

func (r *Uint64RawFile) ReadAt(offset int64, count int64) []uint64 {
	bytes := make([]uint64, 8*count)
	if n, err := r.file.ReadAt(bytes, offset*8); err != nil || n != 8*count {
		return nil
	}
	intArray := make([]uint64, count)
	for n:=int64(0); n < count; n++ {
		intArray[n] = binary.LittleEndian.Uint64(bytes[n*8:(n+1)*8])
	}
	return intArray
}

func (r *Uint64RawFile) WriteAt(offset int64, data []uint64) error {
	bytes := make([]byte, 8*len(data))
	for n:=0; n < len(data); n++ {
		binary.LittleEndian.PutUint64(bytes[n*8:(n+1)*8],data[n])
	}
	if n, err:=r.file.WriteAt(bytes, offset*8); err != nil || n != len(data)*8 {
		return fmt.Errorf("IO Error")
	}
	return nil
}

type ByteStorage interface {
	ReadAt(int64, int64) []byte
	WriteAt(int64, []byte)
	Append([]byte)
	Size() int64
	Close()
}

type IndexStore struct {
	io      ByteStorage
	buckets uint64
}

func NewIndexStore(storage ByteStorage) *IndexStore {
	return &IndexStore{
		io:      storage,
		buckets: 0,
	}
}

func (idx *IndexStore) New


func (us UintStore) Put(v uint64) {
	us.data[0] = byte(v)
	us.data[1] = byte(v >> 8)
	us.data[2] = byte(v >> 16)
	us.data[3] = byte(v >> 24)
	us.data[4] = byte(v >> 32)
	us.data[5] = byte(v >> 40)
	us.data[6] = byte(v >> 48)
	us.data[7] = byte(v >> 56)
	if n, err := us.io.Write(us.data); n != 8 || err != nil {
		return
	}
}

type Position struct {
	Start uint64
	Size  uint16
}

type Storage struct {
	Positions []Position
	Buffer    io.ReadWriteCloser
	Sizes     io.WriteCloser
}

func (s *Storage) ReadNth(n uint64) instructions.Instruction {
	return nil
}

func (s *Storage) Persist(instruction instructions.Instruction) {
	bytes := instruction.Serialize()
	if n, err := s.Buffer.Write(bytes); n != len(bytes) || err != nil {
		return
	}
	blen := make([]byte, 2)
	binary.LittleEndian.PutUint16(blen, uint16(len(bytes)))
	if n, err := s.Sizes.Write(blen); n != 2 || err != nil {
		return
	}
	if len(s.Positions) > 0 {
		last := s.Positions[len(s.Positions)-1]
		newPosition := Position{
			Start: last.Start + uint64(last.Size),
			Size:  uint16(len(bytes)),
		}
		s.Positions = append(s.Positions, newPosition)
	} else {
		newPosition := Position{
			Start: 0,
			Size:  uint16(len(bytes)),
		}
		s.Positions = append(s.Positions, newPosition)
	}
}
