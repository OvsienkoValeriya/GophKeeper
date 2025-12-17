package client

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/OvsienkoValeriya/GophKeeper/internal/crypto"
)

var (
	ErrMasterKeyNotUnlocked = errors.New("master key not unlocked - please enter your master key")
)

// sessionTTL is the time to live for the session
const sessionTTL = 30 * time.Minute

type MasterKeyStore struct {
	mu            sync.RWMutex
	derivedKey    []byte
	cryptoService *crypto.CryptoService
	isUnlocked    bool
	sessionFile   string
}

func NewMasterKeyStore() *MasterKeyStore {
	home, _ := os.UserHomeDir()
	sessionFile := filepath.Join(home, ".gophkeeper", "session.key")

	store := &MasterKeyStore{
		sessionFile: sessionFile,
	}

	store.tryLoadSession()

	return store
}

func (s *MasterKeyStore) tryLoadSession() {
	info, err := os.Stat(s.sessionFile)
	if err != nil {
		return
	}

	if time.Since(info.ModTime()) > sessionTTL {
		os.Remove(s.sessionFile)
		return
	}

	data, err := os.ReadFile(s.sessionFile)
	if err != nil {
		return
	}

	derivedKey, err := hex.DecodeString(string(data))
	if err != nil || len(derivedKey) != 32 {
		os.Remove(s.sessionFile)
		return
	}

	s.derivedKey = derivedKey
	s.cryptoService = crypto.NewCryptoService(derivedKey)
	s.isUnlocked = true
}

func (s *MasterKeyStore) saveSession() error {
	dir := filepath.Dir(s.sessionFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data := hex.EncodeToString(s.derivedKey)
	return os.WriteFile(s.sessionFile, []byte(data), 0600)
}

func (s *MasterKeyStore) clearSession() {
	os.Remove(s.sessionFile)
}

// Unlock unlocks the storage with master key
// Parameters:
//   - masterPassword: master key entered by the user
//   - salt: salt received from the server
//   - verifier: verifier received from the server
func (s *MasterKeyStore) Unlock(masterPassword string, salt, verifier []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	derivedKey, err := crypto.UnlockWithMasterKey(masterPassword, salt, verifier)
	if err != nil {
		return err
	}

	s.derivedKey = derivedKey
	s.cryptoService = crypto.NewCryptoService(derivedKey)
	s.isUnlocked = true

	s.saveSession()

	return nil
}

// SetupAndUnlock sets up a new master key and unlocks the storage
// Parameters:
//   - masterPassword: master key entered by the user
//
// Returns:
//   - salt: salt received from the server
//   - verifier: verifier received from the server
//   - error: error if the master key is invalid
func (s *MasterKeyStore) SetupAndUnlock(masterPassword string) (salt, verifier []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	salt, verifier, derivedKey, err := crypto.SetupMasterKey(masterPassword)
	if err != nil {
		return nil, nil, err
	}

	s.derivedKey = derivedKey
	s.cryptoService = crypto.NewCryptoService(derivedKey)
	s.isUnlocked = true

	s.saveSession()

	return salt, verifier, nil
}

func (s *MasterKeyStore) GetCryptoService() (*crypto.CryptoService, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isUnlocked {
		return nil, ErrMasterKeyNotUnlocked
	}

	return s.cryptoService, nil
}

func (s *MasterKeyStore) IsUnlocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isUnlocked
}

func (s *MasterKeyStore) Lock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clearSession()

	if s.cryptoService != nil {
		s.cryptoService.Clear()
	}

	for i := range s.derivedKey {
		s.derivedKey[i] = 0
	}

	s.derivedKey = nil
	s.cryptoService = nil
	s.isUnlocked = false
}
