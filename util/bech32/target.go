package bech32

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/mraksoll4/btcutil/bech32"
	"math/big"
	"strconv"
	"strings"
)

func GetTargetFromCommitment(commitment string) (float64, error) {
	if len(commitment) <= 6 || !strings.Contains(commitment, "puzzle") {
		return 0, fmt.Errorf("invalid commitment: %s", commitment)
	}
	// Decode the Bech32 string
	_, decoded, err, expected := Decode(commitment)
	if err != nil {
		if expected != "" {
			commitment = commitment[:len(commitment)-6] + expected
			_, decoded, err, _ = Decode(commitment)
			if err != nil {
				return 0, fmt.Errorf("decode again commitment = %s  failed | %v", commitment, err)
			}
		} else {

			return 0, fmt.Errorf("decode commitment = %s failed | %v", commitment, err)
		}
	}

	dataBytes, err := bech32.ConvertBits(decoded, 5, 8, false)
	if err != nil {
		return 0, fmt.Errorf("convertBits(decoded=%v) | %v", decoded, err)
	}

	b := bytes.NewBuffer(dataBytes)

	c := make([]byte, 48)
	_, err = b.Read(c)
	if err != nil {
		fmt.Println("Error reading bytes:", err)
		return 0, fmt.Errorf("error reading bytes:(c=%v) | %v", c, err)
	}

	yIsPositive := c[47]>>7 != 0
	c[47] &= 0x7f

	x := bigintFromBytes(c, binary.LittleEndian)

	resBytes := bigintToBytes(x, binary.LittleEndian)

	if yIsPositive {
		resBytes[len(resBytes)-1] |= 1 << 7
	}

	crypto1 := sha256.Sum256(resBytes)
	crypto2 := sha256.Sum256(crypto1[:])

	result := bigintFromBytes(crypto2[:8], binary.LittleEndian)

	target := new(big.Int).Div(new(big.Int).SetUint64(1<<64-1), result)

	targetNum, err := strconv.ParseFloat(target.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("ParseFloat %s | %v", target.String(), err)
	}

	return targetNum, nil
}

func bigintFromBytes(data []byte, order binary.ByteOrder) *big.Int {
	result := new(big.Int)
	result.SetBytes(data)

	if order == binary.LittleEndian {
		reverseBytes(data)
		result.SetBytes(data)
	}

	return result
}

func bigintToBytes(x *big.Int, order binary.ByteOrder) []byte {
	byteSlice := x.Bytes()

	if len(byteSlice) < 48 {
		paddedSlice := make([]byte, 48)
		copy(paddedSlice[48-len(byteSlice):], byteSlice)
		byteSlice = paddedSlice
	}

	if order == binary.LittleEndian {
		reverseBytes(byteSlice)
	}

	return byteSlice
}

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
