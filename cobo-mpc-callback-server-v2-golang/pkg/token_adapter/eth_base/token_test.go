package eth_base

import (
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"testing"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/stretchr/testify/assert"
)

func TestNewEthToken(t *testing.T) {
	token := NewToken("ETH")
	ethToken, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "ETH", ethToken.tokenID)
	assert.False(t, ethToken.erc20Token)
}

func TestNewErc20Token(t *testing.T) {
	token := NewErc20Token("ETH_USDT")
	erc20Token, ok := token.(*Token)
	assert.True(t, ok)
	assert.Equal(t, "ETH_USDT", erc20Token.tokenID)
	assert.True(t, erc20Token.erc20Token)
}

func TestEthBaseToken_BuildTransaction(t *testing.T) {
	tests := []struct {
		name      string
		txInfo    *token_adapter.TransactionInfo
		wantError bool
	}{
		{
			name: "Valid ETH transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &ethRawTx,
					},
				},
			},
			wantError: false,
		},
		{
			name: "Valid eip1559 ETH transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &ethEip1559RawTx,
					},
				},
			},
			wantError: false,
		},
		{
			name: "Valid erc20 ETH transaction",
			txInfo: &token_adapter.TransactionInfo{
				Transaction: &coboWaaS2.Transaction{
					RawTxInfo: &coboWaaS2.TransactionRawTxInfo{
						UnsignedRawTx: &erc20RawTx,
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

	token := NewToken("ETH")
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
}
