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
	ID       uuid.UUID `db:"id" valid:"uuidv4"`
	Name     string    `db:"name" valid:"stringlength(1|255),required"`
	BundleID string    `db:"bundle_id" valid:"stringlength(1|255),required"`
	Metadata dat.JSON  `db:"metadata" valid:"json"`

	//TODO: Validate dates
	CreatedAt dat.NullTime `db:"created_at" valid:""`
	UpdatedAt dat.NullTime `db:"updated_at" valid:""`
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
func UpsertGame(db runner.Connection, game *Game, mr *MixedMetricsReporter) error {
	return mr.WithSegment(SegmentInsert, func() error {
		return db.
			InsertInto("games").
			Columns("name", "bundle_id").
			Record(game).
			Returning("id", "created_at", "updated_at").
			QueryStruct(game)
	})
}
