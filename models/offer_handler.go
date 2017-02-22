// offers api // https://github.com/topfreegames/offers
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
  "encoding/json"
  uuid "github.com/satori/go.uuid"
)

//Offer represents a tenant in offers API
type Offer struct {
	ID              string `db:"id" valid:"uuidv4,required"`
	GameID          string `db:"game_id" valid:"matches(^[a-z0-9]+(\\-[a-z0-9]+)*$),stringlength(1|255),required"`
	OfferTemplateID string `db:"offer_template_id" valid:"uuidv4,required"`
	PlayerID        string `db:"player_id" valid:"ascii,stringlength(1|1000),required"`
  SeenCounter     int    `db:"seen_counter" valid:""`
  BoughtCounter   int    `db:"bought_counter" valid:""`

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

//Frequency is how many times per unit of time that the offers is shown to player
type Frequency struct {
  Every string
  Max int
}

//Period is how many times per unit of time that the offer can be bought by player
type Period struct {
  Every string
  Max int
}

//GetOfferByID returns a offer by it's pk
func GetOfferByID(db runner.Connection, gameID string, id string, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", "select by id", func() error {
		return db.
			Select("id, game_id, offer_template_id, player_id, created_at, updated_at, claimed_at, last_seen_at, seen_counter").
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

//InsertOffer inserts an offer with the new UUID
func InsertOffer(db runner.Connection, offer *Offer, t time.Time, mr *MixedMetricsReporter) error {
	offer.CreatedAt = dat.NullTimeFrom(t)
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
            RETURNING id`
  err := mr.WithDatastoreSegment("offers", "select by id", func() error {
    return db.SQL(query, t, offerID, playerID, gameID).QueryStruct(&offer)
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

//GetAvailableOffers returns the offers that match the criteria of enabled offer templates
func GetAvailableOffers(db runner.Connection, playerID, gameID string, t time.Time, mr *MixedMetricsReporter) (map[string]*OfferTemplate, error) {
  eot, err := GetEnabledOfferTemplates(db, gameID, mr)
  if err != nil {
    return nil, err
  }
  if len(eot) == 0 {
    return map[string]*OfferTemplate{}, nil
  }

  var trigger TimeTrigger
  filteredOts, err := filterTemplatesByTrigger(&trigger, eot, t)
  if err != nil {
    return nil, err
  }
  if len(filteredOts) == 0 {
    return map[string]*OfferTemplate{}, nil
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
    return map[string]*OfferTemplate{}, nil
  }

  var offerTemplatesByPlacement map[string]*OfferTemplate
  for _, ot := range filteredOts {
    if _, otInMap := offerTemplatesByPlacement[ot.Placement]; !otInMap {
      offerTemplatesByPlacement[ot.Placement] = ot
      o := &Offer{
        ID: uuid.NewV4().String(),
        GameID: ot.GameID,
        OfferTemplateID: ot.ID,
        PlayerID: playerID,
      }
      InsertOffer(db, o, t, mr)
    }
  }

  return offerTemplatesByPlacement, nil
}

func filterTemplatesByTrigger(trigger Trigger, ots []*OfferTemplate, t time.Time) ([]*OfferTemplate, error) {
  var filteredOts []*OfferTemplate
  for _, ot := range ots {
    var times Times
    err := json.Unmarshal(ot.Trigger, &times)
    if err != nil {
      return nil, err
    }

    if trigger.IsTriggered(&times, &t) {
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
      Select("id, offer_template_id").
      From("offers").
      Where("player_id=$1 AND game_id=$2 AND offer_template_id IN ($3)", playerID, gameID, offerTemplateIDs).
      QueryStructs(&offers)
  })

  return offers, err
}

func filterTemplatesByFrequencyAndPeriod(offers []*Offer, ots []*OfferTemplate, t time.Time) ([]*OfferTemplate, error) {
  var offersByOfferTemplateID map[string][]*Offer
  var filteredOts []*OfferTemplate

  for _, offer := range offers {
    if os, ok := offersByOfferTemplateID[offer.OfferTemplateID]; ok {
      offersByOfferTemplateID[offer.OfferTemplateID] = append(os, offer)
    } else {
      offersByOfferTemplateID[offer.OfferTemplateID] = []*Offer{offer}
    }
  }

  for _, offerTemplate := range ots {
    if os, ok := offersByOfferTemplateID[offerTemplate.ID]; ok {
      var f Frequency
      err := json.Unmarshal(offerTemplate.Frequency, &f)
      if err != nil {
        return nil, err
      }
      var p Period
      err = json.Unmarshal(offerTemplate.Period, &p)
      if err != nil {
        return nil, err
      }
      for _, offer := range os {
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
          if offer.LastSeenAt.Time.Add(duration).After(t) {
            continue
          }
        }
        filteredOts = append(filteredOts, offerTemplate)
      }
    } else {
      filteredOts = append(filteredOts, offerTemplate)
    }
  }

  return filteredOts, nil
}
