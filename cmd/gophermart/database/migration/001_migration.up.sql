CREATE SCHEMA IF NOT EXISTS gophermart;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     first_name VARCHAR(50),
     last_name VARCHAR(50),
     login VARCHAR(50) NOT NULL UNIQUE,
     created_at timestamp,
     password text NOT NULL,
     balance FLOAT DEFAULT 0,
     spend FLOAT DEFAULT 0
);

-- CREATE TABLE IF NOT EXISTS orders (
--     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
--     user_id uuid REFERENCES users(id) ON DELETE CASCADE,
--     number VARCHAR(50) NOT NULL UNIQUE,
--     status VARCHAR(50),
--     uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     accrual FLOAT DEFAULT 0
-- );
--
-- CREATE TABLE IF NOT EXISTS withdrawals (
--     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
--     user_id uuid REFERENCES users(id) ON DELETE CASCADE,
--     order_number VARCHAR (50) NOT NULL UNIQUE,
--     status VARCHAR(50) DEFAULT 'NEW',
--     processed_at TIMESTAMP,
--     sum FLOAT DEFAULT 0
-- );