package waas2

import (
	"context"
	"fmt"
	"strings"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/validator"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-api/waas2"
)

func VerifyTxApprovalDetails(ctx context.Context, txApprovalDetails []*TransactionApprovalDetail) error {
	for i, txApprovalDetail := range txApprovalDetails {
		if txApprovalDetail.TransactionId == "" {
			return fmt.Errorf("txApprovalDetail index %d transaction id is empty", i)
		}

		if txApprovalDetail.Transaction == nil {
			return fmt.Errorf("txApprovalDetail index %d tx id %s transaction is nil", i, txApprovalDetail.TransactionId)
		}

		if txApprovalDetail.ApprovalDetail == nil {
			return fmt.Errorf("txApprovalDetail index %d tx id %s approval detail is nil", i, txApprovalDetail.TransactionId)
		}

		if txApprovalDetail.Templates == nil {
			return fmt.Errorf("txApprovalDetail index %d tx id %s templates is nil", i, txApprovalDetail.TransactionId)
		}

		err := verifyTxApprovalDetail(ctx, txApprovalDetail)
		if err != nil {
			return fmt.Errorf("txApprovalDetail index %d tx id %s failed to verify: %w", i, txApprovalDetail.TransactionId, err)
		}
	}
	return nil
}

func verifyTxApprovalDetail(ctx context.Context, txApprovalDetail *TransactionApprovalDetail) error {
	approvalDetail := txApprovalDetail.ApprovalDetail
	transaction := txApprovalDetail.Transaction

	if approvalDetail.TransactionId != nil && *approvalDetail.TransactionId != transaction.TransactionId {
		return fmt.Errorf("tx %s transaction id is not equal to approval detail transaction id", txApprovalDetail.TransactionId)
	}

	transactionType := strings.ToLower(string(*transaction.Type))

	handleUserDetails := func(templateKey string, userDetails []coboWaas2.ApprovalUserDetail) error {
		approveCount, err := verifyUserDetails(templateKey, userDetails, txApprovalDetail)
		if err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify user details: %w", err)
		}

		if approvalDetail.BrokerUser.Threshold == nil {
			return fmt.Errorf("txApprovalDetail broker user threshold is nil")
		}

		if approveCount < int(*approvalDetail.BrokerUser.Threshold) {
			return fmt.Errorf("txApprovalDetail approve count %d is less than threshold %d", approveCount, *approvalDetail.BrokerUser.Threshold)
		}

		if approvalDetail.BrokerUser.Result == nil {
			return fmt.Errorf("txApprovalDetail broker user result is nil")
		}

		if *approvalDetail.BrokerUser.Result != coboWaas2.APPROVALTRANSACTIONRESULT_APPROVED {
			return fmt.Errorf("txApprovalDetail broker user result is not approved")
		}

		return nil
	}

	if approvalDetail.BrokerUser != nil {
		templateKey := "broker_user"
		if err := handleUserDetails(templateKey, approvalDetail.BrokerUser.UserDetails); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify broker user details: %w", err)
		}
		fmt.Println("broker user details verified and approved")
	}

	if approvalDetail.Spender != nil {
		templateKey := transactionType
		if err := handleUserDetails(templateKey, approvalDetail.Spender.UserDetails); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify spender user details: %w", err)
		}
		fmt.Println("spender user details verified and approved")
	}

	if approvalDetail.Approver != nil {
		templateKey := transactionType
		if err := handleUserDetails(templateKey, approvalDetail.Approver.UserDetails); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify approver user details: %w", err)
		}
		fmt.Println("approver user details verified and approved")
	}

	return nil
}

func verifyUserDetails(templateKey string, userDetails []coboWaas2.ApprovalUserDetail, txApprovalDetail *TransactionApprovalDetail) (int, error) {

	approveCount := 0

	for i, userDetail := range userDetails {
		if userDetail.Pubkey == nil {
			return 0, fmt.Errorf("userDetail index %d pubkey is nil", i)
		}

		if userDetail.Signature == nil {
			return 0, fmt.Errorf("userDetail index %d signature is nil", i)
		}

		if userDetail.Result == nil {
			return 0, fmt.Errorf("userDetail index %d result is nil", i)
		}

		if userDetail.TemplateVersion == nil {
			return 0, fmt.Errorf("userDetail index %d template version is nil", i)
		}

		authResult, err := verifyUserDetail(templateKey, userDetail, txApprovalDetail)
		if err != nil {
			return 0, fmt.Errorf("userDetails index %d failed to verify user detail: %w", i, err)
		}

		if authResult {
			approveCount++
		}
	}
	return approveCount, nil
}

func verifyUserDetail(templateKey string, userDetail coboWaas2.ApprovalUserDetail, txApprovalDetail *TransactionApprovalDetail) (bool, error) {
	// get Template
	templates := txApprovalDetail.Templates
	authTemplate := ""
	found := false
	for _, template := range templates {
		if template.TemplateVersion == nil || template.BusinessKey == nil || template.TemplateText == nil {
			continue
		}

		if *template.TemplateVersion == *userDetail.TemplateVersion && *template.BusinessKey == templateKey {
			found = true
			authTemplate = *template.TemplateText
			break
		}
	}

	if !found {
		return false, fmt.Errorf("Template not found, template version: %s, template key: %s", *userDetail.TemplateVersion, templateKey)
	}

	// get biz data

	// get auth result
	authResult := int(*userDetail.Result)

	// get pubkey
	pubkey := *userDetail.Pubkey

	// get signature
	signature := *userDetail.Signature

	// get message
	message := ""

	authData := &validator.AuthData{
		Template:  authTemplate,
		BizData:   "",
		Result:    authResult,
		Pubkey:    pubkey,
		Signature: signature,
		Message:   message,
	}
	authValidator := validator.NewAuthValidator(authData)
	err := authValidator.VerifyAuthData()
	if err != nil {
		return false, fmt.Errorf("failed to verify user detail: %w", err)
	}

	return authResult == int(coboWaas2.APPROVALRESULT_APPROVED), nil
}
