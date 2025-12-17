package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ClientConfig struct {
	ServerAddress string        `json:"server_address"`
	Timeout       time.Duration `json:"timeout"`
	TLSEnabled    bool          `json:"tls_enabled"`
	TLSCertPath   string        `json:"tls_cert_path"`
}

func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		ServerAddress: "localhost:50051",
		Timeout:       10 * time.Second,
		TLSEnabled:    false,
		TLSCertPath:   "server.crt",
	}
}

func LoadConfig() *ClientConfig {
	config := DefaultConfig()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config
	}

	configPath := filepath.Join(homeDir, ".gophkeeper", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}

	if err := json.Unmarshal(data, config); err != nil {
		return DefaultConfig()
	}

	return config
}

func SaveConfig(config *ClientConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".gophkeeper")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, "config.json"), data, 0600)
}
