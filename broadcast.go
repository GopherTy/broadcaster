package broadcaster

import (
	"errors"
	"sync/atomic"
)

// Broadcaster pub/sub patterns
type Broadcaster interface {
	Subscribe(chan interface{}) error
	Unsubscribe(chan interface{}) error
	Publish(v interface{}) error
	Run() error
	Close() error
}

type broadcaster struct {
	input       chan interface{}
	subscribe   chan chan interface{}
	unsubscribe chan chan interface{}
	observers   map[chan interface{}]bool
	done        chan struct{} // a signal to close broadcaster
	l           uint32
}

func (b *broadcaster) Run() error {
	if b == nil {
		return errors.New("broadcaster cannot be  nil")
	}
	go func() {
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
				return
			}
		}
	}()
	return nil
}

func (b *broadcaster) Close() error {
	if b == nil {
		return errors.New("broadcaster cannot be nil")
	}
	select {
	case <-b.done:
		return errors.New("broadcaster already closed")
	default:
		if atomic.CompareAndSwapUint32(&b.l, 0, 1) {
			close(b.done)
		}
	}
	return nil
}

func (b *broadcaster) Publish(message interface{}) error {
	if b == nil {
		return errors.New("broadcaster cannot be nil")
	}
	select {
	case <-b.done:
		return errors.New("broadcaster already closed")
	case b.input <- message: // broadcaster receive message
	}
	return nil
}

func (b *broadcaster) Subscribe(observer chan interface{}) error {
	if b == nil {
		return errors.New("broadcaster cannot be nil")
	}
	select {
	case <-b.done:
		return errors.New("broadcaster already closed")
	case b.subscribe <- observer:
	}
	return nil
}

func (b *broadcaster) Unsubscribe(observer chan interface{}) error {
	if b == nil {
		return errors.New("broadcaster cannot be nil")
	}
	select {
	case <-b.done:
		return errors.New("broadcaster already closed")
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
