package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim0xPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x123456", "123456"},
		{"0xabcdef", "abcdef"},
		{"123456", "123456"},
		{"abcdef", "abcdef"},
		{"0x", ""},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := Trim0xPrefix(test.input)
			if output != test.expected {
				t.Errorf("Trim0xPrefix(%s) = %s; want %s", test.input, output, test.expected)
			}
		})
	}
}

func TestBigIntToHex(t *testing.T) {
	dst, ok := new(big.Int).SetString("1", 16)
	if ok {
		t.Log(dst.String())
	}
	dst2, ok := new(big.Int).SetString("1", 10)
	assert.Equal(t, true, ok)
	s := BigIntToHex(dst2)
	t.Log(s)
}
