package validator

import (
	"testing"
)

func TestSignatureValidator_Verify_Secp256k1(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "test", 1)
	err := validator.Verify()
	if err != nil {
		t.Errorf("Verify() error = %v", err)
	}
}
func TestSignatureValidator_Verify_InvalidSignature(t *testing.T) {
	validator := NewSignatureValidator("test", "test", "invalid", -1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_InvalidPubkey(t *testing.T) {
	validator := NewSignatureValidator("test", "invalid", "test", -1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}

func TestSignatureValidator_Verify_InvalidMessage(t *testing.T) {
	validator := NewSignatureValidator("invalid", "test", "test", 1)
	err := validator.Verify()
	if err == nil {
		t.Errorf("Verify() error = %v", err)
	}
}
