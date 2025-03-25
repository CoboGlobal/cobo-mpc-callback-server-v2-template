package token_adapter

import coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"

// Transaction represents a generic blockchain transaction
type Transaction interface {
	// GetHashes returns a list of transaction hashes
	GetHashes() ([]string, error)

	// GetDestinationAddresses returns a list of destination addresses
	GetDestinationAddresses() ([]string, error)
}

// Token represents a specific blockchain implementation
type Token interface {
	// BuildTransaction builds a transaction from input data
	BuildTransaction(txInfo *TransactionInfo) (Transaction, error)
}

type TransactionInfo struct {
	SourceAddresses []*coboWaaS2.AddressInfo
	Transaction     *coboWaaS2.Transaction
}
