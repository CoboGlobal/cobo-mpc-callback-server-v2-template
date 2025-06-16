package waas2

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"github.com/test-go/testify/assert"
)

const testApiSecret = ""

func TestGetTransactionApprovalDetail(t *testing.T) {
	client := NewClient(testApiSecret)
	detail, err := client.GetTransactionApprovalDetail(context.Background(), "383d10e7-8d3f-40c6-abec-e4ac36a2a998")
	if err != nil {
		t.Fatalf("failed to get transaction approval detail: %v", err)
	}

	detailJson, err := json.Marshal(detail)
	assert.NoError(t, err)
	fmt.Printf("transaction approval detail: %v\n", string(detailJson))
}

func TestCreateAPIKeys(t *testing.T) {
	apiKey, apiSecret, err := crypto.GenerateApiKey()
	assert.NoError(t, err)
	fmt.Printf("api key: %v\n", apiKey)
	fmt.Printf("api secret: %v\n", apiSecret)
}
