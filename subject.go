package broadcaster

import (
	"fmt"
	"sync"
)

// Subject pub/sub patterns
type Subject struct {
	observers map[*Subscription]Observable
	mu        sync.Mutex
}

// Publish publish message to observers
func (s *Subject) Publish(msg interface{}) {
	s.mu.Lock()
	for _, observer := range s.observers {
		observer.Next(msg)
	}
	s.mu.Unlock()
}

// Subscribe add observer to subject
func (s *Subject) Subscribe(obs Observable) *Subscription {
	subs := &Subscription{subj: s}

	s.mu.Lock()
	s.observers[subs] = obs
	s.mu.Unlock()

	return subs
}

func (s *Subject) unsubscribe(subs *Subscription, lock bool) {
	if lock {
		s.mu.Lock()
		delete(s.observers, subs)
		s.mu.Unlock()
	} else {
		delete(s.observers, subs)
	}
}

// Subscription observer unsubscribe object
type Subscription struct {
	subj *Subject
}

// Unsubscribe unsubscribe.
// when you unsubscribe in observer's Next function, lock is false. otherwise is true
func (s *Subscription) Unsubscribe(lock bool) {
	s.subj.unsubscribe(s, lock)
}

// Observable observer to implements
type Observable interface {
	// how to recevie message
	Next(interface{})
	// // error
	// Error(error)
	// // complete
	// Complete()
}

type observer struct {
	tag int
}

func (o *observer) Next(message interface{}) {
	fmt.Printf("observer %v recevied message ---> %v\n", o.tag, message)

}

// NewObserver a example observer implment Observable
func NewObserver() Observable {
	return &observer{}
}
