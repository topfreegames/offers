// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pmylund/go-cache"
	edat "github.com/topfreegames/extensions/dat"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//OfferInstance represents a tenant in offers API it cannot be updated, only inserted
type OfferInstance struct {
	ID           string       `db:"id" json:"id" valid:"uuidv4,required"`
	GameID       string       `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	PlayerID     string       `db:"player_id" json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	OfferID      string       `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	OfferVersion int          `db:"offer_version" json:"offerVersion" valid:"int,required"`
	Contents     dat.JSON     `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	ProductID    string       `db:"product_id" json:"productId" valid:"ascii,stringlength(1|255)"`
	Cost         dat.JSON     `db:"cost" json:"cost" valid:"JSONObject"`
	CreatedAt    dat.NullTime `db:"created_at" json:"createdAt" valid:""`
}

//OfferInstanceOffer is a join of OfferInstance with offer
type OfferInstanceOffer struct {
	ID       string   `db:"id" json:"id" valid:"uuidv4,required"`
	GameID   string   `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	OfferID  string   `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	Contents dat.JSON `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	Enabled  bool     `db:"enabled" json:"enabled"`
}

//OfferToReturn has the fields for the returned offer
type OfferToReturn struct {
	ID        string   `db:"id" json:"id"`
	ProductID string   `db:"product_id" json:"productId,omitempty"`
	Cost      dat.JSON `db:"cost" json:"cost,omitempty" valid:"JSONObject"`
	Contents  dat.JSON `db:"contents" json:"contents"`
	Metadata  dat.JSON `db:"metadata" json:"metadata"`
	ExpireAt  int64    `db:"expire_at" json:"expireAt"`
}

//FrequencyOrPeriod is the struct for basic Frequency and Period types
type FrequencyOrPeriod struct {
	Every string
	Max   int
}

//GetOfferInstanceAndOfferEnabled returns a offer by its pk
func GetOfferInstanceAndOfferEnabled(ctx context.Context, db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferInstanceOffer, error) {
	var offerInstance OfferInstanceOffer
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		builder := db.Select("oi.id, oi.offer_id, oi.contents, o.enabled")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_instances oi JOIN offers o ON (oi.offer_id=o.id)").
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
func GetOfferInstanceByID(ctx context.Context, db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferInstance, error) {
	var offerInstance OfferInstance
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		builder := db.Select("id, offer_id, contents")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_instances").
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
	ctx context.Context,
	db runner.Connection,
	gameID, offerID string,
	claimCounter int,
	t time.Time,
	mr *MixedMetricsReporter,
) (int64, error) {
	offer, err := GetOfferByID(ctx, db, gameID, offerID, mr)
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
	ctx context.Context,
	db runner.Connection,
	gameID, offerInstanceID, playerID, productID, transactionID string,
	timestamp int64,
	t time.Time,
	mr *MixedMetricsReporter,
) (dat.JSON, bool, int64, error) {

	// If an offer instance id is sent
	var offerInstance *OfferInstance
	var previousOfferPlayer bool
	var err error
	var nextAt int64

	if offerInstanceID != "" {
		offerInstance, err = GetOfferInstanceByID(ctx, db, gameID, offerInstanceID, mr)
		if err != nil {
			return nil, false, 0, err
		}
	} else {
		offerInstance, err = GetLastOfferInstanceByPlayerIDAndProductID(ctx, db, gameID, playerID, productID, timestamp, mr)
		if err != nil {
			return nil, false, 0, err
		}
	}
	offerPlayer, err := GetOfferPlayer(ctx, db, gameID, playerID, offerInstance.OfferID, mr)
	if err == nil {
		previousOfferPlayer = true
	} else if !IsNoRowsInResultSetError(err) {
		return nil, false, 0, err
	} else {
		offerPlayer = &OfferPlayer{
			GameID:       gameID,
			PlayerID:     playerID,
			OfferID:      offerInstance.OfferID,
			Transactions: dat.JSON([]byte(`[]`)),
			Impressions:  dat.JSON([]byte(`[]`)),
		}
	}

	isReplay := false
	var transactions []string
	err = offerPlayer.Transactions.Unmarshal(&transactions)

	if err != nil {
		return nil, false, 0, err
	}
	for _, tr := range transactions {
		if transactionID == tr {
			isReplay = true
			break
		}
	}
	if isReplay {
		nextAt, err = getClaimedOfferNextAt(
			ctx, db, gameID, offerInstance.OfferID,
			offerPlayer.ClaimCounter, offerPlayer.ClaimTimestamp.Time, mr)
		if err != nil {
			return nil, false, 0, err
		}
		return offerInstance.Contents, true, nextAt, nil
	}

	if previousOfferPlayer {
		jsonTr, err := dat.NewJSON(append(transactions, transactionID))
		if err != nil {
			return nil, false, 0, err
		}
		offerPlayer.Transactions = *jsonTr
		err = ClaimOfferPlayer(ctx, db, offerPlayer, time.Unix(timestamp, 0), mr)
		if err != nil {
			return nil, false, 0, err
		}
	} else {
		offerPlayer.ClaimCounter = 1
		offerPlayer.ClaimTimestamp = dat.NullTimeFrom(time.Unix(timestamp, 0))
		offerPlayer.Transactions = dat.JSON([]byte(fmt.Sprintf(`["%s"]`, transactionID)))
		err = CreateOfferPlayer(ctx, db, offerPlayer, mr)
		if err != nil {
			return nil, false, 0, err
		}
	}

	nextAt, err = getClaimedOfferNextAt(
		ctx, db, gameID, offerInstance.OfferID,
		offerPlayer.ClaimCounter, time.Unix(timestamp, 0), mr)
	if err != nil {
		return nil, false, 0, err
	}
	return offerInstance.Contents, false, nextAt, nil
}

//GetLastOfferInstanceByPlayerIDAndProductID returns a offer by gameId, playerId and productId
func GetLastOfferInstanceByPlayerIDAndProductID(ctx context.Context, db runner.Connection, gameID, playerID, productID string, timestamp int64, mr *MixedMetricsReporter) (*OfferInstance, error) {
	var offerInstance OfferInstance
	err := mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
		builder := db.SQL("SELECT id, offer_id, contents "+
			"FROM offer_instances "+
			"WHERE game_id=$1 AND player_id=$2 AND product_id=$3 AND created_at < to_timestamp($4) "+
			"ORDER BY created_at DESC FETCH FIRST 1 ROW ONLY", gameID, playerID, productID, timestamp)
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.QueryStruct(&offerInstance)
	})

	err = HandleNotFoundError("offerInstance", map[string]interface{}{
		"GameID":    gameID,
		"PlayerID":  playerID,
		"ProductID": productID,
	}, err)

	return &offerInstance, err
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

//ViewOffer views the offer
func ViewOffer(
	ctx context.Context,
	db runner.Connection,
	gameID, offerInstanceID, playerID, impressionID string,
	t time.Time,
	mr *MixedMetricsReporter,
) (bool, int64, error) {
	var nextAt int64
	var previousOfferPlayer bool

	offerInstance, err := GetOfferInstanceAndOfferEnabled(ctx, db, gameID, offerInstanceID, mr)
	if err != nil {
		return false, 0, err
	}

	// Offer is disabled
	if !offerInstance.Enabled {
		return false, 0, nil
	}

	offerPlayer, err := GetOfferPlayer(ctx, db, gameID, playerID, offerInstance.OfferID, mr)
	if err == nil {
		previousOfferPlayer = true
	} else if !IsNoRowsInResultSetError(err) {
		return false, 0, err
	} else {
		offerPlayer = &OfferPlayer{
			GameID:       gameID,
			PlayerID:     playerID,
			OfferID:      offerInstance.OfferID,
			Transactions: dat.JSON([]byte(`[]`)),
			Impressions:  dat.JSON([]byte(`[]`)),
		}
	}

	isReplay := false
	var impressions []string
	err = offerPlayer.Impressions.Unmarshal(&impressions)
	if err != nil {
		return false, 0, err
	}
	for _, imp := range impressions {
		if impressionID == imp {
			isReplay = true
			break
		}
	}

	if isReplay {
		nextAt, err = getViewedOfferNextAt(ctx, db, gameID, offerInstance.OfferID, offerPlayer.ViewCounter, t, mr)
		if err != nil {
			return false, 0, err
		}
		return true, nextAt, nil
	}

	if previousOfferPlayer {
		jsonTr, err := dat.NewJSON(append(impressions, impressionID))
		if err != nil {
			return false, 0, err
		}
		offerPlayer.Impressions = *jsonTr
		err = ViewOfferPlayer(ctx, db, offerPlayer, t, mr)
		if err != nil {
			return false, 0, err
		}
	} else {
		offerPlayer.ViewCounter = 1
		offerPlayer.ViewTimestamp = dat.NullTimeFrom(t)
		offerPlayer.Impressions = dat.JSON([]byte(fmt.Sprintf(`["%s"]`, impressionID)))
		err = CreateOfferPlayer(ctx, db, offerPlayer, mr)
		if err != nil {
			return false, 0, err
		}
	}

	nextAt, err = getViewedOfferNextAt(ctx, db, gameID, offerInstance.OfferID, offerPlayer.ViewCounter, t, mr)
	if err != nil {
		return false, 0, err
	}
	return false, nextAt, nil
}

//GetAvailableOffers returns the offers that match the criteria of enabled offer templates
func GetAvailableOffers(
	ctx context.Context,
	db runner.Connection,
	offersCache *cache.Cache,
	gameID, playerID string,
	t time.Time,
	expireDuration time.Duration,
	filterAttrs map[string]string,
	allowInefficientQueries bool,
	mr *MixedMetricsReporter,
) (map[string][]*OfferToReturn, error) {
	offersByPlacement := make(map[string][]*OfferToReturn)

	enabledOffers, err := GetEnabledOffers(
		ctx,
		db,
		gameID,
		offersCache,
		expireDuration,
		t,
		filterAttrs,
		allowInefficientQueries,
		mr,
	)
	if err != nil {
		return nil, err
	}
	if len(enabledOffers) == 0 {
		return offersByPlacement, nil
	}

	offersByPlayer, err := GetOffersByPlayer(ctx, db, gameID, playerID, mr)
	if err != nil {
		return nil, err
	}

	filteredOffers, err := filterOffersByFrequencyAndPeriod(playerID, enabledOffers, offersByPlayer, t, mr)
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
			Cost:         offer.Cost,
		})
	}

	// TODO: Change this to use offerVersions
	offerInstances, err = findOrCreateOfferInstance(ctx, db, offerInstances, t, mr)
	if err != nil {
		return nil, err
	}

	for _, offerInstance := range offerInstances {
		offer := offers[offerInstance.OfferID]

		var trigger Times
		json.Unmarshal(offer.Trigger, &trigger)
		offerToReturn := &OfferToReturn{
			ID:        offerInstance.ID,
			ProductID: offer.ProductID,
			Contents:  offer.Contents,
			Cost:      offer.Cost,
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
	ctx context.Context,
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
			fmt.Sprintf("('%s', '%s', '%s', '%d', '%s'::jsonb, '%s', '%s'::jsonb)", o.GameID, o.PlayerID, o.OfferID, o.OfferVersion, o.Contents, o.ProductID, o.Cost))
	}

	query := fmt.Sprintf(`
	WITH
		sel AS (SELECT id, offer_id FROM offer_instances WHERE %s),
		ins AS (INSERT INTO offer_instances(game_id, player_id, offer_id, offer_version, contents, product_id, cost)
				VALUES %s
				ON CONFLICT DO NOTHING
				RETURNING id, offer_id)
	SELECT * FROM ins UNION ALL SELECT * FROM sel
	`, strings.Join(whereClause, " OR "), strings.Join(valueArgs, ","))

	err = mr.WithDatastoreSegment("offer_instances", SegmentInsect, func() error {
		builder := db.SQL(query)
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.QueryStructs(&resOfferInstances)
	})

	return resOfferInstances, err
}

func filterOffersByFrequencyAndPeriod(
	playerID string,
	offers []*Offer,
	playerOffers []*OfferPlayer,
	t time.Time,
	mr *MixedMetricsReporter,
) ([]*Offer, error) {
	playerOffersByOfferID := map[string]*OfferPlayer{}
	for _, playerOffer := range playerOffers {
		playerOffersByOfferID[playerOffer.OfferID] = playerOffer
	}
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

		offerPlayer := &OfferPlayer{}
		if val, ok := playerOffersByOfferID[offer.ID]; ok {
			offerPlayer = val
		}

		if f.Max != 0 && offerPlayer.ViewCounter >= f.Max {
			continue
		}
		if f.Every != "" {
			duration, err := time.ParseDuration(f.Every)
			if err != nil {
				return nil, err
			}
			if offerPlayer.ViewTimestamp.Time.Add(duration).After(t) {
				continue
			}
		}
		if p.Max != 0 && offerPlayer.ClaimCounter >= p.Max {
			continue
		}
		if p.Every != "" {
			duration, err := time.ParseDuration(p.Every)
			if err != nil {
				return nil, err
			}
			if offerPlayer.ClaimTimestamp.Time.Add(duration).After(t) {
				continue
			}
		}
		filteredOffers = append(filteredOffers, offer)
	}

	return filteredOffers, nil
}
