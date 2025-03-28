package eth_base

import (
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	tokenID    string
	erc20Token bool
}

func NewToken(tokenID string) token_adapter.Token {
	return &Token{
		tokenID:    tokenID,
		erc20Token: false,
	}
}

func NewErc20Token(tokenID string) token_adapter.Token {
	return &Token{
		tokenID:    tokenID,
		erc20Token: true,
	}
}

func (t *Token) BuildTransaction(txInfo *token_adapter.TransactionInfo) (token_adapter.Transaction, error) {
	preTxData, err := prepareBuildTransactionData(txInfo)
	if err != nil {
		return nil, fmt.Errorf("prepare build transaction data error: %w", err)
	}

	tx, err := ParseEthTransaction(preTxData.rawTx)
	if err != nil {
		return nil, fmt.Errorf("prepare eth transaction error: %w", err)
	}

	return &Transaction{tx: tx, PrepareTransactionData: preTxData, token: t}, nil
}

func prepareBuildTransactionData(txInfo *token_adapter.TransactionInfo) (data *PrepareTransactionData, err error) {
	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.RawTxInfo == nil || txInfo.Transaction.RawTxInfo.UnsignedRawTx == nil {
		return nil, fmt.Errorf("transaction info raw tx is nil")
	}

	rawTx := *txInfo.Transaction.RawTxInfo.UnsignedRawTx
	rawTxBytes := common.FromHex(rawTx)

	return &PrepareTransactionData{rawTx: rawTxBytes}, nil
}
