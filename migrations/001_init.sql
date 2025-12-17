CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    master_key_salt BYTEA,
    master_key_verifier BYTEA,
    master_key_created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS resources (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,        -- "credentials" | "text" | "binary" | "card"
    storage VARCHAR(20) NOT NULL,      -- "postgres" | "minio"
    object_key VARCHAR(500),           -- object key in MinIO (if storage = "minio")
    size BIGINT DEFAULT 0,
    metadata JSONB,                    -- additional metadata (encrypted)
    data BYTEA,                        -- data (if storage = "postgres")
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_resources_user_id ON resources(user_id);
CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(type);
CREATE INDEX IF NOT EXISTS idx_users_has_master_key ON users((master_key_salt IS NOT NULL));