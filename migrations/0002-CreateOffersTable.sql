CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offers (
  id char(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
  game_id varchar(255) NOT NULL REFERENCES games(id),
  name varchar(255) NOT NULL,
  period JSONB NOT NULL DEFAULT '{}' ::JSONB,
  frequency JSONB NOT NULL DEFAULT '{}' ::JSONB,
  trigger JSONB NOT NULL DEFAULT '{}' ::JSONB,
  placement varchar(255) NOT NULL,
  metadata JSONB DEFAULT '{}' ::JSONB,
  product_id varchar(255) NOT NULL,
  contents JSONB NOT NULL DEFAULT '{}' ::JSONB,
  enabled bool NOT NULL DEFAULT true,
  version integer NOT NULL DEFAULT 1,
  created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX offers_game ON offers (game_id);
