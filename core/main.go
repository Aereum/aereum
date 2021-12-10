package main

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto/ed25519/curve25519"
)

func main() {

	a := make([]byte, 32)
	a_pub, _ := curve25519.X25519(a, curve25519.Basepoint)

	// manda para

	b := make([]byte, 32)
	b_pub, _ := curve25519.X25519(b, curve25519.Basepoint) // -> mando pro a
	shared2, _ := curve25519.X25519(b, a_pub)              // esta é a cifra.... ninguém conhece b

	// a recebe bpub

	shared1, _ := curve25519.X25519(a, b_pub) // calcula

	fmt.Println(shared1)
	fmt.Println(shared2)
}
