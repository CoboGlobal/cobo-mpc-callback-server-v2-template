package waas2

import (
	"context"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type Getter interface {
	GetTransactionApprovalDetail(ctx context.Context, transactionId string) (*coboWaas2.TransactionApprovalDetail, error)
}
