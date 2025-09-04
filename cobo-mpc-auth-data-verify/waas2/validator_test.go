package waas2

import (
	"context"
	"testing"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/stretchr/testify/assert"
)

func TestTxApprovalDetailValidator_Verify(t *testing.T) {
	waas2 := NewWaas2(NewClient(testApiSecret, coboWaas2.DevEnv))
	txApprovalDetails, err := waas2.Build(context.Background(), []string{testTransactionId})
	if err != nil {
		t.Fatalf("failed to build tx approval details: %v", err)
	}

	assert.Equal(t, len(txApprovalDetails), 1)

	validator := NewTxApprovalDetailValidator(txApprovalDetails[0], &Config{})
	err = validator.Verify(context.Background())
	if err != nil {
		t.Fatalf("failed to verify tx approval detail: %v", err)
	}
}
