package waas2

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	waas2 "github.com/CoboGlobal/cobo-waas2-go-api/waas2"
	"github.com/CoboGlobal/cobo-waas2-go-api/waas2/crypto"
	"github.com/test-go/testify/assert"
)

const testApiSecret = "7f7dd12c5c87594e5723b4baea26c9c2d18bd13784255020d00bc4856d8e8013"

// api key: ad8104b0e6f0a4e9ec7d3a60da138e9ee780cbb969ca522ca8d3ff4c9ac1c65d
// api secret: 7f7dd12c5c87594e5723b4baea26c9c2d18bd13784255020d00bc4856d8e8013
func TestListTransactionApprovalDetails(t *testing.T) {
	client := NewClient(testApiSecret, waas2.DevEnv)
	txApprovalDetail, err := client.ListTransactionApprovalDetails(context.Background(), []string{"383d10e7-8d3f-40c6-abec-e4ac36a2a998"})
	if err != nil {
		t.Fatalf("failed to get transaction approval detail: %v", err)
	}

	detailJson, err := json.MarshalIndent(txApprovalDetail, "", "  ")
	assert.NoError(t, err)
	fmt.Printf("transaction approval detail: %v\n", string(detailJson))
}

func TestCreateAPIKeys(t *testing.T) {
	apiKey, apiSecret, err := crypto.GenerateApiKey()
	assert.NoError(t, err)
	fmt.Printf("api key: %v\n", apiKey)
	fmt.Printf("api secret: %v\n", apiSecret)
}
