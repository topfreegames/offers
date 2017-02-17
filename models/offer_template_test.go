// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"github.com/mgutz/dat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
)

var _ = Describe("Offer Template Models", func() {
	Describe("Offer Template instance", func() {
		It("should load a template", func() {

		})

		It("should create an template with valid parameters", func() {
			offerTemplate := &models.OfferTemplate{
				ID:        uuid.NewV4(),
				Name:      "New Awesome Game",
				Pid:       "com.tfg.example",
				GameID:    "nonexisting-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
			}

			err := db.InsertInto("offer_templates").
				Columns("id", "name", "pid", "gameid", "contents", "period", "frequency", "trigger").
				Record(offerTemplate).
				Returning("id").
				QueryStruct(offerTemplate)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error with invalid parameters", func() {
			offerTemplate := &models.OfferTemplate{
				ID:        uuid.NewV4(),
				Name:      "New Awesome Game",
				Pid:       "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
			}

			err := db.InsertInto("offer_templates").
				Columns("id", "name", "pid", "gameid", "contents", "period", "frequency", "trigger").
				Record(offerTemplate).
				Returning("id").
				QueryStruct(offerTemplate)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get all available offers", func() {
		It("Should get all available offers", func() {

		})
	})
})
