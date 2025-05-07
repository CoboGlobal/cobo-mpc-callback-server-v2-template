package validator

import (
	"crypto/ecdsa"
	"crypto/ed25519"
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
	algorithm string
	result    int
}

func NewSignatureValidator(message, pubkey, signature, algorithm string, result int) *SignatureValidator {
	return &SignatureValidator{
		message:   message,
		pubkey:    pubkey,
		signature: signature,
		algorithm: algorithm,
		result:    result,
	}
}

func (v *SignatureValidator) Verify() error {
	signingMessage := getSigningMessage(v.message, v.result)
	messageHash := sha256.Sum256([]byte(signingMessage))

	if v.algorithm == "secp256k1" {
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
	} else if v.algorithm == "ed25519" {
		pubkeyBytes, err := hex.DecodeString(v.pubkey)
		if err != nil {
			return fmt.Errorf("error decoding pubkey: %w", err)
		}
		signature, err := hex.DecodeString(v.signature)
		if err != nil {
			return fmt.Errorf("error decoding signature: %w", err)
		}
		verified := ed25519.Verify(pubkeyBytes, messageHash[:], signature)
		if !verified {
			return fmt.Errorf("signature verification failed")
		}
	} else {
		return fmt.Errorf("unsupported algorithm: %s", v.algorithm)
	}
	return nil
}

func getSigningMessage(message string, result int) string {
	return fmt.Sprintf("%s||%s", message, strconv.Itoa(result))
}
