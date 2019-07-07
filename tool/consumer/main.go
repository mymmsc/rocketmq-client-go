package consumer

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mymmsc/go-rocketmq-client/v1/consumer"
	"github.com/mymmsc/go-rocketmq-client/v1/log"
	"github.com/mymmsc/go-rocketmq-client/v1/message"
)

var (
	namesrvAddrs string
	tags         string
	topic        string
)

func init() {
	flag.StringVar(&namesrvAddrs, "n", "", "name server address")
	flag.StringVar(&topic, "t", "", "topic")
	flag.StringVar(&tags, "g", "", "tags")
}

type messageQueueChanger struct {
	consumer *consumer.PullConsumer
	logger   log.Logger
}

func (qc *messageQueueChanger) Change(topic string, all, divided []*message.Queue) {
	c := qc.consumer
	for _, q := range divided {
		for {
			pr, err := c.PullSync(q, tags, 0, 32)
			if err != nil {
				qc.logger.Errorf("pull error:%v", err)
				break
			}
			qc.logger.Infof("pull %s result:%d", q, len(pr.Messages))
			break
		}
	}
}

func main() {
	flag.Parse()

	if len(namesrvAddrs) == 0 {
		println("bad namesrvAddrs:" + namesrvAddrs)
		return
	}

	if len(topic) == 0 {
		println("bad topic:" + topic)
		return
	}

	if len(tags) == 0 {
		println("bad tags:" + tags)
		return
	}

	logger, err := newLogger()
	if err != nil {
		return
	}

	c := consumer.NewPullConsumer("test-group", strings.Split(namesrvAddrs, ","), logger)

	err = c.Start()
	if err != nil {
		fmt.Printf("start consumer error:%v", err)
		return
	}

	c.Register([]string{topic}, &messageQueueChanger{consumer: c, logger: logger})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signalChan:
		c.Shutdown()
	}
}

func newLogger() (log.Logger, error) {
	file, err := os.Create("pullconsumer.log")
	if err != nil {
		println("create file error", err.Error())
		return nil, err
	}

	return log.New(file, "", log.Ldefault), err
}
