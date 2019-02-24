package remote

import (
	"os"
	"testing"
	"time"

	"github.com/qiniu/log.v1"
	"github.com/zjykzk/rocketmq-client-go/remote/net"
)

func TestBroker(t *testing.T) {
	t.Run("create topic", testCreateTopic)
	t.Run("delete topic", testDeleteTopic)
}

func testCreateTopic(t *testing.T) {
	rpc := NewRPC(NewClient(&net.Config{}, nil, log.New(os.Stderr, "", log.Ldefault)))
	err := rpc.CreateOrUpdateTopic("localhost:10909", &CreateOrUpdateTopicHeader{
		Topic:           "test_create_topic",
		DefaultTopic:    "default_topic",
		ReadQueueNums:   4,
		WriteQueueNums:  4,
		Perm:            0,
		TopicFilterType: "SINGLE_TAG",
		TopicSysFlag:    0,
		Order:           false,
	}, time.Millisecond*100)
	t.Log(err)
}

func testDeleteTopic(t *testing.T) {
	rpc := NewRPC(NewClient(&net.Config{}, nil, log.New(os.Stderr, "", log.Ldefault)))
	err := rpc.DeleteTopicInBroker("localhost:10909", "test_create_topic", time.Millisecond*100)
	t.Log(err)
}
