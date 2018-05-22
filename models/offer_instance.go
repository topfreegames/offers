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
	"time"

	"github.com/pmylund/go-cache"
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
	var offerInstance *OfferVersion
	var previousOfferPlayer bool
	var err error
	var nextAt int64

	if offerInstanceID != "" {
		offerInstance, err = getOfferVersionByID(ctx, db, gameID, offerInstanceID, mr)
		if err != nil {
			return nil, false, 0, err
		}
	} else {
		offerInstance, err = getLastOfferInstanceByPlayerIDAndProductID(ctx, db, gameID, playerID, productID, timestamp, mr)
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

	offerInstance, err := getOfferVersionAndOfferEnabled(ctx, db, gameID, offerInstanceID, mr)
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

	var offerVersions []*OfferVersion
	offers := make(map[string]*Offer)

	for _, offer := range filteredOffers {
		offers[offer.ID] = offer
		offerVersions = append(offerVersions, &OfferVersion{
			GameID:       offer.GameID,
			OfferID:      offer.ID,
			OfferVersion: offer.Version,
		})
	}

	offerVersions, err = findOfferVersions(ctx, db, offerVersions, mr)
	if err != nil {
		return nil, err
	}

	for _, offerInstance := range offerVersions {
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
