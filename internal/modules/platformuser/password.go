package platformuser

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

var passwordAlphabet = []rune("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789")

func GenerateTemporaryPassword() (string, error) {
	parts := make([]string, 0, 3)

	for i := 0; i < 3; i++ {
		part, err := randomString(4)
		if err != nil {
			return "", err
		}

		parts = append(parts, part)
	}

	return strings.Join(parts, "-"), nil
}

func randomString(length int) (string, error) {
	result := make([]rune, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordAlphabet))))
		if err != nil {
			return "", fmt.Errorf("generate random password: %w", err)
		}

		result[i] = passwordAlphabet[index.Int64()]
	}

	return string(result), nil
}
