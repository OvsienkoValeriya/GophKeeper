package client

import (
	"encoding/json"
	"os"
	"path"
	"time"
)

type TokenStore interface {
	SaveTokens(accessToken, refreshToken string) error
	LoadTokens() (accessToken, refreshToken string, err error)
	ClearTokens() error
	IsAccessTokenExpired() (bool, error)
}

type TokenRecord struct {
	UserID       uint      `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	HasMasterKey bool      `json:"has_master_key"`
}

type FileTokenStore struct {
	filePath string
}

func NewFileTokenStore() (*FileTokenStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := path.Join(home, ".gophkeeper")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	path := path.Join(dir, "tokens.json")

	return &FileTokenStore{
		filePath: path,
	}, nil
}

// SaveTokensWithUserID saves tokens (access and refresh) to the file with user ID
// Parameters:
//   - userID: user ID
//   - accessToken: access token
//   - refreshToken: refresh token
//   - expiresAt: expiration time
//
// Returns:
//   - error: error if the tokens saving failed
func (s *FileTokenStore) SaveTokensWithUserID(userID uint, accessToken, refreshToken string, expiresAt time.Time) error {
	existing, _ := s.loadRecord()

	record := TokenRecord{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		HasMasterKey: existing.UserID == userID && existing.HasMasterKey,
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0600)
}

// SetHasMasterKey sets the flag has_master_key to the file
// Parameters:
//   - hasMasterKey: flag has_master_key
//
// Returns:
//   - error: error if the flag saving failed
func (s *FileTokenStore) SetHasMasterKey(hasMasterKey bool) error {
	record, err := s.loadRecord()
	if err != nil {
		return err
	}

	record.HasMasterKey = hasMasterKey

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0600)
}

// GetUserID gets the user ID from the file
// Returns:
//   - uint: user ID
//   - error: error if the user ID retrieval failed
func (s *FileTokenStore) GetUserID() (uint, error) {
	record, err := s.loadRecord()
	if err != nil {
		return 0, err
	}
	return record.UserID, nil
}

// HasMasterKey checks the flag has_master_key from the file
// Returns:
//   - bool: true if the flag has_master_key is true, false otherwise
//   - error: error if the flag retrieval failed
func (s *FileTokenStore) HasMasterKey() (bool, error) {
	record, err := s.loadRecord()
	if err != nil {
		return false, err
	}
	return record.HasMasterKey, nil
}

func (s *FileTokenStore) loadRecord() (TokenRecord, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return TokenRecord{}, err
	}

	var record TokenRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return TokenRecord{}, err
	}

	return record, nil
}

// LoadTokens loads the tokens (access and refresh) from the file
// Returns:
//   - string: access token
//   - string: refresh token
//   - error: error if the tokens loading failed
func (s *FileTokenStore) LoadTokens() (string, string, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return "", "", err
	}

	var record TokenRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return "", "", err
	}

	return record.AccessToken, record.RefreshToken, nil
}

// ClearTokens clears the tokens from the file
// Returns:
//   - error: error if the tokens clearing failed
func (s *FileTokenStore) ClearTokens() error {
	if err := os.Remove(s.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// IsAccessTokenExpired checks if the access token is expired
// Returns:
//   - bool: true if the access token is expired, false otherwise
//   - error: error if the access token expiration check failed
func (s *FileTokenStore) IsAccessTokenExpired() (bool, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return true, err
	}

	var record TokenRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return true, err
	}

	return time.Now().After(record.ExpiresAt), nil
}
