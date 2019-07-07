package producer

import "github.com/mymmsc/go-rocketmq-client/v2/message"

// MessageQueueSelector select the message queue
type MessageQueueSelector interface {
	Select(mqs []*message.Queue, m *message.Message, arg interface{}) *message.Queue
}
