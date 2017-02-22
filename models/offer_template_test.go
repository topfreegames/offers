// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
	"gopkg.in/mgutz/dat.v2/dat"
)

var _ = Describe("Offer Template Models", func() {
	Describe("Offer Template instance", func() {
		It("should load a template by its ID", func() {
			//Given
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"

			//When
			ot, err := models.GetOfferTemplateByID(db, id, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(ot.ProductID).To(Equal("com.tfg.sample"))
			Expect(ot.GameID).To(Equal("offers-game"))
			Expect(ot.Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(ot.Metadata).To(Equal(dat.JSON([]byte(`{}`))))
      Expect(ot.Period).To(Equal(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ot.Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ot.Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486679000, "from": 1486678000}`))))
			Expect(ot.Placement).To(Equal("popup"))
		})

		It("should not load a template from invalid ID", func() {
			//Given
			id := uuid.NewV4().String()

			//When
			_, err := models.GetOfferTemplateByID(db, id, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})

		It("should create an template with valid parameters", func() {
			//Given
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				Placement: "popup",
			}

			//When
			err := models.InsertOfferTemplate(db, offerTemplate, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error if game id does not exist", func() {
			//Given
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "non-existing-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"type": "once"}`)),
				Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				Placement: "popup",
			}

			//When
			err := models.InsertOfferTemplate(db, offerTemplate, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Get enabled offer templates", func() {
		It("Should get all enabled offer templates", func() {
			ots, err := models.GetEnabledOfferTemplates(db, "offers-game", nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(ots).To(HaveLen(2))
			Expect(ots[0].Name).To(Equal("template-1"))
			Expect(ots[1].Name).To(Equal("template-2"))
		})
	})
})
