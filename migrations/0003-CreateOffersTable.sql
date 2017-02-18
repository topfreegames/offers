CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id varchar(255) NOT NULL REFERENCES games(id),
    offer_template_id varchar(255) NOT NULL REFERENCES offer_templates(id),
    player_id varchar(1000) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NULL,
    claimed_at timestamp NULL
);
