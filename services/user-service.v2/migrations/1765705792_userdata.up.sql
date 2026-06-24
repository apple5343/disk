CREATE TABLE IF NOT EXISTS userdata (
    id UUID PRIMARY KEY,
    login VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(128) NOT NULL,
    password_hash BYTEA NOT NULL CHECK (octet_length(password_hash) > 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_userdata_login ON userdata(login);
