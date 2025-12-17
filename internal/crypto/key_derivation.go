package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time    = 1         // number of iterations
	Argon2Memory  = 64 * 1024 // 64 MB
	Argon2Threads = 4         // number of parallel threads
	KeyLength     = 32        // 256 bits
	SaltLength    = 32        // 256 bits
)

// GenerateSalt generates a cryptographically secure random salt
// Returns:
//   - []byte: salt
//   - error: error if the salt generation failed
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// DeriveKey derives an encryption key from the master password and salt
//
// Parameters:
//   - masterPassword: master password (string)
//   - salt: salt (32 bytes)
//
// Returns:
//   - derivedKey: encryption key (32 bytes)
func DeriveKey(masterPassword string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(masterPassword),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		KeyLength,
	)
}
