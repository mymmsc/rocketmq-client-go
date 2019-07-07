package consumer

import (
	"flag"
	"fmt"
	"os"

	"github.com/mymmsc/go-rocketmq-client/v1/admin"
	"github.com/mymmsc/go-rocketmq-client/v1/log"
	"github.com/mymmsc/go-rocketmq-client/v1/tool/command"
)

func init() {
	cmd := &clientIDsOfGroupFetcher{}

	flags := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	flags.StringVar(&cmd.brokerAddr, "a", "", "broker address")
	flags.StringVar(&cmd.group, "g", "", "group")
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", cmd.Name())
		flags.PrintDefaults()
	}

	cmd.flags = flags

	command.RegisterCommand(cmd)
}

type clientIDsOfGroupFetcher struct {
	brokerAddr string
	group      string
	flags      *flag.FlagSet
}

func (c *clientIDsOfGroupFetcher) Name() string {
	return "clientidsOfGroup"
}

func (c *clientIDsOfGroupFetcher) Desc() string {
	return "query client ids"
}

func (c *clientIDsOfGroupFetcher) Run(args []string) {
	c.flags.Parse(args)
	if c.brokerAddr == "" {
		fmt.Printf("bad broker address:%s\n", c.brokerAddr)
		return
	}

	if c.group == "" {
		fmt.Printf("bad group:%s\n", c.group)
		return
	}

	logger := log.Std
	a := admin.New("tool-clientids", []string{"ignore me"}, logger)
	a.Start()

	ids, err := a.GetConsumerIDs(c.brokerAddr, c.group)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return
	}
	fmt.Printf("%v\n", ids)
}

func (c *clientIDsOfGroupFetcher) Usage() {
	c.flags.Usage()
}
