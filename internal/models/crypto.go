package models

type MasterKeySetup struct {
	Salt     []byte // 32 bytes, random
	Verifier []byte // 32 bytes, HMAC from derived key
}

type EncryptedData struct {
	Nonce      []byte
	Ciphertext []byte
}

type DecryptedResource struct {
	ID       int64
	Name     string
	Type     ResourceType
	Data     []byte
	Metadata map[string]string
}
