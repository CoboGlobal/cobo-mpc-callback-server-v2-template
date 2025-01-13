package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Trim0xPrefix removes the "0x" from a hex-string.
func Trim0xPrefix(hexString string) string {
	if len(hexString) >= 2 && hexString[0:2] == "0x" {
		return hexString[2:]
	}
	return hexString
}

func IntToBytes32(inputInt int) [32]byte {
	var result [32]byte
	binary.BigEndian.PutUint32(result[28:], uint32(inputInt))
	return result
}

func HexToBytes32(inputHex string) ([32]byte, error) {
	inputHex = Trim0xPrefix(inputHex)
	if len(inputHex)%2 != 0 {
		inputHex = "0" + inputHex
	}
	var result [32]byte
	inputBytes, err := hex.DecodeString(inputHex)
	if err != nil {
		return result, err
	}
	copy(result[32-len(inputBytes):], inputBytes)
	return result, nil
}

func Bytes32ToHex(input [32]byte) string {
	hexString := hex.EncodeToString(input[:])
	return "0x" + hexString
}

func BytesToBytes32(input []byte) [32]byte {
	var b32 [32]byte
	if len(input) > 32 {
		return b32
	}

	start := 32 - len(input)
	copy(b32[start:], input)

	return b32
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func BigIntToHex(src *big.Int) string {
	return fmt.Sprintf("0x%x", src.Bytes())
}
