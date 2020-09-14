package broadcaster

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	var wg sync.WaitGroup

	bc := NewBroadcaster(0)
	go bc.Run()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		observer := make(chan interface{})
		bc.Subscribe(observer)
		go func(i int) {
			defer wg.Done()
			defer bc.Unsubscribe(observer)
			for {
				select {
				case msg := <-observer:
					fmt.Printf("observer: %v recevied message------>%v \n", i+1, msg)
				case <-bc.Done():
					return
				}
			}
		}(i)
	}

	bc.Publish("Hello ")
	bc.Publish("World")
	// maybe done signal arrived before message
	time.Sleep(100 * time.Millisecond)
	bc.Close()
	wg.Wait()
}

func TestBroadcastClosePublishSubUnsub(t *testing.T) {
	bc := NewBroadcaster(0)
	go bc.Run()
	context.Background()

	err := bc.Close()
	handError(err)
	err = bc.Close()
	handError(err)
	err = bc.Publish(1)
	handError(err)
	observer := make(chan interface{})
	err = bc.Unsubscribe(observer)
	handError(err)
	err = bc.Subscribe(observer)
	handError(err)
}
func handError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func BenchmarkBroadcast(b *testing.B) {
	bc := NewBroadcaster(1000)
	go bc.Run()
	defer bc.Close()

	observer := make(chan interface{})
	bc.Subscribe(observer)
	for i := 0; i < b.N; i++ {
		bc.Publish(1)
		<-observer
	}
}

func BenchmarkParallelBrodcast(b *testing.B) {
	observer := make(chan interface{})

	bc := NewBroadcaster(0)
	go bc.Run()
	defer bc.Close()

	bc.Subscribe(observer)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bc.Publish(1)
			<-observer
		}
	})
}
