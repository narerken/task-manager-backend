package utils

import "testing"

func TestHashPasswordAndValidatePassword(t *testing.T) {
	hash, err := HashPassword("super-secret")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	valid, err := ValidatePassword("super-secret", hash)
	if err != nil {
		t.Fatalf("ValidatePassword() error = %v", err)
	}

	if !valid {
		t.Fatal("expected password to be valid")
	}
}

func TestValidatePasswordReturnsFalseForWrongPassword(t *testing.T) {
	hash, err := HashPassword("super-secret")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	valid, err := ValidatePassword("wrong-password", hash)
	if err != nil {
		t.Fatalf("ValidatePassword() error = %v", err)
	}

	if valid {
		t.Fatal("expected password to be invalid")
	}
}

func TestValidatePasswordReturnsErrorForInvalidHash(t *testing.T) {
	valid, err := ValidatePassword("super-secret", "not-a-bcrypt-hash")
	if err == nil {
		t.Fatal("expected error for invalid hash format")
	}

	if valid {
		t.Fatal("expected invalid result for malformed hash")
	}
}
