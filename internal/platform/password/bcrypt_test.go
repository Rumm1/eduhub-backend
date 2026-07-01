package password

import "testing"

func TestHashAndComparePassword(t *testing.T) {
	plainPassword := "StrongPassword123!"

	hashedPassword, err := Hash(plainPassword)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashedPassword == plainPassword {
		t.Fatal("hashed password should not be equal to plain password")
	}

	if !Compare(hashedPassword, plainPassword) {
		t.Fatal("expected password comparison to be true")
	}

	if Compare(hashedPassword, "WrongPassword") {
		t.Fatal("expected wrong password comparison to be false")
	}
}
