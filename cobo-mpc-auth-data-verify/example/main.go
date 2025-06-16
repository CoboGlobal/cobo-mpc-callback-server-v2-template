package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/ZhaoZheCobo/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/validator"
	// "github.com/CoboGlobal/cobo-mpc-auth-data-verify/waas2"
)

var pubkeyWhitelist = []string{
	"",
}

func main() {
	transactionID := "mock_transaction_id"

	// step 1: get auth data for transaction
	authData, err := getAuthData(transactionID)
	if err != nil {
		log.Printf("error getting auth data: %v\n", err)
		return
	}

	// step 2: ensure pubkey is in whitelist
	if !slices.Contains(pubkeyWhitelist, authData.Pubkey) {
		log.Printf("pubkey is not in whitelist: %s\n", authData.Pubkey)
		return
	}

	// step 3: verify auth data
	err = verifyAuthData(authData)
	if err != nil {
		log.Printf("error verifying auth data: %v\n", err)
		return
	}

	// step 4: verify biz data is valid:
	// 1. biz data and tss callback data are matched
	// 2. biz data and transaction detail from waas2 are matched

}

// getAuthData get auth data from waas2
func getAuthData(transactionID string) (*validator.AuthData, error) {
	// waas2Client := waas2.NewClient(apiSecret)
	// txDetail, err := waas2Client.GetTransactionApprovalDetail(context.Background(), transactionId)
	// if err != nil {
	// 	return fmt.Errorf("failed to get transaction approval detail: %w", err)
	// }
	return nil, fmt.Errorf("not implemented")
}

func verifyAuthData(authData *validator.AuthData) error {
	return validator.NewAuthValidator(authData).Verify()
}
