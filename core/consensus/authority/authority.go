package authority

import (
	"time"

	"github.com/Aereum/aereum/core/instructions"
)

type ProofOfAuthority struct{}

func (poa *ProofOfAuthority) IsLeader(starting time.Time) bool {
	return true
}

func (poa *ProofOfAuthority) IsValidator(starting time.Time) bool {
	return true
}

func (poa *ProofOfAuthority) IsConsensus(block *instructions.Block) bool {

	return true
}

func (poa *ProofOfAuthority) ValidateBlock(block *instructions.Block) bool {
	return true
}
