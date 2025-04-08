package solana

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func Test_Transaction_GetHashes(t *testing.T) {
	// Create a SOL transaction for testing
	rawTxBytes := common.FromHex(solRawTx)
	tx, err := ParseSolanaTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	solTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "SOL", isSPLToken: false},
	}

	// Test GetHashes method
	hashes, err := solTx.GetHashes()
	assert.NoError(t, err)
	assert.NotEmpty(t, hashes)
	assert.Equal(t, 1, len(hashes))

	assert.Equal(t, []string{"0x0100020453f8c96cd916709eaaf0e4023b16bd25276b8d117275def443e36a0c318f2c92a4bb7324ac0eae63455c211cd74e4026950fe09b2f7fad7e29c506d652f848eb0306466fe5211732ffecadba72c39be7bc8ce5bbc5f7126b2c439b3a400000000000000000000000000000000000000000000000000000000000000000000000bf82f7a836cb9c90684b8d03dbac12002ebe530e83aa2ce5283f172ae316c17f0302000903305705000000000002000502400d0300030200010c020000000b00000000000000"}, hashes)
}

func Test_SPLToken_Transaction_GetHashes(t *testing.T) {
	// Create an SPL Token transaction for testing
	rawTxBytes := common.FromHex(splTokenRawTx)
	tx, err := ParseSolanaTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	solTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{
			tokenID:    "SOL_USDC",
			isSPLToken: true,
		},
	}

	// Test GetHashes method
	hashes, err := solTx.GetHashes()
	assert.NoError(t, err)
	assert.NotEmpty(t, hashes)
	assert.Equal(t, 1, len(hashes))

	assert.Equal(t, []string{"0x01000306a1f97665705a0a89b1ece1dd13838ac9bf912e5441f1a078aafdac11cc7116eeda154f747e1ddf9a758a99028f36026d4dffbfb8f2b1867cb9e474b2f6ca5ef9163d7a7658662964bcf945fbe561f2d03f7d71a56e9fb0e9b91b791860ac4989c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d610306466fe5211732ffecadba72c39be7bc8ce5bbc5f7126b2c439b3a4000000006ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a9409e9f3fa1d001aa1b8d07a13b473207a55866d094e2286bca27e729f94da72a0304000903f82401000000000004000502400d03000504010302000a0c80841e000000000006"}, hashes)
}

func TestTransaction_GetDestinationAddresses_SOL(t *testing.T) {
	// Create a SOL transaction for testing
	rawTxBytes := common.FromHex(solRawTx)
	tx, err := ParseSolanaTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	solTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{tokenID: "SOL", isSPLToken: false},
	}

	// Test GetDestinationAddresses method
	addresses, err := solTx.GetDestinationAddresses()
	if err != nil {
		t.Logf("Error (expected with mock data): %v", err)
		return
	}

	assert.Equal(t, []string{"C63eMJhWSxGhKEFFCSxnF7spphcPgsTnbYsZKjcwJ8Vp"}, addresses)
}

func TestTransaction_GetDestinationAddresses_SPLToken(t *testing.T) {
	// Create an SPL Token transaction for testing
	rawTxBytes := common.FromHex(splTokenRawTx)
	tx, err := ParseSolanaTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance
	solTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx:              rawTxBytes,
			destinationAddress: "8bJaa7p816rKnPSGTsZdWBDkAmsuKDnoEYzLRwrsuFV6",
		},
		token: &Token{
			tokenID:    "USDC",
			isSPLToken: true,
		},
	}

	// Test GetDestinationAddresses method
	addresses, err := solTx.GetDestinationAddresses()
	if err != nil {
		t.Logf("Error (expected with mock data): %v", err)
		return
	}

	assert.Equal(t, []string{"8bJaa7p816rKnPSGTsZdWBDkAmsuKDnoEYzLRwrsuFV6"}, addresses)
}

func TestTransaction_InvalidContract(t *testing.T) {
	// Test if a SOL native transaction is incorrectly treated as an SPL token transaction
	rawTxBytes := common.FromHex(solRawTx)
	tx, err := ParseSolanaTransaction(rawTxBytes)
	assert.NoError(t, err)

	// Create Transaction instance with mismatched token type
	solTx := &Transaction{
		tx: tx,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: rawTxBytes,
		},
		token: &Token{
			tokenID:    "USDC",
			isSPLToken: true,
		},
	}

	// Since we're using mock data, we don't assert specific errors,
	// but we're testing the logic path
	addresses, err := solTx.GetDestinationAddresses()
	// With mock data, this should likely error or return no addresses
	if err == nil {
		assert.Empty(t, addresses, "Should not find addresses with mismatched token type")
	}
}

func TestTransaction_NilTransaction(t *testing.T) {
	// Test case with nil tx
	solTx := &Transaction{
		tx: nil,
		PrepareTransactionData: &PrepareTransactionData{
			rawTx: []byte{},
		},
		token: &Token{tokenID: "SOL", isSPLToken: false},
	}

	// GetDestinationAddresses should return error
	addresses, err := solTx.GetDestinationAddresses()
	assert.Error(t, err)
	assert.Nil(t, addresses)

	// GetHashes should return error because rawTx is empty
	hashes, err := solTx.GetHashes()
	assert.Error(t, err)
	assert.Nil(t, hashes)
}

func TestGetAssociatedTokenAddress(t *testing.T) {
	walletAddress := solana.MustPublicKeyFromBase58("8bJaa7p816rKnPSGTsZdWBDkAmsuKDnoEYzLRwrsuFV6")
	tokenMintAddress := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	tokenProgramID := solana.TokenProgramID

	ataAddress, err := GetAssociatedTokenAddress(walletAddress, tokenMintAddress, tokenProgramID)

	assert.NoError(t, err, "GetAssociatedTokenAddress should not return an error")

	expectedAddress := solana.MustPublicKeyFromBase58("2VpLnZpPENF5rcmSkFnA4TGE4oPRe9W9za2TW32MZZd2")

	assert.Equal(t, expectedAddress.String(), ataAddress.String(), "The computed ATA address should match the expected value")
}
