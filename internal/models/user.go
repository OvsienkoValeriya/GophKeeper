package models

import "time"

type User struct {
	ID                 int64      `db:"id"`
	Username           string     `db:"login"`
	Password           []byte     `db:"password_hash"`
	Email              string     `db:"email"`
	CreatedAt          time.Time  `db:"created_at"`
	MasterKeySalt      []byte     `db:"master_key_salt"`
	MasterKeyVerifier  []byte     `db:"master_key_verifier"`
	MasterKeyCreatedAt *time.Time `db:"master_key_created_at"`
}
