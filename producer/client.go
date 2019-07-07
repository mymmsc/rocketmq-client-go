package producer

import (
	"time"

	"github.com/mymmsc/go-rocketmq-client/v1/client"
	"github.com/mymmsc/go-rocketmq-client/v1/client/rpc"
)

type mqClient interface {
	RegisterProducer(p client.Producer) error
	UnregisterProducer(group string)
	SendMessageSync(broker string, body []byte, h *rpc.SendHeader, timeout time.Duration) (*rpc.SendResponse, error)
	UpdateTopicRouterInfoFromNamesrv(topic string) error
}
