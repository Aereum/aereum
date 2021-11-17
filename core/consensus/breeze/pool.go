package breeze

import (
	"sync"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type InstructionPool struct {
	queue  []*instruction.Instruction
	hashes map[crypto.Hash]int
	mu     *sync.Mutex
}

func (pool *InstructionPool) Unqueue() *instruction.Instruction {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if len(pool.queue) > 0 {
		first := pool.queue[0]
		pool.queue = pool.queue[1:]
		return first
	}
	return nil
}

func (pool *InstructionPool) Queue(m *instruction.Instruction, hash crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.queue = append(pool.queue, m)
	pool.hashes[hash] = len(pool.queue) - 1
}

func (pool *InstructionPool) Delete(hash crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	position, ok := pool.hashes[hash]
	if !ok {
		return
	}
	delete(pool.hashes, hash)
	pool.queue = append(pool.queue[0:position-1], pool.queue[position+1:]...)
}

func (pool *InstructionPool) DeleteArray(hashes []crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for _, hash := range hashes {
		position, ok := pool.hashes[hash]
		if !ok {
			return
		}
		delete(pool.hashes, hash)
		pool.queue = append(pool.queue[0:position-1], pool.queue[position+1:]...)
	}
}
