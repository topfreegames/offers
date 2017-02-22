CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offer_templates (
  id char(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
  name varchar(255) NOT NULL,
  product_id varchar(255) NOT NULL,
  game_id varchar(255) NOT NULL REFERENCES games(id),
  contents JSONB NOT NULL DEFAULT '{}' ::JSONB,
  metadata JSONB DEFAULT '{}' ::JSONB,
  period JSONB NOT NULL DEFAULT '{}' ::JSONB,
  frequency JSONB NOT NULL DEFAULT '{}' ::JSONB,
  trigger JSONB NOT NULL DEFAULT '{}' ::JSONB,
  enabled bool NOT NULL DEFAULT true,
  placement varchar(255) NOT NULL,
  created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);
