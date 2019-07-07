package consumer

import (
	"github.com/mymmsc/go-rocketmq-client/v1/client"
)

// ExprType the filter type of the subcription
type ExprType = client.ExprType

const (
	// ExprTypeTag see client.ExprTypeTag
	ExprTypeTag = client.ExprTypeTag
	// ExprTypeSQL92 see client.ExprTypeSQL92
	ExprTypeSQL92 = client.ExprTypeSQL92

	exprAll = client.ExprAll
)
