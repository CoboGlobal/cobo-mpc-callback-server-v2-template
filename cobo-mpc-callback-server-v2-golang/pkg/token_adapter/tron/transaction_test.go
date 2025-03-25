package tron

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_GetHashes(t *testing.T) {
	// Create a TRX transaction for testing
	rawTxBytes := common.FromHex(trxRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	tronTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRON", trc20Token: false},
	}

	// Test GetHashes method
	hashes, err := tronTx.GetHashes()
	assert.NoError(t, err)
	assert.NotEmpty(t, hashes)
	assert.Equal(t, 1, len(hashes))

	fmt.Printf("Hashes: %v\n", hashes)

	// Ensure hash has the correct format (0x prefixed hex string)
	assert.True(t, len(hashes[0]) > 2)
	assert.Equal(t, "0x", hashes[0][:2])
}

func TestTransaction_GetDestinationAddresses_TRX(t *testing.T) {
	// Create a TRX transaction for testing
	rawTxBytes := common.FromHex(trxRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	tronTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRX", trc20Token: false},
	}

	// Test GetDestinationAddresses method
	addresses, err := tronTx.GetDestinationAddresses()
	assert.NoError(t, err)
	assert.NotEmpty(t, addresses)

	// Verify address format and count
	assert.Equal(t, 1, len(addresses))

	// Tron addresses typically start with "T" and are 34 characters long
	assert.True(t, len(addresses[0]) > 0)
	// Note: Specific address format validation depends on your implementation
}

func TestTransaction_GetDestinationAddresses_TRC20(t *testing.T) {
	// Create a TRC20 transaction for testing
	rawTxBytes := common.FromHex(trc20RawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	tronTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRX_USDT", trc20Token: true},
	}

	// Test GetDestinationAddresses method
	addresses, err := tronTx.GetDestinationAddresses()
	assert.NoError(t, err)
	assert.NotEmpty(t, addresses)

	// Verify address format and count
	assert.Equal(t, 1, len(addresses))

	// Tron addresses typically start with "T" and are 34 characters long
	assert.True(t, len(addresses[0]) > 0)
	// Note: Specific address format validation depends on your implementation
}

func TestTransaction_GetDestinationAddresses_InvalidContract(t *testing.T) {
	// Test contract type mismatch scenarios

	// Using TRX transaction but setting trc20Token to true
	rawTxBytes := common.FromHex(trxRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	tronTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRX", trc20Token: true},
	}

	// Expect error since TRX transaction is not TriggerSmartContract type
	addresses, err := tronTx.GetDestinationAddresses()
	assert.Error(t, err)
	assert.Nil(t, addresses)

	// Using TRC20 transaction but setting trc20Token to false
	rawTxBytes = common.FromHex(trc20RawTx)
	tx, err = ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	tronTx = &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRX_USDT", trc20Token: false},
	}

	// Expect error since TRC20 transaction is not TransferContract type
	addresses, err = tronTx.GetDestinationAddresses()
	assert.Error(t, err)
	assert.Nil(t, addresses)
}

func TestTransaction_NilTransaction(t *testing.T) {
	// Test case with nil tx
	tronTx := &Transaction{
		tx: nil,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: []byte{},
		},
		token: &Token{tokenID: "TRX", trc20Token: false},
	}

	// GetDestinationAddresses should return error
	addresses, err := tronTx.GetDestinationAddresses()
	assert.Error(t, err)
	assert.Nil(t, addresses)

	// GetHashes may also return error, depending on implementation
	hashes, err := tronTx.GetHashes()
	// Note: If GetHashes doesn't return error when tx is nil, modify this assertion
	assert.Error(t, err)
	assert.Nil(t, hashes)
}

func TestParseTronTransaction(t *testing.T) {
	// Test valid transaction
	rawTxBytes := common.FromHex(trxRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// Test invalid transaction data
	invalidRawTx := []byte{0x01, 0x02, 0x03}
	tx, err = ParseTronTransaction(invalidRawTx)
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestTronHash(t *testing.T) {
	// Test valid transaction
	rawTxBytes := common.FromHex(trxRawTx)
	hash, err := TronHash(rawTxBytes)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test invalid transaction data
	invalidRawTx := []byte{0x01, 0x02, 0x03}
	hash, err = TronHash(invalidRawTx)
	assert.Error(t, err)
	assert.Empty(t, hash)
}
