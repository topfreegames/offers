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
	Describe("Get offer template by its id", func() {
		It("should load a template from existent id", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			ot, err := models.GetOfferTemplateByID(db, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(ot.Name).To(Equal("template-1"))
			Expect(ot.ProductID).To(Equal("com.tfg.sample"))
			Expect(ot.GameID).To(Equal("offers-game"))
			Expect(ot.Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(ot.Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(ot.Period).To(Equal(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ot.Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ot.Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486679000, "from": 1486678000}`))))
			Expect(ot.Enabled).To(BeTrue())
			Expect(ot.Placement).To(Equal("popup"))
		})

		It("should not load a template from unexistent ID", func() {
			id := uuid.NewV4().String()
			_, err := models.GetOfferTemplateByID(db, id, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferTemplate was not found with specified filters."))
		})
	})

	Describe("Offer Template instance", func() {
		It("should create an template with valid parameters", func() {
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			ot, err := models.InsertOfferTemplate(db, offerTemplate, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(ot.ID).NotTo(Equal(""))
			Expect(ot.Name).To(Equal(offerTemplate.Name))
			Expect(ot.ProductID).To(Equal(offerTemplate.ProductID))
			Expect(ot.GameID).To(Equal(offerTemplate.GameID))
			Expect(ot.Contents).To(Equal(offerTemplate.Contents))
			Expect(ot.Metadata).To(Equal(offerTemplate.Metadata))
			Expect(ot.Period).To(Equal(offerTemplate.Period))
			Expect(ot.Frequency).To(Equal(offerTemplate.Frequency))
			Expect(ot.Trigger).To(Equal(offerTemplate.Trigger))
			Expect(ot.Enabled).To(BeTrue())
			Expect(ot.Placement).To(Equal(offerTemplate.Placement))

			dbOt, err := models.GetOfferTemplateByID(db, ot.ID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbOt.Name).To(Equal(offerTemplate.Name))
			Expect(dbOt.ProductID).To(Equal(offerTemplate.ProductID))
			Expect(dbOt.GameID).To(Equal(offerTemplate.GameID))
			Expect(dbOt.Contents).To(Equal(offerTemplate.Contents))
			Expect(dbOt.Metadata).To(Equal(offerTemplate.Metadata))
			Expect(dbOt.Period).To(Equal(offerTemplate.Period))
			Expect(dbOt.Frequency).To(Equal(offerTemplate.Frequency))
			Expect(dbOt.Trigger).To(BeEquivalentTo(offerTemplate.Trigger))
			Expect(dbOt.Enabled).To(BeTrue())
			Expect(dbOt.Placement).To(Equal(offerTemplate.Placement))
		})

		It("should return error if game with given id does not exist", func() {
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				ProductID: "com.tfg.example",
				GameID:    "non-existing-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "1s"}`)),
				Frequency: dat.JSON([]byte(`{"every": "1s"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				Placement: "popup",
			}
			_, err := models.InsertOfferTemplate(db, offerTemplate, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferTemplate could not be saved due to: insert or update on table \"offer_templates\" violates foreign key constraint \"offer_templates_game_id_fkey\""))

			var ot models.OfferTemplate
			err = conn.
				Select("id").
				From("offer_templates").
				Where("game_id = $1", "non-existing-game-id").
				QueryStruct(&ot)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: no rows in result set"))
		})
	})

	Describe("Get enabled offer templates", func() {
		It("Should get all enabled offer templates", func() {
			ots, err := models.GetEnabledOfferTemplates(db, "offers-game", nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(ots).To(HaveLen(3))
			Expect(ots[0].Name).To(Equal("template-1"))
			Expect(ots[1].ID).To(Equal("d5114990-77d7-45c4-ba5f-462fc86b213f"))
			Expect(ots[1].Name).To(Equal("template-2"))
			Expect(ots[1].ProductID).To(Equal("com.tfg.sample.2"))
			Expect(ots[1].Contents).To(Equal(dat.JSON([]byte(`{"gems": 100, "gold": 5}`))))
			Expect(ots[1].Metadata).To(Equal(dat.JSON([]byte(`{"meta": "data"}`))))
			Expect(ots[1].Period).To(Equal(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ots[1].Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ots[1].Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486679000, "from": 1486678000}`))))
			Expect(ots[1].Placement).To(Equal("store"))
			Expect(ots[2].Name).To(Equal("template-3"))
		})
	})

	Describe("Set enabled offer template", func() {
		It("should disable an enabled offer", func() {
			//Given
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			enabled := false

			//When
			err1 := models.SetEnabledOfferTemplate(db, templateID, enabled, nil)
			ot, err2 := models.GetOfferTemplateByID(db, templateID, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(ot.Enabled).To(BeFalse())
		})

		It("should enable an enabled offer", func() {
			//Given
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			enabled := true

			//When
			err1 := models.SetEnabledOfferTemplate(db, templateID, enabled, nil)
			ot, err2 := models.GetOfferTemplateByID(db, templateID, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(ot.Enabled).To(BeTrue())
		})

		It("should enable a disabled offer", func() {
			//Given
			templateID := "27b0370f-bd61-4346-a10d-50ec052ae125"
			enabled := true

			//When
			err1 := models.SetEnabledOfferTemplate(db, templateID, enabled, nil)
			ot, err2 := models.GetOfferTemplateByID(db, templateID, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(ot.Enabled).To(BeTrue())
		})

		It("should return error if id doesn't exist", func() {
			//Given
			templateID := uuid.NewV4().String()
			enabled := true

			//When
			err1 := models.SetEnabledOfferTemplate(db, templateID, enabled, nil)

			//Then
			Expect(err1).To(HaveOccurred())
		})
	})
})
