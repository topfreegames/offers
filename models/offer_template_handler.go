// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/mgutz/dat"
	runner "github.com/mgutz/dat/sqlx-runner"
	"github.com/topfreegames/offers/errors"
)

//OfferTemplate contains the parameters of a template
type OfferTemplate struct {
	ID        string   `db:"id" valid:"matches(^[^-][a-z0-9-]*$),required"`
	Name      string   `db:"name" valid:"ascii,stringlength(1|255),required"`
	ProductID string   `db:"product_id" valid:"ascii,stringlength(1|255),required"`
	GameID    string   `db:"game_id" valid:"ascii,stringlength(1|255),required"`
	Contents  dat.JSON `db:"contents" valid:"required"`
	Metadata  dat.JSON `db:"metadata" valid:""`
	Period    dat.JSON `db:"period" valid:"required"`
	Frequency dat.JSON `db:"frequency" valid:"required"`
	Trigger   dat.JSON `db:"trigger" valid:"required"`
	Enabled   bool     `db:"enabled" valid:"matches(^(true|false)$),optional"`
  Placement string   `db:"placement" valid:"ascii,stringlength(1|255),required"`
}

const enabledOfferTemplates = `
    WHERE
		ot.game_id = $1
		AND ot.enabled = true
`

//GetOfferTemplateByID returns OfferTemplate by ID
func GetOfferTemplateByID(db runner.Connection, id string, mr *MixedMetricsReporter) (*OfferTemplate, error) {
	var ot OfferTemplate
	err := mr.WithDatastoreSegment("offer_templates", "select by id", func() error {
		return db.
			Select(`
				id, name, product_id, game_id,
				contents, metadata, period,
				frequency, trigger, placement, enabled
			`).
			From("offer_templates").
			Where("id = $1", id).
			QueryStruct(&ot)
	})

	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return nil, errors.NewModelNotFoundError("offer template", map[string]interface{}{
				"ID": id,
			})
		}

		return nil, err
	}

	return &ot, nil
}

//GetEnabledOfferTemplates returns all the enabled offers
func GetEnabledOfferTemplates(db runner.Connection, gameID string, mr *MixedMetricsReporter) ([]*OfferTemplate, error) {
	var ots []*OfferTemplate
	err := mr.WithDatastoreSegment("offer_templates", "select", func() error {
		return db.
			Select(`
				id, name, product_id,
				contents, metadata, period,
				frequency, trigger, placement
			`).
			From("offer_templates ot").
			Scope(enabledOfferTemplates, gameID).
			OrderBy("name asc").
			QueryStructs(&ots)
	})
	if err != nil {
		err = HandleNotFoundError("Offer Template", map[string]interface{}{
			"enabled": true,
		}, err)
		return nil, err
	}
	return ots, nil
}

// InsertOfferTemplate inserts a new offer template into DB
func InsertOfferTemplate(db runner.Connection, ot *OfferTemplate, mr *MixedMetricsReporter) error {
	return mr.WithDatastoreSegment("offer_templates", "insert", func() error {
		return db.
			InsertInto("offer_templates").
			Columns("id", "name", "product_id", "game_id", "contents", "period", "frequency", "trigger", "placement").
			Record(ot).
			Returning("id").
			QueryStruct(ot)
	})
}
