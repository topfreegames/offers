CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offer_templates (
  id varchar(255) PRIMARY KEY,
  name varchar(255) NOT NULL,
  product_id varchar(255) NOT NULL,
  game_id varchar(255) NOT NULL REFERENCES games(id),
  contents JSONB NOT NULL DEFAULT '{}' ::JSONB,
  metadata JSONB NOT NULL DEFAULT '{}' ::JSONB,
  period JSONB NOT NULL DEFAULT '{}' ::JSONB,
  frequency JSONB NOT NULL DEFAULT '{}' ::JSONB,
  trigger JSONB NOT NULL DEFAULT '{}' ::JSONB,
  enabled bool NOT NULL DEFAULT true
);
