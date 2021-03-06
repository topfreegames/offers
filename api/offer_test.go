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
	"io"
	"net/http"
	"net/http/httptest"

	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Template Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("POST /offers", func() {
		It("should return status code 201 for valid parameters", func() {
			name := "New Awesome Game"
			productID := "com.tfg.example"
			gameID := "game-id"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusCreated), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).NotTo(BeEmpty())
			Expect(obj["name"]).To(Equal(name))
			Expect(obj["productId"]).To(Equal(productID))
			Expect(obj["gameId"]).To(Equal(gameID))
			Expect(obj["contents"].(map[string]interface{})["gems"].(float64)).To(BeEquivalentTo(5))
			Expect(obj["contents"].(map[string]interface{})["gold"]).To(BeEquivalentTo(100))
			Expect(obj["period"].(map[string]interface{})["max"].(float64)).To(BeEquivalentTo(1))
			Expect(obj["frequency"].(map[string]interface{})["every"].(string)).To(Equal("24h"))
			Expect(obj["trigger"].(map[string]interface{})["to"].(float64)).To(BeEquivalentTo(1487366964730))
			Expect(obj["trigger"].(map[string]interface{})["from"].(float64)).To(BeEquivalentTo(1487280506875))
			Expect(obj["placement"]).To(Equal(placement))
			Expect(obj["enabled"]).To(BeTrue())
			Expect(int(obj["version"].(float64))).To(Equal(1))
		})

		It("should return status code 201 for valid parameters, including filters", func() {
			name := "New Awesome Game"
			productID := "com.tfg.example"
			gameID := "game-id"
			contents := `{"gems": 5, "gold": 100}`
			period := `{"max": 1}`
			frequency := `{"every": "24h"}`
			trigger := `{"from": 1487280506875, "to": 1487366964730}`
			filters := `{"level": {"eq": "5"}}`
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
				"filters":   dat.JSON([]byte(filters)),
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusCreated), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).NotTo(BeEmpty())
			Expect(obj["name"]).To(Equal(name))
			Expect(obj["productId"]).To(Equal(productID))
			Expect(obj["gameId"]).To(Equal(gameID))
			Expect(obj["contents"].(map[string]interface{})["gems"].(float64)).To(BeEquivalentTo(5))
			Expect(obj["contents"].(map[string]interface{})["gold"]).To(BeEquivalentTo(100))
			Expect(obj["period"].(map[string]interface{})["max"].(float64)).To(BeEquivalentTo(1))
			Expect(obj["frequency"].(map[string]interface{})["every"].(string)).To(Equal("24h"))
			Expect(obj["trigger"].(map[string]interface{})["to"].(float64)).To(BeEquivalentTo(1487366964730))
			Expect(obj["trigger"].(map[string]interface{})["from"].(float64)).To(BeEquivalentTo(1487280506875))
			Expect(obj["filters"].(map[string]interface{})["level"].(map[string]interface{})["eq"].(string)).To(Equal("5"))
			Expect(obj["placement"]).To(Equal(placement))
			Expect(obj["enabled"]).To(BeTrue())
			Expect(int(obj["version"].(float64))).To(Equal(1))
		})

		It("should return status code 201 for valid parameters, including cost", func() {
			name := "New Awesome Game"
			gameID := "game-id"
			contents := `{"gems": 5, "gold": 100}`
			period := `{"max": 1}`
			frequency := `{"every": "24h"}`
			trigger := `{"from": 1487280506875, "to": 1487366964730}`
			cost := `{"gems": 1000}`
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
				"cost":      dat.JSON([]byte(cost)),
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusCreated), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).NotTo(BeEmpty())
			Expect(obj["name"]).To(Equal(name))
			Expect(obj["gameId"]).To(Equal(gameID))
			Expect(obj["contents"].(map[string]interface{})["gems"].(float64)).To(BeEquivalentTo(5))
			Expect(obj["contents"].(map[string]interface{})["gold"]).To(BeEquivalentTo(100))
			Expect(obj["period"].(map[string]interface{})["max"].(float64)).To(BeEquivalentTo(1))
			Expect(obj["frequency"].(map[string]interface{})["every"].(string)).To(Equal("24h"))
			Expect(obj["trigger"].(map[string]interface{})["to"].(float64)).To(BeEquivalentTo(1487366964730))
			Expect(obj["trigger"].(map[string]interface{})["from"].(float64)).To(BeEquivalentTo(1487280506875))
			Expect(obj["cost"].(map[string]interface{})["gems"]).To(BeEquivalentTo(1000))
			Expect(obj).NotTo(HaveKey("productId"))
			Expect(obj["placement"]).To(Equal(placement))
			Expect(obj["enabled"]).To(BeTrue())
			Expect(int(obj["version"].(float64))).To(Equal(1))
		})

		It("should return status code 422 if missing arguments", func() {
			offerReader := JSONFor(JSON{})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;Name: non zero value required;Period: [] does not validate as RequiredJSONObject;;Frequency: [] does not validate as RequiredJSONObject;;Trigger: [] does not validate as RequiredJSONObject;;Placement: non zero value required;Contents: [] does not validate as RequiredJSONObject;;"))
		})

		It("should return status code 422 if missing productID and cost", func() {
			name := "New Awesome Game"
			gameID := "game-id"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Cost and ProductID cannot be both null"))
		})

		It("should return status code 422 if invalid arguments", func() {
			offerReader := JSONFor(JSON{
				"name":      "",
				"productId": "",
				"gameId":    "###",
				"contents":  "{not-a-json}",
				"period":    "{not-a-json}",
				"frequency": "{not-a-json}",
				"trigger":   "{not-a-json}",
				"placement": "",
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: ### does not validate as matches(^[^-][a-zA-Z0-9-_]*$);Name: non zero value required;Period: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Frequency: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Trigger: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Placement: non zero value required;Contents: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;"))
		})

		It("should return status code 422 if game-id doesn't exist", func() {
			offerReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "not-existing-game-id",
				"contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-003"))
			Expect(obj["error"]).To(Equal("InvalidOfferError"))
			Expect(obj["description"]).To(Equal("Offer could not be saved due to: insert or update on table \"offers\" violates foreign key constraint \"offers_game_id_fkey\""))
		})

		It("should return status code 422 if contents is empty", func() {
			offerReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  "",
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Contents: [34 34] does not validate as RequiredJSONObject;;"))
		})

		Describe("Invalid filters format", func() {
			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"eq": 5}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"leq": 5}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"lt": "5", "geq": 2}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"neq": 2}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"geq": "2"}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"lt": 2, "geq": 2}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": 5}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})

			It("should return status code 422 if invalid filter format", func() {
				name := "New Awesome Game"
				productID := "com.tfg.example"
				gameID := "game-id"
				contents := `{"gems": 5, "gold": 100}`
				period := `{"max": 1}`
				frequency := `{"every": "24h"}`
				trigger := `{"from": 1487280506875, "to": 1487366964730}`
				filters := `{"level": {"eq": "arena,1"}}`
				placement := "popup"
				offerReader := JSONFor(JSON{
					"name":      name,
					"productId": productID,
					"gameId":    gameID,
					"contents":  dat.JSON([]byte(contents)),
					"period":    dat.JSON([]byte(period)),
					"frequency": dat.JSON([]byte(frequency)),
					"trigger":   dat.JSON([]byte(trigger)),
					"placement": placement,
					"filters":   dat.JSON([]byte(filters)),
				})

				request, _ := http.NewRequest("POST", "/offers", offerReader)
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-002"))
				Expect(obj["error"]).To(Equal("ValidationFailedError"))
				Expect(obj["description"]).To(ContainSubstring("Filters:"))
			})
		})

		It("returns status code of 500 if database is unavailable", func() {
			offerReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)

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
			Expect(obj["error"]).To(Equal("Insert offer failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code of 401 if no auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			name := "New Awesome Game"
			productID := "com.tfg.example"
			gameID := "game-id"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return status code of 401 if invalid auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			name := "New Awesome Game"
			productID := "com.tfg.example"
			gameID := "game-id"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerReader := JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request, _ := http.NewRequest("POST", "/offers", offerReader)
			request.SetBasicAuth("invaliduser", "invalidpass")

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("PUT /offers/{id}/enable", func() {
		It("should enable an enabled offer", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
		})

		It("should enable a disabled offer", func() {
			id := "27b0370f-bd61-4346-a10d-50ec052ae125"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
		})

		It("returns status code of 500 if database is unavailable", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
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
			Expect(obj["error"]).To(Equal("Update offer failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code 404 if id doesn't exist", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code 422 if invalid parameters", func() {
			id := "not-uuid"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: not-uuid does not validate;"))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/offers//enable?game-id=offers-game", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})

		It("should return status code of 401 if no auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return status code of 401 if invalid auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/enable?game-id=offers-game", id), offerReader)
			request.SetBasicAuth("invaliduser", "invalidpass")

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("PUT /offers/{id}/disable", func() {
		It("should disable an enabled offer", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
		})

		It("should disable a disabled offer", func() {
			id := "27b0370f-bd61-4346-a10d-50ec052ae125"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
		})

		It("should use ID from URI even if a valid one is passed in body", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{
				"id": "aa65a3f2-7cf8-4d76-957f-0a23a1bbbd32",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
		})

		It("returns status code of 500 if database is unavailable", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
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
			Expect(obj["error"]).To(Equal("Update offer failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code 404 if id doesn't exist", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code 422 if invalid parameters", func() {
			id := "not-uuid"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: not-uuid does not validate;"))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/offers//disable?game-id=offers-game", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})

		It("should return status code of 401 if no auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return status code of 401 if invalid auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/disable?game-id=offers-game", id), offerReader)
			request.SetBasicAuth("invaliduser", "invalidpass")

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("PUT /offers/{id}", func() {
		var offerReader io.Reader
		BeforeEach(func() {
			name := "New Awesome Game"
			productID := "com.tfg.example"
			gameID := "offers-game"
			contents := "{\"gems\": 5432}"
			period := "{\"max\": 123}"
			frequency := "{\"every\": \"240h\"}"
			trigger := "{\"from\": 123456789101, \"to\": 123456789111}"
			filters := `{"level": {"geq": 5}}`
			placement := "popup"
			offerReader = JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
				"filters":   dat.JSON([]byte(filters)),
			})
		})

		It("should update offer", func() {
			id := "a411fbcf-dddc-4153-b42b-3f9b2684c965"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
			Expect(int(obj["version"].(float64))).To(Equal(2))
		})

		It("should update offer with cost", func() {
			name := "New Awesome Game"
			gameID := "offers-game"
			contents := "{\"gems\": 5432}"
			cost := "{\"gems\": 1234}"
			period := "{\"max\": 123}"
			frequency := "{\"every\": \"240h\"}"
			trigger := "{\"from\": 123456789101, \"to\": 123456789111}"
			filters := `{"level": {"geq": 5}}`
			placement := "popup"
			offerReader = JSONFor(JSON{
				"name":      name,
				"cost":      dat.JSON([]byte(cost)),
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
				"filters":   dat.JSON([]byte(filters)),
			})
			id := "a411fbcf-dddc-4153-b42b-3f9b2684c965"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(id))
			Expect(int(obj["version"].(float64))).To(Equal(2))
		})

		It("returns status code of 500 if database is unavailable", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
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
			Expect(obj["error"]).To(Equal("Update offer failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code 404 if id doesn't exist", func() {
			id := uuid.NewV4().String()
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code 422 if invalid id", func() {
			id := "not-uuid"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: not-uuid does not validate;"))
		})

		It("should return status code 422 if invalid parameters", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerReader = JSONFor(JSON{
				"contents": "invalid",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;Name: non zero value required;Period: [] does not validate as RequiredJSONObject;;Frequency: [] does not validate as RequiredJSONObject;;Trigger: [] does not validate as RequiredJSONObject;;Placement: non zero value required;Contents: [34 105 110 118 97 108 105 100 34] does not validate as RequiredJSONObject;;"))
		})

		It("should return status code 422 if no productId and cost", func() {
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			name := "New Awesome Game"
			gameID := "offers-game"
			contents := "{\"gems\": 5432}"
			period := "{\"max\": 123}"
			frequency := "{\"every\": \"240h\"}"
			trigger := "{\"from\": 123456789101, \"to\": 123456789111}"
			filters := `{"level": {"geq": 5}}`
			placement := "popup"
			offerReader = JSONFor(JSON{
				"name":      name,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
				"filters":   dat.JSON([]byte(filters)),
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Cost and ProductID cannot be both null"))
		})

		It("should return status code of 401 if no auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return status code of 401 if invalid auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			id := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s", id), offerReader)
			request.SetBasicAuth("invaliduser", "invalidpass")

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("GET /offers", func() {
		It("should return status code of 200 and a list of offers", func() {
			request, _ := http.NewRequest("GET", "/offers?game-id=offers-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())

			offers := obj["offers"].([]interface{})
			Expect(offers).To(HaveLen(5))
			for i := 0; i < len(obj); i++ {
				offer := offers[i].(map[string]interface{})
				Expect(offer).To(HaveKey("id"))
				Expect(offer).To(HaveKey("name"))
				Expect(offer).To(HaveKey("productId"))
				Expect(offer).To(HaveKey("gameId"))
				Expect(offer).To(HaveKey("contents"))
				Expect(offer).To(HaveKey("metadata"))
				Expect(offer).To(HaveKey("enabled"))
				Expect(offer).To(HaveKey("placement"))
				Expect(offer).To(HaveKey("period"))
				Expect(offer).To(HaveKey("frequency"))
				Expect(offer).To(HaveKey("trigger"))
				if offer["id"].(string) == "dd21ec96-2890-4ba0-b8e2-40ea67196990" {
					Expect(offer).To(HaveKey("cost"))
				}
			}

			pages := obj["pages"].(float64)
			Expect(pages).To(Equal(float64(1)))
		})

		It("should return two offers with limit 2 and no offset", func() {
			limit := 2
			url := fmt.Sprintf("/offers?game-id=offers-game&limit=%d", limit)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())

			offers := obj["offers"].([]interface{})
			Expect(offers).To(HaveLen(2))

			pages := obj["pages"].(float64)
			Expect(pages).To(Equal(float64(3)))
		})

		It("should return no offers if limit is 0", func() {
			limit := 0
			url := fmt.Sprintf("/offers?game-id=offers-game&limit=%d", limit)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())

			offers := obj["offers"].([]interface{})
			Expect(offers).To(HaveLen(0))

			pages := obj["pages"].(float64)
			Expect(pages).To(Equal(float64(0)))
		})

		It("should return no offers with no limit and offset 1 (second page)", func() {
			offset := 1
			url := fmt.Sprintf("/offers?game-id=offers-game&offset=%d", offset)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())

			offers := obj["offers"].([]interface{})
			Expect(offers).To(BeEmpty())

			pages := obj["pages"].(float64)
			Expect(pages).To(Equal(float64(1)))
		})

		It("should return error if limit is negative", func() {
			limit := -1
			url := fmt.Sprintf("/offers?game-id=offers-game&limit=%d", limit)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The limit parameter must be an uint."))
			Expect(obj["description"]).To(Equal("strconv.ParseUint: parsing \"-1\": invalid syntax"))
		})

		It("should return error if limit is not a number", func() {
			limit := "qwerty"
			url := fmt.Sprintf("/offers?game-id=offers-game&limit=%s", limit)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The limit parameter must be an uint."))
			Expect(obj["description"]).To(Equal("strconv.ParseUint: parsing \"qwerty\": invalid syntax"))
		})

		It("should return error if offset is negative", func() {
			offset := -1
			url := fmt.Sprintf("/offers?game-id=offers-game&offset=%d", offset)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The offset parameter must be an uint."))
			Expect(obj["description"]).To(Equal("strconv.ParseUint: parsing \"-1\": invalid syntax"))
		})

		It("should return error if offset is not a number", func() {
			offset := "qwerty"
			url := fmt.Sprintf("/offers?game-id=offers-game&offset=%s", offset)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The offset parameter must be an uint."))
			Expect(obj["description"]).To(Equal("strconv.ParseUint: parsing \"qwerty\": invalid syntax"))
		})

		It("should return empty list if no offers", func() {
			request, _ := http.NewRequest("GET", "/offers?game-id=unexistent-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			var obj map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &obj)
			Expect(err).NotTo(HaveOccurred())

			offers := obj["offers"].([]interface{})
			Expect(offers).To(HaveLen(0))

			pages := obj["pages"].(float64)
			Expect(pages).To(Equal(float64(0)))
		})

		It("should return status code of 400 if game-id is not provided", func() {
			request, _ := http.NewRequest("GET", "/offers", nil)

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
			oldDB := app.DB
			defer func() {
				app.DB = oldDB // avoid errors in after each
			}()
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("GET", "/offers?game-id=offers-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("List game offers failed."))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
		})

		It("should return status code of 401 if no auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			request, _ := http.NewRequest("GET", "/offers?game-id=offers-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return status code of 401 if invalid auth provided", func() {
			defer func() {
				config.Set("basicauth.username", "")
				config.Set("basicauth.password", "")
			}()
			config.Set("basicauth.username", "user")
			config.Set("basicauth.password", "pass")
			request, _ := http.NewRequest("GET", "/offers?game-id=offers-game", nil)
			request.SetBasicAuth("invaliduser", "invalidpass")
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
