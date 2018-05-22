// offers api
// https://github.com/topfreeoffers/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Offers <backend@tfgco.com>

package models_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	edat "github.com/topfreegames/extensions/dat"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var _ = Describe("Offer Instance Model", func() {
	const (
		defaultOfferInstanceID string        = "56fc0477-39f1-485c-898e-4909e9155eb1"
		defaultOfferID         string        = "dd21ec96-2890-4ba0-b8e2-40ea67196990"
		defaultGameID          string        = "offers-game"
		defaultPlayerID        string        = "player-1"
		defaultProductID       string        = "com.tfg.sample"
		filterGameID           string        = "another-game-with-filters"
		expireDuration         time.Duration = 300 * time.Second
	)

	Describe("Claim offer", func() {
		var from, to int64 = 1486678000, 148669000
		It("should claim valid offer", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim valid offer without the offer instance id", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "offers-game"
			id := ""
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim valid offer before trigger begins", func() {
			//Given
			currentTime := time.Unix(from-500, 0)
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim valid offer after trigger begins", func() {
			//Given
			currentTime := time.Unix(to+500, 0)
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim and receive 0 nextAt if reached max purchases", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "limited-offers-game"
			id := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			offerID := "aa65a3f2-7cf8-4d76-957f-0a23a1bbbd32"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(int64(0)))
			Expect(err).NotTo(HaveOccurred())

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim and receive 0 nextAt if offer was disabled", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "offers-game"
			id := "b0bffdd6-5cb8-4b54-b250-349b18c07638"
			offerID := "27b0370f-bd61-4346-a10d-50ec052ae125"
			playerID := "player-14"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(int64(0)))
			Expect(err).NotTo(HaveOccurred())

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim and receive nextAt equal to currentTime if no every frequency or period", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "another-game-v8"
			id := "0a90073e-a798-46a8-a4f2-6b32182672ff"
			offerID := "17e42a40-da28-44dc-abd1-0cef8c2dff42"
			playerID := "player-11"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(currentTime.Unix()))
			Expect(err).NotTo(HaveOccurred())

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should claim and receive biggest nextAt considering period and freq", func() {
			//Given
			currentTime := time.Unix(from+500, 0)
			gameID := "another-game"
			id := "29593e1d-792a-4849-8236-9d7b80fc6f6c"
			offerID := "a2539a8c-55f2-4539-a8c0-929b240d8c80"
			playerID := "player-11"
			transactionID := uuid.NewV4().String()

			//When
			contents, alreadyClaimed, nextAt, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(contents).NotTo(BeNil())
			Expect(alreadyClaimed).To(BeFalse())
			Expect(nextAt).To(Equal(currentTime.Unix() + 30))
			Expect(err).NotTo(HaveOccurred())

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
		})

		It("should not claim twice the same offer", func() {
			//Given
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			firstTime := time.Unix(to+500, 0)
			secondTime := time.Unix(to+1000, 0)

			//When
			contents1, alreadyClaimed1, nextAt1, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, firstTime.Unix(), firstTime, nil)
			Expect(err).NotTo(HaveOccurred())
			contents2, alreadyClaimed2, nextAt2, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, secondTime.Unix(), secondTime, nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(contents1).NotTo(BeNil())
			Expect(alreadyClaimed1).To(BeFalse())
			Expect(nextAt1).To(Equal(firstTime.Unix() + 1))

			Expect(contents2).NotTo(BeNil())
			Expect(alreadyClaimed2).To(BeTrue())
			Expect(nextAt2).To(Equal(int64(to + 501)))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(1))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(firstTime.Unix()))
		})

		It("should claim twice the same offer if different transactionIDs", func() {
			//Given
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			playerID := "player-1"
			transactionID1 := uuid.NewV4().String()
			transactionID2 := uuid.NewV4().String()

			firstTime := time.Unix(to+500, 0)
			secondTime := time.Unix(to+1000, 0)

			//When
			contents1, alreadyClaimed1, nextAt1, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID1, firstTime.Unix(), firstTime, nil)
			Expect(err).NotTo(HaveOccurred())
			contents2, alreadyClaimed2, nextAt2, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID2, secondTime.Unix(), secondTime, nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(contents1).NotTo(BeNil())
			Expect(alreadyClaimed1).To(BeFalse())
			Expect(nextAt1).To(Equal(firstTime.Unix() + 1))

			Expect(contents2).NotTo(BeNil())
			Expect(alreadyClaimed2).To(BeFalse())
			Expect(nextAt2).To(Equal(secondTime.Unix() + 1))

			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ClaimCounter).To(Equal(2))
			Expect(offerPlayer.ClaimTimestamp.Time.Unix()).To(Equal(secondTime.Unix()))
		})

		It("should not claim an offer that doesn't exist", func() {
			//Given
			currentTime := time.Unix(to+500, 0)
			gameID := "offers-game"
			id := uuid.NewV4().String()
			playerID := "player-1"
			transactionID := uuid.NewV4().String()
			//When
			_, _, _, err := models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferInstance was not found with specified filters."))
		})

		It("should fail if some error in the database", func() {
			currentTime := time.Unix(from+500, 0)
			gameID := "offers-game"
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			playerID := "player-1"
			transactionID := uuid.NewV4().String()

			oldDB := db
			defer func() {
				db = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable

			_, _, _, err = models.ClaimOffer(nil, db, gameID, id, playerID, defaultProductID, transactionID, currentTime.Unix(), currentTime, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
		})
	})

	Describe("View offer", func() {
		It("should view offer instance", func() {
			//Given
			playerID := "player-1"
			gameID := defaultGameID
			offerInstanceID := defaultOfferInstanceID
			offerID := defaultOfferID
			impressionID := uuid.NewV4().String()
			currentTime := time.Now()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			//Then
			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ViewCounter).To(Equal(1))
			Expect(offerPlayer.ViewTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))
		})

		It("should return 0 nextAt if offer reached max period", func() {
			//Given
			playerID := "player-1"
			gameID := "limited-offers-game"
			currentTime := time.Now()
			offerInstanceID := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			impressionID := uuid.NewV4().String()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAt).To(Equal(int64(0)))
		})

		It("should return 0 nextAt if offer was disabled", func() {
			//Given
			playerID := "player-14"
			gameID := "offers-game"
			currentTime := time.Now()
			offerInstanceID := "b0bffdd6-5cb8-4b54-b250-349b18c07638"
			impressionID := uuid.NewV4().String()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAt).To(Equal(int64(0)))
		})

		It("should return nextAt equal to now if offer has no every in frequency", func() {
			//Given
			playerID := "player-11"
			gameID := "offers-game"
			currentTime := time.Now()
			offerInstanceID := "4407b770-5b24-4ffa-8563-0694d1a10156"
			offerID := "5fed76ab-1fd7-4a91-972d-bca228ce80c4"
			impressionID := uuid.NewV4().String()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			//Then
			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ViewCounter).To(Equal(1))
			Expect(offerPlayer.ViewTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
			Expect(nextAt).To(Equal(currentTime.Unix()))
		})

		It("should return error if non-existing id", func() {
			//Given
			offerInstanceID := uuid.NewV4().String()
			playerID := "player-1"
			gameID := "offers-game"
			impressionID := uuid.NewV4().String()
			currentTime := time.Now()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferInstance was not found with specified filters."))
			Expect(nextAt).To(Equal(int64(0)))
		})

		It("should return error if non-existing game", func() {
			//Given
			offerInstanceID := uuid.NewV4().String()
			playerID := "player-1"
			gameID := "non-existing-offers-game"
			impressionID := uuid.NewV4().String()
			currentTime := time.Now()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferInstance was not found with specified filters."))
			Expect(nextAt).To(Equal(int64(0)))
		})

		It("should return error if non-existing game", func() {
			//Given
			offerInstanceID := uuid.NewV4().String()
			playerID := "player-1"
			gameID := "non-existing-offers-game"
			impressionID := uuid.NewV4().String()
			currentTime := time.Now()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferInstance was not found with specified filters."))
			Expect(nextAt).To(Equal(int64(0)))
		})

		It("should not increment counters if view is a retry", func() {
			//Given
			playerID := "player-1"
			gameID := defaultGameID
			offerInstanceID := defaultOfferInstanceID
			offerID := defaultOfferID
			impressionID := uuid.NewV4().String()
			currentTime := time.Now()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			isReplay, nextAt, err = models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())

			//Then
			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ViewCounter).To(Equal(1))
			Expect(offerPlayer.ViewTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))
		})

		It("should increment twice if different impressionIDs", func() {
			//Given
			playerID := "player-1"
			gameID := defaultGameID
			offerInstanceID := defaultOfferInstanceID
			offerID := defaultOfferID
			currentTime := time.Now()
			impressionID := uuid.NewV4().String()

			//When
			isReplay, nextAt, err := models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			impressionID = uuid.NewV4().String()
			isReplay, nextAt, err = models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, impressionID, currentTime, nil)
			Expect(isReplay).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			//Then
			offerPlayer, err := models.GetOfferPlayer(nil, db, gameID, playerID, offerID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerPlayer.ViewCounter).To(Equal(2))
			Expect(offerPlayer.ViewTimestamp.Time.Unix()).To(Equal(currentTime.Unix()))
			Expect(nextAt).To(Equal(currentTime.Unix() + 1))
		})
	})

	Describe("Get available offers", func() {
		It("should return a list of offer templates for each available placement", func() {
			//Given
			playerID := "player-9"
			gameID := defaultGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(3))
			Expect(offerInstances).To(HaveKey("popup"))
			Expect(offerInstances["popup"]).To(HaveLen(1))
			Expect(offerInstances["popup"][0].ID).To(Equal("56fc0477-39f1-485c-898e-4909e9155eb1"))
			Expect(offerInstances["popup"][0].ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstances["popup"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["popup"][0].Cost).To(Equal(dat.JSON([]byte(`{"gems": 500}`))))
			Expect(offerInstances["popup"][0].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["popup"][0].ExpireAt).To(Equal(int64(1486679000)))

			Expect(offerInstances).To(HaveKey("store"))
			Expect(offerInstances["store"]).To(HaveLen(2))
			i1, i2 := -1, -1
			if offerInstances["store"][0].ID == "35df52e7-3161-446f-975b-92f32871e37c" {
				i1 = 1
				i2 = 0
			} else {
				i1 = 0
				i2 = 1
			}
			Expect(offerInstances["store"][i1].ID).NotTo(BeNil())
			Expect(offerInstances["store"][i1].ProductID).To(Equal("com.tfg.sample.2"))
			Expect(offerInstances["store"][i1].Contents).To(Equal(dat.JSON([]byte(`{"gems": 100, "gold": 5}`))))
			Expect(offerInstances["store"][i1].Metadata).To(Equal(dat.JSON([]byte(`{"meta": "data"}`))))
			Expect(offerInstances["store"][i1].ExpireAt).To(Equal(int64(1486679200)))

			Expect(offerInstances["store"][i2].ID).To(Equal("35df52e7-3161-446f-975b-92f32871e37c"))
			Expect(offerInstances["store"][i2].ProductID).To(Equal("com.tfg.sample.3"))
			Expect(offerInstances["store"][i2].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["store"][i2].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["store"][i2].ExpireAt).To(Equal(int64(1486679100)))
		})

		It("should return a list of offer templates for each available placement when with interval filters and allow inneficient queries", func() {
			//Given
			playerID := "player-13"
			gameID := filterGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := map[string]string{
				"level": "1",
			}

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, true, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(2))
			Expect(offerInstances).To(HaveKey("popup"))
			Expect(offerInstances["popup"]).To(HaveLen(1))
			Expect(offerInstances["popup"][0].ID).To(Equal("33f9bbc1-5e9e-4a80-ae95-8d74d8774629"))
			Expect(offerInstances["popup"][0].ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstances["popup"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["popup"][0].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["popup"][0].ExpireAt).To(Equal(int64(1486679000)))
			Expect(offerInstances).To(HaveKey("popup"))
			Expect(offerInstances["store"]).To(HaveLen(1))
			Expect(offerInstances["store"][0].ID).To(Equal("0eebb309-753c-4736-98f1-5be851e1ac4d"))
			Expect(offerInstances["store"][0].ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstances["store"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["store"][0].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["store"][0].ExpireAt).To(Equal(int64(1486679000)))
		})

		It("should return a list of offer templates for each available placement when with interval filters and not allow inneficient queries", func() {
			//Given
			playerID := "player-13"
			gameID := filterGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := map[string]string{
				"level": "1",
			}

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(0))
		})

		It("should return a list of offer templates for each available placement when with filters not matching", func() {
			//Given
			playerID := "player-13"
			gameID := filterGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := map[string]string{
				"level": "3",
			}

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(0))
		})

		It("should return a list of offer templates for each available placement when with filters but offer has no filters and inneficient queries allowed", func() {
			//Given
			playerID := "player-9"
			gameID := defaultGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := map[string]string{
				"key1": "value1",
				"key2": "1.2",
				"key3": "1",
			}

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, true, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(3))
			Expect(offerInstances).To(HaveKey("popup"))
			Expect(offerInstances["popup"]).To(HaveLen(1))
			Expect(offerInstances["popup"][0].ID).To(Equal("56fc0477-39f1-485c-898e-4909e9155eb1"))
			Expect(offerInstances["popup"][0].ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstances["popup"][0].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["popup"][0].Cost).To(Equal(dat.JSON([]byte(`{"gems": 500}`))))
			Expect(offerInstances["popup"][0].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["popup"][0].ExpireAt).To(Equal(int64(1486679000)))

			Expect(offerInstances).To(HaveKey("store"))
			Expect(offerInstances["store"]).To(HaveLen(2))
			i1, i2 := -1, -1
			if offerInstances["store"][0].ID == "35df52e7-3161-446f-975b-92f32871e37c" {
				i1 = 1
				i2 = 0
			} else {
				i1 = 0
				i2 = 1
			}
			Expect(offerInstances["store"][i1].ID).NotTo(BeNil())
			Expect(offerInstances["store"][i1].ProductID).To(Equal("com.tfg.sample.2"))
			Expect(offerInstances["store"][i1].Contents).To(Equal(dat.JSON([]byte(`{"gems": 100, "gold": 5}`))))
			Expect(offerInstances["store"][i1].Metadata).To(Equal(dat.JSON([]byte(`{"meta": "data"}`))))
			Expect(offerInstances["store"][i1].ExpireAt).To(Equal(int64(1486679200)))

			Expect(offerInstances["store"][i2].ID).To(Equal("35df52e7-3161-446f-975b-92f32871e37c"))
			Expect(offerInstances["store"][i2].ProductID).To(Equal("com.tfg.sample.3"))
			Expect(offerInstances["store"][i2].Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstances["store"][i2].Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstances["store"][i2].ExpireAt).To(Equal(int64(1486679100)))
		})

		It("should return no offer templates when with filters but offer has no filters and inneficient queries not allowed", func() {
			//Given
			playerID := "player-1"
			gameID := defaultGameID
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := map[string]string{
				"key1": "value1",
				"key2": "1.2",
				"key3": "1",
			}

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(0))
		})

		It("should return empty offer template list if gameID doesn't exist", func() {
			//Given
			playerID := "player-1"
			gameID := "non-existing-game"
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(offerInstances).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not return offer-template-1 if player-1 saw it a while ago", func() {
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)
			nextTime := time.Unix(1486678000, 100)
			filterAttrs := make(map[string]string)
			impressionID := uuid.NewV4().String()

			//When
			_, _, err := models.ViewOffer(nil, db, gameID, defaultOfferInstanceID, playerID, impressionID, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, nextTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			Expect(offerInstances).To(HaveLen(2))
			Expect(offerInstances).To(HaveKey("store"))
			Expect(offerInstances).To(HaveKey("unique-place"))
		})

		It("should not return offer-template-1 if player-1 bought it a while ago", func() {
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)
			nextTime := time.Unix(1486678000, 100)
			filterAttrs := make(map[string]string)
			offerInstanceID := "4407b770-5b24-4ffa-8563-0694d1a10156"

			//When
			_, _, _, err := models.ClaimOffer(nil, db, gameID, offerInstanceID, "", "", "", currentTime.Unix(), currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, nextTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())

			//Then
			for _, offerInstance := range offerInstances["store"] {
				Expect(offerInstance.ID).NotTo(Equal(offerInstanceID))
			}
		})

		It("should not return template if it has empty trigger", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-empty-trigger"
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			//When
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(BeEmpty())
		})

		It("should not return template if it reached max frequency", func() {
			playerID := "player-1"
			gameID := "offers-game-max-freq"
			productID := "com.tfg.sample"
			transactionID := uuid.NewV4().String()
			currentTime := time.Unix(1486678000, 0)
			claimTime := int64(1486678000)
			filterAttrs := make(map[string]string)

			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(1))
			Expect(offerInstances).To(HaveKey("store"))
			offerInstanceID := offerInstances["store"][0].ID

			_, alreadyClaimed, _, err := models.ClaimOffer(nil, db, gameID, offerInstanceID, playerID, productID, transactionID, claimTime, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(alreadyClaimed).To(BeFalse())

			offerInstances, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(0))
		})

		It("should not return template if it reached max period", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-max-period"
			productID := "com.tfg.sample"
			transactionID := uuid.NewV4().String()
			currentTime := time.Unix(1486678000, 0)
			claimTime := int64(1486678000)
			filterAttrs := make(map[string]string)

			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(1))
			Expect(offerInstances).To(HaveKey("store"))
			offerInstanceID := offerInstances["store"][0].ID

			_, alreadyClaimed, _, err := models.ClaimOffer(nil, db, gameID, offerInstanceID, playerID, productID, transactionID, claimTime, currentTime, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(alreadyClaimed).To(BeFalse())

			offerInstances, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(0))
		})

		It("should fail if template has invalid frequency", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-invalid-every-freq"
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			//When
			_, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("time: invalid duration invalid"))
		})

		It("should fail if template has invalid period", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game-invalid-every-period"
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			//When
			_, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("time: invalid duration invalid"))
		})

		It("should fail if some error in the database", func() {
			playerID := "player-1"
			gameID := "offers-game"
			currentTime := time.Unix(1486678000, 0)
			filterAttrs := make(map[string]string)

			oldDB := db
			defer func() {
				db = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable

			_, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
		})
	})

	Describe("Get offer info", func() {
		It("should return offer info for an already seen offer", func() {
			//Given
			gameID := defaultGameID
			offerInstanceID := "eb7e8d2a-2739-4da3-aa31-7970b63bdad7"

			//When
			offerInstance, err := models.GetOfferInfo(nil, db, gameID, offerInstanceID, expireDuration, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstance.ID).To(Equal(offerInstanceID))
			Expect(offerInstance.ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstance.Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstance.Cost).To(Equal(dat.JSON([]byte(`{"gems": 500}`))))
			Expect(offerInstance.Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstance.ExpireAt).To(Equal(int64(1486679000)))
		})

		It("should return offer info for an offer instance not in offer versions", func() {
			//Given
			gameID := defaultGameID
			offerInstanceID := "abcd8d2a-2739-4da3-aa31-8970b63bdad7"

			//When
			offerInstance, err := models.GetOfferInfo(nil, db, gameID, offerInstanceID, expireDuration, nil)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstance.ID).To(Equal(offerInstanceID))
			Expect(offerInstance.ProductID).To(Equal("com.tfg.sample"))
			Expect(offerInstance.Contents).To(Equal(dat.JSON([]byte(`{"gems": 5, "gold": 100}`))))
			Expect(offerInstance.Cost).To(Equal(dat.JSON([]byte(`{"gems": 500}`))))
			Expect(offerInstance.Metadata).To(Equal(dat.JSON([]byte(`{}`))))
			Expect(offerInstance.ExpireAt).To(Equal(int64(1486679000)))
		})

		It("should error if gameID doesn't exist", func() {
			//Given
			gameID := "non-existing-game"
			offerInstanceID := "eb7e8d2a-2739-4da3-aa31-7970b63bdad7"

			//When
			_, err := models.GetOfferInfo(nil, db, gameID, offerInstanceID, expireDuration, nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("OfferInstance was not found with specified filters."))
		})

		It("should fail if some error in the database", func() {
			//Given
			gameID := "non-existing-game"
			offerInstanceID := "eb7e8d2a-2739-4da3-aa31-7970b63bdad7"

			oldDB := db
			defer func() {
				db = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			db.(*runner.DB).DB.Close() // make DB connection unavailable

			//When
			_, err = models.GetOfferInfo(nil, db, gameID, offerInstanceID, expireDuration, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("sql: database is closed"))
		})
	})

	Describe("Claim and GetAvailableOffers integrated", func() {
		It("should not return consumed offer after it has been updated", func() {
			offerID := "a2539a8c-55f2-4539-a8c0-929b240d8c80"
			playerID := "player-1"
			gameID := "another-game"
			currentTime := time.Unix(1486678000, 0)
			place := "unique-place"
			transactionID := uuid.NewV4().String()
			filterAttrs := make(map[string]string)

			// Get fot the first time
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveLen(1))
			Expect(offerInstances[place]).To(HaveLen(2))

			// Claim the offer instance
			_, alreadyClaimed, _, err := models.ClaimOffer(nil, db, gameID, offerInstances[place][0].ID, playerID, "", transactionID, currentTime.Unix(), currentTime, nil)
			Expect(alreadyClaimed).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())

			// Get offer to update it
			offer := new(models.Offer)
			builder := db.SQL("SELECT * FROM offers WHERE id = $1 AND game_id = $2", offerID, gameID)
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.QueryStruct(offer)
			Expect(err).NotTo(HaveOccurred())

			// Update its contents and insert with same key
			offer.Contents = dat.JSON([]byte(`{ "somethingNew": 100 }`))
			offer, err = models.UpdateOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			// Should not return the popup offer, since it was claimed for the first time
			offerInstances, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveKey(place))
			Expect(offerInstances[place]).To(HaveLen(1))
		})

		It("should return updated offer with one remaining view", func() {
			playerID := "player-1"
			gameID := "another-game-2"
			place := "unique-place"
			offerID := "f1f74fcd-17ae-4ccd-a248-f77c60e78c82"
			currentTime := time.Unix(1486678000, 0)
			nextTime := func(currentTime time.Time) time.Time {
				return time.Unix(currentTime.Unix()+10, 0)
			}
			filterAttrs := make(map[string]string)

			// Get offer instances
			offerInstances, err := models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances[place]).To(HaveLen(1))
			offerInstanceID := offerInstances[place][0].ID

			// View once
			_, _, err = models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, uuid.NewV4().String(), currentTime, nil)
			Expect(err).NotTo(HaveOccurred())

			// Update Offer
			offer := new(models.Offer)
			builder := db.SQL("SELECT * FROM offers WHERE id = $1 AND game_id = $2", offerID, gameID)
			builder.Execer = edat.NewExecer(builder.Execer)
			err = builder.QueryStruct(offer)
			Expect(err).NotTo(HaveOccurred())
			offer.Contents = dat.JSON([]byte(`{ "somethingNew": 100 }`))
			_, err = models.UpdateOffer(nil, db, offer, offersCache, nil)
			Expect(err).NotTo(HaveOccurred())

			// Get offer
			currentTime = nextTime(currentTime)
			offerInstances, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).To(HaveKey(place))
			Expect(offerInstances[place]).To(HaveLen(1))
			offerInstanceID = offerInstances[place][0].ID

			// Sees twice
			_, _, err = models.ViewOffer(nil, db, gameID, offerInstanceID, playerID, uuid.NewV4().String(), currentTime, nil)
			Expect(err).NotTo(HaveOccurred())

			// Get offer, expect unique-place to not be returned
			currentTime = nextTime(currentTime)
			offerInstances, err = models.GetAvailableOffers(nil, db, offersCache, gameID, playerID, currentTime, expireDuration, filterAttrs, false, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offerInstances).NotTo(HaveKey(place))
		})
	})
})
