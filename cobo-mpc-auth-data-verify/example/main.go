package main

import (
	"context"
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-auth-data-verify/validator"
	"github.com/CoboGlobal/cobo-mpc-auth-data-verify/waas2"
)

func main() {

}

func getVerifier(apiSecret string, transactionId string) error {
	//
	waas2Client := waas2.NewClient(apiSecret)
	txDetail, err := waas2Client.GetTransactionApprovalDetail(context.Background(), transactionId)
	if err != nil {
		return fmt.Errorf("failed to get transaction approval detail: %w", err)
	}

	// 	message   string - 不需要，计算出来的
	// 	result    string
	// 	pubkey    string  - int
	// 	signature string
	// 	algorithm string - api 有相应字段
	//  templateContent string - 审核模版 api 有相应字段
	//  bizData string - api 有相应字段
	//  callbackData string - callback 数据

	if txDetail == nil {
		return fmt.Errorf("transaction approval detail is nil")
	}
	if txDetail.Spender == nil {
		return fmt.Errorf("spender is nil")
	}
	txSpender := txDetail.Spender

	// parse total result
	txResult := txSpender.Result
	if txResult == nil {
		return fmt.Errorf("spender result is nil")
	}
	result := *txResult
	if result != "Approved" {
		return fmt.Errorf("spender result is not approved")
	}

	// parse each user detail
	for _, userDetail := range txSpender.UserDetails {
		// parse message
		txMessage := userDetail.Message
		if txMessage == nil {
			return fmt.Errorf("message is nil")
		}
		message := *txMessage

		// parse pubkey
		txPubkey := userDetail.Pubkey
		if txPubkey == nil {
			return fmt.Errorf("pubkey is nil")
		}
		pubkey := *txPubkey

		// parse signature
		txSignature := userDetail.Signature
		if txSignature == nil {
			return fmt.Errorf("signature is nil")
		}
		signature := *txSignature

		// parse result
		txResult := userDetail.Result
		if txResult == nil {
			return fmt.Errorf("result is nil")
		}
		result := *txResult

		// template content
		templateContent := ""
		v := validator.NewValidator(-1, pubkey, signature, "secp256k1", templateContent, bizData)
		err = v.Verify()
		if err != nil {
			return fmt.Errorf("error verifying: %w", err)
		}
		// pubkey is valid in whitelist

		// compare callback data with biz data
		equal, err := verifier.CompareCallbackData(callbackData, bizData)
		if err != nil {
			return fmt.Errorf("error comparing callback data: %w", err)
		}
		if !equal {
			return fmt.Errorf("callback data is not equal to biz data")
		}

		// verify callback data and signing hash
		// TODO: implement this

	}

	return nil
}
