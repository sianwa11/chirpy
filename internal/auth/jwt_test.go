package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	secret := "mangumi-nyangumi"
	wrongSecret := "wrong-secret"
	expiresIn := 1 * time.Hour

    // Test successful token creation and validation
    t.Run("Valid token", func(t *testing.T) {
        token, err := MakeJWT(userID, secret, expiresIn)
        if err != nil {
            t.Fatalf("Failed to create token: %v", err)
        }

        // Validate the token
        gotUserID, err := ValidateJWT(token, secret)
        if err != nil {
            t.Fatalf("Failed to validate token: %v", err)
        }

        // Check user ID matches
        if gotUserID != userID {
            t.Errorf("User ID mismatch: got %v, want %v", gotUserID, userID)
        }
    })

    // Test expired token rejection
    t.Run("Expired token", func(t *testing.T) {
        // Create a token that expires almost immediately
        token, err := MakeJWT(userID, secret, 1*time.Millisecond)
        if err != nil {
            t.Fatalf("Failed to create token: %v", err)
        }

        // Wait for token to expire
        time.Sleep(10 * time.Millisecond)

        // Attempt to validate the expired token
        _, err = ValidateJWT(token, secret)
        if err == nil {
            t.Error("Expected error for expired token, got nil")
        }
    })

    // Test wrong secret rejection
    t.Run("Wrong secret", func(t *testing.T) {
        // Create a token with the correct secret
        token, err := MakeJWT(userID, secret, expiresIn)
        if err != nil {
            t.Fatalf("Failed to create token: %v", err)
        }

        // Attempt to validate with wrong secret
        _, err = ValidateJWT(token, wrongSecret)
        if err == nil {
            t.Error("Expected error for wrong secret, got nil")
        }
    })

    // Test malformed token
    t.Run("Malformed token", func(t *testing.T) {
        _, err := ValidateJWT("not-a-valid-token", secret)
        if err == nil {
            t.Error("Expected error for malformed token, got nil")
        }
    })
}