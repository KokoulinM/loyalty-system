-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    login VARCHAR(50) NOT NULL UNIQUE,
    created_at timestamp,
    password text NOT NULL
);

-- +goose Down
DROP TABLE users;
