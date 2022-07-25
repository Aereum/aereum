package main

import (
	"net"
	"sync"
	"time"

	"github.com/Aereum/aereum/core/instructions"
	"github.com/gobwas/ws/wsutil"
)

type InstructionQueue struct {
	send  chan instructions.Instruction
	close chan struct{}
}

func (s *InstructionQueue) Send(instruction instructions.Instruction) {
	s.send <- instruction
}

func (s *InstructionQueue) Close() {
	s.close <- struct{}{}
}

type SocketQueue struct {
	send  chan jsoner
	close chan struct{}
}

func (s *SocketQueue) Send(msg jsoner) {
	s.send <- msg
}

func (s *SocketQueue) Close() {
	s.close <- struct{}{}
}

func NewSocketQueue(conn net.Conn) *SocketQueue {
	queueData := make([]jsoner, 0)
	queue := make(chan jsoner)
	terminate := make(chan struct{})
	terminated := false
	lock := sync.Mutex{}
	// queuer goroutine
	go func() {
		for {
			select {
			case data := <-queue:
				if data != nil {
					lock.Lock()
					queueData = append(queueData, data)
					lock.Unlock()
				}
			case <-terminate:
				close(queue)
				terminated = true
				return
			}
		}
	}()
	// worker goroutine
	go func() {
		hundredmili, _ := time.ParseDuration("100ms")
		next := time.After(hundredmili)
		for {
			<-next
			for {
				lock.Lock()
				var data jsoner
				if len(queueData) != 0 {
					data = queueData[0]
					queueData = queueData[1:]
				}
				lock.Unlock()
				if data != nil {
					bytes := data.JSON()
					if bytes != nil {
						wsutil.WriteClientText(conn, bytes)
					}
				} else {
					break
				}
			}
			if terminated {
				return
			}
			next = time.After(hundredmili)
		}
	}()
	return &SocketQueue{send: queue, close: terminate}
}
