// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pmylund/go-cache"
	edat "github.com/topfreegames/extensions/dat"
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
	ProductID string    `db:"product_id" json:"productId,omitempty" valid:"ascii,stringlength(1|255)"`
	Contents  dat.JSON  `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	Enabled   bool      `db:"enabled" json:"enabled" valid:"matches(^(true|false)$),optional"`
	Version   int       `db:"version" json:"version" valid:"int,optional"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" valid:"optional"`
	Filters   dat.JSON  `db:"filters" json:"filters" valid:"FilterJSONObject"`
	Cost      dat.JSON  `db:"cost" json:"cost,omitempty" valid:"JSONObject"`
}

const enabledOffers = `
    WHERE
		offers.game_id = $1
		AND offers.enabled = true
		AND (trigger->>'to')::int >= $2
		AND (trigger->>'from')::int <= $2
`

var isValidString = regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`).MatchString

//ValidateString validates the string contains only valid characters for filters
func ValidateString(s string) bool {
	return isValidString(s)
}

func buildInefficientScope(enabledOffers string, filterAttrs map[string]string) string {
	subQueries := []string{enabledOffers}
	for k, v := range filterAttrs {
		// TODO: Possible SQL injection
		if !ValidateString(k) || !ValidateString(v) {
			subQueries = []string{enabledOffers}
			break
		}
		rawSubQuery := `
		AND (
			NOT (filters ? '%[1]s') OR
			filters @> '{"%[1]s": {"eq": "%[2]s"}}' OR
			filters @> '{"%[1]s": {"neq": "%[2]s"}}'
		)`
		subQuery := fmt.Sprintf(rawSubQuery, k, v)
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			rawSubQuery = `
			AND (
				NOT (filters ? '%[1]s') OR
				filters @> '{"%[1]s": {"eq": "%[2]s"}}' OR
				filters @> '{"%[1]s": {"neq": "%[2]s"}}' OR
				(((filters::json#>>'{"%[1]s",geq}') IS NOT NULL OR (filters::json#>>'{"%[1]s",geq}') IS NOT NULL) AND
				((filters::json#>>'{"%[1]s",geq}') IS NULL OR %[3]f >= (filters::json#>>'{"%[1]s",geq}')::float) AND
				((filters::json#>>'{"%[1]s",lt}') IS NULL OR %[3]f < (filters::json#>>'{"%[1]s",lt}')::float))
			)`
			subQuery = fmt.Sprintf(rawSubQuery, k, v, f)
		}
		subQueries = append(subQueries, subQuery)
	}
	query := strings.Join(subQueries, " ")
	return query
}

func buildEfficientScope(enabledOffers string, filterAttrs map[string]string) string {
	subQueries := []string{enabledOffers}
	for k, v := range filterAttrs {
		// TODO: Possible SQL injection
		if !ValidateString(k) || !ValidateString(v) {
			subQueries = []string{enabledOffers}
			break
		}
		rawSubQuery := `
		AND (
			filters @> '{"%[1]s": {"eq": "%[2]s"}}' OR
			filters @> '{"%[1]s": {"neq": "%[2]s"}}'
		)`
		subQuery := fmt.Sprintf(rawSubQuery, k, v)
		subQueries = append(subQueries, subQuery)
	}
	query := strings.Join(subQueries, " ")
	return query
}

//GetOfferByID returns Offer by ID
func GetOfferByID(ctx context.Context, db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*Offer, error) {
	var offer Offer
	err := mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		builder := db.Select("id, frequency, period, version, enabled")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offers").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offer)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":     id,
		"GameID": gameID,
	}, err)
	return &offer, err
}

//GetEnabledOffers returns all the enabled offers and matching offers
func GetEnabledOffers(ctx context.Context, db runner.Connection, gameID string, offersCache *cache.Cache, expireDuration time.Duration, currentTime time.Time, filterAttrs map[string]string, allowInefficientQueries bool, mr *MixedMetricsReporter) ([]*Offer, error) {
	var offers []*Offer
	var err error

	enabledOffersKey := GetEnabledOffersKey(gameID)
	if len(filterAttrs) == 0 {
		offersInterface, found := offersCache.Get(enabledOffersKey)
		if found {
			offers = offersInterface.([]*Offer)
			return offers, err
		}
	}

	var scope string
	if allowInefficientQueries {
		scope = buildInefficientScope(enabledOffers, filterAttrs)
	} else {
		scope = buildEfficientScope(enabledOffers, filterAttrs)
	}

	err = mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		// TODO: Add a configurable limit to this query

		builder := db.
			Select(`
		id, game_id, name, period, frequency,
		trigger, placement, metadata,
		product_id, contents, version, cost
		`)
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offers").
			Scope(scope, gameID, currentTime.Unix()).
			QueryStructs(&offers)
	})
	err = HandleNotFoundError("Offer", map[string]interface{}{"enabled": true}, err)

	if err == nil && len(filterAttrs) == 0 {
		offersCache.Set(enabledOffersKey, offers, expireDuration)
	}

	return offers, err
}

//ListOffers returns all the offer templates for a given game
//return the number of pages using the number of offers and given the limit for each page
func ListOffers(
	ctx context.Context,
	db runner.Connection,
	gameID string,
	limit, offset uint64,
	mr *MixedMetricsReporter,
) ([]*Offer, int, error) {
	offers := []*Offer{}
	var numberOffers int
	err := mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
		builder := db.Select("COUNT(*)")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offers").
			Where("game_id = $1", gameID).
			QueryScalar(&numberOffers)
	})
	if err != nil {
		return offers, 0, err
	}

	var pages int
	if limit != 0 {
		pages = numberOffers / int(limit)
		if numberOffers%int(limit) != 0 {
			pages = pages + 1
		}

		start := offset * limit
		err = mr.WithDatastoreSegment("offers", SegmentSelect, func() error {
			builder := db.Select("*")
			builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
			return builder.From("offers").
				Where("game_id = $1", gameID).
				OrderBy("created_at").
				Limit(limit).
				Offset(start).
				QueryStructs(&offers)
		})
		if err != nil {
			return offers, 0, err
		}
	}

	return offers, pages, nil
}

// InsertOffer inserts a new offer template into DB
func InsertOffer(ctx context.Context, db runner.Connection, offer *Offer, offersCache *cache.Cache, mr *MixedMetricsReporter) (*Offer, error) {
	if offer.Metadata == nil {
		offer.Metadata = dat.JSON([]byte(`{}`))
	}
	if offer.Filters == nil {
		offer.Filters = dat.JSON([]byte(`{}`))
	}
	if offer.Cost == nil {
		offer.Cost = dat.JSON([]byte(`{}`))
	}
	err := mr.WithDatastoreSegment("offers", SegmentInsert, func() error {
		builder := db.InsertInto("offers")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.Columns("game_id", "name", "period", "frequency", "trigger", "placement", "metadata", "product_id", "contents", "filters", "cost").
			Record(offer).
			Returning("id, enabled, version").
			QueryStruct(offer)
	})

	foreignKeyErr := HandleForeignKeyViolationError("Offer", err)
	if err == nil {
		enabledOffersKey := GetEnabledOffersKey(offer.GameID)
		offersCache.Delete(enabledOffersKey)
	}
	return offer, foreignKeyErr
}

// UpdateOffer updates a given offer
func UpdateOffer(ctx context.Context, db runner.Connection, offer *Offer, offersCache *cache.Cache, mr *MixedMetricsReporter) (*Offer, error) {
	prevOffer, err := GetOfferByID(ctx, db, offer.GameID, offer.ID, mr)
	if err != nil {
		return nil, err
	}
	if offer.Metadata == nil {
		offer.Metadata = dat.JSON([]byte(`{}`))
	}
	if offer.Filters == nil {
		offer.Filters = dat.JSON([]byte(`{}`))
	}
	if offer.Cost == nil {
		offer.Cost = dat.JSON([]byte(`{}`))
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
		"filters":    offer.Filters,
		"cost":       offer.Cost,
		"version":    prevOffer.Version + 1,
	}
	offer.Version = prevOffer.Version + 1
	err = mr.WithDatastoreSegment("offers", SegmentUpdate, func() error {
		builder := db.Update("offers")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.SetMap(offersMap).
			Where("id = $1 AND game_id = $2", offer.ID, offer.GameID).
			Returning("id, version").
			QueryStruct(offer)
	})
	if err == nil {
		enabledOffersKey := GetEnabledOffersKey(offer.GameID)
		offersCache.Delete(enabledOffersKey)
	}
	return offer, err
}

//SetEnabledOffer can enable or disable an offer template
func SetEnabledOffer(ctx context.Context, db runner.Connection, gameID, id string, enabled bool, offersCache *cache.Cache, mr *MixedMetricsReporter) error {
	var offerTemplate Offer
	err := mr.WithDatastoreSegment("offers", SegmentUpdate, func() error {
		builder := db.Update("offers")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.Set("enabled", enabled).
			Where("id=$1 AND game_id=$2", id, gameID).
			Returning("id").
			QueryStruct(&offerTemplate)
	})

	err = HandleNotFoundError("Offer", map[string]interface{}{
		"ID":     id,
		"GameID": gameID,
	}, err)
	if err == nil {
		enabledOffersKey := GetEnabledOffersKey(gameID)
		offersCache.Delete(enabledOffersKey)
	}
	return err
}
