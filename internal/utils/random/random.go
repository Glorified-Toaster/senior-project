package random

import (
	"crypto/rand"
	"math/big"
)

const (
	letters    = "abcdefghijklmnopqrstuvwxyz"
	uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers    = "0123456789"
	allAllowed = letters + uppercase + numbers
)

func String(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allAllowed))))
		if err != nil {
			return "", err
		}
		result[i] = allAllowed[num.Int64()]
	}
	return string(result), nil
}

func Password(n int) (string, error) {
	if n < 8 {
		n = 8
	}

	result := make([]byte, n)

	// Ensure at least one uppercase
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(uppercase))))
	if err != nil {
		return "", err
	}
	result[0] = uppercase[num.Int64()]

	// Ensure at least one number
	num, err = rand.Int(rand.Reader, big.NewInt(int64(len(numbers))))
	if err != nil {
		return "", err
	}
	result[1] = numbers[num.Int64()]

	// Fill the rest
	for i := 2; i < n; i++ {
		num, err = rand.Int(rand.Reader, big.NewInt(int64(len(allAllowed))))
		if err != nil {
			return "", err
		}
		result[i] = allAllowed[num.Int64()]
	}

	return string(result), nil
}
