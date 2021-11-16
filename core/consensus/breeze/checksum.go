package breeze

import (
	"encoding/binary"
	"math/big"
	"sort"

	"github.com/Aereum/aereum/core/crypto"
)

const StakeRounding = 1000

type Node struct {
	Token []byte
	Stake uint64
}

type Slot struct {
	Token []byte
	Hash  *big.Int
}

type Slots []Slot

func (s Slots) Less(i, j int) bool {
	return s[i].Hash.Cmp(s[j].Hash) == -1
}

func (s Slots) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Slots) Len() int {
	return len(s)
}

func SortSlots(checksum []byte, nodes []Node) Slots {
	slots := make(Slots, 0)
	for _, node := range nodes {
		for n := uint64(0); n < node.Stake/StakeRounding; n++ {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, n)
			data = append(data, node.Token...)
			hash := crypto.Hasher(data)
			hashBigInt := big.NewInt(0)
			hashBigInt.SetBytes(hash[:])
			slots = append(slots, Slot{Token: node.Token, Hash: hashBigInt})
		} //
	}
	sort.Sort(slots)
	return slots
}
