// offers api // https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models

import (
	"encoding/json"
	"time"

	"github.com/topfreegames/offers/errors"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//Offer represents a tenant in offers API
type Offer struct {
	ID              string `db:"id" json:"id" valid:"uuidv4,required"`
	GameID          string `db:"game_id" json:"gameId" valid:"matches(^[^-][a-z0-9-]*$),stringlength(1|255),required"`
	OfferTemplateID string `db:"offer_template_id" json:"offerTemplateId" valid:"uuidv4,required"`
	PlayerID        string `db:"player_id" json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	SeenCounter     int    `db:"seen_counter" json:"seenCounter" valid:""`
	BoughtCounter   int    `db:"bought_counter" json:"boughtCounter" valid:""`

	CreatedAt  dat.NullTime `db:"created_at" json:"createdAt" valid:""`
	UpdatedAt  dat.NullTime `db:"updated_at" json:"updatedAt" valid:""`
	ClaimedAt  dat.NullTime `db:"claimed_at" json:"claimedAt" valid:""`
	LastSeenAt dat.NullTime `db:"last_seen_at" json:"lastSeenAt" valid:""`
}

//OfferToUpdate has required fields for claiming an offer
type OfferToUpdate struct {
	ID       string `db:"id" valid:"uuidv4,required"`
	GameID   string `db:"game_id" valid:"matches(^[^-][a-z0-9-]*$),stringlength(1|255),required"`
	PlayerID string `db:"player_id" valid:"ascii,stringlength(1|1000),required"`
}

//OfferToReturn has the fields for the returned offer
type OfferToReturn struct {
	ID                   string   `json:"id"`
	ProductID            string   `json:"productId"`
	Contents             dat.JSON `json:"contents"`
	Metadata             dat.JSON `json:"metadata"`
	RemainingPurchases   int      `json:"remainingPurchases,omitempty"`
	RemainingImpressions int      `json:"remainingImpressions,omitempty"`
}

//FrequencyOrPeriod is the struct for basic Frequecy and Period types
type FrequencyOrPeriod struct {
	Every string
	Max   int
}

//GetOfferByID returns a offer by it's pk
func GetOfferByID(db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, game_id, offer_template_id, player_id, created_at, updated_at, claimed_at, last_seen_at, seen_counter, bought_counter").
			From("offers").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offer)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"GameID": gameID,
		"ID":     id,
	}, err)

	return &offer, err
}

//InsertOffer inserts an offer with the new UUID
func InsertOffer(db runner.Connection, offer *Offer, t time.Time, mr *MixedMetricsReporter) error {
	err := mr.WithDatastoreSegment("offers", "insect", func() error {
		return db.
			InsertInto("offers").
			Columns("game_id", "offer_template_id", "player_id").
			Record(offer).
			Returning("id").
			QueryStruct(offer)
	})

	if err != nil {
		if pqErr, ok := IsForeignKeyViolationError(err); ok {
			return errors.NewInvalidModelError("Offer", pqErr.Message)
		}
		return err
	}

	return nil
}

//ClaimOffer sets claimed_at to time
func ClaimOffer(db runner.Connection, offerID, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) (dat.JSON, bool, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, claimed_at, offer_template_id, bought_counter").
			From("offers").
			Where("id=$1 AND player_id=$2 AND game_id=$3", offerID, playerID, gameID).
			QueryStruct(&offer)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":       offerID,
		"GameID":   gameID,
		"PlayerID": playerID,
	}, err)

	if err != nil {
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
			Set("bought_counter", offer.BoughtCounter+1).
			Where("id=$1", offer.ID).
			Returning("claimed_at").
			QueryStruct(&offer)
	})

	return ot.Contents, false, err
}

//UpdateOfferLastSeenAt updates last seen timestamp of an offer
func UpdateOfferLastSeenAt(db runner.Connection, offerID, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) error {
	var offer Offer

	query := `UPDATE offers
            SET
              last_seen_at = $1,
              seen_counter = seen_counter + 1
            WHERE
              id = $2 AND
              player_id = $3 AND
              game_id = $4
            RETURNING id, last_seen_at`
	err := mr.WithDatastoreSegment("offers", "update", func() error {
		return db.SQL(query, t, offerID, playerID, gameID).QueryStruct(&offer)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":       offerID,
		"GameID":   gameID,
		"PlayerID": playerID,
	}, err)

	return err
}

//GetAvailableOffers returns the offers that match the criteria of enabled offer templates
func GetAvailableOffers(db runner.Connection, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) (map[string][]*OfferToReturn, error) {
	eot, err := GetEnabledOfferTemplates(db, gameID, mr)
	if err != nil {
		return nil, err
	}
	if len(eot) == 0 {
		return map[string][]*OfferToReturn{}, nil
	}

	var trigger TimeTrigger
	filteredOts, err := filterTemplatesByTrigger(trigger, eot, t)
	if err != nil {
		return nil, err
	}
	if len(filteredOts) == 0 {
		return map[string][]*OfferToReturn{}, nil
	}

	offerTemplateIDs := make([]string, len(filteredOts))
	for idx, ot := range filteredOts {
		offerTemplateIDs[idx] = ot.ID
	}
	playerOffers, err := getPlayerOffersByOfferTemplateIDs(db, gameID, playerID, offerTemplateIDs, mr)
	if err != nil {
		return nil, err
	}
	filteredOts, err = filterTemplatesByFrequencyAndPeriod(playerOffers, filteredOts, t)
	if err != nil {
		return nil, err
	}
	if len(filteredOts) == 0 {
		return map[string][]*OfferToReturn{}, nil
	}

	playerOffersByOfferTemplateID := map[string]*Offer{}
	for _, o := range playerOffers {
		playerOffersByOfferTemplateID[o.OfferTemplateID] = o
	}
	offerTemplatesByPlacement := make(map[string][]*OfferToReturn)
	for _, ot := range filteredOts {
		offerToReturn := &OfferToReturn{
			ProductID: ot.ProductID,
			Contents:  ot.Contents,
			Metadata:  ot.Metadata,
		}
		var f FrequencyOrPeriod
		var p FrequencyOrPeriod
		json.Unmarshal(ot.Frequency, &f)
		json.Unmarshal(ot.Period, &p)
		if f.Max > 0 {
			offerToReturn.RemainingImpressions = f.Max
		}
		if p.Max > 0 {
			offerToReturn.RemainingPurchases = p.Max
		}
		o := &Offer{
			GameID:          ot.GameID,
			OfferTemplateID: ot.ID,
			PlayerID:        playerID,
		}
		playerOffer, playerHasOffer := playerOffersByOfferTemplateID[ot.ID]
		if playerHasOffer {
			offerToReturn.ID = playerOffer.ID
			if offerToReturn.RemainingImpressions > 0 {
				offerToReturn.RemainingImpressions = offerToReturn.RemainingImpressions - playerOffer.SeenCounter
			}
			if offerToReturn.RemainingPurchases > 0 {
				offerToReturn.RemainingPurchases = offerToReturn.RemainingPurchases - playerOffer.BoughtCounter
			}
		} else {
			err := InsertOffer(db, o, t, mr)
			if err != nil {
				return nil, err
			}
			offerToReturn.ID = o.ID
		}
		if _, otInMap := offerTemplatesByPlacement[ot.Placement]; !otInMap {
			offerTemplatesByPlacement[ot.Placement] = []*OfferToReturn{offerToReturn}
		} else {
			offerTemplatesByPlacement[ot.Placement] = append(offerTemplatesByPlacement[ot.Placement], offerToReturn)
		}
	}

	return offerTemplatesByPlacement, nil
}

func filterTemplatesByTrigger(trigger Trigger, ots []*OfferTemplate, t time.Time) ([]*OfferTemplate, error) {
	var (
		filteredOts []*OfferTemplate
		times       Times
	)
	for _, ot := range ots {
		if err := json.Unmarshal(ot.Trigger, &times); err != nil {
			return nil, err
		}
		if trigger.IsTriggered(times, t) {
			filteredOts = append(filteredOts, ot)
		}
	}
	return filteredOts, nil
}

func getPlayerOffersByOfferTemplateIDs(
	db runner.Connection,
	gameID string,
	playerID string,
	offerTemplateIDs []string,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	var offers []*Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, offer_template_id, game_id, last_seen_at, claimed_at, seen_counter, bought_counter").
			From("offers").
			Where("player_id=$1 AND game_id=$2 AND offer_template_id IN $3", playerID, gameID, offerTemplateIDs).
			QueryStructs(&offers)
	})
	return offers, err
}

func filterTemplatesByFrequencyAndPeriod(offers []*Offer, ots []*OfferTemplate, t time.Time) ([]*OfferTemplate, error) {
	var filteredOts []*OfferTemplate
	offerByOfferTemplateID := make(map[string]*Offer)
	for _, offer := range offers {
		offerByOfferTemplateID[offer.OfferTemplateID] = offer
	}

	for _, offerTemplate := range ots {
		if offer, ok := offerByOfferTemplateID[offerTemplate.ID]; ok {
			var (
				f FrequencyOrPeriod
				p FrequencyOrPeriod
			)
			if err := json.Unmarshal(offerTemplate.Frequency, &f); err != nil {
				return nil, err
			}
			if err := json.Unmarshal(offerTemplate.Period, &p); err != nil {
				return nil, err
			}
			if f.Max != 0 && offer.SeenCounter >= f.Max {
				continue
			}
			if f.Every != "" {
				duration, err := time.ParseDuration(f.Every)
				if err != nil {
					return nil, err
				}
				if offer.LastSeenAt.Time.Add(duration).After(t) {
					continue
				}
			}
			if p.Max != 0 && offer.BoughtCounter >= p.Max {
				continue
			}
			if p.Every != "" {
				duration, err := time.ParseDuration(p.Every)
				if err != nil {
					return nil, err
				}
				if offer.ClaimedAt.Time.Add(duration).After(t) {
					continue
				}
			}
			filteredOts = append(filteredOts, offerTemplate)
		} else {
			filteredOts = append(filteredOts, offerTemplate)
		}
	}

	return filteredOts, nil
}
