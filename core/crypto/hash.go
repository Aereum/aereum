package crypto

import (
	"bytes"
	"crypto/sha256"
)

type Hash [Size]byte

var ZeroHash Hash = Hasher([]byte{})

func (hash Hash) ToInt64() int64 {
	return int64(hash[0]) + (int64(hash[1]) << 8) + (int64(hash[2]) << 16) + (int64(hash[3]) << 24)
}

func BytesToHash(bytes []byte) Hash {
	var hash Hash
	if len(bytes) != Size {
		return hash
	}
	copy(hash[:], bytes)
	return hash
}

func (h Hash) Equal(another Hash) bool {
	return h == another
}

func (h Hash) Equals(another []byte) bool {
	if len(another) < Size {
		return false
	}
	return bytes.Equal(h[:], another[:Size])
}

func Hasher(data []byte) Hash {
	return Hash(sha256.Sum256(data))
}

func HashToken(token Token) Hash {
	return Hash(sha256.Sum256(token[:]))
}
