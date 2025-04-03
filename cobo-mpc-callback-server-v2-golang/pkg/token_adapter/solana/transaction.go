package solana

import (
	"encoding/hex"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// Transaction structure
type Transaction struct {
	token *Token
	*PrepareTransactionData
	tx *solana.Transaction
}

// PrepareTransactionData contains the raw transaction data
type PrepareTransactionData struct {
	rawTx              []byte
	destinationAddress string
}

// GetHashes implements Transaction interface for Solana
func (t *Transaction) GetHashes() ([]string, error) {
	if t.tx == nil {
		return nil, fmt.Errorf("transaction data is nil")
	}

	var hashes []string
	rawTxBytes, err := t.tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal binary transaction: %w", err)
	}
	hashStr := hex.EncodeToString(rawTxBytes)
	hashes = append(hashes, "0x"+hashStr)

	return hashes, nil
}

// GetDestinationAddresses implements Transaction interface for Solana
func (t *Transaction) GetDestinationAddresses() ([]string, error) {
	if t.tx == nil {
		return nil, fmt.Errorf("transaction data is nil")
	}

	var addresses []string

	// Process instructions
	for idx, inst := range t.tx.Message.Instructions {
		programID := t.tx.Message.AccountKeys[inst.ProgramIDIndex]

		// Check if this is a Solana System Program transfer
		if programID.Equals(solana.SystemProgramID) && !t.token.isSPLToken {
			// For system transfers,
			// 1 source recipient
			// 2 destination recipient
			if len(inst.Accounts) < 2 {
				return nil, fmt.Errorf("parse system instruction index %v account length %v less than 2", idx, len(inst.Accounts))
			}
			recipientIdx := inst.Accounts[1]
			recipient := t.tx.Message.AccountKeys[recipientIdx]
			addresses = append(addresses, recipient.String())

			// parse inst.Data

			// accounts := make([]*solana.AccountMeta, len(inst.Accounts))
			// for i, idx := range inst.Accounts {
			// 	if int(idx) >= len(t.tx.Message.AccountKeys) {
			// 		continue
			// 	}
			// 	accounts[i] = &solana.AccountMeta{
			// 		PublicKey: t.tx.Message.AccountKeys[idx],
			// 		IsSigner:  false, // todo
			// 		IsWritable: false, // todo
			// 	}
			// }
			// decoded, err := system.DecodeInstruction(accounts, inst.Data)
			// if err != nil {
			// 	return nil, err
			// }

			// // check it is Transfer
			// transfer, ok := decoded.Impl.(*system.Transfer)
			// if !ok {
			// 	return nil, err
			// }
			// _ = *transfer.Lamports

		} else if programID.Equals(solana.TokenProgramID) && t.token.isSPLToken {
			// For SPL token transfers,
			// 0 - source
			// 1 - mint = token address
			// 2 - destination
			// 3 - The source account's owner/delegate = wallet address
			if len(inst.Accounts) < 4 {
				return nil, fmt.Errorf("parse spl token instruction index %v account length %v less than 4", idx, len(inst.Accounts))
			}
			// recipientIdx := inst.Accounts[0]
			// sourceAccount := t.tx.Message.AccountKeys[recipientIdx]
			recipientIdx := inst.Accounts[1]
			mintAccount := t.tx.Message.AccountKeys[recipientIdx]
			recipientIdx = inst.Accounts[2]
			destinationAccount := t.tx.Message.AccountKeys[recipientIdx]
			// recipientIdx = inst.Accounts[3]
			//s ourceOwnerAccount := t.tx.Message.AccountKeys[recipientIdx]

			if t.destinationAddress == "" {
				return nil, fmt.Errorf("tx destination address is nil")
			}
			destinationOwnerAccount := solana.MustPublicKeyFromBase58(t.destinationAddress)

			desTokenAccount, err := GetAssociatedTokenAddress(
				destinationOwnerAccount,
				mintAccount,
				programID,
			)
			if err != nil {
				return nil, fmt.Errorf("fail to get associated token address: %w", err)
			}
			if desTokenAccount.Equals(destinationAccount) {
				addresses = append(addresses, destinationOwnerAccount.String())
			} else {
				return nil, fmt.Errorf(
					"parse spl token instruction index %v destination token address mismatch %s with %s",
					idx, desTokenAccount, destinationAccount)
			}

			// parse inst.Data
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no destination addresses parse in transaction")
	}

	return addresses, nil
}

// ParseSolanaTransaction parses a raw transaction bytes into a Solana Transaction
func ParseSolanaTransaction(rawTx []byte) (*solana.Transaction, error) {
	tx, err := solana.TransactionFromBase64(string(rawTx))
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction message from base64: %w", err)
	}

	return tx, nil
}

func GetAssociatedTokenAddress(
	walletAddress solana.PublicKey,
	tokenMintAddress solana.PublicKey,
	tokenProgramID solana.PublicKey,
) (solana.PublicKey, error) {

	if tokenProgramID.IsZero() {
		tokenProgramID = solana.TokenProgramID
	}

	seeds := [][]byte{
		walletAddress.Bytes(),
		tokenProgramID.Bytes(),
		tokenMintAddress.Bytes(),
	}

	address, _, err := solana.FindProgramAddress(
		seeds,
		solana.SPLAssociatedTokenAccountProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to find program address: %w", err)
	}

	return address, nil
}
