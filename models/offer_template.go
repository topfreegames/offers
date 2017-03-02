// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/topfreegames/offers/errors"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//OfferTemplate contains the parameters of a template
type OfferTemplate struct {
	ID        string   `db:"id" json:"id" valid:"uuidv4"`
	Name      string   `db:"name" json:"name" valid:"ascii,stringlength(1|255),required"`
	ProductID string   `db:"product_id" json:"productId" valid:"ascii,stringlength(1|255),required"`
	GameID    string   `db:"game_id" json:"gameId" valid:"matches(^[^-][a-z0-9-]*$),stringlength(1|255),required"`
	Contents  dat.JSON `db:"contents" json:"contents" valid:"RequiredJSONObject"`
	Metadata  dat.JSON `db:"metadata" json:"metadata" valid:"JSONObject"`
	Period    dat.JSON `db:"period" json:"period" valid:"RequiredJSONObject"`
	Frequency dat.JSON `db:"frequency" json:"frequency" valid:"RequiredJSONObject"`
	Trigger   dat.JSON `db:"trigger" json:"trigger" valid:"RequiredJSONObject"`
	Enabled   bool     `db:"enabled" json:"enabled" valid:"matches(^(true|false)$),optional"`
	Placement string   `db:"placement" json:"placement" valid:"ascii,stringlength(1|255),required"`
}

//OfferTemplateToUpdate is used by api/app.go to call setEnabled
type OfferTemplateToUpdate struct {
	ID      string `db:"id" valid:"uuidv4,required"`
	Enabled bool   `db:"enabled" valid:"bool,optional"`
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

	err = HandleNotFoundError("OfferTemplate", map[string]interface{}{"ID": id}, err)
	return &ot, err
}

//GetOfferTemplateByNameAndGame returns OfferTemplate by Name
func GetOfferTemplateByNameAndGame(db runner.Connection, name, gameID string, mr *MixedMetricsReporter) (*OfferTemplate, error) {
	var ot OfferTemplate
	err := mr.WithDatastoreSegment("offer_templates", "select by name", func() error {
		return db.
			Select(`
				id, name, product_id, game_id,
				contents, metadata, period,
				frequency, trigger, placement, enabled
			`).
			From("offer_templates").
			Where("name = $1 AND game_id = $2 AND enabled = true", name, gameID).
			QueryStruct(&ot)
	})

	err = HandleNotFoundError("OfferTemplate", map[string]interface{}{"Name": name}, err)
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
	err = HandleNotFoundError("OfferTemplate", map[string]interface{}{"enabled": true}, err)
	return ots, err
}

// InsertOfferTemplate inserts a new offer template into DB
func InsertOfferTemplate(db runner.Connection, ot *OfferTemplate, mr *MixedMetricsReporter) (*OfferTemplate, error) {
	_, err := GetOfferTemplateByNameAndGame(db, ot.Name, ot.GameID, mr)

	if err != nil {
		notFoundErr := HandleNotFoundError("OfferTemplate", map[string]interface{}{"Name": ot.Name}, err)
		if err != notFoundErr {
			return nil, err
		}
	} else {
		msg := "An offer template with name " + ot.Name + " already exist and is enabled"
		return ot, errors.NewConflictedModelError("OfferTemplate", msg)
	}

	if ot.Metadata == nil {
		ot.Metadata = dat.JSON([]byte(`{}`))
	}
	err = mr.WithDatastoreSegment("offer_templates", SegmentInsert, func() error {
		return db.
			InsertInto("offer_templates").
			Columns("name", "product_id", "game_id", "contents", "period", "frequency", "trigger", "placement", "metadata").
			Record(ot).
			Returning("id, enabled").
			QueryStruct(ot)
	})

	foreignKeyErr := HandleForeignKeyViolationError("OfferTemplate", err)
	return ot, foreignKeyErr
}

//SetEnabledOfferTemplate can enable or disable an offer template
func SetEnabledOfferTemplate(db runner.Connection, id string, enabled bool, mr *MixedMetricsReporter) error {
	var offerTemplate OfferTemplate
	err := mr.WithDatastoreSegment("offer_templates", SegmentUpdate, func() error {
		return db.
			Update("offer_templates").
			Set("enabled", enabled).
			Where("id=$1", id).
			Returning("id").
			QueryStruct(&offerTemplate)
	})

	err = HandleNotFoundError("OfferTemplate", map[string]interface{}{
		"ID": id,
	}, err)

	return err
}
