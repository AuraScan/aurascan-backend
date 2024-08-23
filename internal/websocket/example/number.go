package example

import (
	"aurascan-backend/internal/websocket/pubsub"
	"sync/atomic"
	"time"
)

// Number 每秒+1 返回最新int
type Number struct {
	n int64
}

func (*Number) Topic() pubsub.Topic {
	return "number"
}
func (c *Number) Cache() interface{} {
	return &pubsub.Message{Topic: c.Topic(), Data: atomic.LoadInt64(&c.n)}
}

func (c *Number) Publish(pub chan *pubsub.Message, stop <-chan struct{}) {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			atomic.AddInt64(&c.n, 1)
			pub <- &pubsub.Message{Topic: c.Topic(), Data: atomic.LoadInt64(&c.n)}
			timer.Reset(time.Second)
		case <-stop:
			return
		}
	}
}
