ALTER TABLE offers ADD COLUMN filters JSONB NOT NULL DEFAULT '{}' ::JSONB;