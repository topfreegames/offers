// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models

import (
	"github.com/mgutz/dat"
	runner "github.com/mgutz/dat/sqlx-runner"
	"github.com/topfreegames/offers/errors"
	"time"
)

//Offer represents a tenant in offers API
type Offer struct {
	ID              string `db:"id" valid:"uuidv4,required"`
	GameID          string `db:"game_id" valid:"matches(^[a-z0-9]+(\\-[a-z0-9]+)*$),stringlength(1|255),required"`
	OfferTemplateID string `db:"offer_template_id" valid:"uuidv4,required"`
	PlayerID        string `db:"player_id" valid:"ascii,stringlength(1|1000),required"`

	CreatedAt  dat.NullTime `db:"created_at" valid:""`
	UpdatedAt  dat.NullTime `db:"updated_at" valid:""`
	ClaimedAt  dat.NullTime `db:"claimed_at" valid:""`
  LastSeenAt dat.NullTime `db:"last_seen_at" valid:""`
}

//OfferToUpdate has required fields for claiming an offer
type OfferToUpdate struct {
	ID              string `db:"id" valid:"uuidv4,required"`
	GameID          string `db:"game_id" valid:"matches(^[a-z0-9]+(\\-[a-z0-9]+)*$),stringlength(1|255),required"`
	PlayerID        string `db:"player_id" valid:"ascii,stringlength(1|1000),required"`
}

const playerSeenOffersScope = `
	WHERE
		o.game_id = $1
	AND o.player_id = $2
  AND o.last_seen_at IS NOT NULL
`
const playerUnseenOffersScope = `
	WHERE
		o.game_id = $1
	AND o.player_id = $2
  AND o.last_seen_at IS NULL
`

//GetOfferByID returns a offer by it's pk
func GetOfferByID(db runner.Connection, gameID string, id string, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, game_id, offer_template_id, player_id, created_at, updated_at, claimed_at").
			From("offers").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offer)
	})

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, errors.NewModelNotFoundError("Offer", map[string]interface{}{
				"GameID": gameID,
				"ID":     id,
			})
		}
		return nil, err
	}

	return &offer, nil
}

// GetPlayerSeenOffers returns all the offers player has seen, but only brings the
// ID, offer_template_id and claimed_at fields
func GetPlayerSeenOffers(
	db runner.Connection,
	gameID string,
	playerID string,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	return getPlayerOffers(db, gameID, playerID, mr, playerSeenOffersScope)
}

// GetPlayerUnseenOffers returns all the offers player has not seen, but only brings the
// ID, offer_template_id and claimed_at fields
func GetPlayerUnseenOffers(
	db runner.Connection,
	gameID string,
	playerID string,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	return getPlayerOffers(db, gameID, playerID, mr, playerUnseenOffersScope)
}

func getPlayerOffers(
	db runner.Connection,
	gameID string,
	playerID string,
	mr *MixedMetricsReporter,
	scope string,
) ([]*Offer, error) {
	params := []interface{}{
		gameID,
		playerID,
	}

	var offers []*Offer
	err := mr.WithDatastoreSegment("offers", "select seen offers", func() error {
		return db.
			Select("o.id, o.offer_template_id, o.claimed_at").
			From("offers o").
			Scope(scope, params...).
			QueryStructs(&offers)
	})

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, errors.NewModelNotFoundError("Offer", map[string]interface{}{
				"GameID":   gameID,
				"PlayerID": playerID,
			})
		}
		return nil, err
	}

	return offers, nil
}

//UpsertOffer updates a offer with new meta or insert with the new UUID
func UpsertOffer(db runner.Connection, offer *Offer, t time.Time, mr *MixedMetricsReporter) error {
	offer.CreatedAt = dat.NullTimeFrom(t)
	err := mr.WithDatastoreSegment("offers", "upsert", func() error {
		return db.
			Upsert("offers").
			Columns("game_id", "offer_template_id", "player_id", "created_at").
			Record(offer).
			Where("id=$1", offer.ID).
			Returning("id", "created_at", "updated_at", "claimed_at").
			QueryStruct(offer)
	})

	if pqErr, ok := IsForeignKeyViolationError(err); ok {
		return errors.NewInvalidModelError("Offer", pqErr.Message)
	}

	if err != nil {
		return err
	}

	return nil
}

//ClaimOffer sets claimed_at to time
func ClaimOffer(db runner.Connection, offerID, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) (dat.JSON, bool, error) {
  
	var offer Offer

  err := mr.WithDatastoreSegment("offers", "select by id", func() error {
    return db.
      Select("id, claimed_at, offer_template_id").
      From("offers").
      Where("id=$1 AND player_id=$2 AND game_id=$3", offerID, playerID, gameID).
      QueryStruct(&offer)
  })

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, false, errors.NewModelNotFoundError("Offer", map[string]interface{}{
        "ID": offerID,
				"GameID":   gameID,
				"PlayerID": playerID,
			})
		}
		return nil, false, err
	}

  ot, err := GetOfferTemplateByID(db, offer.OfferTemplateID, mr)
  if err != nil {
    return nil, false, err
  }

  if offer.ClaimedAt.Valid {
    return ot.Contents, true, nil
  }

	err = mr.WithDatastoreSegment("offers", "update", func() error {
		return db.
			Update("offers").
			Set("claimed_at", t).
			Where("id=$1", offer.ID).
			Returning("claimed_at").
			QueryStruct(&offer)
	})

	return ot.Contents, false, err
}

//UpdateLastSeenAt updates last seen timestamp of an offer
func UpdateLastSeenAt(db runner.Connection, offerID, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) error {
  var offer Offer
  err := mr.WithDatastoreSegment("offers", "update", func() error {
		return db.
			Update("offers").
			Set("last_seen_at", t).
      Where("id=$1 AND player_id=$2 AND game_id=$3", offerID, playerID, gameID).
			Returning("id").
      QueryStruct(&offer)
	})

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return errors.NewModelNotFoundError("Offer", map[string]interface{}{
        "ID": offerID,
				"GameID":   gameID,
				"PlayerID": playerID,
			})
		}
		return err
  }

  return nil
}
