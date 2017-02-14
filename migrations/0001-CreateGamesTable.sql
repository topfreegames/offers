CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE games (
    id UUID PRIMARY KEY,
    name varchar(255) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at bigint NULL
);