package password

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
)

func Hash(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("password is empty")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	sum := checksum(raw, salt)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(sum[:]), nil
}

func Compare(encoded, raw string) bool {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	actual := checksum(raw, salt)
	return subtle.ConstantTimeCompare(expected, actual[:]) == 1
}

func checksum(raw string, salt []byte) [32]byte {
	payload := append([]byte(raw), salt...)
	return sha256.Sum256(payload)
}
