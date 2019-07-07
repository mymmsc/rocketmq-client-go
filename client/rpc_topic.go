package client

import (
	"time"

	"github.com/mymmsc/go-rocketmq-client/v2/client/rpc"
)

// CreateOrUpdateTopic create topic from broker
func (c *MQClient) CreateOrUpdateTopic(
	addr string, header *rpc.CreateOrUpdateTopicHeader, to time.Duration,
) error {
	return rpc.CreateOrUpdateTopic(c.Client, addr, header, to)
}

// DeleteTopicInBroker delete topic in the broker
func (c *MQClient) DeleteTopicInBroker(addr, topic string, to time.Duration) error {
	return rpc.DeleteTopicInBroker(c.Client, addr, topic, to)
}
