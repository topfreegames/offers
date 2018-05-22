CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offer_versions (
    id char(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id varchar(255) NOT NULL REFERENCES games(id),
    offer_id char(36) NOT NULL REFERENCES offers(id),
    offer_version integer NOT NULL,
    contents JSONB NOT NULL DEFAULT '{}' ::JSONB,
    product_id varchar(255) NOT NULL,
    cost JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX game_id_offer_id_offer_version ON offer_versions (game_id, offer_id, offer_version);
