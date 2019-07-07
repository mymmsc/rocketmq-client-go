package producer

import (
	"testing"
	"time"

	"github.com/mymmsc/go-rocketmq-client/v2/message"

	"github.com/stretchr/testify/assert"
)

type fakeTopicRouter struct {
	writeCount     int
	nextQueueIndex uint32
}

func (tr *fakeTopicRouter) SelectOneQueue() *message.Queue {
	return &message.Queue{
		BrokerName: "SEL",
		Topic:      "SEL",
	}
}
func (tr *fakeTopicRouter) NextQueueIndex() uint32 {
	tr.nextQueueIndex++
	return tr.nextQueueIndex
}
func (tr *fakeTopicRouter) MessageQueues() []*message.Queue {
	return []*message.Queue{
		&message.Queue{
			BrokerName: "b1",
			Topic:      "b1",
			QueueID:    0,
		},
		&message.Queue{
			BrokerName: "b1",
			Topic:      "b1",
			QueueID:    1,
		},
		&message.Queue{
			BrokerName: "b2",
			Topic:      "b2",
			QueueID:    1,
		},
	}
}
func (tr *fakeTopicRouter) WriteQueueCount(broker string) int {
	tr.writeCount++
	if tr.writeCount == 1 {
		return 0
	}
	return tr.writeCount
}
func (tr *fakeTopicRouter) SelectOneQueueNotOf(lastBroker string) *message.Queue {
	return &message.Queue{
		BrokerName: "HINT",
		Topic:      "HINT",
	}
}

func TestLatency(t *testing.T) {
	fs := NewMQFaultStrategy(true)

	fs.UpdateFault("b1", 1*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.UpdateFault("b1", 99*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.UpdateFault("b1", 49*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.UpdateFault("b1", 50*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.UpdateFault("b1", 100*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.sendLatencyFaultEnable = false
	fs.UpdateFault("b1", 1200*time.Millisecond, false)
	assert.True(t, fs.Available("b1"))

	fs.sendLatencyFaultEnable = true
	fs.UpdateFault("b1", 600*time.Millisecond, false)
	assert.False(t, fs.Available("b1"))
}

func TestSelectOneQueue(t *testing.T) {
	fs, tp := NewMQFaultStrategy(false), &fakeTopicRouter{}

	q := fs.SelectOneQueue(tp, "b1")
	assert.Equal(t, "HINT", q.BrokerName)
	assert.Equal(t, "HINT", q.Topic)

	fs.sendLatencyFaultEnable = true
	q = fs.SelectOneQueue(tp, "b1")
	assert.Equal(t, "b2", q.BrokerName)
	assert.Equal(t, "b2", q.Topic)
	assert.Equal(t, uint8(1), q.QueueID)
	fs.UpdateFault("b1", 0, true)
	fs.UpdateFault("b2", 0, true)

	q = fs.SelectOneQueue(tp, "not exist")
	assert.Equal(t, "HINT", q.BrokerName)
	assert.Equal(t, "HINT", q.Topic)
	assert.Equal(t, 1, tp.writeCount)

	q = fs.SelectOneQueue(tp, "not exist")
	assert.True(t, q.BrokerName == "b1" || q.BrokerName == "b2")
	assert.Equal(t, "SEL", q.Topic)
	assert.Equal(t, 2, tp.writeCount)

	fs.UpdateFault("not exist", 0, false)
	q = fs.SelectOneQueue(tp, "not exist")
	assert.Equal(t, "not exist", q.BrokerName)
	assert.Equal(t, "SEL", q.Topic)
	assert.Equal(t, uint8(0), q.QueueID)
	assert.Equal(t, 3, tp.writeCount)
}
