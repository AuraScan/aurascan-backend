package pubsub

import (
	"ch-common-package/logger"
	"github.com/google/uuid"
	"sync"
)

type Topic string
type Publisher interface {
	Publish(chan *Message, <-chan struct{})
	Topic() Topic
	Cache() interface{}
}

type Message struct {
	Topic Topic       `json:"topic"`
	Data  interface{} `json:"data"`
}

var ps pubSub

func Start() {
	ps.pubChan = make(chan *Message)
	ps.unpub = make(chan Publisher)
	ps.publisher = make(chan Publisher)
	ps.pubs = make(map[Topic]Publisher)
	ps.pubStop = make(map[Topic]chan struct{})
	ps.subs = make(map[*Sub]struct{})
	ps.sub = make(chan *Sub, 1)
	ps.unsub = make(chan *Sub, 1)
	ps.done = make(chan struct{})

	go ps.Start()
	logger.Info("starting PubSub!")
}

func RegisterPublisher(p Publisher) {
	ps.RegisterPublisher(p)
}

func Unpub(p Publisher) {
	ps.Unpub(p)
}

func Subscribe(topic ...Topic) *Sub {
	return ps.Sub(topic...)
}

func Unsub(sub *Sub) {
	ps.unsub <- sub
}

func (ps *pubSub) RegisterPublisher(p Publisher) {
	ps.publisher <- p
}

func (ps *pubSub) Unpub(p Publisher) {
	ps.unpub <- p
}

type pubSub struct {
	pubChan   chan *Message
	unpub     chan Publisher
	publisher chan Publisher
	pubs      map[Topic]Publisher
	pubStop   map[Topic]chan struct{}
	subs      map[*Sub]struct{}
	sub       chan *Sub
	unsub     chan *Sub
	done      chan struct{}
	// timeout   time.Duration
}

type Sub struct {
	C      chan interface{}
	uid    uuid.UUID
	topics map[Topic]struct{}
	mu     sync.Mutex
}

func (s *Sub) HasTopic(t Topic) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, has := s.topics[t]
	return has
}

func (s *Sub) AllTopics() []Topic {
	s.mu.Lock()
	defer s.mu.Unlock()
	var tps []Topic
	for t := range s.topics {
		tps = append(tps, t)
	}
	return tps
}

func (s *Sub) AddTopic(t Topic) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics[t] = struct{}{}
}

func (s *Sub) DelTopic(t Topic) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.topics, t)
}

func (s *Sub) UUID() uuid.UUID {
	return s.uid
}

func (ps *pubSub) Sub(topic ...Topic) *Sub {
	uid, _ := uuid.NewUUID()
	sub := &Sub{
		C:      make(chan interface{}, 1),
		uid:    uid,
		topics: make(map[Topic]struct{}),
		mu:     sync.Mutex{},
	}

	for _, t := range topic {
		sub.topics[t] = struct{}{}
	}

	ps.sub <- sub
	return sub
}

func (ps *pubSub) Start() {
	for {
		select {
		case msg := <-ps.pubChan:
			// 发布 非阻塞
			for sub := range ps.subs {
				if sub.HasTopic(msg.Topic) {
					select {
					case sub.C <- msg:
					default:
						logger.Errorf("cannot send message, out of buffer uid=%s, topics=%+v", sub.uid.String(), sub.topics)
					}
				}
			}
			logger.Debugf("received message from publisher | type=%s", msg.Topic)

		case sub := <-ps.sub:
			// 添加新的订阅并返回最新的订阅内容
			// 如果未指定订阅 默认订阅所有主题
			ps.subs[sub] = struct{}{}

			if len(sub.topics) == 0 {
				for _, pub := range ps.pubs {
					sub.topics[pub.Topic()] = struct{}{}
				}
			}

			for topic, _ := range sub.topics {
				if pub, has := ps.pubs[topic]; has {
					if cache := pub.Cache(); cache != nil {
						select {
						case sub.C <- cache:
						default:
							logger.Warnf("cannot send message, out of buffer uid=%s, topics=%+v", sub.uid.String(), sub.topics)
						}
					}
				}
			}
			logger.Infof("pubSub.Start added sub chan | uid is %s, topics is%+v", sub.uid.String(), sub.topics)

		case sub := <-ps.unsub:
			// 取消订阅
			close(sub.C)
			delete(ps.subs, sub)
			logger.Debugf("closed sub chan | uid=%s", sub.uid.String())

		case pub := <-ps.publisher:
			// 添加新的发布者 并运行发布者函数
			_, has := ps.pubs[pub.Topic()]
			if !has {
				stop := make(chan struct{})
				ps.pubs[pub.Topic()] = pub
				ps.pubStop[pub.Topic()] = stop
				go pub.Publish(ps.pubChan, stop)
				logger.Infof("added publish func | type=%v", pub.Topic())
			} else {
				logger.Warnf("already published func | type=%v", pub.Topic())
			}

		case pub := <-ps.unpub:
			// 停止发布者 并移除
			_, ok := ps.pubs[pub.Topic()]
			if ok {
				close(ps.pubStop[pub.Topic()])
				delete(ps.pubs, pub.Topic())
				delete(ps.pubStop, pub.Topic())
				logger.Infof("closed publish chan | type=%v", pub.Topic())
			}

		case <-ps.done:
			// 关闭所有订阅者和发布者
			for sub := range ps.subs {
				close(sub.C)
			}
			for _, pub := range ps.pubStop {
				close(pub)
			}
			return
		}
	}
}
