package admin

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mymmsc/go-rocketmq-client/v1/admin"
	"github.com/mymmsc/go-rocketmq-client/v1/log"
	"github.com/mymmsc/go-rocketmq-client/v1/message"
	"github.com/mymmsc/go-rocketmq-client/v1/tool/command"
)

func init() {
	cmd := &maxOffset{}
	flags := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	flags.IntVar(&cmd.queueID, "q", -1, "queue id")
	flags.StringVar(&cmd.topic, "t", "", "topic")
	flags.StringVar(&cmd.broker, "b", "", "broker name")
	flags.StringVar(&cmd.namesrvAddrs, "n", "", "name servers")

	cmd.flags = flags

	command.RegisterCommand(cmd)
}

type maxOffset struct {
	broker       string
	topic        string
	queueID      int
	namesrvAddrs string

	flags *flag.FlagSet
}

func (mo *maxOffset) Name() string {
	return "maxoffset"
}

func (mo *maxOffset) Desc() string {
	return "query the max offset"
}

func (mo *maxOffset) Run(args []string) {
	mo.flags.Parse(args)

	if len(mo.broker) == 0 {
		fmt.Println("empty broker:[" + mo.broker + "]")
		mo.Usage()
		return
	}

	if len(mo.topic) == 0 {
		fmt.Println("empty topic: [" + mo.topic + "]")
		mo.Usage()
		return
	}

	if len(mo.namesrvAddrs) == 0 {
		fmt.Println("empty namesrv: [" + mo.namesrvAddrs + "]")
		mo.Usage()
		return
	}

	if mo.queueID < 0 {
		fmt.Printf("bad queueID: [%d]\n", mo.queueID)
		mo.Usage()
		return
	}

	logger := log.Std
	a := admin.New("tool-max-offset", strings.Split(mo.namesrvAddrs, ","), logger)
	a.Start()

	offset, err := a.MaxOffset(&message.Queue{
		BrokerName: mo.broker,
		Topic:      mo.topic,
		QueueID:    uint8(mo.queueID),
	})
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("%d\n", offset)
}

func (mo *maxOffset) Usage() {
	mo.flags.Usage()
}
