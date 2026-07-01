package password

import "golang.org/x/crypto/bcrypt"

const defaultCost = bcrypt.DefaultCost

func Hash(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), defaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Compare(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
