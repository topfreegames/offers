// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"time"

	"github.com/pmylund/go-cache"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//Offer contains the parameters of an offer
type Offer struct {
	ID        string    `db:"id" json:"id" valid:"uuidv4"`
	GameID    string    `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	Name      string    `db:"name" json:"name" valid:"ascii,stringlength(1|255),required"`
	Period    dat.JSON  `db:"period" json:"period" valid:"RequiredJSONObject"`
	Frequency dat.JSON  `db:"frequency" json:"frequency" valid:"RequiredJSONObject"`
	Trigger   dat.JSON  `db:"trigger" json:"trigger" valid:"RequiredJSONObject"`
	Placement string    `db:"placement" json:"placement" valid:"ascii,stringlength(1|255),required"`
	Metadata  dat.JSON  `db:"metadata" json:"metadata" valid:"JSONObject"`
	ProductID string    `db:"product_id" json:"productId" valid:"ascii,stringlength(1|255),required"`
	Contents  dat.JSON  `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	Enabled   bool      `db:"enabled" json:"enabled" valid:"matches(^(true|false)$),optional"`
	Version   int       `db:"version" json:"version" valid:"int,optional"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" valid:"optional"`
}

const enabledOffers = `
    WHERE
		offers.game_id = $1
		AND offers.enabled = true
`

//GetOfferByID returns Offer by ID
func GetOfferByID(db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		return db.
			Select("id, frequency, period, version").
			From("offers").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offer)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":     id,
		"GameID": gameID,
	}, err)
	return &offer, err
}

//GetEnabledOffers returns all the enabled offers
func GetEnabledOffers(db runner.Connection, gameID string, offersCache *cache.Cache, expireDuration time.Duration, mr *MixedMetricsReporter) ([]*Offer, error) {
	var offers []*Offer
	var err error

	enabledOffersKey := GetEnabledOffersKey(gameID)
	offersInterface, found := offersCache.Get(enabledOffersKey)

	if found {
		//fmt.Println("Offers Cache Hit")
		offers = offersInterface.([]*Offer)
		return offers, err
	}

	//fmt.Println("Offers Cache Miss")
	err = mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		return db.
			Select(`
		id, game_id, name, period, frequency,
		trigger, placement, metadata,
		product_id, contents, version
		`).
			From("offers").
			Scope(enabledOffers, gameID).
			QueryStructs(&offers)
	})
	err = HandleNotFoundError("Offer", map[string]interface{}{"enabled": true}, err)

	if err == nil {
		offersCache.Set(enabledOffersKey, offers, expireDuration)
	}

	return offers, err
}

//ListOffers returns all the offer templates for a given game
func ListOffers(db runner.Connection, gameID string, mr *MixedMetricsReporter) ([]*Offer, error) {
	var offers []*Offer
	err := mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		return db.
			Select("*").
			From("offers").
			Where("game_id = $1", gameID).
			QueryStructs(&offers)
	})
	return offers, err
}

// InsertOffer inserts a new offer template into DB
func InsertOffer(db runner.Connection, offer *Offer, mr *MixedMetricsReporter) (*Offer, error) {
	if offer.Metadata == nil {
		offer.Metadata = dat.JSON([]byte(`{}`))
	}

	err := mr.WithDatastoreSegment("offers", SegmentInsert, func() error {
		return db.
			InsertInto("offers").
			Columns("game_id", "name", "period", "frequency", "trigger", "placement", "metadata", "product_id", "contents").
			Record(offer).
			Returning("id, enabled, version").
			QueryStruct(offer)
	})

	foreignKeyErr := HandleForeignKeyViolationError("Offer", err)
	return offer, foreignKeyErr
}

// UpdateOffer updates a given offer
func UpdateOffer(db runner.Connection, offer *Offer, mr *MixedMetricsReporter) (*Offer, error) {
	prevOffer, err := GetOfferByID(db, offer.GameID, offer.ID, mr)
	if err != nil {
		return nil, err
	}
	if offer.Metadata == nil {
		offer.Metadata = dat.JSON([]byte(`{}`))
	}
	offersMap := map[string]interface{}{
		"name":       offer.Name,
		"period":     offer.Period,
		"frequency":  offer.Frequency,
		"trigger":    offer.Trigger,
		"placement":  offer.Placement,
		"metadata":   offer.Metadata,
		"product_id": offer.ProductID,
		"contents":   offer.Contents,
		"version":    prevOffer.Version + 1,
	}
	offer.Version = prevOffer.Version + 1
	err = mr.WithDatastoreSegment("offers", SegmentUpdate, func() error {
		return db.
			Update("offers").
			SetMap(offersMap).
			Where("id = $1 AND game_id = $2", offer.ID, offer.GameID).
			Returning("id, version").
			QueryStruct(offer)
	})
	return offer, err
}

//SetEnabledOffer can enable or disable an offer template
func SetEnabledOffer(db runner.Connection, gameID, id string, enabled bool, mr *MixedMetricsReporter) error {
	var offerTemplate Offer
	err := mr.WithDatastoreSegment("offers", SegmentUpdate, func() error {
		return db.
			Update("offers").
			Set("enabled", enabled).
			Where("id=$1 AND game_id=$2", id, gameID).
			Returning("id").
			QueryStruct(&offerTemplate)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":     id,
		"GameID": gameID,
	}, err)

	return err
}
