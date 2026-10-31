package internal

import "sync"

// DataQueue is a thread-safe queue for storing messages.
type DataQueue struct {
	mux      sync.Mutex
	messages []any
	Closed   bool
}

func NewDataQueue() *DataQueue {
	return &DataQueue{
		messages: make([]any, 0),
	}
}

// Close marks the queue as closed, preventing further enqueues.
func (q *DataQueue) Close() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.Closed = true
}
func (q *DataQueue) Enqueue(message any) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.Closed {
		return
	}
	q.messages = append(q.messages, message)
}

// Dequeue removes and returns all messages from the queue for the given key.
// Returns the messages and true if successful, or nil and false if the queue is empty or the key does not exist.
func (q *DataQueue) Dequeue() ([]any, bool) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(q.messages) == 0 {
		return nil, false
	}
	messages := q.messages
	q.messages = make([]any, 0)
	return messages, true
}
