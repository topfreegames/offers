// offers api
// https://github.com/topfreeoffers/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
	"time"
)

var _ = Describe("Offers Model", func() {
	var defaultOfferTemplateID string
	defaultOfferTemplateID = "dd21ec96-2890-4ba0-b8e2-40ea67196990"

	Describe("Offer Instance", func() {
		It("Shoud load a offer", func() {
			//Given
			offerID := "56fc0477-39f1-485c-898e-4909e9155eb1"

			//When
			var offer models.Offer
			err := db.
				Select("*").
				From("offers").
				Where("id = $1", offerID).
				QueryStruct(&offer)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).To(Equal(offerID))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID).To(Equal(defaultOfferTemplateID))
			Expect(offer.CreatedAt.Valid).To(BeTrue())
		})

		It("Should create offer", func() {
			//Given
			offer := &models.Offer{
				GameID:          "offers-game",
				OfferTemplateID: defaultOfferTemplateID,
				PlayerID:        "player-3",
			}

			//When
			err := db.
				InsertInto("offers").
				Columns("game_id", "offer_template_id", "player_id").
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

			Expect(offer2.ID).To(Equal(offer.ID))
		})
	})

	Describe("Get offer by id", func() {
		It("Should load offer by id", func() {
			//Given
			offerID := "56fc0477-39f1-485c-898e-4909e9155eb1"

			//When
			offer, err := models.GetOfferByID(db, "offers-game", offerID, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).To(Equal(offerID))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID).To(Equal(defaultOfferTemplateID))
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
			Expect(offer).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Get player's seen offers", func() {
		var enabledOfferTemplates []*models.OfferTemplate
		const playerID = "player-seen-offers"
		var offer *models.Offer

		BeforeEach(func() {
			//Given
			var err error
			enabledOfferTemplates, err = models.GetEnabledOfferTemplates(db, "offers-game", nil)
			Expect(err).NotTo(HaveOccurred())

			offer = &models.Offer{
				GameID:          "offers-game",
				OfferTemplateID: enabledOfferTemplates[0].ID,
				PlayerID:        playerID,
			}
			err = models.UpsertOffer(db, offer, time.Now(), nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should get all offers the player has seen and are enabled", func() {
			//When
			offers, err := models.GetPlayerSeenOffers(db, "offers-game", playerID, enabledOfferTemplates, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(1))
			Expect(offers[0].OfferTemplateID).To(Equal(enabledOfferTemplates[0].ID))
		})

		It("Should return empty list if invalid game", func() {
			//When
			offers, err := models.GetPlayerSeenOffers(db, "invalid-game", playerID, enabledOfferTemplates, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(0))
		})

	})

	Describe("Upsert offer", func() {
		It("should insert offer with new id", func() {
			//Given
			offer := &models.Offer{
				GameID:          "offers-game",
				OfferTemplateID: defaultOfferTemplateID,
				PlayerID:        "player-3",
			}

			//When
			err := models.UpsertOffer(db, offer, time.Now(), nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			offerFromDB, err := models.GetOfferByID(db, "offers-game", offer.ID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerFromDB.ID).To(Equal(offer.ID))
		})

		It("should update offer with existing id", func() {
			//Given
			offerID := "35df52e7-3161-446f-975b-92f32871e37c"
			offer := &models.Offer{
				ID:              offerID,
				GameID:          "offers-game-2",
				OfferTemplateID: defaultOfferTemplateID,
				PlayerID:        "player-4",
			}

			//When
			err := models.UpsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).NotTo(HaveOccurred())

			offerFromDB, err := models.GetOfferByID(db, "offers-game-2", offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerFromDB.ID).To(Equal(offerID))
			Expect(offerFromDB.GameID).To(Equal("offers-game-2"))
			Expect(offerFromDB.PlayerID).To(Equal("player-4"))
			Expect(offerFromDB.OfferTemplateID).To(Equal(defaultOfferTemplateID))
		})

		It("should fail if game does not exist", func() {
			//Given
			offer := &models.Offer{
				GameID:          "invalid-offers-game",
				OfferTemplateID: defaultOfferTemplateID,
				PlayerID:        "player-3",
			}
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert or update on table \"offers\" violates foreign key constraint \"offers_game_id_fkey\"",
			)

			//When
			err := models.UpsertOffer(db, offer, time.Now(), nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))

			//Test that after error our connection is still usable
			offerID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			//Must use CONN and not db here to skip transaction
			dbOffer, err := models.GetOfferByID(conn, "offers-game", offerID, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.ID).To(Equal(offerID))
		})

		It("should fail if offer template does not exist", func() {
			//Given
			offer := &models.Offer{
				GameID:          "offers-game-2",
				OfferTemplateID: uuid.NewV4().String(),
				PlayerID:        "player-3",
			}
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert or update on table \"offers\" violates foreign key constraint \"offers_offer_template_id_fkey\"",
			)

			//When
			err := models.UpsertOffer(db, offer, time.Now(), nil)

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
	})
})
