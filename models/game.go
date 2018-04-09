// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"context"
	"time"

	edat "github.com/topfreegames/extensions/dat"
	dat "gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//Game represents a tenant in offers API
type Game struct {
	ID       string   `db:"id" json:"id" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255)"`
	Name     string   `db:"name" json:"name" valid:"ascii,stringlength(1|255),required"`
	Metadata dat.JSON `db:"metadata" json:"metadata" valid:"JSONObject"`

	//TODO: Validate dates
	CreatedAt dat.NullTime `db:"created_at" json:"createdAt" valid:""`
	UpdatedAt dat.NullTime `db:"updated_at" json:"updatedAt" valid:""`
}

//GetMetadata for game
func (g *Game) GetMetadata() (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := g.Metadata.Unmarshal(&obj)
	return obj, err
}

//GetGameByID returns a game by it's pk
func GetGameByID(ctx context.Context, db runner.Connection, id string, mr *MixedMetricsReporter) (*Game, error) {
	var game Game
	err := mr.WithDatastoreSegment("games", SegmentInsert, func() error {
		builder := db.Select("*")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("games").
			Where("id = $1", id).
			QueryStruct(&game)
	})

	err = HandleNotFoundError("Game", map[string]interface{}{"ID": id}, err)
	return &game, err
}

//ListGames returns a the full list of games
func ListGames(ctx context.Context, db runner.Connection, mr *MixedMetricsReporter) ([]*Game, error) {
	var games []*Game
	err := mr.WithDatastoreSegment("games", "select all", func() error {
		builder := db.Select("*")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("games").
			QueryStructs(&games)
	})
	return games, err
}

//UpsertGame updates a game with new meta or insert with the new UUID
//func UpsertGame(db runner.Connection, game *Game, t time.Time, mr *MixedMetricsReporter) error {
func UpsertGame(ctx context.Context, db runner.Connection, game *Game, t time.Time, mr *MixedMetricsReporter) error {
	if game.Metadata == nil {
		game.Metadata = dat.JSON([]byte(`{}`))
	}
	game.UpdatedAt = dat.NullTimeFrom(t)
	return mr.WithDatastoreSegment("games", SegmentUpsert, func() error {
		builder := db.Upsert("games")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.Columns("id", "name", "updated_at", "metadata").
			Record(game).
			Where("id=$1", game.ID).
			Returning("created_at", "updated_at").
			QueryStruct(game)
	})
}
