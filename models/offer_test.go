// offers api
// https://github.com/topfreeoffers/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models_test

import (
	"time"

	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offers Model", func() {
	defaultOfferTemplateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
	defaultOfferTemplateKey := "da700673-0415-43c3-a8e0-18331b794482"
	defaultOfferID := "56fc0477-39f1-485c-898e-4909e9155eb1"

	Describe("Offer Instance", func() {
		It("Should load a offer", func() {
			//Given

			//When
			var offer models.Offer
			err := db.
				Select("*").
				From("offers").
				Where("id = $1", defaultOfferID).
				QueryStruct(&offer)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).To(Equal(defaultOfferID))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID).To(Equal(defaultOfferTemplateID))
			Expect(offer.OfferTemplateKey).To(Equal(defaultOfferTemplateKey))
			Expect(offer.CreatedAt.Valid).To(BeTrue())
		})

		It("Should create offer", func() {
			//Given
			offer := &models.Offer{
				GameID:           "offers-game",
				OfferTemplateID:  defaultOfferTemplateID,
				OfferTemplateKey: defaultOfferTemplateKey,
				PlayerID:         "player-3",
			}

			//When
			err := db.
				InsertInto("offers").
				Columns("game_id", "offer_template_id", "offer_template_key", "player_id").
				Record(offer).
				Returning("id", "claimed_at", "created_at", "updated_at").
				QueryStruct(offer)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).NotTo(Equal(""))

			var offer2 models.Offer
			err = db.
				Select("*").
				From("offers").
				Where("id = $1", offer.ID).
				QueryStruct(&offer2)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer2.ID).To(Equal(offer.ID))
		})
	})

	Describe("Get offer by id", func() {
		It("Should load offer by id", func() {
			//Given

			//When
			offer, err := models.GetOfferByID(db, "offers-game", defaultOfferID, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).To(Equal(defaultOfferID))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID).To(Equal(defaultOfferTemplateID))
			Expect(offer.OfferTemplateKey).To(Equal(defaultOfferTemplateKey))
			Expect(offer.CreatedAt.Valid).To(BeTrue())
		})

		It("Should return error if offer not found", func() {
			//Given
			offerID := uuid.NewV4().String()
			expectedError := errors.NewModelNotFoundError("Offer", map[string]interface{}{
				"GameID": "offers-game",
				"ID":     offerID,
			})

			//When
			offer, err := models.GetOfferByID(db, "offers-game", offerID, nil)

			//Then
			Expect(offer.ID).To(Equal(""))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Insert offer", func() {
		It("should insert offer with new id", func() {
			//Given
			gameID := "offers-game"
			offer := &models.Offer{
				GameID:           gameID,
				OfferTemplateID:  defaultOfferTemplateID,
				OfferTemplateKey: defaultOfferTemplateKey,
				PlayerID:         "player-3",
			}

			//When
			err := models.InsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).NotTo(BeEmpty())
		})

		It("should fail if game does not exist", func() {
			//Given
			offer := &models.Offer{
				GameID:           "non-existing-game",
				OfferTemplateID:  defaultOfferTemplateID,
				OfferTemplateKey: defaultOfferTemplateKey,
				PlayerID:         "player-3",
			}
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert on table \"offers\" violates constraint \"offer_templates_key\" \"offer_templates_id\" \"offer_templated_game_id\"",
			)

			//When
			err := models.InsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))

			//Test that after error our connection is still usable
			//Must use CONN and not db here to skip transaction
			dbOffer, err := models.GetOfferByID(conn, "offers-game", defaultOfferID, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.ID).To(Equal(defaultOfferID))
		})

		It("should fail if some error in the database", func() {
			gameID := "offers-game"
			offer := &models.Offer{
				GameID:          gameID,
				OfferTemplateID: defaultOfferTemplateID,
				PlayerID:        "player-3",
			}
			oldDB := db
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable
			err = models.InsertOffer(db, offer, time.Now(), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
			db = oldDB // avoid errors in after each
		})

		It("should fail if no offer template has specified id", func() {
			//Given
			offer := &models.Offer{
				GameID:           "offers-game",
				OfferTemplateID:  uuid.NewV4().String(),
				OfferTemplateKey: defaultOfferTemplateKey,
				PlayerID:         "player-3",
			}
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert on table \"offers\" violates constraint \"offer_templates_key\" \"offer_templates_id\" \"offer_templated_game_id\"",
			)

			//When
			err := models.InsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})

		It("should fail if no offer template has specified key", func() {
			//Given
			offer := &models.Offer{
				GameID:           "offers-game-2",
				OfferTemplateID:  defaultOfferTemplateID,
				OfferTemplateKey: uuid.NewV4().String(),
				PlayerID:         "player-3",
			}
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert on table \"offers\" violates constraint \"offer_templates_key\" \"offer_templates_id\" \"offer_templated_game_id\"",
			)

			//When
			err := models.InsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Claim offer", func() {
		var from, to int64 = 1486678000, 148669000
		It("should claim valid offer", func() {
			//Given
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(from+500, 0)

			//When
			contents, alreadyClaimed, err := models.ClaimOffer(db, id, playerID, gameID, currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			claimedOffer, err := models.GetOfferByID(db, "offers-game", defaultOfferID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(claimedOffer.ClaimedAt.Valid).To(BeTrue())
			Expect(claimedOffer.ClaimedAt.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim valid offer before trigger begins", func() {
			//Given
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			currentTime := time.Unix(from-500, 0)
			playerID := "player-1"
			gameID := "offers-game"

			//When
			contents, alreadyClaimed, err := models.ClaimOffer(db, id, playerID, gameID, currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should claim valid offer after trigger begins", func() {
			//Given
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			currentTime := time.Unix(to+500, 0)
			playerID := "player-1"
			gameID := "offers-game"

			//When
			contents, alreadyClaimed, err := models.ClaimOffer(db, id, playerID, gameID, currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not claim twice the same offer", func() {
			//Given
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			playerID := "player-1"
			gameID := "offers-game"
			firstTime := time.Unix(to+500, 0)
			secondTime := time.Unix(to+1000, 0)

			//When
			contents1, alreadyClaimed1, err1 := models.ClaimOffer(db, id, playerID, gameID, firstTime, nil)
			contents2, alreadyClaimed2, err2 := models.ClaimOffer(db, id, playerID, gameID, secondTime, nil)

			//Then
			Expect(contents1).NotTo(BeNil())
			Expect(alreadyClaimed1).To(BeFalse())
			Expect(err1).NotTo(HaveOccurred())

			Expect(contents2).NotTo(BeNil())
			Expect(alreadyClaimed2).To(BeTrue())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("should not claim an offer that doesn't exist", func() {
			//Given
			id := uuid.NewV4().String()
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(to+500, 0)

			//When
			_, _, err := models.ClaimOffer(db, id, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})

		It("should fail if some error in the database", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(from+500, 0)

			oldDB := db
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable

			_, _, err = models.ClaimOffer(db, id, playerID, gameID, currentTime, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
			db = oldDB // avoid errors in after each
		})
	})

	Describe("Update offer last seen at", func() {
		It("should update last seen offer at now and increment seen counter", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Now()

			//When
			offerBefore, err1 := models.GetOfferByID(db, "offers-game", defaultOfferID, nil)
			err2 := models.UpdateOfferLastSeenAt(db, defaultOfferID, playerID, gameID, currentTime, nil)
			offerAfter, err3 := models.GetOfferByID(db, "offers-game", defaultOfferID, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(offerAfter.LastSeenAt.Time.Unix()).To(Equal(currentTime.Unix()))
			Expect(offerAfter.LastSeenAt.Valid).To(BeTrue())
			Expect(offerBefore.SeenCounter).To(Equal(0))
			Expect(offerAfter.SeenCounter).To(Equal(1))
		})

		It("should return status code 422 if invalid id", func() {
			//Given
			id := uuid.NewV4().String()
			playerID := "player-1"
			gameID := "offers-game"

			//When
			err := models.UpdateOfferLastSeenAt(db, id, playerID, gameID, time.Now(), nil)

			//Then
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Get available offers", func() {
		It("should return a list of offer templates for each available placement", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(3))
			Expect(templates).To(HaveKey("popup"))
			Expect(templates["popup"]).To(HaveLen(1))
			Expect(templates["popup"][0].ID).To(Equal("56fc0477-39f1-485c-898e-4909e9155eb1"))
			Expect(templates["popup"][0].ProductID).To(Equal("com.tfg.sample"))
			Expect(templates["popup"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(templates["popup"][0].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(templates["popup"][0].RemainingPurchases).To(Equal(0))
			Expect(templates["popup"][0].RemainingImpressions).To(Equal(0))
			Expect(templates["popup"][0].ExpireAt).To(Equal(int64(1486679000)))

			Expect(templates).To(HaveKey("store"))
			Expect(templates["store"]).To(HaveLen(2))
			Expect(templates["store"][0].ID).NotTo(BeNil())
			Expect(templates["store"][0].ProductID).To(Equal("com.tfg.sample.2"))
			Expect(templates["store"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 100, "gold": 5}`))))
			Expect(templates["store"][0].Metadata).To(Equal(dat.JSON([]byte(`{"meta": "data"}`))))
			Expect(templates["store"][0].RemainingPurchases).To(Equal(0))
			Expect(templates["store"][0].RemainingImpressions).To(Equal(0))
			Expect(templates["store"][0].ExpireAt).To(Equal(int64(1486679200)))

			Expect(templates["store"][1].ID).To(Equal("6c4a79f2-24b8-4be9-93d4-12413b789823"))
			Expect(templates["store"][1].ProductID).To(Equal("com.tfg.sample.3"))
			Expect(templates["store"][1].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(templates["store"][1].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(templates["store"][1].RemainingPurchases).To(Equal(1))
			Expect(templates["store"][1].RemainingImpressions).To(Equal(1))
			Expect(templates["store"][1].ExpireAt).To(Equal(int64(1486679100)))
		})

		It("should return offers for two different players of game offers-game", func() {
			//Given
			playerID1 := "player-1"
			playerID2 := "player-2"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates1, err1 := models.GetAvailableOffers(db, playerID1, gameID, currentTime, nil)
			templates2, err2 := models.GetAvailableOffers(db, playerID2, gameID, currentTime, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(templates1).To(HaveLen(3))
			Expect(templates2).To(HaveLen(3))
		})

		It("should return empty offer template list if gameID doesn't exist", func() {
			//Given
			playerID := "player-1"
			gameID := "non-existing-game"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(templates).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not return offer-template-1 if last_seen_at is not long ago", func() {
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)

			//When
			err1 := models.UpdateOfferLastSeenAt(db, defaultOfferID, playerID, gameID, currentTime, nil)
			templates, err2 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(2))
			Expect(templates).To(HaveKey("store"))
		})

		It("should not return limited-template more than once", func() {
			//Given
			playerID := "player-1"
			gameID := "limited-offers-game"
			offerID := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			currentTime := time.Unix(1486678000, 0)

			//When
			templatesBefore, err1 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			err2 := models.UpdateOfferLastSeenAt(db, offerID, playerID, gameID, currentTime, nil)
			templatesAfter, err3 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(templatesBefore).To(HaveLen(1))
			Expect(templatesAfter).To(HaveLen(0))
		})

		It("should not return limited-template after claim", func() {
			//Given
			playerID := "player-1"
			gameID := "limited-offers-game"
			offerID := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			currentTime := time.Unix(1486678000, 0)

			//When
			templatesBefore, err1 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			_, alreadyClaimed, err2 := models.ClaimOffer(db, offerID, playerID, gameID, currentTime, nil)
			templatesAfter, err3 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(templatesBefore).To(HaveLen(1))
			Expect(templatesAfter).To(HaveLen(0))
		})

		It("should not return offer-template-1 if last_seen_at is not long ago", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game"
			offerID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			currentTime := time.Unix(1486678000, 0)

			//When
			templatesBefore, err1 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			_, alreadyClaimed, err2 := models.ClaimOffer(db, offerID, playerID, gameID, currentTime, nil)
			templatesAfter, err3 := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(templatesBefore).To(HaveLen(3))
			Expect(templatesAfter).To(HaveLen(2))
		})

		It("should not return template if it has empty trigger", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-empty-trigger"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(BeEmpty())
		})

		It("should not return template if it reached max frequency", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-max-freq"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(BeEmpty())
		})

		It("should not return template if it reached max period", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-max-period"
			currentTime := time.Unix(1486678000, 0)

			//When
			templates, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(BeEmpty())
		})

		It("should fail if template has invalid frequency", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-invalid-every-freq"
			currentTime := time.Unix(1486678000, 0)

			//When
			_, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("time: invalid duration invalid"))
		})

		It("should fail if some error in the database", func() {
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)

			oldDB := db
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable

			_, err = models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
			db = oldDB // avoid errors in after each
		})

		It("should fail if template has invalid frequency", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-invalid-every-period"
			currentTime := time.Unix(1486678000, 0)

			//When
			_, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("time: invalid duration invalid"))
		})
	})

	Describe("Claim and GetAvailableOffers integrated", func() {
		It("should not return consumed offer after it has been updated", func() {
			offerTemplateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)

			// Get fot the first time
			offers, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())

			// Disable offer template
			err = models.SetEnabledOfferTemplate(db, offerTemplateID, false, nil)
			Expect(err).NotTo(HaveOccurred())

			// Claim the offer
			_, alreadyClaimed, err := models.ClaimOffer(db, offers["popup"][0].ID, playerID, gameID, currentTime, nil)
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			// Get offer template to update it
			offerTemplate, err := models.GetOfferTemplateByID(db, offerTemplateID, nil)
			Expect(err).NotTo(HaveOccurred())

			// Update its contents and insert with same key
			offerTemplate.Contents = dat.JSON([]byte(`{ "somethingNew": 100 }`))
			offerTemplate, err = models.InsertOfferTemplate(db, offerTemplate, nil)
			Expect(err).NotTo(HaveOccurred())

			// Should not return the popup offer, since it was claimed for the first time
			offers, err = models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).NotTo(HaveKey("popup"))
		})

		It("should return updated offer with one remaining view", func() {
			playerID := "player-1"
			gameID := "offers-game"
			place := "unique-place"
			offerTemplateID := "5fed76ab-1fd7-4a91-972d-bca228ce80c4"
			currentTime := time.Unix(1486678000, 0)

			// Get offer
			offers, err := models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			offerID := offers[place][0].ID

			// Sees once
			err = models.UpdateOfferLastSeenAt(db, offerID, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())

			// Disable offer template
			err = models.SetEnabledOfferTemplate(db, offerTemplateID, false, nil)
			Expect(err).NotTo(HaveOccurred())

			// Update Offer template
			ot, err := models.GetOfferTemplateByID(db, offerTemplateID, nil)
			Expect(err).NotTo(HaveOccurred())
			ot.Contents = dat.JSON([]byte(`{ "somethingNew": 100 }`))
			ot, err = models.InsertOfferTemplate(db, ot, nil)
			Expect(err).NotTo(HaveOccurred())

			// Get offer
			offers, err = models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			offerID = offers[place][0].ID

			// Sees twice
			err = models.UpdateOfferLastSeenAt(db, offerID, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())

			// Get offer, expect unique-place to not be returned
			offers, err = models.GetAvailableOffers(db, playerID, gameID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).NotTo(HaveKey(place))
		})
	})
})
