package eth_base

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// hashPool holds LegacyKeccak256 hash for rlpHash.
var hashPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

type Transaction struct {
	token *Token
	*PrepareTransactionData
	tx *types.Transaction
}

type PrepareTransactionData struct {
	rawTx []byte
}

// GetHashes implements Transaction interface for Ethereum
func (t *Transaction) GetHashes() ([]string, error) {
	var hashes []string
	h, err := EthHash(t.rawTx)
	if err != nil {
		return nil, fmt.Errorf("calc eth hash error: %w", err)
	}
	hashes = append(hashes, h.String())
	return hashes, nil
}

// GetDestinationAddresses implements Transaction interface for Ethereum
func (t *Transaction) GetDestinationAddresses() ([]string, error) {
	if t.tx == nil {
		return nil, fmt.Errorf("transaction is nil")
	}

	var addresses []string

	if t.token.erc20Token {
		// parse ERC20 transfer destination addresses
		// ERC20 transfer method ID: 0xa9059cbb
		input := t.tx.Data()
		if len(input) < 4+32 { // method(4) + address(32)
			return nil, fmt.Errorf("invalid ERC20 transfer data length")
		}

		// check method ID: transfer
		methodID := input[:4]
		if common.Bytes2Hex(methodID) != "a9059cbb" {
			return nil, fmt.Errorf("not an ERC20 transfer method")
		}

		// parse destination addresses
		addressBytes := input[4:36]
		address := common.BytesToAddress(addressBytes)
		addresses = append(addresses, address.Hex())

	} else {
		// eth transfer
		if to := t.tx.To(); to != nil {
			addresses = append(addresses, to.Hex())
		} else {
			return []string{}, nil
		}
	}

	return addresses, nil
}

func ParseEthTransaction(rawTx []byte) (*types.Transaction, error) {
	if len(rawTx) < 2 {
		return nil, fmt.Errorf("parse raw tx length too short")
	}
	newRawTx := rawTx
	if rawTx[0] == 0x02 {
		// eip25519
		firstByte := rawTx[0]

		var fields []interface{}
		err := rlp.DecodeBytes(rawTx[1:], &fields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode Bytes raw tx %x : %w", rawTx[1:], err)
		}
		fields = append(fields, big.NewInt(0), big.NewInt(0), big.NewInt(0))
		newRawTx, err = rlp.EncodeToBytes(fields)
		if err != nil {
			return nil, fmt.Errorf("failed to encode Bytes to raw tx: %w", err)
		}
		newRawTx = append([]byte{firstByte}, newRawTx...)
	}

	tx := &types.Transaction{}
	if err := tx.UnmarshalBinary(newRawTx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw tx %x : %w", rawTx, err)
	}
	return tx, nil
}

func EthHash(x []byte) (h common.Hash, err error) {
	sha, ok := hashPool.Get().(crypto.KeccakState)
	if !ok {
		err = fmt.Errorf("failed to get Keccak in hasher pool")
		return
	}

	defer hashPool.Put(sha)

	sha.Reset()
	if _, err = sha.Write(x); err != nil {
		return
	}
	if _, err = sha.Read(h[:]); err != nil {
		return
	}
	return
}
