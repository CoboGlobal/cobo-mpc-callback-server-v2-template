package eth_base

import (
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"testing"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// test data
var (
	// ETH transfer raw transaction data
	ethRawTx = "f502843b9ee0418252089447f95183e513da1c599e76611dcb00304d07b6a087038d7ea4c6800088112233445566778883aa36a78080"

	// ETH eip1559 transfer raw transaction data
	ethEip1559RawTx = "02f283aa36a70181de6f8252089447f95183e513da1c599e76611dcb00304d07b6a087038d7ea4c68000881122334455667788c0"

	// ERC20 transfer raw transaction data (method ID: 0xa9059cbb)
	erc20RawTx = "02f8b00180843b9aca00825209944a57e687b9126435a9b19e4a802113e266adebdb80b844a9059cbb000000000000000000000000b4c79dab8f259c7aee6e5b2aa729821864227e84000000000000000000000000000000000000000000000000000000174876e800808401546d4080c001a05d429559c67fc853d656fe02f0b2baadc634e0a9c4cc6ac24b89aa0d8fe3e714a02b886e1bfef3c45b0779034890aa72b68ca2a5ea7c354ba89c1ff8963109925f"
)

func TestNewEthToken(t *testing.T) {
	token := NewEthToken("ETH")
	ethToken, ok := token.(*EthBaseToken)
	assert.True(t, ok)
	assert.Equal(t, "ETH", ethToken.tokenID)
	assert.False(t, ethToken.erc20Token)
}

func TestNewErc20Token(t *testing.T) {
	token := NewErc20Token("USDT")
	erc20Token, ok := token.(*EthBaseToken)
	assert.True(t, ok)
	assert.Equal(t, "USDT", erc20Token.tokenID)
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
						RawTx: &ethRawTx,
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
						RawTx: &ethEip1559RawTx,
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
						RawTx: &erc20RawTx,
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
						RawTx: nil,
					},
				},
			},
			wantError: true,
		},
	}

	token := NewEthToken("ETH")
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

func TestEthBaseTransaction_GetHashes(t *testing.T) {
	rawTxBytes := common.FromHex(ethRawTx)
	ethTx, err := ParseEthTransaction(rawTxBytes)
	assert.NoError(t, err)

	transaction := &EthBaseTransaction{
		token: &EthBaseToken{tokenID: "ETH"},
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		tx: ethTx,
	}

	hashes, err := transaction.GetHashes()
	assert.NoError(t, err)
	assert.Len(t, hashes, 1)
	assert.NotEmpty(t, hashes[0])
}

func TestEthBaseTransaction_GetDestinationAddresses(t *testing.T) {
	tests := []struct {
		name          string
		rawTx         string
		isERC20       bool
		expectedAddrs []string
		wantError     bool
	}{
		{
			name:          "ETH transfer",
			rawTx:         ethRawTx,
			isERC20:       false,
			expectedAddrs: []string{"0x47f95183e513da1c599E76611DCB00304d07b6A0"},
			wantError:     false,
		},
		{
			name:          "ERC20 transfer",
			rawTx:         erc20RawTx,
			isERC20:       true,
			expectedAddrs: []string{"0xb4c79dab8f259c7aee6e5b2aa729821864227e84"},
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTxBytes := common.FromHex(tt.rawTx)
			tx, err := ParseEthTransaction(rawTxBytes)
			assert.NoError(t, err)

			transaction := &EthBaseTransaction{
				token: &EthBaseToken{
					tokenID:    "TEST",
					erc20Token: tt.isERC20,
				},
				PrepareTransactionData: &PrepareTransactionData{
					rawTx: rawTxBytes,
				},
				tx: tx,
			}

			addresses, err := transaction.GetDestinationAddresses()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAddrs, addresses)
			}
		})
	}
}

func TestParseEthTransaction(t *testing.T) {
	tests := []struct {
		name      string
		rawTx     string
		wantError bool
	}{
		{
			name:      "Valid ETH transaction",
			rawTx:     ethRawTx,
			wantError: false,
		},
		{
			name:      "Empty raw tx",
			rawTx:     "",
			wantError: true,
		},
		{
			name:      "Invalid raw tx",
			rawTx:     "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTxBytes := common.FromHex(tt.rawTx)
			tx, err := ParseEthTransaction(rawTxBytes)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.IsType(t, &types.Transaction{}, tx)
			}
		})
	}
}

func TestEthHash(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantError bool
	}{
		{
			name:      "Valid input",
			input:     []byte("test data"),
			wantError: false,
		},
		{
			name:      "Empty input",
			input:     []byte{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := EthHash(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash.String())
			}
		})
	}
}
