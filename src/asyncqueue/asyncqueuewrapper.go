package asyncqueue

import (
	"sync"

	queuing "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
)

type AsyncQueueWrapper interface {
	Offer(message queuing.QueueMessage)
	Poll() (queuing.QueueMessage, bool)
}

type AsyncQueue struct {
	queue        queuing.Queue
	channel      chan queuing.QueueMessage
	closeChannel chan int
	muLock       *sync.Mutex
}

func NewAsyncQueue(queue queuing.Queue) *AsyncQueue {
	queueWrapper := AsyncQueue{
		queue:        queue,
		channel:      make(chan queuing.QueueMessage),
		closeChannel: make(chan int),
		muLock:       new(sync.Mutex),
	}

	go queueWrapper.startRunningAsync()
	return &queueWrapper
}

func (queueWrapper *AsyncQueue) startRunningAsync() {
	isChannelClosed := false
	for !isChannelClosed {
		select {
		case message := <-queueWrapper.channel:
			queueWrapper.muLock.Lock()
			queueWrapper.queue.Offer(message)
			queueWrapper.muLock.Unlock()
		case <-queueWrapper.closeChannel:
			isChannelClosed = true
			close(queueWrapper.channel)
		default:
		}
	}
}

func (queueWrapper *AsyncQueue) Offer(message queuing.QueueMessage) {
	queueWrapper.channel <- message
}

func (queueWrapper *AsyncQueue) Poll() (queuing.QueueMessage, bool) {
	queueWrapper.muLock.Lock()
	defer queueWrapper.muLock.Unlock()

	return queueWrapper.queue.Poll()
}
