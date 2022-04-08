package actors

import (
	"github.com/Aereum/aereum/core/crypto"
)

type Author struct {
	Caption    string
	PrivateKey crypto.PrivateKey
	Attorney   crypto.PrivateKey
}
