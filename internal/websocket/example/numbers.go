package example

import (
	"aurascan-backend/internal/websocket/pubsub"
	"time"
)

// Numbers 每秒+1 使用cache 返回最新int组成的切片
type Numbers []int

const numberLength = 5

var cache = pubsub.New(numberLength)

func (*Numbers) Topic() pubsub.Topic {
	return "numbers"
}
func (n *Numbers) Cache() interface{} {
	return &pubsub.Message{Topic: n.Topic(), Data: cache.Peek(numberLength)}
}

func (n *Numbers) Publish(pub chan *pubsub.Message, stop <-chan struct{}) {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			var number int = 0
			c := cache.Peek(1)
			if c != nil {
				number = c[0].Data.(int) + 1
			}
			cache.Push(&pubsub.Message{Topic: n.Topic(), Data: number})
			pub <- &pubsub.Message{Topic: n.Topic(), Data: cache.Peek(numberLength)}
			timer.Reset(time.Second)
		case <-stop:
			return
		}
	}
}
