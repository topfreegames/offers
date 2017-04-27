// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pmylund/go-cache"
	"github.com/topfreegames/offers/util"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
	redis "gopkg.in/redis.v5"
)

//OfferInstance represents a tenant in offers API it cannot be updated, only inserted
type OfferInstance struct {
	ID           string       `db:"id" json:"id" valid:"uuidv4,required"`
	GameID       string       `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	PlayerID     string       `db:"player_id" json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	OfferID      string       `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	OfferVersion int          `db:"offer_version" json:"offerVersion" valid:"int,required"`
	Contents     dat.JSON     `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	ProductID    string       `db:"product_id" json:"productId" valid:"ascii,stringlength(1|255),required"`
	CreatedAt    dat.NullTime `db:"created_at" json:"createdAt" valid:""`
}

//OfferInstanceOffer is a join of OfferInstance with
type OfferInstanceOffer struct {
	ID       string   `db:"id" json:"id" valid:"uuidv4,required"`
	GameID   string   `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	OfferID  string   `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	Contents dat.JSON `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	Enabled  bool     `db:"enabled" json:"enabled"`
}

//OfferToReturn has the fields for the returned offer
type OfferToReturn struct {
	ID        string   `json:"id"`
	ProductID string   `json:"productId"`
	Contents  dat.JSON `json:"contents"`
	Metadata  dat.JSON `json:"metadata"`
	ExpireAt  int64    `json:"expireAt"`
}

//FrequencyOrPeriod is the struct for basic Frequency and Period types
type FrequencyOrPeriod struct {
	Every string
	Max   int
}

//GetOfferInstanceAndOfferEnabled returns a offer by its pk
func GetOfferInstanceAndOfferEnabled(db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferInstanceOffer, error) {
	var offerInstance OfferInstanceOffer
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		return db.
			Select("oi.id, oi.offer_id, oi.contents, o.enabled").
			From("offer_instances oi JOIN offers o ON (oi.offer_id=o.id)").
			Where("oi.id=$1 AND oi.game_id=$2", id, gameID).
			QueryStruct(&offerInstance)
	})

	err = HandleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID": gameID,
		"ID":     id,
	}, err)

	return &offerInstance, err
}

//GetOfferInstanceByID returns a offer by its pk
func GetOfferInstanceByID(db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferInstance, error) {
	var offerInstance OfferInstance
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		return db.
			Select("id, offer_id, contents").
			From("offer_instances").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offerInstance)
	})

	err = HandleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID": gameID,
		"ID":     id,
	}, err)

	return &offerInstance, err
}

func getClaimedOfferNextAt(
	db runner.Connection,
	gameID, offerID string,
	claimCounter int,
	t time.Time,
	mr *MixedMetricsReporter,
) (int64, error) {
	offer, err := GetOfferByID(db, gameID, offerID, mr)
	if err != nil {
		return 0, err
	}
	if !offer.Enabled {
		return 0, nil
	}

	var p FrequencyOrPeriod
	var f FrequencyOrPeriod
	json.Unmarshal(offer.Period, &p)
	json.Unmarshal(offer.Frequency, &f)

	if p.Max != 0 && claimCounter >= p.Max {
		return 0, nil
	}

	if p.Every == "" && f.Every == "" {
		return t.Unix(), nil
	}

	var duration time.Duration
	var nextAt int64
	if p.Every != "" {
		duration, _ = time.ParseDuration(p.Every)
		nextAt = t.Add(duration).Unix()
	}

	if f.Every != "" {
		duration, _ = time.ParseDuration(f.Every)
		if t.Add(duration).Unix() > nextAt {
			nextAt = t.Add(duration).Unix()
		}
	}
	return nextAt, nil
}

//ClaimOffer claims the offer
func ClaimOffer(
	db runner.Connection,
	redisClient *util.RedisClient,
	gameID, offerInstanceID, playerID, productID, transactionID string,
	timestamp int64,
	t time.Time,
	mr *MixedMetricsReporter,
) (dat.JSON, bool, int64, error) {

	// If an offer instance id is sent
	var offerInstance *OfferInstance
	var err error
	var isReplay bool
	var claimCount, claimTimestamp int64
	var nextAt int64

	if offerInstanceID != "" {
		offerInstance, err = GetOfferInstanceByID(db, gameID, offerInstanceID, mr)
		if err != nil {
			return nil, false, 0, err
		}
	} else {
		offerInstance, err = GetLastOfferInstanceByPlayerIDAndProductID(db, gameID, playerID, productID, timestamp, mr)
		if err != nil {
			return nil, false, 0, err
		}
	}

	transactionsKey := GetTransactionsKey(playerID, gameID)
	claimCounterKey := GetClaimCounterKey(playerID, offerInstance.OfferID)
	claimTimestampKey := GetClaimTimestampKey(playerID, offerInstance.OfferID)

	err = mr.WithRedisSegment(SegmentSIsMember, func() error {
		isReplay, err = redisClient.Client.SIsMember(transactionsKey, transactionID).Result()
		return err
	})
	if err != nil {
		return nil, false, 0, err
	}

	if isReplay {
		pipe := redisClient.Client.TxPipeline()

		claimCountGetOp := pipe.Get(claimCounterKey)
		claimTimestampGetOp := pipe.Get(claimTimestampKey)

		err = mr.WithRedisSegment(SegmentGet, func() error {
			_, err = pipe.Exec()
			return err
		})
		if err != nil {
			return nil, false, 0, err
		}

		claimCount, err = claimCountGetOp.Int64()
		if err != nil {
			return nil, false, 0, err
		}

		claimTimestamp, err = claimTimestampGetOp.Int64()
		if err != nil {
			return nil, false, 0, err
		}

		nextAt, err = getClaimedOfferNextAt(db, gameID, offerInstance.OfferID, int(claimCount), time.Unix(claimTimestamp, 0), mr)
		if err != nil {
			return nil, false, 0, err
		}
		return offerInstance.Contents, true, nextAt, nil
	}

	pipe := redisClient.Client.TxPipeline()

	claimCountIncrOp := pipe.Incr(claimCounterKey)
	pipe.Set(claimTimestampKey, timestamp, 0)
	pipe.SAdd(transactionsKey, transactionID)

	err = mr.WithRedisSegment(SegmentRedis, func() error {
		_, err = pipe.Exec()
		return err
	})
	if err != nil {
		return nil, false, 0, err
	}

	claimCount, err = claimCountIncrOp.Result()
	if err != nil {
		return nil, false, 0, err
	}

	nextAt, err = getClaimedOfferNextAt(db, gameID, offerInstance.OfferID, int(claimCount), time.Unix(timestamp, 0), mr)
	if err != nil {
		return nil, false, 0, err
	}
	return offerInstance.Contents, false, nextAt, nil
}

//GetLastOfferInstanceByPlayerIDAndProductID returns a offer by gameId, playerId and productId
func GetLastOfferInstanceByPlayerIDAndProductID(db runner.Connection, gameID, playerID, productID string, timestamp int64, mr *MixedMetricsReporter) (*OfferInstance, error) {
	var offerInstance OfferInstance
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		return db.SQL("SELECT id, offer_id, contents "+
			"FROM offer_instances "+
			"WHERE game_id=$1 AND player_id=$2 AND product_id=$3 AND created_at < to_timestamp($4) "+
			"ORDER BY created_at DESC FETCH FIRST 1 ROW ONLY", gameID, playerID, productID, timestamp).
			QueryStruct(&offerInstance)
	})

	err = HandleNotFoundError("offerInstance", map[string]interface{}{
		"GameID":    gameID,
		"PlayerID":  playerID,
		"ProductID": productID,
	}, err)

	return &offerInstance, err
}

func getViewedOfferNextAt(
	db runner.Connection,
	gameID, offerID string,
	viewCounter int,
	t time.Time,
	mr *MixedMetricsReporter,
) (int64, error) {
	offer, err := GetOfferByID(db, gameID, offerID, mr)
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

//ViewOffer views the offer
func ViewOffer(
	db runner.Connection,
	redisClient *util.RedisClient,
	gameID, offerInstanceID, playerID, impressionID string,
	t time.Time,
	mr *MixedMetricsReporter,
) (bool, int64, error) {

	offerInstance, err := GetOfferInstanceAndOfferEnabled(db, gameID, offerInstanceID, mr)
	if err != nil {
		return false, 0, err
	}

	// Offer is disabled
	if !offerInstance.Enabled {
		return false, 0, nil
	}

	impressionsKey := GetImpressionsKey(playerID, gameID)
	viewCounterKey := GetViewCounterKey(playerID, offerInstance.OfferID)
	viewTimestampKey := GetViewTimestampKey(playerID, offerInstance.OfferID)

	var isReplay bool
	var viewCount int64
	var nextAt int64
	err = mr.WithRedisSegment(SegmentSIsMember, func() error {
		isReplay, err = redisClient.Client.SIsMember(impressionsKey, impressionID).Result()
		return err
	})
	if err != nil {
		return false, 0, err
	}

	if isReplay {
		err = mr.WithRedisSegment(SegmentGet, func() error {
			viewCount, err = redisClient.Client.Get(viewCounterKey).Int64()
			return err
		})
		if err != nil {
			return false, 0, err
		}

		nextAt, err = getViewedOfferNextAt(db, gameID, offerInstance.OfferID, int(viewCount), t, mr)
		if err != nil {
			return false, 0, err
		}
		return true, nextAt, nil
	}

	pipe := redisClient.Client.TxPipeline()

	viewCounterIncrOp := pipe.Incr(viewCounterKey)
	pipe.Set(viewTimestampKey, t.Unix(), 0)
	pipe.SAdd(impressionsKey, impressionID)

	err = mr.WithRedisSegment(SegmentRedis, func() error {
		_, err = pipe.Exec()
		return err
	})
	if err != nil {
		return false, 0, err
	}

	viewCount, err = viewCounterIncrOp.Result()
	if err != nil {
		return false, 0, err
	}

	nextAt, err = getViewedOfferNextAt(db, gameID, offerInstance.OfferID, int(viewCount), t, mr)
	if err != nil {
		return false, 0, err
	}
	return false, nextAt, nil
}

//GetAvailableOffers returns the offers that match the criteria of enabled offer templates
func GetAvailableOffers(
	db runner.Connection,
	redisClient *util.RedisClient,
	offersCache *cache.Cache,
	gameID, playerID string,
	t time.Time,
	expireDuration time.Duration,
	filterAttrs map[string]string,
	mr *MixedMetricsReporter,
) (map[string][]*OfferToReturn, error) {
	offersByPlacement := make(map[string][]*OfferToReturn)

	enabledOffers, err := GetEnabledOffers(db, gameID, offersCache, expireDuration, filterAttrs, mr)

	if err != nil {
		return nil, err
	}
	if len(enabledOffers) == 0 {
		return offersByPlacement, nil
	}

	var trigger TimeTrigger
	filteredOffers, err := filterTemplatesByTrigger(trigger, enabledOffers, t)
	if err != nil {
		return nil, err
	}
	if len(filteredOffers) == 0 {
		return offersByPlacement, nil
	}

	filteredOffers, err = filterOffersByFrequencyAndPeriod(redisClient, playerID, filteredOffers, t, mr)
	if err != nil {
		return nil, err
	}
	if len(filteredOffers) == 0 {
		return offersByPlacement, nil
	}

	var offerInstances []*OfferInstance
	offers := make(map[string]*Offer)

	for _, offer := range filteredOffers {
		offers[offer.ID] = offer
		offerInstances = append(offerInstances, &OfferInstance{
			GameID:       offer.GameID,
			PlayerID:     playerID,
			OfferID:      offer.ID,
			OfferVersion: offer.Version,
			Contents:     offer.Contents,
			ProductID:    offer.ProductID,
		})
	}

	offerInstances, err = findOrCreateOfferInstance(db, offerInstances, t, mr)

	for _, offerInstance := range offerInstances {
		offer := offers[offerInstance.OfferID]

		var trigger Times
		json.Unmarshal(offer.Trigger, &trigger)
		offerToReturn := &OfferToReturn{
			ID:        offerInstance.ID,
			ProductID: offer.ProductID,
			Contents:  offer.Contents,
			Metadata:  offer.Metadata,
			ExpireAt:  trigger.To,
		}

		if _, offerInMap := offersByPlacement[offer.Placement]; !offerInMap {
			offersByPlacement[offer.Placement] = []*OfferToReturn{offerToReturn}
		} else {
			offersByPlacement[offer.Placement] = append(offersByPlacement[offer.Placement], offerToReturn)
		}
	}

	return offersByPlacement, nil
}

func findOrCreateOfferInstance(
	db runner.Connection,
	offerInstances []*OfferInstance,
	t time.Time,
	mr *MixedMetricsReporter,
) ([]*OfferInstance, error) {
	resOfferInstances := make([]*OfferInstance, 0, len(offerInstances))
	var err error

	whereClause := make([]string, 0, len(offerInstances))
	valueArgs := make([]string, 0, len(offerInstances))

	for _, o := range offerInstances {
		whereClause = append(whereClause,
			fmt.Sprintf("(game_id='%s' AND player_id='%s' AND offer_id='%s' AND offer_version=%d)",
				o.GameID,
				o.PlayerID,
				o.OfferID,
				o.OfferVersion))
		valueArgs = append(valueArgs,
			fmt.Sprintf("('%s', '%s', '%s', '%d', '%s'::jsonb, '%s')", o.GameID, o.PlayerID, o.OfferID, o.OfferVersion, o.Contents, o.ProductID))
	}

	query := fmt.Sprintf(`
	WITH 
		sel AS (SELECT id, offer_id FROM offer_instances WHERE %s),
		ins AS (INSERT INTO offer_instances(game_id, player_id, offer_id, offer_version, contents, product_id)
						VALUES %s
						ON CONFLICT DO NOTHING
						RETURNING id, offer_id)
	SELECT * FROM ins UNION ALL SELECT * FROM sel
	`, strings.Join(whereClause, " OR "), strings.Join(valueArgs, ","))

	err = mr.WithDatastoreSegment("offer_instances", SegmentInsect, func() error {
		err := db.SQL(query).QueryStructs(&resOfferInstances)
		return err
	})

	return resOfferInstances, err
}

func filterTemplatesByTrigger(trigger Trigger, offers []*Offer, t time.Time) ([]*Offer, error) {
	var (
		filteredOffers []*Offer
		times          Times
	)
	for _, ot := range offers {
		if err := json.Unmarshal(ot.Trigger, &times); err != nil {
			return nil, err
		}
		if trigger.IsTriggered(times, t) {
			filteredOffers = append(filteredOffers, ot)
		}
	}
	return filteredOffers, nil
}

func filterOffersByFrequencyAndPeriod(
	redisClient *util.RedisClient,
	playerID string,
	offers []*Offer,
	t time.Time,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	var err error
	var filteredOffers []*Offer
	for _, offer := range offers {
		var (
			f FrequencyOrPeriod
			p FrequencyOrPeriod
		)
		if err = json.Unmarshal(offer.Frequency, &f); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(offer.Period, &p); err != nil {
			return nil, err
		}

		claimCounterKey := GetClaimCounterKey(playerID, offer.ID)
		claimTimestampKey := GetClaimTimestampKey(playerID, offer.ID)
		viewCounterKey := GetViewCounterKey(playerID, offer.ID)
		viewTimestampKey := GetViewTimestampKey(playerID, offer.ID)

		var claimCounter int64
		var claimTimestamp int64
		var viewCounter int64
		var viewTimestamp int64

		pipe := redisClient.Client.TxPipeline()
		claimCounterGetOp := pipe.Get(claimCounterKey)
		claimTimestampGetOp := pipe.Get(claimTimestampKey)
		viewCounterGetOp := pipe.Get(viewCounterKey)
		viewTimestampGetOp := pipe.Get(viewTimestampKey)

		err = mr.WithRedisSegment(SegmentRedis, func() error {
			_, err = pipe.Exec()
			return err
		})
		// If err == redis.Nil, then Get didn't found claimCounter fot that key
		// Either player doesn't exist, or it was never inserted yet.
		// Since claimCounter already is 0 and this key will be crated in Incr, just go on.
		if err != nil && err != redis.Nil {
			return nil, err
		}

		claimCounter, err = claimCounterGetOp.Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		claimTimestamp, err = claimTimestampGetOp.Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		lastClaimAt := time.Unix(claimTimestamp, 0)

		viewCounter, err = viewCounterGetOp.Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		viewTimestamp, err = viewTimestampGetOp.Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		lastViewAt := time.Unix(viewTimestamp, 0)

		if f.Max != 0 && int(viewCounter) >= f.Max {
			continue
		}
		if f.Every != "" {
			duration, err := time.ParseDuration(f.Every)
			if err != nil {
				return nil, err
			}
			if lastViewAt.Add(duration).After(t) {
				continue
			}
		}
		if p.Max != 0 && int(claimCounter) >= p.Max {
			continue
		}
		if p.Every != "" {
			duration, err := time.ParseDuration(p.Every)
			if err != nil {
				return nil, err
			}
			if lastClaimAt.Add(duration).After(t) {
				continue
			}
		}
		filteredOffers = append(filteredOffers, offer)
	}

	return filteredOffers, nil
}
