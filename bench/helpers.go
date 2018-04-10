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

	edat "github.com/topfreegames/extensions/dat"
	"github.com/topfreegames/offers/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

func getRoute(url string) string {
	return fmt.Sprintf("http://localhost:8889%s", url)
}

func getGames(db *runner.Connection) ([]*models.Game, error) {
	var games []*models.Game

	query := `SELECT id, name FROM games`
	builder := (*db).SQL(query)
	builder.Execer = edat.NewExecer(builder.Execer)
	err := builder.QueryStructs(&games)

	return games, err
}

func getOffers(db *runner.Connection) ([]*models.Offer, error) {
	var offers []*models.Offer

	query := "SELECT * FROM offers"
	builder := (*db).SQL(query)
	builder.Execer = edat.NewExecer(builder.Execer)
	err := builder.QueryStructs(&offers)

	return offers, err
}

func getOfferInstances(gameID, playerID string) ([]*models.OfferInstance, error) {
	var offerInstances []*models.OfferInstance
	route := getRoute(fmt.Sprintf("/available-offers?game-id=%s&player-id=%s", gameID, playerID))
	res, err := get(route)
	validateResp(res, err)

	var offersPerPlacement map[string][]*models.OfferToReturn
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	err = json.Unmarshal(buf.Bytes(), &offersPerPlacement)
	if err != nil {
		panic(err.Error())
	}
	res.Body.Close()

	for _, offers := range offersPerPlacement {
		for _, offer := range offers {
			offerInstances = append(offerInstances, &models.OfferInstance{
				GameID:   gameID,
				PlayerID: playerID,
				ID:       offer.ID,
			})
		}
	}

	return offerInstances, err
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
