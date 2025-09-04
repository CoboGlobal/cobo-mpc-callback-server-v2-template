package waas2

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/validator"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type TxApprovalDetail struct {
	TransactionId  string
	Transaction    *coboWaas2.Transaction
	ApprovalDetail *coboWaas2.ApprovalDetail
	Templates      []coboWaas2.ApprovalTemplate
}
type Config struct {
	PubkeyWhitelist []string
}

type Validator interface {
	Verify(ctx context.Context) error
}

type TxApprovalDetailValidator struct {
	tad    *TxApprovalDetail
	config *Config
}

func NewTxApprovalDetailValidator(tad *TxApprovalDetail, config *Config) *TxApprovalDetailValidator {
	return &TxApprovalDetailValidator{
		tad:    tad,
		config: config,
	}
}

func (t *TxApprovalDetailValidator) Verify(ctx context.Context) error {
	if t.tad == nil {
		return fmt.Errorf("txApprovalDetail is nil")
	}

	if t.tad.TransactionId == "" {
		return fmt.Errorf("txApprovalDetail transaction id is empty")
	}

	if t.tad.Transaction == nil {
		return fmt.Errorf("txApprovalDetail transaction is nil")
	}

	if t.tad.ApprovalDetail == nil {
		return fmt.Errorf("txApprovalDetail approval detail is nil")
	}

	if t.tad.Templates == nil {
		return fmt.Errorf("txApprovalDetail templates is nil")
	}

	err := t.verifyTxApprovalDetail(ctx)
	if err != nil {
		return fmt.Errorf("txApprovalDetail failed to verify: %w", err)
	}
	return nil
}

func (t *TxApprovalDetailValidator) verifyTxApprovalDetail(ctx context.Context) error {
	approvalDetail := t.tad.ApprovalDetail
	transaction := t.tad.Transaction

	if approvalDetail.TransactionId != nil && *approvalDetail.TransactionId != transaction.TransactionId {
		return fmt.Errorf("tx %s transaction id is not equal to approval detail transaction id", t.tad.TransactionId)
	}

	transactionType := strings.ToLower(string(*transaction.Type))

	handleUserDetails := func(templateKey string, roleDetail *coboWaas2.RoleDetail) error {
		if len(roleDetail.UserDetails) == 0 {
			return fmt.Errorf("role detail user details is empty")
		}

		approveCount, err := t.verifyUserDetails(templateKey, roleDetail.UserDetails)
		if err != nil {
			return fmt.Errorf("user details failed to verify: %w", err)
		}

		if roleDetail.Threshold == nil {
			return fmt.Errorf("role detail threshold is nil")
		}

		if approveCount < int(*roleDetail.Threshold) {
			return fmt.Errorf("user detail approve count %d is less than threshold %d", approveCount, *roleDetail.Threshold)
		}

		if roleDetail.Result == nil {
			return fmt.Errorf("role detail result is nil")
		}

		if *roleDetail.Result != coboWaas2.APPROVALTRANSACTIONRESULT_APPROVED {
			return fmt.Errorf("role detail result is not approved")
		}

		return nil
	}

	if approvalDetail.AddressOwner != nil {
		templateKey := "address_owner"
		if err := handleUserDetails(templateKey, approvalDetail.AddressOwner); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify address owner details: %w", err)
		}
		fmt.Println("address owner details verified and approved")
	}

	if approvalDetail.Spender != nil {
		templateKey := transactionType
		if err := handleUserDetails(templateKey, approvalDetail.Spender); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify spender user details: %w", err)
		}
		fmt.Println("spender user details verified and approved")
	}

	if approvalDetail.Approver != nil {
		templateKey := transactionType
		if err := handleUserDetails(templateKey, approvalDetail.Approver); err != nil {
			return fmt.Errorf("txApprovalDetail failed to verify approver user details: %w", err)
		}
		fmt.Println("approver user details verified and approved")
	}

	return nil
}

func (t *TxApprovalDetailValidator) verifyUserDetails(templateKey string, userDetails []coboWaas2.ApprovalUserDetail) (int, error) {
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

		authResult, err := t.verifyUserDetail(templateKey, userDetail)
		if err != nil {
			return 0, fmt.Errorf("userDetails index %d failed to verify user detail: %w", i, err)
		}

		if authResult {
			approveCount++
		}
	}
	return approveCount, nil
}

func (t *TxApprovalDetailValidator) verifyUserDetail(templateKey string, userDetail coboWaas2.ApprovalUserDetail) (bool, error) {
	// get Template
	templates := t.tad.Templates
	authTemplate := ""
	found := false
	for _, template := range templates {
		if template.TemplateVersion == nil || template.BusinessKey == nil || template.TemplateText == nil {
			continue
		}

		if *template.TemplateVersion == *userDetail.TemplateVersion && businessKeyToTemplateKey(*template.BusinessKey) == templateKey {
			found = true
			authTemplate = *template.TemplateText
			break
		}
	}

	if !found {
		return false, fmt.Errorf("Template not found, template version: %s, template key: %s", *userDetail.TemplateVersion, templateKey)
	}

	// get biz data
	bizData, err := buildBizData(t.tad.Transaction, userDetail)
	if err != nil {
		return false, fmt.Errorf("failed to merge transaction and user detail: %w", err)
	}

	// get auth result
	authResult := int(*userDetail.Result)

	// get pubkey
	pubkey := *userDetail.Pubkey
	if t.config != nil && t.config.PubkeyWhitelist != nil {
		if !slices.Contains(t.config.PubkeyWhitelist, pubkey) {
			return false, fmt.Errorf("pubkey %s is not in whitelist", pubkey)
		}
	}

	// get signature
	signature := *userDetail.Signature

	// // get message // api not get message
	// message := ""

	authData := &validator.AuthData{
		Template:  authTemplate,
		BizData:   bizData,
		Result:    authResult,
		Pubkey:    pubkey,
		Signature: signature,
		// Message:   message,
	}
	authValidator := validator.NewAuthValidator(authData)
	err = authValidator.VerifyAuthData()
	if err != nil {
		return false, fmt.Errorf("failed to verify user detail: %w", err)
	}

	return authResult == int(coboWaas2.APPROVALRESULT_APPROVED), nil
}

// buildBizData merges the properties of Transaction and ApprovalUserDetail into a single JSON string
func buildBizData(transaction *coboWaas2.Transaction, userDetail coboWaas2.ApprovalUserDetail) (string, error) {
	// Create a map to hold the merged data
	mergedData := make(map[string]interface{})

	// Add transaction properties
	if transaction != nil {
		// Convert transaction to map
		txMap, err := transaction.ToMap()
		if err != nil {
			return "", fmt.Errorf("failed to convert transaction to map: %w", err)
		}

		// Add transaction properties with "tx_" prefix to avoid conflicts
		for key, value := range txMap {
			mergedData[key] = value
		}
	}
	userDetailMap, err := userDetail.ToMap()
	if err != nil {
		return "", fmt.Errorf("failed to convert user detail to map: %w", err)
	}
	for key, value := range userDetailMap {
		mergedData[key] = value
	}

	// Convert merged data to JSON string
	jsonData, err := json.Marshal(mergedData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged data to JSON: %w", err)
	}

	return string(jsonData), nil
}
