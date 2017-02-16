// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/mgutz/dat"
	runner "github.com/mgutz/dat/sqlx-runner"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
)

//Game represents a tenant in offers API
type Game struct {
	ID        uuid.UUID    `db:"id"`
	Name      string       `db:"name"`
	BundleID  string       `db:"bundle_id"`
	Metadata  dat.JSON     `db:"metadata"`
	CreatedAt dat.NullTime `db:"created_at"`
	UpdatedAt dat.NullTime `db:"updated_at"`
}

//GetMetadata for game
func (g *Game) GetMetadata() (interface{}, error) {
	var obj interface{}
	err := g.Metadata.Unmarshal(&obj)
	return obj, err
}

//GetGameByID returns a game by it's pk
func GetGameByID(db runner.Connection, id uuid.UUID) (*Game, error) {
	var game Game
	err := db.
		Select("*").
		From("games").
		Where("id = $1", id).
		QueryStruct(&game)

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, errors.NewGameNotFoundError(map[string]interface{}{
				"ID": id,
			})
		}
		return nil, err
	}

	return &game, nil
}

//UpsertGame updates a game with new meta or insert with the new UUID
func UpsertGame(db runner.Connection, game *Game) error {
	return db.
		InsertInto("games").
		Columns("name", "bundle_id").
		Record(game).
		Returning("id", "created_at", "updated_at").
		QueryStruct(game)
}
