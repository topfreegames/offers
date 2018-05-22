// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"context"
	"encoding/json"
	"time"

	edat "github.com/topfreegames/extensions/dat"
	dat "gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//OfferPlayer represents an offer seen by a player
type OfferPlayer struct {
	ID             string       `db:"id" json:"id" valid:"uuidv4,required"`
	GameID         string       `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	PlayerID       string       `db:"player_id" json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	OfferID        string       `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	ClaimCounter   int          `db:"claim_counter" json:"claimCounter" valid:"int"`
	ClaimTimestamp dat.NullTime `db:"claim_timestamp" json:"claimTimestamp" valid:""`
	ViewCounter    int          `db:"view_counter" json:"viewCounter" valid:"int"`
	ViewTimestamp  dat.NullTime `db:"view_timestamp" json:"viewTimestamp" valid:""`
	Transactions   dat.JSON     `db:"transactions" json:"transactions" valid:""`
	Impressions    dat.JSON     `db:"impressions" json:"impressions" valid:""`
}

//GetOfferPlayer returns an offer player
func GetOfferPlayer(ctx context.Context, db runner.Connection, gameID, playerID, offerID string, mr *MixedMetricsReporter) (*OfferPlayer, error) {
	var offerPlayer OfferPlayer
	err := mr.WithDatastoreSegment("offer_players", SegmentSelect, func() error {
		builder := db.Select("*")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_players").
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", gameID, playerID, offerID).
			QueryStruct(&offerPlayer)
	})

	return &offerPlayer, err
}

//GetOffersByPlayer returns all offers by player
func GetOffersByPlayer(ctx context.Context, db runner.Connection, gameID, playerID string, mr *MixedMetricsReporter) ([]*OfferPlayer, error) {
	var offersByPlayer []*OfferPlayer
	err := mr.WithDatastoreSegment("offer_players", "select all", func() error {
		builder := db.Select("*")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_players").
			Where("game_id = $1 AND player_id = $2", gameID, playerID).
			QueryStructs(&offersByPlayer)
	})

	return offersByPlayer, err
}

//CreateOfferPlayer creates an offer player
func CreateOfferPlayer(ctx context.Context, db runner.Connection, offerPlayer *OfferPlayer, mr *MixedMetricsReporter) error {
	if offerPlayer.Transactions == nil {
		offerPlayer.Transactions = dat.JSON([]byte(`[]`))
	}
	if offerPlayer.Impressions == nil {
		offerPlayer.Impressions = dat.JSON([]byte(`[]`))
	}
	return mr.WithDatastoreSegment("offer_players", SegmentInsert, func() error {
		builder := db.InsertInto("offer_players")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.
			Columns("game_id", "player_id", "offer_id", "claim_counter", "claim_timestamp", "view_counter", "view_timestamp", "transactions", "impressions").
			Record(offerPlayer).
			Returning("*").
			QueryStruct(offerPlayer)
	})
}

//ClaimOfferPlayer increments the claim counter and updates the timestamp
func ClaimOfferPlayer(ctx context.Context, db runner.Connection, offerPlayer *OfferPlayer, t time.Time, mr *MixedMetricsReporter) error {
	return mr.WithDatastoreSegment("offer_players", SegmentUpdate, func() error {
		const incrCounter = dat.UnsafeString("claim_counter + 1")
		builder := db.Update("offer_players")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.Set("claim_counter", incrCounter).
			Set("claim_timestamp", t).
			Set("transactions", offerPlayer.Transactions).
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", offerPlayer.GameID, offerPlayer.PlayerID, offerPlayer.OfferID).
			Returning("claim_counter, claim_timestamp, transactions").
			QueryStruct(offerPlayer)
	})
}

//ViewOfferPlayer increments the view counter and updates the timestamp
func ViewOfferPlayer(ctx context.Context, db runner.Connection, offerPlayer *OfferPlayer, t time.Time, mr *MixedMetricsReporter) error {
	return mr.WithDatastoreSegment("offer_players", SegmentUpdate, func() error {
		const incrCounter = dat.UnsafeString("view_counter + 1")
		builder := db.Update("offer_players")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.Set("view_counter", incrCounter).
			Set("view_timestamp", t).
			Set("impressions", offerPlayer.Impressions).
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", offerPlayer.GameID, offerPlayer.PlayerID, offerPlayer.OfferID).
			Returning("view_counter, view_timestamp, impressions").
			QueryStruct(offerPlayer)
	})
}

func getViewedOfferNextAt(
	ctx context.Context,
	db runner.Connection,
	gameID, offerID string,
	viewCounter int,
	t time.Time,
	mr *MixedMetricsReporter,
) (int64, error) {
	offer, err := GetOfferByID(ctx, db, gameID, offerID, mr)
	if err != nil {
		return 0, err
	}
	var f FrequencyOrPeriod

	json.Unmarshal(offer.Frequency, &f)
	if f.Max != 0 && viewCounter >= f.Max {
		return 0, nil
	}

	if f.Every != "" {
		duration, err := time.ParseDuration(f.Every)
		if err != nil {
			return 0, err
		}
		return t.Add(duration).Unix(), nil
	}
	return t.Unix(), nil
}
