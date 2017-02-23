CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE offers (
    id char(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id varchar(255) NOT NULL REFERENCES games(id),
    offer_template_id char(36) NOT NULL REFERENCES offer_templates(id),
    player_id varchar(1000) NOT NULL,
    seen_counter integer NOT NULL DEFAULT 0,
    bought_counter integer NOT NULL DEFAULT 0,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NULL,
    claimed_at timestamp WITH TIME ZONE NULL,
    last_seen_at timestamp WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX game_id_player_id_offer_template_id
ON offers (game_id, player_id, offer_template_id)
