// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"time"

	dat "gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//Game represents a tenant in offers API
type Game struct {
	ID       string   `db:"id" json:"id" valid:"matches(^[^-][a-z0-9-]*$),stringlength(1|255),required"`
	Name     string   `db:"name" json:"name" valid:"ascii,stringlength(1|255),required"`
	BundleID string   `db:"bundle_id" json:"bundleId" valid:"stringlength(1|255),required"`
	Metadata dat.JSON `db:"metadata" json:"metadata" valid:"json"`

	//TODO: Validate dates
	CreatedAt dat.NullTime `db:"created_at" json:"createdAt" valid:""`
	UpdatedAt dat.NullTime `db:"updated_at" json:"updatedAt" valid:""`
}

//GetMetadata for game
func (g *Game) GetMetadata() (interface{}, error) {
	var obj interface{}
	err := g.Metadata.Unmarshal(&obj)
	return obj, err
}

//GetGameByID returns a game by it's pk
func GetGameByID(db runner.Connection, id string, mr *MixedMetricsReporter) (*Game, error) {
	var game Game
	err := mr.WithDatastoreSegment("games", SegmentInsert, func() error {
		return db.
			Select("*").
			From("games").
			Where("id = $1", id).
			QueryStruct(&game)
	})

	err = HandleNotFoundError("Game", map[string]interface{}{"ID": id}, err)
	return &game, err
}

//ListGames returns a the full list of games
func ListGames(db runner.Connection, mr *MixedMetricsReporter) ([]*Game, error) {
	var games []*Game
	err := mr.WithDatastoreSegment("games", "select all", func() error {
		return db.
			Select("*").
			From("games").
			QueryStructs(&games)
	})
	return games, err
}

//UpsertGame updates a game with new meta or insert with the new UUID
//func UpsertGame(db runner.Connection, game *Game, t time.Time, mr *MixedMetricsReporter) error {
func UpsertGame(db runner.Connection, game *Game, t time.Time, mr *MixedMetricsReporter) error {
	game.UpdatedAt = dat.NullTimeFrom(t)
	return mr.WithDatastoreSegment("games", SegmentUpsert, func() error {
		return db.
			Upsert("games").
			Columns("id", "name", "bundle_id", "updated_at").
			Record(game).
			Where("id=$1", game.ID).
			Returning("created_at", "updated_at").
			QueryStruct(game)
	})
}
