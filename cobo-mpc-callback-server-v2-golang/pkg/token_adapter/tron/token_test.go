package tron

import (
	"testing"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/stretchr/testify/assert"
)

// Sample transaction raw data, in real tests, replace with valid Tron transaction data
var (
	// TRON transfer transaction example
	trxRawTx = "0x0a024e6e220825a91be86ea2bb1e40f0d39ab3d8325a68080112640a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412330a1541d9bfdf13be2f3d0409a7c25bf0e0d314f722ba1a1215412ea8adf8a3b2f667f340268d2622ae06b8b4fde51880ade20470d9b4d19ed8329001809f49"

	// TRON_USDT token transfer transaction example
	trc20RawTx = "0x0a024bcc220802c981bea16eca8a40f89bafa0cd325aae01081f12a9010a31747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e54726967676572536d617274436f6e747261637412740a1541d9bfdf13be2f3d0409a7c25bf0e0d314f722ba1a121541a614f803b6fd780986a42c78ec9c7f77e6ded13c2244a9059cbb000000000000000000000000508f34bf91b6a63a6713723f4df8c74669a285630000000000000000000000000000000000000000000000000000000000989680708fcdfb9ecd329001e0c9a711"
)

func TestNewToken(t *testing.T) {
	token := NewToken("TRON")
	tronToken, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "TRON", tronToken.tokenID)
	assert.False(t, tronToken.trc20Token)
}

func TestNewTrc20Token(t *testing.T) {
	token := NewTrc20Token("TRON_USDT")
	trc20Token, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "TRON_USDT", trc20Token.tokenID)
	assert.True(t, trc20Token.trc20Token)
}

func TestToken_BuildTransaction(t *testing.T) {
	tests := []struct {
		name      string
		txInfo    *token_adapter.TransactionInfo
		wantError bool
	}{
		{
			name: "Valid TRON transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &trxRawTx,
					},
				},
			},
			wantError: false,
		},
		{
			name: "Valid TRC20 transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &trc20RawTx,
					},
				},
			},
			wantError: false,
		},
		{
			name:      "Nil transaction info",
			txInfo:    nil,
			wantError: true,
		},
		{
			name: "Nil raw tx",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: nil,
					},
				},
			},
			wantError: true,
		},
	}

	// Test TRON token
	t.Run("TRON Tests", func(t *testing.T) {
		token := NewToken("TRON")
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tx, err := token.BuildTransaction(tt.txInfo)
				if tt.wantError {
					assert.Error(t, err)
					assert.Nil(t, tx)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, tx)
				}
			})
		}
	})

	// Test TRC20 token
	t.Run("TRC20 Token Tests", func(t *testing.T) {
		token := NewTrc20Token("TRON_USDT")
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tx, err := token.BuildTransaction(tt.txInfo)
				if tt.wantError {
					assert.Error(t, err)
					assert.Nil(t, tx)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, tx)
				}
			})
		}
	})
}

func TestPrepareBuildTransactionData(t *testing.T) {
	// Test normal case
	t.Run("Valid transaction info", func(t *testing.T) {
		rawTx := trxRawTx
		txInfo := &token_adapter.TransactionInfo{
			Transaction: &coboWaaS2.Transaction{
				RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
					UnsignedRawTx: &rawTx,
				},
			},
		}

		data, err := prepareBuildTransactionData(txInfo)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.NotEmpty(t, data.rawTx)
	})

	// Test various error cases
	t.Run("Nil transaction info", func(t *testing.T) {
		data, err := prepareBuildTransactionData(nil)
		assert.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("Nil Transaction field", func(t *testing.T) {
		txInfo := &token_adapter.TransactionInfo{
			Transaction: nil,
		}
		data, err := prepareBuildTransactionData(txInfo)
		assert.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("Nil RawTxInfo field", func(t *testing.T) {
		txInfo := &token_adapter.TransactionInfo{
			Transaction: &coboWaaS2.Transaction{
				RawTxInfo: nil,
			},
		}
		data, err := prepareBuildTransactionData(txInfo)
		assert.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("Nil UnsignedRawTx field", func(t *testing.T) {
		txInfo := &token_adapter.TransactionInfo{
			Transaction: &coboWaaS2.Transaction{
				RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
					UnsignedRawTx: nil,
				},
			},
		}
		data, err := prepareBuildTransactionData(txInfo)
		assert.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("Invalid hex in UnsignedRawTx", func(t *testing.T) {
		invalidHex := "invalid hex string"
		txInfo := &token_adapter.TransactionInfo{
			Transaction: &coboWaaS2.Transaction{
				RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
					UnsignedRawTx: &invalidHex,
				},
			},
		}

		data, err := prepareBuildTransactionData(txInfo)
		// Note: common.FromHex doesn't return errors for invalid hex strings
		// It tries to skip invalid characters, so there's no error here
		// But the result will be an empty or incomplete byte array
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Empty(t, data.rawTx)
	})
}
