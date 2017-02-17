CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE games (
    id varchar(255) PRIMARY KEY,
    name varchar(255) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    bundle_id varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NULL
);
