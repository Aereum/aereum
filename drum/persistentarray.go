package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/Aereum/aereum/core/util"
)

// append only persistent bytearray with log.Fatal on io error logic
type PersistentByteArray struct {
	filename string
	io       io.ReadWriteCloser
}

func OpenParsistentByteArray(fileName string) *PersistentByteArray {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Could not open instructions file.")
	}
	return &PersistentByteArray{
		filename: fileName,
		io:       file,
	}
}

func (p *PersistentByteArray) Append(data []byte) {
	prependData := make([]byte, 8+len(data))
	util.PutUint64(uint64(len(data)), &data)
	prependData = append(prependData, data...)
	if n, err := p.io.Write(data); err != nil || n != len(prependData) {
		log.Fatal("Could not save instruction to disk.")
	}
}

func (p *PersistentByteArray) Read() []byte {
	length := make([]byte, 8)
	if n, _ := p.io.Read(length); n != 8 {
		log.Fatalf("Could not parse file %v", p.filename)
	}
	data := make([]byte, int(binary.LittleEndian.Uint64(length)))
	if n, err := p.io.Read(length); n != len(data) {
		log.Fatalf("Could not parse file %v", p.filename)
	} else if err == io.EOF {
		p.io.Close()
		if file, err := os.OpenFile(p.filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err != nil {
			log.Fatalf("Could not open file %v in append mode.", p.filename)
		} else {
			p.io = file
		}
	} else if err != nil {
		log.Fatalf("Could not parse file %v", p.filename)
	}
	return data
}
