package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
)

const VerifierContext = "gophkeeper-master-key-verifier-v1"

// CreateVerifier creates a verifier for checking the correctness of the master key
//
// Parameters:
//   - derivedKey: derived key (result of DeriveKey)
//
// Returns:
//   - verifier: 32 bytes that are saved on the server
func CreateVerifier(derivedKey []byte) []byte {
	h := hmac.New(sha256.New, derivedKey)
	h.Write([]byte(VerifierContext))
	return h.Sum(nil)
}

// ValidateVerifier checks if the entered master key matches the saved verifier
//
// Parameters:
//   - derivedKey: derived key from the entered master key
//   - storedVerifier: verifier saved on the server
//
// Returns:
//   - true if the master key is correct
//
// Important: use subtle.ConstantTimeCompare to protect from timing attacks
func ValidateVerifier(derivedKey, storedVerifier []byte) bool {
	computedVerifier := CreateVerifier(derivedKey)
	return subtle.ConstantTimeCompare(computedVerifier, storedVerifier) == 1
}
