CREATE TABLE offer_players (
    id char(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id varchar(255) NOT NULL REFERENCES games(id),
    player_id varchar(1000) NOT NULL,
    offer_id varchar(255) NOT NULL REFERENCES offers(id),
    claim_counter integer,
    claim_timestamp timestamp WITH TIME ZONE,
    view_counter integer,
    view_timestamp timestamp WITH TIME ZONE,
    transactions JSONB NOT NULL DEFAULT '[]'::JSONB,
    impressions JSONB NOT NULL DEFAULT '[]'::JSONB
);

CREATE UNIQUE INDEX game_id_player_id_offer_id ON offer_players (game_id, player_id, offer_id)
