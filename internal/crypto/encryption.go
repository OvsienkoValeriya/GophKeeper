package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

const NonceSize = 12 // 96 bits - standard for AES-GCM

// Encrypt encrypts data using AES-256-GCM
// Parameters:
//   - plaintext: plaintext data to encrypt
//   - key: encryption key (32 bytes, result of DeriveKey)
//
// Returns:
//   - ciphertext: encrypted data in format nonce + encrypted_data
//   - error: error if the encryption failed
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, NonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts data encrypted by the Encrypt function
//
// Parameters:
//   - ciphertext: encrypted data (nonce + encrypted_data)
//   - key: encryption key (32 bytes)
//
// Returns:
//   - plaintext: decrypted data
//   - error: error if the decryption failed (including wrong key)
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := ciphertext[:NonceSize]
	encryptedData := ciphertext[NonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong key or corrupted data): %w", err)
	}

	return plaintext, nil
}
