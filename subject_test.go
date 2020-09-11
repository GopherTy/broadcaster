package broadcaster

import (
	"sync"
	"testing"
)

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
		s.Publish("Hello World")
	}
}

func BenchmarkParallelSubject(b *testing.B) {
	s := &Subject{
		observers: make(map[*Subscription]Observable),
	}

	const count = 10000
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			observer := &observer{tag: i}
			subs := s.Subscribe(observer)
			for i := 0; i < count; i++ {
				go func(s *Subscription) {
					s.Unsubscribe(true)
					wg.Done()
				}(subs)
			}
		}(i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Publish("Hello World")
		}
	})
	wg.Wait()
}
