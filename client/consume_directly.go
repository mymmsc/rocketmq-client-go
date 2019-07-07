package client

import (
	"github.com/mymmsc/go-rocketmq-client/v2/client/rpc"
)

// ConsumeDirectlyResult the flag of the consuming directly
type ConsumeDirectlyResult = rpc.ConsumeDirectlyResult

// predefined `ConsumeDirectlyResult` values
const (
	Success ConsumeDirectlyResult = iota
	Later
	Rollback
	Commit
	Error
	ReturnNil
)

var consumeDirectlyResultStrings = []string{"Success", "Later", "Rollback", "Commit", "Error"}

// ConsumeMessageDirectlyResult consume the message directly, sending by the broker
type ConsumeMessageDirectlyResult = rpc.ConsumeMessageDirectlyResult
