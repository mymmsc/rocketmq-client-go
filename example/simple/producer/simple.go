package main

import (
	"time"

	"github.com/mymmsc/go-rocketmq-client/v2/message"
	"github.com/mymmsc/go-rocketmq-client/v2/producer"
)

func sendSimple(p *producer.Producer) (*producer.SendResult, error) {
	now := time.Now()
	m := &message.Message{Topic: topic, Body: []byte(now.String())}
	m.SetTags("simple")

	return p.SendSync(m)
}
