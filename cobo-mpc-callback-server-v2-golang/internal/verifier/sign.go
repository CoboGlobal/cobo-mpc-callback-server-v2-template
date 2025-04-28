package verifier

import (
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/utils"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

func (v *TssVerifier) verifySign(detail *coboWaaS2.TSSKeySignRequest, extra *coboWaaS2.TSSKeySignExtra) error {
	if detail == nil || extra == nil {
		return fmt.Errorf("detail or request info is nil")
	}
	if extra.Transaction == nil || extra.Transaction.TokenId == nil {
		return fmt.Errorf("transaction or token id is nil")
	}

	// get token
	tokenID := *extra.Transaction.TokenId
	token, err := token_adapter.NewToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	tx, err := token.BuildTransaction(&token_adapter.TransactionInfo{
		SourceAddresses: extra.SourceAddresses,
		Transaction:     extra.Transaction,
	})
	if err != nil {
		return fmt.Errorf("failed to build transaction: %w", err)
	}

	hashes, err := tx.GetHashes()
	if err != nil {
		return fmt.Errorf("failed to get hashes: %w", err)
	}

	// check hashes
	if !utils.IsSubset(detail.MsgHashList, hashes) {
		return fmt.Errorf("msg hash list %v is not part of hashes %v", detail.MsgHashList, hashes)
	}

	// check destination addresses
	toAddresses, err := tx.GetDestinationAddresses()
	if err != nil {
		return fmt.Errorf("failed to get destination addresses: %w", err)
	}

	if len(v.addressWhitelist) > 0 {
		if !utils.IsSubset(toAddresses, v.addressWhitelist) {
			return fmt.Errorf("destination addresses %v is not part of address whitelist", toAddresses)
		}
	}

	return nil
}
