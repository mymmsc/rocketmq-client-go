package main

import (
	"os"

	"github.com/mymmsc/go-rocketmq-client/v2/tool/command"

	_ "github.com/mymmsc/go-rocketmq-client/v2/tool/admin"
	_ "github.com/mymmsc/go-rocketmq-client/v2/tool/consumer"
	_ "github.com/mymmsc/go-rocketmq-client/v2/tool/message"
)

func main() {
	command.Run(os.Args[1:])
}
