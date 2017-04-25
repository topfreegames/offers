// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("Invalid route", func() {
		It("should return status code 404 if invalid route", func() {
			url := "/invalid"
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /available-offers", func() {
		It("should return available offers", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).To(HaveKey("popup"))
			Expect(jsonBody).To(HaveKey("store"))
			popup := jsonBody["popup"]
			Expect(popup).To(HaveLen(1))
			Expect(popup[0]).To(HaveKey("id"))
			Expect(popup[0]).To(HaveKey("productId"))
			Expect(popup[0]).To(HaveKey("contents"))
			Expect(popup[0]).To(HaveKey("metadata"))
			store := jsonBody["store"]
			Expect(store).To(HaveLen(2))
			Expect(store[0]).To(HaveKey("id"))
			Expect(store[0]).To(HaveKey("productId"))
			Expect(store[0]).To(HaveKey("contents"))
			Expect(store[0]).To(HaveKey("metadata"))
			Expect(store[0]).To(HaveKey("expireAt"))
			Expect(store[1]).To(HaveKey("id"))
			Expect(store[1]).To(HaveKey("productId"))
			Expect(store[1]).To(HaveKey("contents"))
			Expect(store[1]).To(HaveKey("metadata"))
			Expect(store[1]).To(HaveKey("expireAt"))
			maxAge := app.MaxAge
			Expect(recorder.Header().Get("Cache-Control")).To(Equal(fmt.Sprintf("max-age=%d", maxAge)))
		})

		It("should return filtered available offers", func() {
			playerID := "player-13"
			gameID := "another-game-with-filters"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s&level=1", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).To(HaveKey("popup"))
			Expect(jsonBody).To(HaveKey("store"))
			popup := jsonBody["popup"]
			Expect(popup).To(HaveLen(1))
			Expect(popup[0]).To(HaveKey("id"))
			Expect(popup[0]).To(HaveKey("productId"))
			Expect(popup[0]).To(HaveKey("contents"))
			Expect(popup[0]).To(HaveKey("metadata"))
			store := jsonBody["store"]
			Expect(store).To(HaveLen(1))
			Expect(store[0]).To(HaveKey("id"))
			Expect(store[0]).To(HaveKey("productId"))
			Expect(store[0]).To(HaveKey("contents"))
			Expect(store[0]).To(HaveKey("metadata"))
			Expect(store[0]).To(HaveKey("expireAt"))
			maxAge := app.MaxAge
			Expect(recorder.Header().Get("Cache-Control")).To(Equal(fmt.Sprintf("max-age=%d", maxAge)))
		})

		It("should return filtered available offers not being received", func() {
			playerID := "player-13"
			gameID := "another-game-with-filters"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s&level=3", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).NotTo(HaveKey("popup"))
			Expect(jsonBody).To(HaveKey("store"))
			store := jsonBody["store"]
			Expect(store).To(HaveLen(1))
			Expect(store[0]).To(HaveKey("id"))
			Expect(store[0]).To(HaveKey("productId"))
			Expect(store[0]).To(HaveKey("contents"))
			Expect(store[0]).To(HaveKey("metadata"))
			Expect(store[0]).To(HaveKey("expireAt"))
			maxAge := app.MaxAge
			Expect(recorder.Header().Get("Cache-Control")).To(Equal(fmt.Sprintf("max-age=%d", maxAge)))
		})

		It("should return game cacheMaxAge if available", func() {
			playerID := "player-1"
			gameID := "offers-game-maxage"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Header().Get("Cache-Control")).To(Equal("max-age=123"))
		})

		It("should return empty list of available offers", func() {
			playerID := "player-1"
			gameID := "non-existing-offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody).To(BeEmpty())
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return empty list if current time is before all triggers", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Clock = MockClock{
				CurrentTime: 100,
			}
			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).To(BeEmpty())
			app.Clock = MockClock{
				CurrentTime: 1486678000,
			}
		})

		It("should return empty list if current time is after all triggers", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Clock = MockClock{
				CurrentTime: 2 * 1000 * 1000 * 1000,
			}
			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).To(BeEmpty())
		})

		It("should return offer again if currentTime - lastView >= frequency", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			offerInstanceID := "eb7e8d2a-2739-4da3-aa31-7970b63bdad7"
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string]interface{}

			app.Router.ServeHTTP(recorder, request)

			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     100,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})

			request, _ = http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()

			app.Router.ServeHTTP(recorder, request)

			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			contents := jsonBody["contents"].(map[string]interface{})

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody["nextAt"]).To(BeEquivalentTo(101))
			Expect(contents["gems"]).To(BeEquivalentTo(5))
			Expect(contents["gold"]).To(BeEquivalentTo(100))

			app.Router.ServeHTTP(recorder, request)

			offerReader = JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     101,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID})

			request, _ = http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()

			app.Router.ServeHTTP(recorder, request)
			err = json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			contents = jsonBody["contents"].(map[string]interface{})

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody["nextAt"]).To(BeEquivalentTo(102))
			Expect(contents["gems"]).To(BeEquivalentTo(5))
			Expect(contents["gold"]).To(BeEquivalentTo(100))
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should not return offer again if currentTime - lastView < frequency", func() {
			playerID := "player-1"
			gameID := "offers-game"
			offerInstanceID := "4407b770-5b24-4ffa-8563-0694d1a10156"
			var jsonBody map[string]interface{}

			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     100,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)
			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			contents := jsonBody["contents"].(map[string]interface{})

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody["nextAt"]).To(BeEquivalentTo(43300))
			Expect(contents["gems"]).To(BeEquivalentTo(5))
			Expect(contents["gold"]).To(BeEquivalentTo(100))

			offerReader = JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     100,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ = http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			err = json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			contents = jsonBody["contents"].(map[string]interface{})

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody["nextAt"]).To(BeEquivalentTo(43300))
			Expect(contents["gems"]).To(BeEquivalentTo(5))
			Expect(contents["gold"]).To(BeEquivalentTo(100))
		})

		It("should return status code 400 if player-id is not informed available offers", func() {
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?game-id=%s", gameID)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The player-id parameter cannot be empty."))
			Expect(obj["description"]).To(Equal("The player-id parameter cannot be empty"))

		})

		It("should return status code 400 if game-id is not informed available offers", func() {
			playerID := "player-1"
			url := fmt.Sprintf("/available-offers?player-id=%s", playerID)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The game-id parameter cannot be empty."))
			Expect(obj["description"]).To(Equal("The game-id parameter cannot be empty"))
		})

		It("should return status code of 500 if some error occurred", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("Failed to retrieve offer for player"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should not return instance after claim if offer period has max 1", func() {
			// Create Offer by requesting it
			gameID := "limited-offers-game"
			playerID := "player-1"
			place := "store"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var body map[string][]*models.OfferToReturn

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			err := json.Unmarshal(recorder.Body.Bytes(), &body)
			Expect(err).ToNot(HaveOccurred())

			// Claim the Offer
			id := body[place][0].ID
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     time.Now().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			})
			request, _ = http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Body.String()).To(Equal(`{"contents":{"gems":5,"gold":100}}`))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must not be returned again in next Get
			request, _ = http.NewRequest("GET", url, nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var newBody map[string][]*models.OfferToReturn
			err = json.Unmarshal(recorder.Body.Bytes(), &newBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(newBody).NotTo(HaveKey(place))
		})

		It("should return offer again if max is 2 and currentTime - lastClaim >= period", func() {
			// Create Offer by requesting it
			gameID := "another-game"
			playerID := "player-12"
			place := "unique-place"
			offerInstanceID := "f0e3ce8c-bf9b-4da2-886b-e1b5a62a18cb"

			// Claim the Offer
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     1486678000,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)

			var jsonBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			contents := jsonBody["contents"]
			Expect(contents).To(Equal(map[string]interface{}{
				"gold": float64(100),
				"gems": float64(5),
			}))
			Expect(jsonBody["nextAt"]).To(Equal(float64(1486678010)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must not be returned again in next Get
			app.Clock = MockClock{CurrentTime: 1486678010}
			request, _ = http.NewRequest("GET", fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID), nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var newBody map[string][]*models.OfferToReturn
			err := json.Unmarshal(recorder.Body.Bytes(), &newBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(newBody).To(HaveKey(place))
		})

		It("should return offer again if max is 2 and currentTime - lastClaim < period", func() {
			// Create Offer by requesting it
			gameID := "another-game"
			playerID := "player-12"
			offerInstanceID := "f0e3ce8c-bf9b-4da2-886b-e1b5a62a18cb"

			// Claim the Offer
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     1486678000,
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)

			var jsonBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			contents := jsonBody["contents"]
			Expect(contents).To(Equal(map[string]interface{}{
				"gold": float64(100),
				"gems": float64(5),
			}))
			Expect(jsonBody["nextAt"]).To(Equal(float64(1486678010)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must not be returned again in next Get
			app.Clock = MockClock{CurrentTime: 148667805}
			request, _ = http.NewRequest("GET", fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID), nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var newBody map[string][]*models.OfferToReturn
			err := json.Unmarshal(recorder.Body.Bytes(), &newBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(newBody).To(BeEmpty())
		})

		It("should return instance after one view if offer frequency has max 2", func() {
			// Create Offer by requesting it
			gameID := "offers-game"
			playerID := "player-1"
			id := "6c4a79f2-24b8-4be9-93d4-12413b789823"
			var body map[string]int64

			// View the Offer
			offerReader := JSONFor(JSON{
				"gameId":       gameID,
				"playerId":     playerID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			json.Unmarshal(recorder.Body.Bytes(), &body)
			Expect(body["nextAt"]).To(Equal(int64(1486678001)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must be returned again in next Get
			app.Clock = MockClock{CurrentTime: 1486678002}
			request, _ = http.NewRequest("GET", fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID), nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))

			var offersToReturn map[string][]models.OfferToReturn
			json.Unmarshal(recorder.Body.Bytes(), &offersToReturn)
			offer := models.OfferToReturn{
				ID:        "6c4a79f2-24b8-4be9-93d4-12413b789823",
				ProductID: "com.tfg.sample.3",
				Contents:  dat.JSON([]byte(`{"gems":5,"gold":100}`)),
				Metadata:  dat.JSON([]byte("{}")),
				ExpireAt:  1486679100,
			}
			Expect(offersToReturn["store"]).To(ContainElement(offer))
		})

		It("should return instance after one view if offer frequency has max 2", func() {
			// Create Offer by requesting it
			gameID := "offers-game"
			playerID := "player-1"
			id := "6c4a79f2-24b8-4be9-93d4-12413b789823"
			var body map[string]int64

			// View the Offer for the first time
			offerReader := JSONFor(JSON{
				"gameId":       gameID,
				"playerId":     playerID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			json.Unmarshal(recorder.Body.Bytes(), &body)
			Expect(body["nextAt"]).To(Equal(int64(1486678001)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// View the Offer for the second time
			app.Clock = MockClock{CurrentTime: 1486678001}
			offerReader = JSONFor(JSON{
				"gameId":       gameID,
				"playerId":     playerID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ = http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			json.Unmarshal(recorder.Body.Bytes(), &body)
			Expect(body["nextAt"]).To(Equal(int64(1486678001))) // No nextAt, so keeps the same as before
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must be returned again in next Get
			app.Clock = MockClock{CurrentTime: 1486678002}
			request, _ = http.NewRequest("GET", fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID), nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))

			var offersToReturn map[string][]models.OfferToReturn
			json.Unmarshal(recorder.Body.Bytes(), &offersToReturn)
			offer := models.OfferToReturn{
				ID:        "6c4a79f2-24b8-4be9-93d4-12413b789823",
				ProductID: "com.tfg.sample.3",
				Contents:  dat.JSON([]byte(`{"gems":5,"gold":100}`)),
				Metadata:  dat.JSON([]byte("{}")),
				ExpireAt:  1486679100,
			}
			Expect(offersToReturn["store"]).NotTo(ContainElement(offer))
		})

		It("should return available offers from cache on second request", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/available-offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)

			start := time.Now().UnixNano()
			app.Router.ServeHTTP(recorder, request)
			dbElapsedTime := time.Now().UnixNano() - start

			recorder = httptest.NewRecorder()
			start = time.Now().UnixNano()
			app.Router.ServeHTTP(recorder, request)
			cacheElapsedTime := time.Now().UnixNano() - start

			Expect(dbElapsedTime).To(BeNumerically(">", cacheElapsedTime))
		})
	})

	Describe("PUT /offers/claim", func() {
		It("should claim valid offer", func() {
			offerInstanceID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			gameID := "offers-game"
			playerID := "player-1"
			claimCounterKey := models.GetClaimCounterKey(playerID, offerID)
			claimTimestampKey := models.GetClaimTimestampKey(playerID, offerID)
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "product-id",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+1)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			claimCount, err := app.RedisClient.Client.Get(claimCounterKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimCount).To(Equal(int64(1)))

			claimTimestamp, err := app.RedisClient.Client.Get(claimTimestampKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should claim valid offer even if id is not passed", func() {
			gameID := "offers-game"
			playerID := "player-1"
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			claimCounterKey := models.GetClaimCounterKey(playerID, offerID)
			claimTimestampKey := models.GetClaimTimestampKey(playerID, offerID)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+1)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			claimCount, err := app.RedisClient.Client.Get(claimCounterKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimCount).To(Equal(int64(1)))

			claimTimestamp, err := app.RedisClient.Client.Get(claimTimestampKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should not claim a claimed offer", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			gameID := "offers-game"
			offerJSON := JSON{
				"gameId":        gameID,
				"playerId":      "player-1",
				"productId":     "product-id",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			}
			offerReader1 := JSONFor(offerJSON)
			offerReader2 := JSONFor(offerJSON)
			request1, _ := http.NewRequest("PUT", "/offers/claim", offerReader1)
			request2, _ := http.NewRequest("PUT", "/offers/claim", offerReader2)

			app.Router.ServeHTTP(recorder, request1)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+1)))
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request2)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+1)))
			Expect(recorder.Code).To(Equal(http.StatusConflict))
		})

		It("should return nextAt to lastClaim + every if replay request", func() {
			offerID := "5fed76ab-1fd7-4a91-972d-bca228ce80c4"
			offerInstanceID := "4407b770-5b24-4ffa-8563-0694d1a10156"
			gameID := "offers-game"
			playerID := "player-1"
			claimCounterKey := models.GetClaimCounterKey(playerID, offerID)
			claimTimestampKey := models.GetClaimTimestampKey(playerID, offerID)
			transactionID := uuid.NewV4().String()

			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "product-id",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": transactionID,
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+12*60*60)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			offerReader = JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "product-id",
				"timestamp":     app.Clock.GetTime().Unix() + 100,
				"transactionId": transactionID,
				"id":            offerInstanceID,
			})
			request, _ = http.NewRequest("PUT", "/offers/claim", offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+12*60*60)))
			Expect(recorder.Code).To(Equal(http.StatusConflict))

			claimCount, err := app.RedisClient.Client.Get(claimCounterKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimCount).To(Equal(int64(1)))

			claimTimestamp, err := app.RedisClient.Client.Get(claimTimestampKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should not return nextAt if max is achieved", func() {
			offerInstanceID := "fe528bb0-dab6-4f5a-b6cd-347422fd9817"
			offerID := "9456f6c4-f9f1-4dd9-8841-9e5770c8e62c"
			gameID := "offers-game-max-freq"
			playerID := "player-1"
			claimCounterKey := models.GetClaimCounterKey(playerID, offerID)
			claimTimestampKey := models.GetClaimTimestampKey(playerID, offerID)
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "product-id",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            offerInstanceID,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(`{"contents":{"gems":5,"gold":100}}`))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			claimCount, err := app.RedisClient.Client.Get(claimCounterKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimCount).To(Equal(int64(1)))

			claimTimestamp, err := app.RedisClient.Client.Get(claimTimestampKey).Int64()
			Expect(err).ToNot(HaveOccurred())
			Expect(claimTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should return 422 if invalid OfferID", func() {
			id := "invalid-offer-id"
			offerReader := JSONFor(JSON{
				"gameId":        "offers-game",
				"playerId":      "player-1",
				"productId":     "product-id",
				"timestamp":     time.Now().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal(fmt.Sprintf("OfferInstanceID: %s does not validate as uuidv4;", id)))
		})

		It("should return 422 if missing parameters", func() {
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;PlayerID: non zero value required;ProductID: non zero value required;Timestamp: non zero value required;TransactionID: non zero value required;"))
		})

		It("should return 404 if non existing OfferID", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{
				"gameId":        "offers-game",
				"playerId":      "player-1",
				"productId":     "product-id",
				"timestamp":     time.Now().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 if non existing GameID", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"gameId":        "non-existing-offers-game",
				"playerId":      "player-1",
				"productId":     "product-id",
				"timestamp":     time.Now().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferInstanceNotFoundError"))
			Expect(obj["description"]).To(Equal("OfferInstance was not found with specified filters."))
		})

		It("should return status code of 500 if some error occurred", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"gameId":        "offers-game",
				"playerId":      "player-1",
				"productId":     "product-id",
				"timestamp":     time.Now().Unix(),
				"transactionId": uuid.NewV4().String(),
				"id":            id,
			})

			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("sql: database is closed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code 400 if invalid json is sent", func() {
			invalidJSON := `"gameId:   "offers-g "Player-1"`
			request, _ := http.NewRequest("PUT", "/offers/claim", strings.NewReader(invalidJSON))

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("json: cannot unmarshal string into Go value of type models.ClaimOfferPayload"))
		})

		It("should return 404 if claim on inexistent offer when offerInstanceID is not passed", func() {
			gameID := "non-existing-offers-game"
			playerID := "player-1"
			offerReader := JSONFor(JSON{
				"gameId":        gameID,
				"playerId":      playerID,
				"productId":     "com.tfg.sample",
				"timestamp":     app.Clock.GetTime().Unix(),
				"transactionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", "/offers/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var jsonBody map[string]string
			json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(jsonBody["code"]).To(Equal("OFF-001"))
			Expect(jsonBody["description"]).To(Equal("offerInstance was not found with specified filters."))
			Expect(jsonBody["error"]).To(Equal("offerInstanceNotFoundError"))
		})
	})

	Describe("PUT /offers/{id}/impressions", func() {
		It("should increment view counter and update ", func() {
			offerInstanceID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			gameID := "offers-game"
			playerID := "player-1"
			viewCounterKey := models.GetViewCounterKey(playerID, offerID)
			viewTimestampKey := models.GetViewTimestampKey(playerID, offerID)
			offerReader := JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix() + 1))

			viewCount, err := app.RedisClient.Client.Get(viewCounterKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewCount).To(Equal(int64(1)))

			viewTimestamp, err := app.RedisClient.Client.Get(viewTimestampKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should return the current timestamp as nextAt if offer reached max period", func() {
			id := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			gameID := "limited-offers-game"
			offerReader := JSONFor(JSON{
				"playerId":     "player-1",
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(HaveLen(0))
		})

		It("should return nextAt equal to now if offer has no every in frequency", func() {
			id := "4407b770-5b24-4ffa-8563-0694d1a10156"
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId":     "player-11",
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should increment view counter twice after seeing twice", func() {
			offerInstanceID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			gameID := "offers-game"
			playerID := "player-1"
			viewCounterKey := models.GetViewCounterKey(playerID, offerID)
			viewTimestampKey := models.GetViewTimestampKey(playerID, offerID)

			// View for the first time
			offerReader := JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix() + 1))

			// View for the second time
			app.Clock = MockClock{CurrentTime: app.Clock.GetTime().Unix() + 1}
			offerReader = JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ = http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix() + 1))

			viewCount, err := app.RedisClient.Client.Get(viewCounterKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewCount).To(Equal(int64(2)))

			viewTimestamp, err := app.RedisClient.Client.Get(viewTimestampKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewTimestamp).To(Equal(app.Clock.GetTime().Unix()))
		})

		It("should return nextAt zero after seeing twice offer with max period 2", func() {
			offerInstanceID := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			playerID := "player-1"
			offerID := "aa65a3f2-7cf8-4d76-957f-0a23a1bbbd32"
			gameID := "limited-offers-game"
			viewCounterKey := models.GetViewCounterKey(playerID, offerID)
			err := app.RedisClient.Client.Set(viewCounterKey, 1, 0).Err()
			Expect(err).ToNot(HaveOccurred())

			offerReader := JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]int64
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["nextAt"]).To(Equal(int64(0)))
		})

		It("should not increment view counter if impressionID is the same", func() {
			offerInstanceID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			gameID := "offers-game"
			playerID := "player-1"
			viewCounterKey := models.GetViewCounterKey(playerID, offerID)
			viewTimestampKey := models.GetViewTimestampKey(playerID, offerID)
			impressionID := uuid.NewV4().String()
			timestamp := app.Clock.GetTime().Unix()

			// View for the first time
			offerReader := JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": impressionID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(timestamp + 1))

			// View for the second time
			app.Clock = MockClock{CurrentTime: timestamp + 1}
			offerReader = JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": impressionID,
			})
			request, _ = http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusConflict))
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(timestamp + 2))

			viewCount, err := app.RedisClient.Client.Get(viewCounterKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewCount).To(Equal(int64(1)))

			viewTimestamp, err := app.RedisClient.Client.Get(viewTimestampKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewTimestamp).To(Equal(timestamp))
		})

		It("should not increment when is a retry request and rechead max view", func() {
			offerInstanceID := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			gameID := "offers-game"
			playerID := "player-1"
			viewCounterKey := models.GetViewCounterKey(playerID, offerID)
			viewTimestampKey := models.GetViewTimestampKey(playerID, offerID)
			impressionKey := models.GetImpressionsKey(playerID, gameID)
			impressionID := uuid.NewV4().String()
			timestamp := app.Clock.GetTime().Unix()

			// already seen once
			err := app.RedisClient.Client.Set(viewCounterKey, 1, 0).Err()
			Expect(err).ToNot(HaveOccurred())
			err = app.RedisClient.Client.Set(viewTimestampKey, timestamp, 0).Err()
			Expect(err).ToNot(HaveOccurred())
			err = app.RedisClient.Client.SAdd(impressionKey, impressionID).Err()
			Expect(err).ToNot(HaveOccurred())

			offerReader := JSONFor(JSON{
				"playerId":     playerID,
				"gameId":       gameID,
				"impressionId": impressionID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", offerInstanceID), offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusConflict))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(timestamp + 1))

			viewCount, err := app.RedisClient.Client.Get(viewCounterKey).Int64()
			Expect(err).NotTo(HaveOccurred())
			Expect(viewCount).To(Equal(int64(1)))
		})

		It("should return status code 422 if invalid parameters", func() {
			offerReader := JSONFor(JSON{
				"playerId":     "player-1",
				"gameId":       "offers-game",
				"impressionId": uuid.NewV4().String(),
			})
			url := "/offers/invalid-uuid/impressions"
			request, _ := http.NewRequest("PUT", url, offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: invalid-uuid does not validate;"))
		})

		It("should return status code 422 if missing body parameters", func() {
			offerReader := JSONFor(JSON{})
			id := uuid.NewV4().String()
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;PlayerID: non zero value required;ImpressionID: non zero value required;"))
		})

		It("should return status code 404 if offer with given ID does not exist", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{
				"playerId":     "player-1",
				"gameId":       "offers-game",
				"impressionId": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferInstanceNotFoundError"))
			Expect(obj["description"]).To(Equal("OfferInstance was not found with specified filters."))
		})

		It("should return status code of 500 if some error occurred", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"playerId":     "player-1",
				"gameId":       "offers-game",
				"impressionId": uuid.NewV4().String(),
			})

			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("sql: database is closed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code 409 if impression was already received", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerJSON := JSON{
				"playerId":     "player-1",
				"gameId":       "offers-game",
				"impressionId": uuid.NewV4().String(),
			}

			offerReader1 := JSONFor(offerJSON)
			offerReader2 := JSONFor(offerJSON)
			request1, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader1)
			request2, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/impressions", id), offerReader2)

			app.Router.ServeHTTP(recorder, request1)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request2)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusConflict))

			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix() + 1))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/offers//impressions", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})
})
