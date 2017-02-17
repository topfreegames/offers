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
	"github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
)

//Offer represents a tenant in offers API
type Offer struct {
	ID              uuid.UUID `db:"id" valid:"uuidv4,required"`
	GameID          string    `db:"game_id" valid:"matches(^[a-z0-9]+(\\-[a-z0-9]+)*$),stringlength(1|255),required"`
	OfferTemplateID uuid.UUID `db:"offer_template_id" valid:"uuidv4,required"`
	PlayerID        string    `db:"player_id" valid:"ascii,stringlength(1|1000),required"`

	CreatedAt dat.NullTime `db:"created_at" valid:""`
	UpdatedAt dat.NullTime `db:"updated_at" valid:""`
	ClaimedAt dat.NullTime `db:"claimed_at" valid:""`
}

const playerSeenOffersScope = `
	WHERE
		o.game_id = $1
		o.player_id = $2
	AND o.offer_template_id in $3
`

//GetOfferByID returns a offer by it's pk
func GetOfferByID(db runner.Connection, gameID string, id uuid.UUID, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, game_id, offer_template_id, player_id, created_at, updated_at, claimed_at").
			From("offers").
			Where("id = $1 AND game_id=$2", id, gameID).
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
	offerTemplates []*OfferTemplate,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	if len(offerTemplates) == 0 {
		return []*Offer{}, nil
	}

	offerTemplateIDs := make([]uuid.UUID, len(offerTemplates))
	for i, offerTemplate := range offerTemplates {
		offerTemplateIDs[i] = offerTemplate.ID
	}

	var offers []*Offer
	err := mr.WithDatastoreSegment("offers", "select seen offers", func() error {
		return db.
			Select("id, offer_template_id, claimed_at").
			From("offers").
			Scope(playerSeenOffersScope, gameID, playerID, offerTemplateIDs).
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
func UpsertOffer(db runner.Connection, offer *Offer, mr *MixedMetricsReporter) error {
	err := mr.WithDatastoreSegment("offers", "upsert", func() error {
		return db.
			Upsert("offers").
			Columns("game_id", "offer_template_id", "player_id").
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
