package asyncqueue

import (
	"sync"

	queuing "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
)

type AsyncQueueWrapper interface {
	Offer(queuing.QueueMessage)
	Poll() (queuing.QueueMessage, bool)
	BatchOffer([]queuing.QueueMessage)
	PollBatch(int) []queuing.QueueMessage
}

type AsyncQueue struct {
	queue   queuing.Queue
	channel chan queuing.QueueMessage
	muLock  *sync.Mutex // TODO investigate this mutex lock under heavy load. Looking for how quickly/efficiently reads are handled
}

func NewAsyncQueue(queue queuing.Queue) *AsyncQueue {
	queueWrapper := AsyncQueue{
		queue:   queue,
		channel: make(chan queuing.QueueMessage),
		muLock:  new(sync.Mutex),
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

func (queueWrapper *AsyncQueue) BatchOffer(messages []queuing.QueueMessage) {
	// TODO investigate this more for potential bottleneck as well as non durable store while message is waiting to be processed in the channel
	for _, message := range messages {
		queueWrapper.channel <- message
	}
}

func (queueWrapper *AsyncQueue) PollBatch(batchSize int) []queuing.QueueMessage {
	queueWrapper.muLock.Lock()
	defer queueWrapper.muLock.Unlock()

	batch := make([]queuing.QueueMessage, batchSize)
	for i := 0; i < batchSize; i++ {
		message, messageFoundInQueue := queueWrapper.queue.Poll()

		if !messageFoundInQueue {
			break
		}

		batch = append(batch, message)
	}

	return batch
}
