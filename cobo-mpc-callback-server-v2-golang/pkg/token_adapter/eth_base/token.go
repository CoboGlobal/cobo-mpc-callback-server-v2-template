package eth_base

import (
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/ethereum/go-ethereum/common"
)

type EthBaseChain struct {
}

func NewEthBaseToken() token_adapter.Token {
	return &EthBaseChain{}
}

func (e *EthBaseChain) BuildTransaction(txInfo *token_adapter.TransactionInfo) (token_adapter.Transaction, error) {

	preTxData, err := prepareBuildTransactionData(txInfo)
	if err != nil {
		return nil, fmt.Errorf("prepare build transaction data error: %w", err)
	}

	tx, err := ParseEthTransaction(preTxData.rawTx)
	if err != nil {
		return nil, fmt.Errorf("prepare eth transaction error: %w", err)
	}

	return &EthBaseTransaction{tx: tx, PrepareTransactionData: preTxData}, nil
}

func prepareBuildTransactionData(txInfo *token_adapter.TransactionInfo) (data *PrepareTransactionData, err error) {
	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.RawTxInfo == nil || txInfo.Transaction.RawTxInfo.RawTx == nil {
		return nil, fmt.Errorf("transaction info raw tx is nil")
	}

	rawTx := *txInfo.Transaction.RawTxInfo.RawTx
	rawTxBytes := common.FromHex(rawTx)

	return &PrepareTransactionData{rawTx: rawTxBytes}, nil
}
