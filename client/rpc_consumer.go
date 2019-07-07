package client

import (
	"time"

	"github.com/mymmsc/go-rocketmq-client/v1/client/rpc"
)

// GetConsumerIDs get the client id from the broker wraper
func (c *MQClient) GetConsumerIDs(addr, group string, to time.Duration) (ids []string, err error) {
	return rpc.GetConsumerIDs(c.Client, addr, group, to)
}
