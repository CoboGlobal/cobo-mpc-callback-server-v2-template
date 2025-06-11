package validator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

type SignatureValidator struct {
	message   string
	pubkey    string
	signature string
	result    int
}

func NewSignatureValidator(message, pubkey, signature string, result int) *SignatureValidator {
	return &SignatureValidator{
		message:   message,
		pubkey:    pubkey,
		signature: signature,
		result:    result,
	}
}

func (v *SignatureValidator) Verify() error {
	signingMessage := getSigningMessage(v.message, v.result)
	messageHash := sha256.Sum256([]byte(signingMessage))

	pubkeyBytes, err := hex.DecodeString(v.pubkey)
	if err != nil {
		return fmt.Errorf("error decoding pubkey: %w", err)
	}
	signatureBytes, err := hex.DecodeString(v.signature)
	if err != nil {
		return fmt.Errorf("error decoding signature: %w", err)
	}
	pubkey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(pubkeyBytes[:32]),
		Y:     new(big.Int).SetBytes(pubkeyBytes[32:64]),
	}
	signatureR := new(big.Int).SetBytes(signatureBytes[:32])
	signatureS := new(big.Int).SetBytes(signatureBytes[32:64])
	verified := ecdsa.Verify(pubkey, messageHash[:], signatureR, signatureS)
	if !verified {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func getSigningMessage(message string, result int) string {
	return fmt.Sprintf("%s||%s", message, strconv.Itoa(result))
}
