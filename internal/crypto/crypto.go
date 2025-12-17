package crypto

import (
	"encoding/json"
	"fmt"
)

type CryptoService struct {
	derivedKey []byte // key derived from the master password
}

func NewCryptoService(derivedKey []byte) *CryptoService {
	return &CryptoService{
		derivedKey: derivedKey,
	}
}

// SetupMasterKey sets up a new master key from the master password
// Parameters:
//   - masterPassword: master password
//
// Returns:
//   - salt: salt for saving to the server
//   - verifier: verifier for saving to the server
//   - derivedKey: key for using in the current session
//   - error: error if the master key setup failed
func SetupMasterKey(masterPassword string) (salt, verifier, derivedKey []byte, err error) {
	salt, err = GenerateSalt()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	derivedKey = DeriveKey(masterPassword, salt)

	verifier = CreateVerifier(derivedKey)

	return salt, verifier, derivedKey, nil
}

// UnlockWithMasterKey unlocks access to data using the master key
//
// Parameters:
//   - masterPassword: master password
//   - salt: salt received from the server
//   - storedVerifier: verifier received from the server
//
// Returns:
//   - derivedKey: key for encryption/decryption
//   - error: error if the master key is invalid
func UnlockWithMasterKey(masterPassword string, salt, storedVerifier []byte) ([]byte, error) {

	derivedKey := DeriveKey(masterPassword, salt)

	if !ValidateVerifier(derivedKey, storedVerifier) {
		return nil, fmt.Errorf("Invalid master key")
	}

	return derivedKey, nil
}

// EncryptData encrypts data of any type
// Parameters:
//   - data: data to encrypt
//
// Returns:
//   - []byte: encrypted data
//   - error: error if the data encryption failed
func (s *CryptoService) EncryptData(data []byte) ([]byte, error) {
	return Encrypt(data, s.derivedKey)
}

// DecryptData decrypts data
// Parameters:
//   - encryptedData: encrypted data
//
// Returns:
//   - []byte: decrypted data
//   - error: error if the data decryption failed
func (s *CryptoService) DecryptData(encryptedData []byte) ([]byte, error) {
	return Decrypt(encryptedData, s.derivedKey)
}

// EncryptJSON encrypts a structure serialized to JSON
// Parameters:
//   - v: structure to encrypt
//
// Returns:
//   - []byte: encrypted data
//   - error: error if the data encryption failed
func (s *CryptoService) EncryptJSON(v interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.EncryptData(jsonData)
}

// DecryptJSON decrypts data and deserializes from JSON
// Parameters:
//   - encryptedData: encrypted data
//   - v: structure to decrypt
//
// Returns:
//   - error: error if the data decryption failed
func (s *CryptoService) DecryptJSON(encryptedData []byte, v interface{}) error {
	jsonData, err := s.DecryptData(encryptedData)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

// Clear clears the key from memory (when logging out)
func (s *CryptoService) Clear() {

	for i := range s.derivedKey {
		s.derivedKey[i] = 0
	}
	s.derivedKey = nil
}
