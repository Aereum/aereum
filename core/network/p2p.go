package network

/*type MessageQueueRequest struct {
	message  message.Message
	response chan bool
}

type MessageUnqueueRequest struct {
	response chan message.Message
}

type ValidMessageQueue struct {
	queue   chan MessageQueueRequest
	unqueue chan MessageUnqueueRequest
}

func NewValidMessageQueue() *ValidMessageQueue {
	queue := make(chan MessageQueueRequest)
	unqueue := make(chan MessageUnqueueRequest)
	//messages := make([]message.Message, 0)
	//hashes := make(map[hashdb.Hash]struct{})
	go func() {
		for {
			select {
			case <-queue:
				//
			case <-unqueue:
				//
			}
		}
	}()
	return &ValidMessageQueue{
		queue:   queue,
		unqueue: unqueue,
	}
}

func NewMessageServer() {
	//messages := NewValidMessageQueue()
	service := fmt.Sprintf(":%v", messageReceiveConnectionPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	logger.MustOrPanic(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	logger.MustOrPanic(err)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

		}
	}()
}

func NewValidatorServer(nodeAddress string) {
	service := fmt.Sprintf(":%v", validationNodePort)
	tcpAdrr, err := net.ResolveTCPAddr("tcp", service)
	logger.MustOrPanic(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	logger.MustOrPanic(err)
	go func() {

	}

}

type MessageReceiveConnection struct {
	conn     *net.TCPConn
	msgqueue *ValidMessageQueue
}
*/
