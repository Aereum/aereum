package hashdb

/*import (
	"crypto/sha256"
	"sync"
)

const (
	size       = sha256.Size
	bucketSize = 6
)

type Hash [size]byte

var zeroHash = Hash{}

type account struct {
	hash    Hash
	balance int
}

type bucket struct {
	acc      [bucketSize]account
	overflow int
}

type walletStore struct {
	bitsForBucket int
	mask          int
	data          []bucket
	bucketCount   []uint8
}

type Wallet struct {
	wallet *walletStore
	*sync.Mutex
}

func NewWallet(bitsForBucket int) *Wallet {
	if bitsForBucket > 4*8 || bitsForBucket < 0 {
		panic("invalid wallet parameters")
	}
	return &Wallet{
		wallet: &walletMap{
			bitsForBucket: bitsForBucket,
			mask:          1<<bitsForBucket - 1,
			bucketCount:   make([]uint8, 1<<bitsForBucket),
			data:          make([]account, bucketSize*(1<<(bitsForBucket))),
		},
	}
}

func (w *Wallet) Withdraw(hash Hash, value int) boll {
	position := int(hash[0]) + (int(hash[1]) << 8) + (int(hash[2]) << 16) + (int(hash[3]) << 24)
	position = bucket & w.wallet.mask
	bucket = w.wallet.data[position]
	for n := 0; n < bucketCount; n++ {
		bucket[n].has
	}
}
*/
