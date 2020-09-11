package broadcaster

import (
	"fmt"
	"sync"
	"testing"
)

func TestBroadcast(t *testing.T) {
	var wg sync.WaitGroup

	bc := NewBroadcaster(0)
	bc.Run()
	defer bc.Close()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		observer := make(chan interface{})
		bc.Subscribe(observer)
		go func(i int) {
			defer wg.Done()
			defer bc.Unsubscribe(observer)
			for j := 0; j < 2; j++ {
				if msg, ok := <-observer; ok {
					fmt.Printf("observer: %v recevied  %v message------>%v \n", i+1, j+1, msg)
				}
			}
		}(i)
	}

	bc.Publish("Hello ")
	bc.Publish("World")
	wg.Wait()
}

func TestBroadcastClosePublishSubUnsub(t *testing.T) {
	bc := NewBroadcaster(0)
	bc.Run()

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
	bc.Run()
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
	bc.Run()
	defer bc.Close()

	bc.Subscribe(observer)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bc.Publish(1)
			<-observer
		}
	})
}
