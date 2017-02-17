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
)

var _ = Describe("Offers Model", func() {
	Describe("Offer Instance", func() {
		It("Shoud load a offer", func() {
			offerID, _ := uuid.FromString("56fc0477-39f1-485c-898e-4909e9155eb1")
			offerTemplateID, _ := uuid.FromString("4118141e-1d20-4839-8ce8-ead92b298a86")

			var offer models.Offer
			err := db.
				Select("*").
				From("offers").
				Where("id = $1", offerID).
				QueryStruct(&offer)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID.String()).To(Equal(offerID.String()))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID.String()).To(Equal(offerTemplateID.String()))
			Expect(offer.CreatedAt.Valid).To(BeTrue())
		})

		It("Should create offer", func() {
			offer := &models.Offer{
				GameID:          "offers-game",
				OfferTemplateID: uuid.NewV4(),
				PlayerID:        "player-3",
			}
			err := db.
				InsertInto("offers").
				Columns("game_id", "offer_template_id", "player_id").
				Record(offer).
				Returning("id", "claimed_at", "created_at", "updated_at").
				QueryStruct(offer)

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
			offerID, _ := uuid.FromString("56fc0477-39f1-485c-898e-4909e9155eb1")
			offerTemplateID, _ := uuid.FromString("4118141e-1d20-4839-8ce8-ead92b298a86")
			offer, err := models.GetOfferByID(db, offerID, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID.String()).To(Equal(offerID.String()))
			Expect(offer.GameID).To(Equal("offers-game"))
			Expect(offer.PlayerID).To(Equal("player-1"))
			Expect(offer.OfferTemplateID.String()).To(Equal(offerTemplateID.String()))
			Expect(offer.CreatedAt.Valid).To(BeTrue())
		})

		It("Should return error if offer not found", func() {
			offerID := uuid.NewV4()
			expectedError := errors.NewModelNotFoundError("Offer", map[string]interface{}{
				"ID": offerID,
			})
			offer, err := models.GetOfferByID(db, offerID, nil)
			Expect(offer).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Upsert offer", func() {
		It("should insert offer with new id", func() {
			offer := &models.Offer{
				GameID:          "offers-game",
				OfferTemplateID: uuid.NewV4(),
				PlayerID:        "player-3",
			}

			err := models.UpsertOffer(db, offer, nil)
			Expect(err).NotTo(HaveOccurred())

			offerFromDB, err := models.GetOfferByID(db, offer.ID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerFromDB.ID).To(Equal(offer.ID))
		})

		It("should update offer with existing id", func() {
			offerID, _ := uuid.FromString("35df52e7-3161-446f-975b-92f32871e37c")
			offerTemplateID := uuid.NewV4()
			offer := &models.Offer{
				ID:              offerID,
				GameID:          "offers-game-2",
				OfferTemplateID: offerTemplateID,
				PlayerID:        "player-4",
			}
			err := models.UpsertOffer(db, offer, nil)
			Expect(err).NotTo(HaveOccurred())

			offerFromDB, err := models.GetOfferByID(db, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerFromDB.ID).To(Equal(offerID))
			Expect(offerFromDB.GameID).To(Equal("offers-game-2"))
			Expect(offerFromDB.PlayerID).To(Equal("player-4"))
			Expect(offerFromDB.OfferTemplateID.String()).To(Equal(offerTemplateID.String()))
		})

		It("should fail if game does not exist", func() {
			offer := &models.Offer{
				GameID:          "invalid-game",
				OfferTemplateID: uuid.NewV4(),
				PlayerID:        "player-3",
			}

			err := models.UpsertOffer(db, offer, nil)
			Expect(err).To(HaveOccurred())
			expectedError := errors.NewInvalidModelError(
				"Offer",
				"insert or update on table \"offers\" violates foreign key constraint \"offers_game_id_fkey\"",
			)
			Expect(err).To(MatchError(expectedError))

			//Test that after error our connection is still usable
			offerID, _ := uuid.FromString("56fc0477-39f1-485c-898e-4909e9155eb1")
			dbOffer, err := models.GetOfferByID(db, offerID, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.ID.String()).To(Equal(offerID.String()))
		})
	})
})
