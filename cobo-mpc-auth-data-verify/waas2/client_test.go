package waas2

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"github.com/test-go/testify/assert"
)

const (
	testApiSecret       = ""
	testTransactionId   = ""
	testTemplateKey     = "withdrawal" // transaction type from transaction detail
	testTemplateVersion = "1.0.0"      // template version from approval detail
)

func TestListTransactions(t *testing.T) {
	client := NewClient(testApiSecret, coboWaas2.DevEnv)
	txs, err := client.ListTransactions(context.Background(), []string{testTransactionId})
	if err != nil {
		t.Fatalf("failed to get transactions: %v", err)
	}
	detailJson, err := json.MarshalIndent(txs, "", "  ")
	fmt.Printf("transactions: %v\n", string(detailJson))
}

func TestListTransactionApprovalDetails(t *testing.T) {
	client := NewClient(testApiSecret, coboWaas2.DevEnv)
	txApprovalDetail, err := client.ListTransactionApprovalDetails(context.Background(), []string{testTransactionId})
	if err != nil {
		t.Fatalf("failed to get transaction approval detail: %v", err)
	}

	detailJson, err := json.MarshalIndent(txApprovalDetail, "", "  ")
	assert.NoError(t, err)
	fmt.Printf("transaction approval detail: %v\n", string(detailJson))
}

func TestListTransactionTemplates(t *testing.T) {
	client := NewClient(testApiSecret, coboWaas2.DevEnv)
	txTemplates, err := client.ListTransactionTemplates(context.Background(), []TemplateName{{TemplateKey: testTemplateKey, TemplateVersion: testTemplateVersion}})
	if err != nil {
		t.Fatalf("failed to get transaction templates: %v", err)
	}
	detailJson, err := json.MarshalIndent(txTemplates, "", "  ")
	fmt.Printf("transaction templates: %v\n", string(detailJson))
}

func TestCreateAPIKeys(t *testing.T) {
	apiKey, apiSecret, err := crypto.GenerateApiKey()
	assert.NoError(t, err)
	fmt.Printf("api key: %v\n", apiKey)
	fmt.Printf("api secret: %v\n", apiSecret)
}
