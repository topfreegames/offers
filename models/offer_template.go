// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//OfferTemplate contains the parameters of a template
type OfferTemplate struct {
	ID        string   `db:"id" valid:"uuidv4,required"`
	Name      string   `db:"name" valid:"ascii,stringlength(1|255),required"`
	ProductID string   `db:"product_id" valid:"ascii,stringlength(1|255),required"`
	GameID    string   `db:"game_id" valid:"matches(^[^-][a-z0-9-]*$),stringlength(1|255),required"`
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

	err = HandleNotFoundError("Offer Template", map[string]interface{}{"ID": id}, err)
	return &ot, err
}

//GetEnabledOfferTemplates returns all the enabled offers
func GetEnabledOfferTemplates(db runner.Connection, gameID string, mr *MixedMetricsReporter) ([]*OfferTemplate, error) {
	var ots []*OfferTemplate
	err := mr.WithDatastoreSegment("offer_templates", "select", func() error {
		return db.
			Select(`
				id, name, game_id, product_id,
				contents, metadata, period,
				frequency, trigger, placement
			`).
			From("offer_templates ot").
			Scope(enabledOfferTemplates, gameID).
			OrderBy("name asc").
			QueryStructs(&ots)
	})
	err = HandleNotFoundError("Offer Template", map[string]interface{}{"enabled": true}, err)
	return ots, err
}

// InsertOfferTemplate inserts a new offer template into DB
func InsertOfferTemplate(db runner.Connection, ot *OfferTemplate, mr *MixedMetricsReporter) (*OfferTemplate, error) {
	if ot.Metadata == nil {
		ot.Metadata = dat.JSON([]byte(`{}`))
	}
	err := mr.WithDatastoreSegment("offer_templates", "insert", func() error {
		return db.
			InsertInto("offer_templates").
			Columns("name", "product_id", "game_id", "contents", "period", "frequency", "trigger", "placement", "metadata").
			Record(ot).
			Returning("id, enabled").
			QueryStruct(ot)
	})
	return ot, err
}
