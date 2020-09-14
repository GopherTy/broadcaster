# broadcaster
Go  subscribe/publish patterns

# Install 

`go get github.com/gopherty/broadcaster`

# Usage

### Broadcaster

```go
import "github.com/gopherty/broadcaster"

func main(){
    bc := NewBroadcaster(0) // 0 is buffer size
	go bc.Run()
	defer bc.Close()

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
```

### Subject 

```go
import "github.com/gopherty/broadcaster"

type observer struct {
	tag int
}

func (o *observer) Next(message interface{}) {
	fmt.Printf("observer %v recevied message ---> %v\n", o.tag, message)
}

func main(){
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

	s.Publish("hello ")
	s.Publish("world")
	s.Complete()
	s.Publish("你好,世界")
}
```

# License

[MIT](https://github.com/GopherTy/v2ray-web/blob/master/LICENSE) © GopherTy

