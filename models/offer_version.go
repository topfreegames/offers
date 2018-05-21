package models

import (
	"context"
	"time"

	edat "github.com/topfreegames/extensions/dat"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//OfferVersion represents a tenant in offers API it cannot be updated, only inserted
type OfferVersion struct {
	ID           string       `db:"id" json:"id" valid:"uuidv4,required"`
	GameID       string       `db:"game_id" json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	OfferID      string       `db:"offer_id" json:"offerId" valid:"uuidv4,required"`
	OfferVersion int          `db:"offer_version" json:"offerVersion" valid:"int,required"`
	Contents     dat.JSON     `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	ProductID    string       `db:"product_id" json:"productId" valid:"ascii,stringlength(1|255)"`
	Cost         dat.JSON     `db:"cost" json:"cost" valid:"JSONObject"`
	CreatedAt    dat.NullTime `db:"created_at" json:"createdAt" valid:""`
}

func offerVersionFromOffer(offer *Offer) *OfferVersion {
	return &OfferVersion{
		GameID:       offer.GameID,
		OfferID:      offer.ID,
		OfferVersion: offer.Version,
		Contents:     offer.Contents,
		ProductID:    offer.ProductID,
		Cost:         offer.Cost,
	}
}

func getOfferToReturn(
	ctx context.Context,
	db runner.Connection,
	gameID, offerID string,
	mr *MixedMetricsReporter,
) (*OfferToReturn, error) {
	var offerVersion OfferToReturn

	err := mr.WithDatastoreSegment("offer_versions", SegmentSelect, func() error {
		builder := db.
			Select("oi.id, oi.product_id, oi.contents, oi.cost, o.metadata, o.trigger#>>'{to}' AS expire_at")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_versions oi JOIN offers o ON (oi.offer_id=o.id)").
			Where("oi.id=$1 AND oi.game_id=$2", offerID, gameID).
			QueryStruct(&offerVersion)
	})

	err = HandleNotFoundError("OfferVersion", map[string]interface{}{
		"GameID": gameID,
		"ID":     offerID,
	}, err)

	return &offerVersion, err
}

//GetOfferInfo returns the offers that match the criteria of enabled offer templates
func GetOfferInfo(
	ctx context.Context,
	db runner.Connection,
	gameID, offerInstanceID string,
	expireDuration time.Duration,
	mr *MixedMetricsReporter,
) (*OfferToReturn, error) {
	offer, err := getOfferToReturn(ctx, db, gameID, offerInstanceID, mr)

	if err != nil {
		return nil, err
	}

	return offer, nil
}
