package validator

import (
	"testing"
)

func TestSignatureValidator_Verify_Secp256k1(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "test", "secp256k1", 1)
	err := validator.Verify()
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_Ed25519(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "test", "ed25519", 1)
	err := validator.Verify()
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_UnsupportedAlgorithm(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "test", "ed25519", 1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_InvalidSignature(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "test", "invalid", -1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_InvalidPubkey(t *testing.T) {
	validator := NewSignatureValidator("test", "invalid", "test", "secp256k1", -1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_InvalidMessage(t *testing.T) {
	validator := NewSignatureValidator("invalid", "test", "test", "secp256k1", 1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}
