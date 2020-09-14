package broadcaster

import (
	"errors"
	"sync"
)

// Broadcaster pub/sub patterns
type Broadcaster interface {
	Subscribe(chan interface{}) error
	Unsubscribe(chan interface{}) error
	Publish(v interface{}) error
	Run() error
	Close() error
	Done() chan struct{}
}

// ErrCanceled .
var ErrCanceled = errors.New("broadcaster already closed")

type broadcaster struct {
	input       chan interface{}
	subscribe   chan chan interface{}
	unsubscribe chan chan interface{}
	observers   map[chan interface{}]bool
	done        chan struct{} // a signal to close broadcaster
	mu          sync.Mutex
}

func (b *broadcaster) Run() error {
	for {
		select {
		case msg := <-b.input:
			// publish message to observers
			for observer := range b.observers {
				observer <- msg
			}
		case observer := <-b.subscribe:
			b.observers[observer] = true
		case observer := <-b.unsubscribe:
			delete(b.observers, observer)
		case <-b.done:
			return nil
		}
	}
}

func (b *broadcaster) Close() error {
	select {
	case <-b.done:
		return ErrCanceled
	default:
		b.mu.Lock()
		close(b.done)
		b.mu.Unlock()
	}
	return nil
}

func (b *broadcaster) Done() chan struct{} {
	b.mu.Lock()
	if b.done == nil {
		b.done = make(chan struct{})
	}
	d := b.done
	b.mu.Unlock()
	return d
}

func (b *broadcaster) Publish(message interface{}) error {
	select {
	case <-b.done:
		return ErrCanceled
	case b.input <- message: // broadcaster receive message
	}
	return nil
}

func (b *broadcaster) Subscribe(observer chan interface{}) error {
	select {
	case <-b.done:
		return ErrCanceled
	case b.subscribe <- observer:
	}
	return nil
}

func (b *broadcaster) Unsubscribe(observer chan interface{}) error {
	select {
	case <-b.done:
		return ErrCanceled
	case b.unsubscribe <- observer:
	}
	return nil
}

// NewBroadcaster create a broadcaster, buf is message buffer.
func NewBroadcaster(buf int) Broadcaster {
	b := &broadcaster{
		input:       make(chan interface{}, buf),
		subscribe:   make(chan chan interface{}),
		unsubscribe: make(chan chan interface{}),
		observers:   make(map[chan interface{}]bool),
		done:        make(chan struct{}),
	}
	return b
}
