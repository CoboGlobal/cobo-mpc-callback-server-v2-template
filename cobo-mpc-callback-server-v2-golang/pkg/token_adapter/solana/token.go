package solana

import (
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	tokenID    string
	isSPLToken bool
}

func NewToken(tokenID string) token_adapter.Token {
	return &Token{
		tokenID:    tokenID,
		isSPLToken: false,
	}
}

func NewSPLToken(tokenID string) token_adapter.Token {
	return &Token{
		tokenID:    tokenID,
		isSPLToken: true,
	}
}

func (t *Token) BuildTransaction(txInfo *token_adapter.TransactionInfo) (token_adapter.Transaction, error) {
	preTxData, err := prepareBuildTransactionData(txInfo)
	if err != nil {
		return nil, fmt.Errorf("prepare build transaction data error: %w", err)
	}

	tx, err := ParseSolanaTransaction(preTxData.rawTx)
	if err != nil {
		return nil, fmt.Errorf("prepare solana transaction error: %w", err)
	}

	return &Transaction{tx: tx, PrepareTransactionData: preTxData, token: t}, nil
}

func prepareBuildTransactionData(txInfo *token_adapter.TransactionInfo) (data *PrepareTransactionData, err error) {
	if txInfo == nil || txInfo.Transaction == nil {
		return nil, fmt.Errorf("transaction info raw tx is nil")
	}

	var rawTxBytes []byte
	if txInfo.Transaction.RawTxInfo != nil && txInfo.Transaction.RawTxInfo.UnsignedRawTx != nil {
		rawTx := *txInfo.Transaction.RawTxInfo.UnsignedRawTx
		rawTxBytes = common.FromHex(rawTx)
	}

	var destinationAddress string
	if txInfo.Transaction.Destination.GetActualInstance() != nil {
		destination := txInfo.Transaction.Destination.TransactionTransferToAddressDestination
		if destination != nil && destination.AccountOutput != nil {
			destinationAddress = *destination.AccountOutput.Address
		}
	}

	return &PrepareTransactionData{rawTx: rawTxBytes[:], destinationAddress: destinationAddress}, nil
}
