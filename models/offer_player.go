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
func GetOfferPlayer(db runner.Connection, gameID, playerID, offerID string, mr *MixedMetricsReporter) (*OfferPlayer, error) {
	var offerPlayer OfferPlayer
	err := mr.WithDatastoreSegment("offer_players", SegmentSelect, func() error {
		return db.
			Select("*").
			From("offer_players").
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", gameID, playerID, offerID).
			QueryStruct(&offerPlayer)
	})

	return &offerPlayer, err
}

//GetOffersByPlayer returns all offers by player
func GetOffersByPlayer(db runner.Connection, gameID, playerID string, mr *MixedMetricsReporter) ([]*OfferPlayer, error) {
	var offersByPlayer []*OfferPlayer
	err := mr.WithDatastoreSegment("offer_players", "select all", func() error {
		return db.
			Select("*").
			From("offer_players").
			Where("game_id = $1 AND player_id = $2", gameID, playerID).
			QueryStructs(&offersByPlayer)
	})

	return offersByPlayer, err
}

//CreateOfferPlayer creates an offer player
func CreateOfferPlayer(db runner.Connection, offerPlayer *OfferPlayer, mr *MixedMetricsReporter) error {
	if offerPlayer.Transactions == nil {
		offerPlayer.Transactions = dat.JSON([]byte(`[]`))
	}
	if offerPlayer.Impressions == nil {
		offerPlayer.Impressions = dat.JSON([]byte(`[]`))
	}
	return mr.WithDatastoreSegment("offer_players", SegmentInsert, func() error {
		return db.
			InsertInto("offer_players").
			Columns("game_id", "player_id", "offer_id", "claim_counter", "claim_timestamp", "view_counter", "view_timestamp", "transactions", "impressions").
			Record(offerPlayer).
			Returning("*").
			QueryStruct(offerPlayer)
	})
}

//ClaimOfferPlayer increments the claim counter and updates the timestamp
func ClaimOfferPlayer(db runner.Connection, offerPlayer *OfferPlayer, t time.Time, mr *MixedMetricsReporter) error {
	return mr.WithDatastoreSegment("offer_players", SegmentUpdate, func() error {
		const incrCounter = dat.UnsafeString("claim_counter + 1")
		return db.
			Update("offer_players").
			Set("claim_counter", incrCounter).
			Set("claim_timestamp", t).
			Set("transactions", offerPlayer.Transactions).
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", offerPlayer.GameID, offerPlayer.PlayerID, offerPlayer.OfferID).
			Returning("claim_counter, claim_timestamp, transactions").
			QueryStruct(offerPlayer)
	})
}

//ViewOfferPlayer increments the view counter and updates the timestamp
func ViewOfferPlayer(db runner.Connection, offerPlayer *OfferPlayer, t time.Time, mr *MixedMetricsReporter) error {
	return mr.WithDatastoreSegment("offer_players", SegmentUpdate, func() error {
		const incrCounter = dat.UnsafeString("view_counter + 1")
		return db.
			Update("offer_players").
			Set("view_counter", incrCounter).
			Set("view_timestamp", t).
			Set("impressions", offerPlayer.Impressions).
			Where("game_id = $1 AND player_id = $2 AND offer_id = $3", offerPlayer.GameID, offerPlayer.PlayerID, offerPlayer.OfferID).
			Returning("view_counter, view_timestamp, impressions").
			QueryStruct(offerPlayer)
	})
}
