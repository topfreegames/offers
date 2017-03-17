// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package bench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

func getRoute(url string) string {
	return fmt.Sprintf("http://localhost:8888%s", url)
}

func createGames(db *runner.Connection, numberOfGames int) ([]*models.Game, error) {
	var games []*models.Game

	for i := 0; i < numberOfGames; i++ {
		now := time.Now()
		game := &models.Game{
			Name: fmt.Sprintf("game-%d", i),
			ID:   uuid.NewV4().String(),
		}
		if err := models.UpsertGame(*db, game, now, nil); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func createOffers(db *runner.Connection, game *models.Game, enabled bool, numberOfOffers int) ([]*models.Offer, error) {
	var err error
	offers := getOffers(game, enabled, numberOfOffers)

	for _, offer := range offers {
		offer, err = models.InsertOffer(*db, offer, nil)
		if err != nil {
			return nil, err
		}
	}

	return offers, nil
}

func getOffers(game *models.Game, enabled bool, numberOfOffers int) []*models.Offer {
	var offers []*models.Offer

	for i := 0; i < numberOfOffers; i++ {
		to := time.Now().Unix() + 5*60
		trigger := fmt.Sprintf("{\"from\": 1486678000, \"to\": %d}", to)
		offer := &models.Offer{
			GameID:    game.ID,
			Name:      fmt.Sprintf("offer-%d", i),
			Period:    dat.JSON([]byte(`{"every": "1h"}`)),
			Frequency: dat.JSON([]byte(`{"every": "1h"}`)),
			Trigger:   dat.JSON([]byte(trigger)),
			Placement: "popup",
			ProductID: "tfg.com.sample",
			Contents:  dat.JSON([]byte(`{"x": 1}`)),
			Enabled:   enabled,
		}
		offers = append(offers, offer)
	}

	return offers
}

func get(url string) (*http.Response, error) {
	return sendTo("GET", url, nil)
}

func postTo(url string, payload map[string]interface{}) (*http.Response, error) {
	return sendTo("POST", url, payload)
}

func putTo(url string, payload map[string]interface{}) (*http.Response, error) {
	return sendTo("PUT", url, payload)
}

func sendTo(method, url string, payload map[string]interface{}) (*http.Response, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func validateResp(res *http.Response, err error) {
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		bts, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Request failed with status code %d\n", res.StatusCode)
		panic(string(bts))
	}
}
