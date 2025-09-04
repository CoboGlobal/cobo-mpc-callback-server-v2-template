package main

import (
	"context"
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/waas2"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

var (
	pubkeyWhitelist = []string{
		"",
	}

	apiSecret = ""
	env       = coboWaas2.DevEnv
)

func main() {
	// init waas2 client
	transactionIds := []string{"mock_transaction_id"}

	client := waas2.NewClient(apiSecret, env)

	waas2Client := waas2.NewWaas2(client)

	// build transaction and approval details
	txApprovalDetails, err := waas2Client.Build(context.Background(), transactionIds)
	if err != nil {
		panic(fmt.Errorf("failed to build transaction approval details: %w", err))
	}

	config := waas2.Config{
		PubkeyWhitelist: pubkeyWhitelist,
	}

	for _, txApprovalDetail := range txApprovalDetails {
		// verify transaction approval detail
		validator := waas2.NewTxApprovalDetailValidator(txApprovalDetail, &config)
		err = validator.Verify(context.Background())
		if err != nil {
			panic(fmt.Errorf("failed to verify transaction approval detail: %w", err))
		}

		// verify txApprovalDetail (transaction and approval detail)
		// txApprovalDetail and tss callback data are matched
	}
}
