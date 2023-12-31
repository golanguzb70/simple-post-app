CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    fullname VARCHAR(128) NOT NULL DEFAULT '',
    user_name VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(64) NOT NULL  UNIQUE,
    hashed_password TEXT NOT NULL,
    refresh_token TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);