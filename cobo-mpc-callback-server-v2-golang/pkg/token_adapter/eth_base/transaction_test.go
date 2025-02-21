package eth_base

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// test data
var (
	// ETH transfer raw transaction data
	ethRawTx        = "0xe980842ffee68f825208940f76f604fd7762bd94b48ca2523f69ab9665c97f865af3107a400080018080"
	ethRawTxHash    = "0x1acf0342f77d22389a9cf6524a2af14dc420928d2942efe0dc0764be2ea5321b"
	ethRawTxDesAddr = "0x0f76F604fd7762Bd94B48CA2523F69ab9665c97f"

	// ETH eip1559 transfer raw transaction data
	ethEip1559RawTx = "02f283aa36a70181de6f8252089447f95183e513da1c599e76611dcb00304d07b6a087038d7ea4c68000881122334455667788c0"

	// ERC20 transfer raw transaction data (method ID: 0xa9059cbb)
	erc20RawTx        = "0xf86a04850127efef2283016f5b94dac17f958d2ee523a2206206994597c13d831ec780b844a9059cbb0000000000000000000000008b45b84e2cf29e5f826797df7e1aa93fc71a2bfd0000000000000000000000000000000000000000000000000000000002faf080018080"
	erc20RawTxHash    = "0xf8d61123554ede2f342c56b4da018ec640c7ded1574a7417f0a0396aba664007"
	erc20RawTxDesAddr = "0x8B45b84e2cF29E5F826797dF7e1Aa93FC71a2bfd"
)

func TestEthBaseTransaction_GetHashes(t *testing.T) {
	tests := []struct {
		name      string
		rawTx     string
		hashHex   string
		isErc20   bool
		wantError bool
	}{
		{
			name:      "ETH transfer transaction",
			rawTx:     ethRawTx,
			hashHex:   ethRawTxHash,
			isErc20:   false,
			wantError: false,
		},
		{
			name:      "ERC20 transfer transaction",
			rawTx:     erc20RawTx,
			hashHex:   erc20RawTxHash,
			isErc20:   true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTxBytes := common.FromHex(tt.rawTx)
			tx, err := ParseEthTransaction(rawTxBytes)
			assert.NoError(t, err, "Failed to parse transaction")

			transaction := &EthBaseTransaction{
				token: &EthBaseToken{
					tokenID:    "TEST",
					erc20Token: tt.isErc20,
				},
				PrepareTransactionData: &PrepareTransactionData{
					rawTx: rawTxBytes,
				},
				tx: tx,
			}

			hashes, err := transaction.GetHashes()
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, hashes, 1, "Should return exactly one hash")
			assert.Equal(t, tt.hashHex, hashes[0], "Transaction hash does not match expected value")
		})
	}
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
			expectedAddrs: []string{ethRawTxDesAddr},
			wantError:     false,
		},
		{
			name:          "ERC20 transfer",
			rawTx:         erc20RawTx,
			isERC20:       true,
			expectedAddrs: []string{erc20RawTxDesAddr},
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
		wantHash  []string
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
