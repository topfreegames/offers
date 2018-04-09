// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	edat "github.com/topfreegames/extensions/dat"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
	oTesting "github.com/topfreegames/offers/testing"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

const (
	defaultOfferID string        = "dd21ec96-2890-4ba0-b8e2-40ea67196990"
	defaultGameID  string        = "offers-game"
	expireDuration time.Duration = 300 * time.Second
)

var _ = Describe("Offer Models", func() {
	Describe("Get offer id", func() {
		It("should load an offer from existent id", func() {
			id := defaultOfferID
			gameID := defaultGameID

			offer, err := models.GetOfferByID(nil, db, gameID, id, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).To(Equal(id))
			Expect(offer.Period).To(Equal(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(offer.Frequency).To(BeEquivalentTo(dat.JSON([]byte(`{"every": "1s"}`))))
			Expect(offer.Version).To(Equal(1))
		})

		It("should not load an offer from nonexistent ID", func() {
			id := uuid.NewV4().String()
			gameID := defaultGameID

			_, err := models.GetOfferByID(nil, db, gameID, id, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Offer was not found with specified filters."))
		})

		It("should not load an offer from nonexistent GameID", func() {
			id := uuid.NewV4().String()
			gameID := defaultGameID

			_, err := models.GetOfferByID(nil, db, gameID, id, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Offer was not found with specified filters."))
		})
	})

	Describe("Insert Offer", func() {
		It("should create an offer with valid parameters", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			offer, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).NotTo(Equal(""))
			Expect(offer.Enabled).To(BeTrue())
			Expect(offer.Version).To(Equal(1))
		})

		It("should succeed and reset offers cache", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			enabledOffersKey := models.GetEnabledOffersKey(offer.GameID)
			offersCache.Set(enabledOffersKey, []*models.Offer{offer}, time.Minute)

			_, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeFalse())
		})

		It("should create an offer with valid parameters, including filters and metadata", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
				Metadata:  dat.JSON([]byte(`{"Awesome": "offer"}`)),
				Filters:   dat.JSON([]byte(`{"level": {"geq": 1.0, "lt": 3.0}}`)),
			}

			offer, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).NotTo(Equal(""))
			Expect(offer.Enabled).To(BeTrue())
			Expect(offer.Version).To(Equal(1))
		})

		It("should create an offer with valid parameters, with cost instead of productId", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				Cost:      dat.JSON([]byte(`{"gems": 500}`)),
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			offer, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(offer.ID).NotTo(Equal(""))
			Expect(offer.Enabled).To(BeTrue())
			Expect(offer.Version).To(Equal(1))
		})

		It("should return error if game with given id does not exist", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "non-existing-game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			offer, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`Offer could not be saved due to: insert or update on table "offers" violates foreign key constraint "offers_game_id_fkey"`))

			err = conn.
				Select("id").
				From("offers").
				Where("game_id = $1", "non-existing-game-id").
				QueryStruct(&offer)
			Expect(offer.ID).To(BeEmpty())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: no rows in result set"))
		})

		It("should return error if inserting offer template with missing parameters", func() {
			//Given
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
			}

			//When
			_, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`pq: null value in column "period" violates not-null constraint`))
		})

		It("should return error if DB is closed", func() {
			oldDB := db
			defer func() {
				db = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}

			_, err = models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
		})

		It("should fail and not reset offers cache", func() {
			//Given
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
			}
			enabledOffersKey := models.GetEnabledOffersKey(offer.GameID)
			offersCache.Set(enabledOffersKey, "cached-stuff", time.Minute)

			//When
			_, err := models.InsertOffer(nil, db, offer, offersCache, nil)

			//Then
			Expect(err).To(HaveOccurred())

			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeTrue())
		})
	})

	Describe("List offers", func() {
		It("Should return the full list of offers for a game", func() {
			var limit uint64 = 5
			var offset uint64 = 0
			games, pages, err := models.ListOffers(nil, db, "offers-game", limit, offset, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(5))
			Expect(pages).To(Equal(1))
		})

		It("should return empty list if non-existing game id", func() {
			var limit uint64 = 5
			var offset uint64 = 0
			games, pages, err := models.ListOffers(nil, db, "non-existing-game", limit, offset, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(0))
			Expect(pages).To(Equal(0))
		})

		It("should return two offers with limit 2 and offset 0", func() {
			var limit uint64 = 2
			var offset uint64 = 0
			var pages int = 5/2 + 1
			games, pages, err := models.ListOffers(nil, db, "offers-game", limit, offset, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(2))
			Expect(pages).To(Equal(pages))
		})

		It("should return no offer with limit 2 and offset 3", func() {
			var limit uint64 = 2
			var offset uint64 = 3
			var pages int = 5/2 + 1
			games, pages, err := models.ListOffers(nil, db, "offers-game", limit, offset, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(BeEmpty())
			Expect(pages).To(Equal(pages))
		})

		It("should return 0 pages if limit is 0", func() {
			var limit uint64 = 0
			var offset uint64 = 0
			var pages int = 0
			games, pages, err := models.ListOffers(nil, db, "offers-game", limit, offset, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(0))
			Expect(pages).To(Equal(pages))
		})

		It("should return error if db isn't connected", func() {
			db, err := oTesting.GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			err = db.(*runner.DB).DB.Close()
			Expect(err).NotTo(HaveOccurred())

			var limit uint64 = 10
			var offset uint64 = 0
			_, _, err = models.ListOffers(nil, db, "offers-game", limit, offset, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Get enabled offers", func() {
		It("should get the enabled offers for the given game", func() {
			expectedIDs := []string{
				"dd21ec96-2890-4ba0-b8e2-40ea67196990",
				"d5114990-77d7-45c4-ba5f-462fc86b213f",
				"a411fbcf-dddc-4153-b42b-3f9b2684c965",
				"5fed76ab-1fd7-4a91-972d-bca228ce80c4",
			}
			gameID := defaultGameID

			filterAttrs := make(map[string]string)
			offers, err := models.GetEnabledOffers(nil, db, gameID, offersCache, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(4))
			for i := 0; i < len(offers); i++ {
				Expect(expectedIDs).To(ContainElement(offers[i].ID))
			}
		})

		It("should return an empty list if there are no enabled offers", func() {
			gameID := uuid.NewV4().String()
			filterAttrs := make(map[string]string)
			offers, err := models.GetEnabledOffers(nil, db, gameID, offersCache, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(0))
		})

		It("should get enabled offers from cache on the second time", func() {
			expectedIDs := []string{
				"dd21ec96-2890-4ba0-b8e2-40ea67196990",
				"d5114990-77d7-45c4-ba5f-462fc86b213f",
				"a411fbcf-dddc-4153-b42b-3f9b2684c965",
				"5fed76ab-1fd7-4a91-972d-bca228ce80c4",
			}
			gameID := defaultGameID
			start := time.Now().UnixNano()
			filterAttrs := make(map[string]string)
			offers, err := models.GetEnabledOffers(nil, db, gameID, offersCache, expireDuration, filterAttrs, false, nil)
			dbElapsedTime := time.Now().UnixNano() - start
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(4))

			start = time.Now().UnixNano()
			offers, err = models.GetEnabledOffers(nil, db, gameID, offersCache, expireDuration, filterAttrs, false, nil)
			cacheElapsedTime := time.Now().UnixNano() - start
			Expect(err).NotTo(HaveOccurred())
			_, found := offersCache.Get(models.GetEnabledOffersKey(gameID))
			Expect(found).To(BeTrue())

			Expect(dbElapsedTime).To(BeNumerically(">", cacheElapsedTime))
			Expect(offers).To(HaveLen(4))
			for i := 0; i < len(offers); i++ {
				Expect(expectedIDs).To(ContainElement(offers[i].ID))
			}
		})
	})

	Describe("Set enabled offer template", func() {
		It("should disable an enabled offer", func() {
			//Given
			offerID := defaultOfferID
			gameID := defaultGameID
			enabled := false
			var offer models.Offer

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			err = db.SQL("SELECT enabled FROM offers WHERE game_id=$1 AND id=$2", gameID, offerID).QueryStruct(&offer)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(offer.Enabled).To(BeFalse())
		})

		It("should enable an enabled offer", func() {
			//Given
			offerID := defaultOfferID
			gameID := defaultGameID
			enabled := true
			var offer models.Offer

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			err = db.SQL("SELECT enabled FROM offers WHERE game_id=$1 AND id=$2", gameID, offerID).QueryStruct(&offer)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(offer.Enabled).To(BeTrue())
		})

		It("should enable a disabled offer", func() {
			//Given
			offerID := "27b0370f-bd61-4346-a10d-50ec052ae125"
			gameID := defaultGameID
			enabled := true
			var offer models.Offer

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			err = db.SQL("SELECT enabled FROM offers WHERE game_id=$1 AND id=$2", gameID, offerID).QueryStruct(&offer)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(offer.Enabled).To(BeTrue())
		})

		It("should succeed and reset offers cache", func() {
			//Given
			offerID := defaultOfferID
			gameID := defaultGameID
			enabled := true

			enabledOffersKey := models.GetEnabledOffersKey(gameID)
			offersCache.Set(enabledOffersKey, "something", time.Minute)

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeFalse())
		})

		It("should return error if id doesn't exist", func() {
			//Given
			offerID := uuid.NewV4().String()
			gameID := defaultGameID
			enabled := true

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})

		It("should return error if game id doesn't exist", func() {
			//Given
			offerID := defaultOfferID
			gameID := "non-existing-game-id"
			enabled := true

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)

			//Then
			Expect(err).To(HaveOccurred())
		})

		It("should fail and not reset offers cache", func() {
			//Given
			offerID := uuid.NewV4().String()
			gameID := defaultGameID
			enabled := true

			enabledOffersKey := models.GetEnabledOffersKey(gameID)
			offersCache.Set(enabledOffersKey, "cached-stuff", time.Minute)

			//When
			err := models.SetEnabledOffer(nil, db, gameID, offerID, enabled, offersCache, nil)

			//Then
			Expect(err).To(HaveOccurred())

			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeTrue())
		})
	})

	Describe("Update Offer", func() {
		It("should update the offer and increment the version with valid parameters", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}

			updatedOffer, err := models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedOffer.ID).To(Equal(offerUpdate.ID))
			Expect(updatedOffer.Version).To(Equal(createdOffer.Version + 1))

			var dbOffer models.Offer
			builder := db.Select("*")
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.From("offers").Where("id=$1", offerUpdate.ID).QueryStruct(&dbOffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.GameID).To(Equal(offerUpdate.GameID))
			Expect(dbOffer.Name).To(Equal(offerUpdate.Name))
			Expect(dbOffer.ProductID).To(Equal(offerUpdate.ProductID))
			Expect(dbOffer.Contents).To(Equal(offerUpdate.Contents))
			Expect(dbOffer.Period).To(Equal(offerUpdate.Period))
			Expect(dbOffer.Frequency).To(Equal(offerUpdate.Frequency))
			Expect(dbOffer.Trigger).To(Equal(offerUpdate.Trigger))
			Expect(dbOffer.Placement).To(Equal(offerUpdate.Placement))
			Expect(dbOffer.Version).To(Equal(createdOffer.Version + 1))
		})

		It("should update the offer and increment the version with valid parameters, including filters and metadata", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Filters:   dat.JSON([]byte(`{"level": {"eq": "5"}}`)),
				Metadata:  dat.JSON([]byte(`{"Cool": "offer"}`)),
				Placement: "store",
			}

			updatedOffer, err := models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedOffer.ID).To(Equal(offerUpdate.ID))
			Expect(updatedOffer.Version).To(Equal(createdOffer.Version + 1))

			var dbOffer models.Offer
			builder := db.Select("*")
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.From("offers").Where("id=$1", offerUpdate.ID).QueryStruct(&dbOffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.GameID).To(Equal(offerUpdate.GameID))
			Expect(dbOffer.Name).To(Equal(offerUpdate.Name))
			Expect(dbOffer.ProductID).To(Equal(offerUpdate.ProductID))
			Expect(dbOffer.Contents).To(Equal(offerUpdate.Contents))
			Expect(dbOffer.Period).To(Equal(offerUpdate.Period))
			Expect(dbOffer.Frequency).To(Equal(offerUpdate.Frequency))
			Expect(dbOffer.Trigger).To(Equal(offerUpdate.Trigger))
			Expect(dbOffer.Placement).To(Equal(offerUpdate.Placement))
			Expect(dbOffer.Filters).To(Equal(offerUpdate.Filters))
			Expect(dbOffer.Metadata).To(Equal(offerUpdate.Metadata))
			Expect(dbOffer.Version).To(Equal(createdOffer.Version + 1))
		})

		It("should update the offer and increment the version with valid parameters, including cost", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "game-id",
				Name:      "offer-2",
				Cost:      dat.JSON([]byte(`{"gems": 500}`)),
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Filters:   dat.JSON([]byte(`{"level": {"eq": "5"}}`)),
				Metadata:  dat.JSON([]byte(`{"Cool": "offer"}`)),
				Placement: "store",
			}

			updatedOffer, err := models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedOffer.ID).To(Equal(offerUpdate.ID))
			Expect(updatedOffer.Version).To(Equal(createdOffer.Version + 1))

			var dbOffer models.Offer
			builder := db.Select("*")
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.From("offers").Where("id=$1", offerUpdate.ID).QueryStruct(&dbOffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.ProductID).To(BeEmpty())
			Expect(dbOffer.Cost).To(Equal(offerUpdate.Cost))
			Expect(dbOffer.Version).To(Equal(createdOffer.Version + 1))
		})

		It("should succeed and reset offers cache", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			enabledOffersKey := models.GetEnabledOffersKey(offer.GameID)
			offersCache.Set(enabledOffersKey, []*models.Offer{offer}, time.Minute)

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}

			_, err = models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeFalse())
		})

		It("should return error if game with given id does not exist", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "non-existing-game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}

			_, err = models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Offer was not found with specified filters."))

			var dbOffer models.Offer
			builder := db.Select("*")
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.From("offers").Where("id=$1", offerUpdate.ID).QueryStruct(&dbOffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbOffer.Version).To(Equal(createdOffer.Version))
		})

		It("should return error if offer with given id does not exist", func() {
			id := uuid.NewV4().String()
			offerUpdate := &models.Offer{
				ID:        id,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}

			_, err := models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Offer was not found with specified filters."))
		})

		It("should return error if inserting offer template with missing parameters", func() {
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
			}

			_, err = models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`pq: null value in column "period" violates not-null constraint`))
		})

		It("should return error if DB is closed", func() {
			oldDB := db
			defer func() {
				db = oldDB // avoid errors in after each
			}()
			offer := &models.Offer{
				Name:      "offer-1",
				ProductID: "com.tfg.example",
				GameID:    "game-id",
				Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				Period:    dat.JSON([]byte(`{"every": "10m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "24h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1487280506875}`)),
				Placement: "popup",
			}
			createdOffer, err := models.InsertOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable
			offerUpdate := &models.Offer{
				ID:        createdOffer.ID,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}

			_, err = models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
		})

		It("should fail and not reset offers cache", func() {

			id := uuid.NewV4().String()
			offerUpdate := &models.Offer{
				ID:        id,
				GameID:    "game-id",
				Name:      "offer-2",
				ProductID: "com.tfg.example2",
				Contents:  dat.JSON([]byte(`{"gems": 5}`)),
				Period:    dat.JSON([]byte(`{"every": "1m"}`)),
				Frequency: dat.JSON([]byte(`{"every": "2h"}`)),
				Trigger:   dat.JSON([]byte(`{"from": 1111111111111}`)),
				Placement: "store",
			}
			enabledOffersKey := models.GetEnabledOffersKey(offerUpdate.GameID)
			offersCache.Set(enabledOffersKey, "cached-stuff", time.Minute)

			_, err := models.UpdateOffer(nil, db, offerUpdate, offersCache, nil)
			Expect(err).To(HaveOccurred())

			_, found := offersCache.Get(enabledOffersKey)
			Expect(found).To(BeTrue())
		})
	})
})
