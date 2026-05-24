package auth

import (
	"enactus/internal/models/inputs"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTSecretGenerateTokenAndValidateToken(t *testing.T) {
	jwtSecret := &JWTSecret{Secret: []byte("test-secret")}
	authUser := &inputs.AuthUser{Id: 42, Email: "user@example.com"}

	tokenStr, err := jwtSecret.GenerateToken(authUser, "admin")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	token, err := jwtSecret.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected jwt.MapClaims, got %T", token.Claims)
	}

	if claims["user_id"] != float64(42) {
		t.Fatalf("expected user_id claim 42, got %v", claims["user_id"])
	}

	if claims["role"] != "admin" {
		t.Fatalf("expected role claim admin, got %v", claims["role"])
	}
}

func TestJWTSecretValidateTokenRejectsInvalidToken(t *testing.T) {
	jwtSecret := &JWTSecret{Secret: []byte("test-secret")}

	token, err := jwtSecret.ValidateToken("not-a-jwt")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	if token != nil {
		t.Fatalf("expected nil token, got %v", token)
	}
}

func TestJWTSecretValidateTokenRejectsWrongSecret(t *testing.T) {
	issuer := &JWTSecret{Secret: []byte("issuer-secret")}
	validator := &JWTSecret{Secret: []byte("validator-secret")}
	authUser := &inputs.AuthUser{Id: 7}

	tokenStr, err := issuer.GenerateToken(authUser, "user")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	token, err := validator.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for token signed with another secret")
	}

	if token != nil {
		t.Fatalf("expected nil token, got %v", token)
	}
}
