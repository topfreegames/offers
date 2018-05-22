package models

import (
	"context"
	"fmt"
	"strings"
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

	// This part was left to be backwards compatible with previously existing offer instances
	if err != nil && IsNoRowsInResultSetError(err) {
		err = mr.WithDatastoreSegment("offer_instances", SegmentSelect, func() error {
			builder := db.
				Select("oi.id, oi.product_id, oi.contents, oi.cost, o.metadata, o.trigger#>>'{to}' AS expire_at")
			builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
			return builder.From("offer_instances oi JOIN offers o ON (oi.offer_id=o.id)").
				Where("oi.id=$1 AND oi.game_id=$2", offerID, gameID).
				QueryStruct(&offerVersion)
		})
	}
	err = handleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID": gameID,
		"ID":     offerID,
	}, err)

	return &offerVersion, err
}

func findOfferVersions(
	ctx context.Context,
	db runner.Connection,
	offerVersions []*OfferVersion,
	mr *MixedMetricsReporter,
) ([]*OfferVersion, error) {
	resOfferInstances := make([]*OfferVersion, 0, len(offerVersions))
	var err error

	whereClause := make([]string, 0, len(offerVersions))
	for _, o := range offerVersions {
		whereClause = append(whereClause, fmt.Sprintf("(offer_id='%s' AND offer_version=%d)",
			o.OfferID, o.OfferVersion))
	}

	query := fmt.Sprintf(`
	SELECT * FROM (SELECT id, offer_id FROM offer_versions
	WHERE game_id=$1 AND (%s)) AS sel
	`, strings.Join(whereClause, " OR "))

	err = mr.WithDatastoreSegment("offer_versions", SegmentInsect, func() error {
		builder := db.SQL(query, offerVersions[0].GameID)
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.QueryStructs(&resOfferInstances)
	})

	return resOfferInstances, err
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

func getOfferVersionAndOfferEnabled(ctx context.Context, db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferInstanceOffer, error) {
	var offerInstance OfferInstanceOffer
	err := mr.WithDatastoreSegment("offer_versions", SegmentSelect, func() error {
		builder := db.Select("oi.id, oi.offer_id, oi.contents, o.enabled")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_versions oi JOIN offers o ON (oi.offer_id=o.id)").
			Where("oi.id=$1 AND oi.game_id=$2", id, gameID).
			QueryStruct(&offerInstance)
	})

	err = handleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID": gameID,
		"ID":     id,
	}, err)

	return &offerInstance, err
}

func getOfferVersionByID(ctx context.Context, db runner.Connection, gameID, id string, mr *MixedMetricsReporter) (*OfferVersion, error) {
	var offerInstance OfferVersion
	err := mr.WithDatastoreSegment("offer_versions", SegmentSelect, func() error {
		builder := db.Select("id, offer_id, contents")
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.From("offer_versions").
			Where("id=$1 AND game_id=$2", id, gameID).
			QueryStruct(&offerInstance)
	})

	err = handleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID": gameID,
		"ID":     id,
	}, err)

	return &offerInstance, err
}

func getLastOfferInstanceByPlayerIDAndProductID(ctx context.Context, db runner.Connection, gameID, playerID, productID string, timestamp int64, mr *MixedMetricsReporter) (*OfferVersion, error) {
	var offerInstance OfferVersion
	err := mr.WithDatastoreSegment("offer_players", SegmentSelect, func() error {
		builder := db.SQL(`
		SELECT ov.id, ov.offer_id, ov.contents
			FROM offer_players op JOIN offer_versions ov ON ov.offer_id=op.offer_id
			WHERE op.game_id=$1 AND op.player_id=$2 AND op.view_timestamp < to_timestamp($4)
				AND ov.game_id=$1 AND ov.product_id=$3
			ORDER BY op.view_timestamp DESC FETCH FIRST 1 ROW ONLY`,
			gameID, playerID, productID, timestamp,
		)
		builder.Execer = edat.NewExecer(builder.Execer).WithContext(ctx)
		return builder.QueryStruct(&offerInstance)
	})

	err = handleNotFoundError("OfferInstance", map[string]interface{}{
		"GameID":    gameID,
		"PlayerID":  playerID,
		"ProductID": productID,
	}, err)

	return &offerInstance, err
}
