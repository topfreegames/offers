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

	//TODO: Validate dates
	CreatedAt dat.NullTime `db:"created_at" valid:""`
	UpdatedAt dat.NullTime `db:"updated_at" valid:""`
	ClaimedAt dat.NullTime `db:"claimed_at" valid:""`
}

//GetOfferByID returns a offer by it's pk
func GetOfferByID(db runner.Connection, id uuid.UUID, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("*").
			From("offers").
			Where("id = $1", id).
			QueryStruct(&offer)
	})

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, errors.NewModelNotFoundError("Offer", map[string]interface{}{
				"ID": id,
			})
		}
		return nil, err
	}

	return &offer, nil
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
