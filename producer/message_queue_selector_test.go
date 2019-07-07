package producer

import "github.com/mymmsc/go-rocketmq-client/v2/message"

type fakeMessageQueueSelector struct {
	selectRet *message.Queue
}

func (s *fakeMessageQueueSelector) Select(mqs []*message.Queue, m *message.Message, arg interface{}) *message.Queue {
	return s.selectRet
}
