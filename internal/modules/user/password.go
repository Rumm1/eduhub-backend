package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

var temporaryPasswordAlphabet = []rune("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789")

func generateTemporaryPassword() (string, error) {
	parts := make([]string, 0, 3)

	for i := 0; i < 3; i++ {
		part, err := randomPasswordPart(4)
		if err != nil {
			return "", err
		}

		parts = append(parts, part)
	}

	return strings.Join(parts, "-"), nil
}

func randomPasswordPart(length int) (string, error) {
	result := make([]rune, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(temporaryPasswordAlphabet))))
		if err != nil {
			return "", fmt.Errorf("generate temporary password: %w", err)
		}

		result[i] = temporaryPasswordAlphabet[index.Int64()]
	}

	return string(result), nil
}
