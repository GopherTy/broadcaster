package broadcaster

import (
	"fmt"
	"testing"
)

// usage
type observer struct {
	tag int
}

func (o *observer) Next(message interface{}) {
	fmt.Printf("observer %v recevied message ---> %v\n", o.tag, message)
}

func TestSubject(t *testing.T) {
	s := &Subject{
		observers: make(map[*Subscription]Observable),
	}

	const count = 4
	for i := 1; i < count; i++ {
		observer := &observer{
			tag: i,
		}
		subscription := s.Subscribe(observer)
		if i == 2 {
			subscription.Unsubscribe(true)
		}
	}

	s.Publish(1)
	s.Publish(2)
	s.Complete()
	s.Publish(3)
}

func TestHandleFuncSubject(t *testing.T) {
	s := NewSubject()

	const count = 4
	for i := 1; i < count; i++ {
		var subscription *Subscription
		func(i int) {
			subscription = s.HandleFunc(func(msg interface{}) {
				fmt.Printf("observer %v recieve message ---> %v \n", i, msg)
			})
		}(i)

		if i == 2 {
			subscription.Unsubscribe(true)
		}
	}

	s.Publish("hello")
	s.Publish("world")
	s.Complete()
	s.Publish(3)
}

func BenchmarkSubject(b *testing.B) {
	s := &Subject{
		observers: make(map[*Subscription]Observable),
	}
	const count = 10
	for i := 0; i < count; i++ {
		observerT := &observer{tag: i}
		s.Subscribe(observerT)
	}

	for i := 0; i < b.N; i++ {
		s.Publish(i)
	}
}

func BenchmarkParallelSubject(b *testing.B) {
	s := &Subject{
		observers: make(map[*Subscription]Observable),
	}

	const count = 2
	for i := 0; i < count; i++ {
		observer := &observer{tag: i}
		s.Subscribe(observer)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Publish("Hello World")
		}
	})
}
