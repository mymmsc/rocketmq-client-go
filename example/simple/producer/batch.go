package main

import (
	"fmt"
	"time"

	"github.com/mymmsc/go-rocketmq-client/v2/message"
	"github.com/mymmsc/go-rocketmq-client/v2/producer"
)

func sendBatch(p *producer.Producer) (*producer.SendResult, error) {
	data := message.Data{Body: []byte(time.Now().String())}
	m := &message.Batch{Topic: topic, Datas: []message.Data{data}}

	r, err := p.SendBatchSync(m)

	if err == nil {
		fmt.Printf("sub message id:%s\n", message.GetUniqID(m.Datas[0].Properties))
	}

	return r, err
}
