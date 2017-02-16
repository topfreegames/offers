CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offer_templates (
  id char(36) PRIMARY KEY,
  name varchar(255) NOT NULL,
  pid varchar(255) NOT NULL,
  gameid varchar(255) NOT NULL,
  contents JSONB NOT NULL DEFAULT '{}' ::JSONB,
  metadata JSONB NOT NULL DEFAULT '{}' ::JSONB,
  period JSONB NOT NULL DEFAULT '{}' ::JSONB,
  frequency JSONB NOT NULL DEFAULT '{}' ::JSONB,
  trigger JSONB NOT NULL DEFAULT '{}' ::JSONB
);
