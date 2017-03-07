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
	. "github.com/topfreegames/offers/testing"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var _ = Describe("Offer Template Models", func() {
	Describe("Get offer template by its id", func() {
		It("should load a template from existent id", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			ot, err := models.GetOfferTemplateByID(db, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(ot.Name).To(Equal("template-1"))
			Expect(ot.Key).To(Equal("da700673-0415-43c3-a8e0-18331b794482"))
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
		It("should create a template with valid parameters", func() {
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				Key:       uuid.NewV4().String(),
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
				Key:       uuid.NewV4().String(),
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

		It("should return error if inserting offer template with repeated key and game with an enabled offer template", func() {
			//Given
			offerTemplate := &models.OfferTemplate{
				Name:      "template-1",
				Key:       "da700673-0415-43c3-a8e0-18331b794482",
				ProductID: "com.tfg.example",
				GameID:    "offers-game",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			//When
			_, err := models.InsertOfferTemplate(db, offerTemplate, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`OfferTemplate could not be saved due to: An offer template with key da700673-0415-43c3-a8e0-18331b794482 already exist and is enabled`))
		})

		It("should insert offer template with repeated key and game with a disabled offer template", func() {
			//Given
			offerTemplate := &models.OfferTemplate{
				Name:      "template-1",
				Key:       "da700673-0415-43c3-a8e0-18331b794482",
				ProductID: "com.tfg.example",
				GameID:    "offers-game",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			//When
			err1 := models.SetEnabledOfferTemplate(db, "dd21ec96-2890-4ba0-b8e2-40ea67196990", false, nil)
			_, err2 := models.InsertOfferTemplate(db, offerTemplate, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("should return error if inserting offer template with missing parameters", func() {
			//Given
			offerTemplate := &models.OfferTemplate{
				Name:      "non-existing-template",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
			}

			//When
			_, err := models.InsertOfferTemplate(db, offerTemplate, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})

		It("should return error if DB is closed", func() {
			oldDB := db
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable
			offerTemplate := &models.OfferTemplate{
				Name:      "New Awesome Game",
				Key:       uuid.NewV4().String(),
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			_, err = models.InsertOfferTemplate(db, offerTemplate, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
			db = oldDB // avoid errors in after each
		})
	})

	Describe("List offer templates", func() {
		It("Should return the full list of offer templates for the given game", func() {
			games, err := models.ListOfferTemplates(db, "offers-game", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(5))
		})
	})

	Describe("Get enabled offer templates", func() {
		It("Should get all enabled offer templates", func() {
			ots, err := models.GetEnabledOfferTemplates(db, "offers-game", nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(ots).To(HaveLen(4))
			Expect(ots[0].Name).To(Equal("template-1"))
			Expect(ots[1].ID).To(Equal("5fed76ab-1fd7-4a91-972d-bca228ce80c4"))
			Expect(ots[1].Key).To(Equal("471c730b-b8ed-4caa-a245-f46822914c8c"))
			Expect(ots[1].Name).To(Equal("template-10"))
			Expect(ots[1].ProductID).To(Equal("com.tfg.sample"))
			Expect(ots[1].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(ots[1].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(ots[1].Period).To(Equal(dat.JSON([]byte(`{"every": "12h"}`))))
			Expect(ots[1].Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"max": 2}`))))
			Expect(ots[1].Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486679000, "from": 1486678000}`))))
			Expect(ots[1].Placement).To(Equal("unique-place"))
			Expect(ots[2].ID).To(Equal("d5114990-77d7-45c4-ba5f-462fc86b213f"))
			Expect(ots[2].Key).To(Equal("bd36b563-8fd4-4f08-82bb-d9344717b50c"))
			Expect(ots[2].Name).To(Equal("template-2"))
			Expect(ots[2].ProductID).To(Equal("com.tfg.sample.2"))
			Expect(ots[2].Contents).To(Equal(dat.JSON([]byte(`{"gems": 100, "gold": 5}`))))
			Expect(ots[2].Metadata).To(Equal(dat.JSON([]byte(`{"meta": "data"}`))))
			Expect(ots[2].Period).To(Equal(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ots[2].Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(ots[2].Trigger).To(Equal(dat.JSON([]byte(`{"to": 1486679200, "from": 1486678000}`))))
			Expect(ots[2].Placement).To(Equal("store"))
			Expect(ots[3].Name).To(Equal("template-3"))
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

	Describe("Get offer template by key", func() {
		It("should return offer template by its name", func() {
			//Given
			key := "da700673-0415-43c3-a8e0-18331b794482"
			gameID := "offers-game"
			expectedOt := &models.OfferTemplate{
				ID:        "dd21ec96-2890-4ba0-b8e2-40ea67196990",
				Key:       key,
				Name:      "template-1",
				ProductID: "com.tfg.sample",
				GameID:    gameID,
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Metadata:  dat.JSON([]byte(`{}`)),
				Period:    dat.JSON([]byte(`{"every": "1s"}`)),
				Frequency: dat.JSON([]byte(`{"every": "1s"}`)),
				Trigger:   dat.JSON([]byte(`{"to": 1486679000, "from": 1486678000}`)),
				Placement: "popup",
				Enabled:   true,
			}

			//When
			ot, err := models.GetEnabledOfferTemplateByKeyAndGame(db, key, gameID, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(ot).To(Equal(expectedOt))
		})

		It("should return empty if key not found", func() {
			//Given
			gameID := "offers-game"
			key := uuid.NewV4().String()

			//When
			_, err := models.GetEnabledOfferTemplateByKeyAndGame(db, key, gameID, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: no rows in result set"))
		})
	})
})
