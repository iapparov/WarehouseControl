CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);