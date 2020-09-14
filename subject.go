package broadcaster

import (
	"sync"
)

// Subject pub/sub patterns
type Subject struct {
	observers map[*Subscription]Observable
	complete  bool
	mu        sync.Mutex
}

// Publish publish message to observers
func (s *Subject) Publish(msg interface{}) {
	if s.complete {
		return
	}
	s.mu.Lock()
	for _, observer := range s.observers {
		observer.Next(msg)
	}
	s.mu.Unlock()
}

// Complete complete publish message
func (s *Subject) Complete() {
	if s.complete {
		return
	}
	s.mu.Lock()
	s.complete = true
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
	// Next recevie message
	Next(interface{})
}

// NewSubject return subject
func NewSubject() *Subject {
	return &Subject{observers: make(map[*Subscription]Observable)}
}
