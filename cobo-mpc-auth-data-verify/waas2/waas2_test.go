package waas2

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

func TestWaas2_Build(t *testing.T) {
	waas2 := NewWaas2(NewClient(testApiSecret, coboWaas2.DevEnv))
	txApprovalDetails, err := waas2.Build(context.Background(), []string{testTransactionId})
	if err != nil {
		t.Fatalf("failed to build tx approval details: %v", err)
	}
	detailJson, err := json.MarshalIndent(txApprovalDetails, "", "  ")
	fmt.Printf("tx approval details: %v\n", string(detailJson))
}

func TestWaas2_getTemplatesByApprovalDetail(t *testing.T) {
	// Create a mock client for testing
	waas2 := NewWaas2(NewClient(testApiSecret, coboWaas2.DevEnv))

	// Create mock approval detail with BrokerUser, Spender, and Approver
	approvalDetail := coboWaas2.ApprovalDetail{
		TransactionId: stringPtr(testTransactionId),
		AddressOwner: &coboWaas2.RoleDetail{
			UserDetails: []coboWaas2.ApprovalUserDetail{
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
			},
		},
		Spender: &coboWaas2.RoleDetail{
			UserDetails: []coboWaas2.ApprovalUserDetail{
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
			},
		},
		Approver: &coboWaas2.RoleDetail{
			UserDetails: []coboWaas2.ApprovalUserDetail{
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
				{
					TemplateVersion: stringPtr(testTemplateVersion),
				},
			},
		},
	}

	// Test with withdrawal transaction type
	transactionType := "withdrawal"
	templates, err := waas2.getTemplatesByApprovalDetail(context.Background(), approvalDetail, transactionType)
	if err != nil {
		t.Fatalf("failed to get templates by approval detail: %v", err)
	}

	expectedTemplateCount := 2
	if len(templates) != expectedTemplateCount {
		t.Errorf("expected %d templates, got %d", expectedTemplateCount, len(templates))
	}

	// Verify template types and versions
	templateMap := make(map[string]bool)
	for _, template := range templates {
		if template.BusinessKey == nil || template.TemplateVersion == nil {
			t.Errorf("template has nil BusinessKey or TemplateVersion")
			continue
		}

		key := getTemplateKey(TemplateName{
			TemplateKey:     businessKeyToTemplateKey(*template.BusinessKey),
			TemplateVersion: *template.TemplateVersion,
		})
		templateMap[key] = true
	}

	// Check for expected templates
	expectedTemplates := []string{
		"address_owner_1.0.0",
		"withdrawal_1.0.0", // from spender
	}

	for _, expected := range expectedTemplates {
		if !templateMap[expected] {
			t.Errorf("expected template %s not found", expected)
		}
	}

	// Verify template caching
	if len(waas2.templateMapCache) != expectedTemplateCount {
		t.Error("template cache should be 2 after fetching templates")
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}
