package example

import (
	"aurascan-backend/internal/websocket/pubsub"
	"ch-common-package/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url         = "http://172.30.56.103:3030/testnet/block/latest"
	cacheLength = 5
)

var blockCache = pubsub.New(cacheLength)

type Block struct {
	Header struct {
		Metadata struct {
			Height         uint64 `json:"height"`
			CoinbaseTarget uint64 `json:"coinbase_target"`
			ProofTarget    uint64 `json:"proof_target"`
			Timestamp      uint64 `json:"timestamp"`
		} `json:"metadata"`
	} `json:"header"`
}

func (*Block) Topic() pubsub.Topic {
	return "latest-block"
}
func (*Block) Cache() interface{} {
	x := blockCache.Peek(1)
	if x == nil {
		return nil
	}
	return x[0]
}

func (b *Block) Publish(pub chan *pubsub.Message, stop <-chan struct{}) {
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			body, err := HttpGet(url)
			if err == nil {
				var block Block
				err = json.Unmarshal([]byte(body), &block)
				if err == nil {
					c := blockCache.Peek(1)
					if c == nil || block.Header.Metadata.Height > c[0].Data.(*Block).Header.Metadata.Height {
						msg := &pubsub.Message{Topic: b.Topic(), Data: &block}
						pub <- msg
						blockCache.Push(msg)
						logger.Infof("update new block | height=%v", block.Header.Metadata.Height)
					}
				} else {
					logger.Errorf("failed to parse block | err=%v", err)
				}
			} else {
				logger.Errorf("failed to get latest block | err=%v", err)
			}
			timer.Reset(time.Second * 3)

		case <-stop:
			return
		}
	}
}

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
