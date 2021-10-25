package hashdb

/*import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

func (one Hash) Equal(another Hash) bool {
	for n := 0; n < size; n++ {
		if one[n] != another[n] {
			return false
		}
	}
	return true
}

const accountSize = size + 8
const bucketBytes = bucketSize * accountSize

type walletMap struct {
	bitsForBucket int
	mask          int
	data          []account
	bucketCount   []uint8
	overflows     map[Hash]int
}

func (b *Wallet) Withdraw(hash Hash, value int) bool {
	if balance, ok := b.wallet.overflows[hash]; ok {
		if balance > value {
			b.wallet.overflows[hash] = balance - value
			return true
		}
		return false
	}
	bucket := int(hash[0]) + int(hash[1])<<8 + int(hash[2])<<16 + int(hash[3])<<32
	bucket = bucket & b.wallet.mask
	position := bucket * bucketSize
	for n := 0; n < bucketSize; n++ {
		if hash.Equal(b.wallet.data[position+n].token) {
			if b.wallet.data[bucket+n].balance > value {
				return true
			}
			return false
		}
	}
	return false
}

func (b *Wallet) Credit(hash Hash, value int) bool {
	if balance, ok := b.wallet.overflows[hash]; ok {
		b.wallet.overflows[hash] = balance + value
		return true
	}
	bucket := int(hash[0]) + int(hash[1])<<8 + int(hash[2])<<16 + int(hash[3])<<32
	bucket = bucket & b.wallet.mask
	position := bucket * bucketSize
	for n := 0; n < bucketSize; n++ {
		if hash.Equal(b.wallet.data[position+n].token) {
			b.wallet.data[bucket+n].balance += value
			return true
		}
	}
	if count := b.wallet.bucketCount[bucket]; count < bucketSize {
		b.wallet.data[bucket+int(count)] = account{
			token:   hash,
			balance: value,
		}
		b.wallet.bucketCount[bucket] += 1
	} else {
		b.wallet.overflows[hash] = value
		if len(b.wallet.overflows) > 3*(1<<b.wallet.bitsForBucket) {
			b.wallet = doubleWallet(b.wallet)
		}
	}
	return true
}

func doubleWallet(w *walletMap) *walletMap {
	start := time.Now()
	nw := NewWallet(w.bitsForBucket + 1).wallet
	stat := make([]float64, len(nw.bucketCount))
	for n := 0; n < len(w.data); n++ {
		hash := w.data[n].token
		bucket := int(hash[0]) + int(hash[1])<<8 + int(hash[2])<<16 + int(hash[3])<<32
		bucket = bucket & nw.mask
		stat[bucket] += 1.0
		count := int(nw.bucketCount[bucket])
		if count < bucketSize {
			position := bucket*bucketSize + count
			nw.bucketCount[bucket] += 1
			nw.data[position] = w.data[n]
		} else {
			acc := w.data[n]
			nw.overflows[acc.token] = acc.balance
		}
	}
	for hash, balance := range w.overflows {
		bucket := int(hash[0]) + int(hash[1])<<8 + int(hash[2])<<16 + int(hash[3])<<32
		bucket = bucket & nw.mask
		stat[bucket] += 1.0
		count := int(nw.bucketCount[bucket])
		if count < bucketSize {
			position := bucket*bucketSize + count
			nw.bucketCount[bucket] += 1
			nw.data[position] = account{token: hash, balance: balance}
		} else {
			nw.overflows[hash] = balance
		}
	}
	avg := 0.0
	for _, c := range stat {
		avg += c
	}
	tot := avg
	avg = avg / float64(len(stat))
	count := 0
	count2 := 0
	for _, c := range stat {
		if c > bucketSize {
			count++
		}
		if c > 2*bucketSize {
			count2++
		}

	}
	fmt.Println(len(nw.overflows), tot, avg, float64(count)/float64(len(stat)), float64(count2)/float64(len(stat)), time.Since(start))
	return nw
}

func DumpWallet(w *Wallet) {
	data := make([]account, len(w.wallet.data))
	copy(data, w.wallet.data)
	for hash, balance := range w.wallet.overflows {
		data = append(data, account{token: hash, balance: balance})
	}
}

func main() {

	w := NewWallet(4)
	for n := 1; n < 10e6; n++ {
		hashed := make([]byte, 32)
		rand.Read(hashed)
		//shashed := sha256.Sum256(hashed[:])
		var hash Hash
		for n := 0; n < 32; n++ {
			hash[n] = hashed[n]
		}
		w.Credit(hash, 1)
	}

	start := time.Now()
	fmt.Println(time.Since(start))

}
*/
