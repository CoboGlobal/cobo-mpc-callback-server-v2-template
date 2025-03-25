package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// hashPool holds sha256 instances for transaction hashing
var hashPool = sync.Pool{
	New: func() interface{} { return sha256.New() },
}

// Transaction structure with core.TransactionRaw instead of core.Transaction
type Transaction struct {
	token *Token
	*PrepareTransactionData
	tx *core.TransactionRaw
}

// PrepareTransactionData contains the raw transaction data
type PrepareTransactionData struct {
	rawTx []byte
}

// GetHashes implements Transaction interface for Tron
func (t *Transaction) GetHashes() ([]string, error) {
	if len(t.rawTx) == 0 {
		return nil, fmt.Errorf("transaction raw data is empty")
	}
	var hashes []string
	sha, ok := hashPool.Get().(hash.Hash)
	if !ok {
		return nil, fmt.Errorf("failed to get SHA256 from pool")
	}
	defer hashPool.Put(sha)

	sha.Reset()
	sha.Write(t.rawTx)
	hash := sha.Sum(nil)
	hashStr := hex.EncodeToString(hash)
	hashes = append(hashes, "0x"+hashStr)
	return hashes, nil
}

// GetDestinationAddresses implements Transaction interface for Tron
func (t *Transaction) GetDestinationAddresses() ([]string, error) {
	if t.tx == nil {
		return nil, fmt.Errorf("transaction raw data is nil")
	}

	var addresses []string

	// Get transaction contract
	if len(t.tx.GetContract()) == 0 {
		return nil, fmt.Errorf("transaction contract is empty")
	}

	contract := t.tx.GetContract()[0]
	contractType := contract.GetType()

	if t.token.trc20Token {
		// TRC20 transfer
		if contractType != core.Transaction_Contract_TriggerSmartContract {
			return nil, fmt.Errorf("not a TRC20 transfer contract")
		}

		parameter := new(core.TriggerSmartContract)
		if err := proto.Unmarshal(contract.GetParameter().GetValue(), parameter); err != nil {
			return nil, fmt.Errorf("unmarshal trigger smart contract error: %w", err)
		}

		// Decode TRC20 transfer data to get the recipient address
		data := parameter.GetData()
		if len(data) < 4+32 { // method(4) + address(32)
			return nil, fmt.Errorf("invalid TRC20 transfer data length")
		}

		// Check if it's a transfer function (method ID: a9059cbb)
		methodID := hex.EncodeToString(data[:4])
		if methodID != "a9059cbb" {
			return nil, fmt.Errorf("not a TRC20 transfer method")
		}

		// Extract address (note: Tron addresses need special handling)
		addressBytes := data[16:36]
		// Convert to Tron address format
		addressBytes = append([]byte{0x41}, addressBytes...)

		tronAddress := address.Address(addressBytes)
		addresses = append(addresses, tronAddress.String())
	} else {
		// Native TRON transfer
		if contractType != core.Transaction_Contract_TransferContract {
			return nil, fmt.Errorf("not a TRX transfer contract")
		}

		parameter := new(core.TransferContract)
		if err := proto.Unmarshal(contract.GetParameter().GetValue(), parameter); err != nil {
			return nil, fmt.Errorf("unmarshal transfer contract error: %w", err)
		}
		tronAddress := address.Address(parameter.GetToAddress())
		addresses = append(addresses, tronAddress.String())
	}

	return addresses, nil
}

// ParseTronTransaction parses a raw transaction bytes into a Tron TransactionRaw
func ParseTronTransaction(rawTx []byte) (*core.TransactionRaw, error) {
	tx := new(core.TransactionRaw)
	if err := proto.Unmarshal(rawTx, tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return tx, nil
}
