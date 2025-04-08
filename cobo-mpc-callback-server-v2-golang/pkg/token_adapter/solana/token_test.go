package solana

import (
	"testing"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/stretchr/testify/assert"
)

// Sample transaction raw data, in real tests, replace with valid Solana transaction data
var (
	// Native SOL transfer transaction example
	solRawTx = "41414541416752542b4d6c7332525a776e717277354149374672306c4a32754e45584a313376524434326f4d4d5938736b715337637953734471356a52567768484e644f51436156442b43624c332b7466696e4642745a532b456a7241775a47622b5568467a4c2f374b323663734f623537794d3562764639784a724c454f624f6b41414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414c2b4339366732793579516145754e4139757345674175766c4d4f67366f733553672f4679726a4673462f4177494143514d7756775541414141414141494142514a4144514d4141774941415177434141414143774141414141414141413d"

	// SPL token SOL_USDC transfer transaction example
	splTokenRawTx = "0x41414541417761682b585a6c63466f4b69624873346430546734724a76354575564548786f4869712f6177527a48455737746f565433522b48642b616459715a416f3832416d314e2f372b3438724747664c6e6b644c4c32796c3735466a3136646c686d4b5753382b5558373557487930443939636156756e374470755274354747437353596e472b6e727a767475744f6a316c38327172795851787362766b77744c32344f5238706749445253396459514d47526d2f6c495263792f2b7974756e4c446d2b65386a4f573778666353617978446d7a704141414141427433323464646c6f5a505a792b46477a7574357242793068653166577a65524f6f7a316858372f414b6c416e70382f6f6441427168754e42364537527a49487056686d304a54694b47764b4a2b63702b55326e4b674d4541416b442b435142414141414141414541415543514130444141554541514d4341416f4d674951654141414141414147"
)

func TestNewToken(t *testing.T) {
	token := NewToken("SOL")
	solToken, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "SOL", solToken.tokenID)
	assert.False(t, solToken.isSPLToken)
}

func TestNewSPLToken(t *testing.T) {
	token := NewSPLToken("SOL_USDC")
	splToken, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "SOL_USDC", splToken.tokenID)
	assert.True(t, splToken.isSPLToken)
}

func TestToken_BuildTransaction(t *testing.T) {
	tests := []struct {
		name      string
		txInfo    *token_adapter.TransactionInfo
		wantError bool
	}{
		{
			name: "Valid SOL transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &solRawTx,
					},
				},
			},
			wantError: false,
		},
		{
			name: "Valid SPL Token transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &splTokenRawTx,
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

	// Test SOL token
	t.Run("SOL Tests", func(t *testing.T) {
		token := NewToken("SOL")
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

	// Test SPL token
	t.Run("SPL Token Tests", func(t *testing.T) {
		token := NewSPLToken("SOL_USDC")
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
		rawTx := solRawTx
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
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Empty(t, data.rawTx)
	})
}
