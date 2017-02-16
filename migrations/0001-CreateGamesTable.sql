CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    name varchar(255) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    bundle_id varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at bigint NULL
);
