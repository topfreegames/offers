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
		It("should load a template by its ID", func() {
			var ot models.OfferTemplate
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			err := db.
				Select("*").
				From("offer_templates").
				Where("id = $1", id).
				QueryStruct(&ot)

			Expect(err).NotTo(HaveOccurred())

			Expect(ot.ProductID).To(Equal("com.tfg.sample"))
			Expect(ot.GameID).To(Equal("awesome game"))
			Expect(ot.Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(ot.Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(ot.Period).To(Equal(dat.JSON([]byte(`{"type": "once"}`))))
			Expect(ot.Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"unit": "hour", "every": 12}`))))
			Expect(ot.Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486678079, "from": 1486678078}`))))
		})

		It("should not load a template from invalid ID", func() {
			var ot models.OfferTemplate
			id := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
			err := db.
				Select("*").
				From("offer_templates").
				Where("id = $1", id).
				QueryStruct(&ot)

			Expect(err).To(HaveOccurred())
		})

		It("should create an template with valid parameters", func() {
			offerTemplate := &models.OfferTemplate{
				ID:        uuid.NewV4(),
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "nonexisting-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
			}

			err := db.InsertInto("offer_templates").
				Columns("id", "name", "product_id", "game_id", "contents", "period", "frequency", "trigger").
				Record(offerTemplate).
				Returning("id").
				QueryStruct(offerTemplate)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error if game id does not exist", func() {
			offerTemplate := &models.OfferTemplate{
				ID:        uuid.NewV4(),
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "nonexisting-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
			}

			err := db.InsertInto("offer_templates").
				Columns("id", "name", "product_id", "game_id", "contents", "period", "frequency", "trigger").
				Record(offerTemplate).
				Returning("id").
				QueryStruct(offerTemplate)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get all available offers", func() {
		It("Should get all available offers", func() {
			ots, err := models.GetEnabledOfferTemplates(db, "awesome game", nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(ots).To(HaveLen(2))

			Expect(ots[0].Name).To(Equal("ot-1"))
			Expect(ots[1].Name).To(Equal("ot-2"))
		})
	})
})
