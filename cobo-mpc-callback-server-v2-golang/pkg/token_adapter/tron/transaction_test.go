package tron

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func Test_Transaction_GetHashes(t *testing.T) {
	// Create a TRON transaction for testing
	rawTxBytes := common.FromHex(tronRawTx)
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

	assert.Equal(t, []string{"0xf872c594b7a2d2cc6647a6a090f870655e4856b877ff3e125644e2be27dd8b8d"}, hashes)
}

func Test_TRC20_Transaction_GetHashes(t *testing.T) {
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
		token: &Token{tokenID: "TRON_USDT", trc20Token: false},
	}

	// Test GetHashes method
	hashes, err := tronTx.GetHashes()
	assert.NoError(t, err)
	assert.NotEmpty(t, hashes)
	assert.Equal(t, 1, len(hashes))

	assert.Equal(t, []string{"0x0f6da9739ddfef69237e1e89adc3a4cbd9f11cc7e7451719bd720dba6949430c"}, hashes)
}

func TestTransaction_GetDestinationAddresses_TRON(t *testing.T) {
	// Create a TRON transaction for testing
	rawTxBytes := common.FromHex(tronRawTx)
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

	// Test GetDestinationAddresses method
	addresses, err := tronTx.GetDestinationAddresses()
	assert.NoError(t, err)
	assert.Equal(t, []string{"TEDv9wo5epcVi7pW3rEZUa7wbtPuAKZCer"}, addresses)
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
		token: &Token{tokenID: "TRON_USDT", trc20Token: true},
	}

	// Test GetDestinationAddresses method
	addresses, err := tronTx.GetDestinationAddresses()
	assert.NoError(t, err)
	assert.Equal(t, []string{"THKAcY3fvSyfkzbYxj2aAgxC5R6YAPMJqa"}, addresses)
}

func TestTransaction_GetDestinationAddresses_InvalidContract(t *testing.T) {
	// Test contract type mismatch scenarios

	// Using TRX transaction but setting trc20Token to true
	rawTxBytes := common.FromHex(tronRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)

	tronTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "TRX_USDT", trc20Token: true},
	}

	// Expect error since TRON transaction is not TriggerSmartContract type
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
		token: &Token{tokenID: "TRON", trc20Token: false},
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
		token: &Token{tokenID: "TRON", trc20Token: false},
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
	rawTxBytes := common.FromHex(tronRawTx)
	tx, err := ParseTronTransaction(rawTxBytes)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// Test invalid transaction data
	invalidRawTx := []byte{0x01, 0x02, 0x03}
	tx, err = ParseTronTransaction(invalidRawTx)
	assert.Error(t, err)
	assert.Nil(t, tx)
}
