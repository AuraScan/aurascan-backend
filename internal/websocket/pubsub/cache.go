package pubsub

import "sync"

// Cache LIFO
type Cache struct {
	data  []*Message
	mu    sync.Mutex
	limit int
}

const defaultCacheLength = 10

func New(limit int) *Cache {
	if limit < 1 {
		limit = defaultCacheLength
	}
	return &Cache{
		data:  make([]*Message, 0),
		mu:    sync.Mutex{},
		limit: limit,
	}
}

// Push 写入元素 超出限制就移除最早的元素
func (c *Cache) Push(msg *Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.data) >= c.limit {
		c.data = c.data[1:]
	}
	c.data = append(c.data, msg)
}

// Pop 推出最新元素 并删除
func (c *Cache) Pop() *Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.data) == 0 {
		return nil
	}
	result := c.data[len(c.data)-1]
	c.data = c.data[:len(c.data)-1]
	return result
}

// Peek 以LIFO模式查看length个最新元素
func (c *Cache) Peek(length int) []*Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) == 0 {
		return nil
	}

	if length < 1 || length > c.limit {
		length = c.limit
	}
	if length > len(c.data) {
		length = len(c.data)
	}

	var result []*Message
	for i := 1; i <= length; i++ {
		result = append(result, c.data[len(c.data)-i])
	}
	return result
}
