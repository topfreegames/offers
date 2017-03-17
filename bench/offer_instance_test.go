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
	"github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
	"testing"
	"time"
)

func BenchmarkAvailableOffers(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	for _, game := range games {
		_, err = createOffers(&db, game, true, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		playerID := fmt.Sprintf("player-%d", i)
		route := getRoute(fmt.Sprintf("/available-offers?game-id=%s&player-id=%s", game.ID, playerID))
		res, err := get(route)
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkClaimOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	for _, game := range games {
		_, err := createOffers(&db, game, false, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
	}

	var offerInstances []*models.OfferToReturn

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		playerID := fmt.Sprintf("player-%d", i)
		route := getRoute(fmt.Sprintf("/available-offers?game-id=%s&player-id=%s", game.ID, playerID))
		res, err := get(route)
		validateResp(res, err)

		var offersPerPlacement map[string][]*models.OfferToReturn
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = json.Unmarshal(buf.Bytes(), &offersPerPlacement)
		if err != nil {
			panic(err.Error())
		}

		for _, offers := range offersPerPlacement {
			offerInstances = append(offerInstances, offers...)
		}

		res.Body.Close()
	}

	b.ResetTimer()

	for i, offerInstance := range offerInstances {
		gamePosition := i / NumberOfOffersPerGame
		game := games[gamePosition%NumberOfGames]
		playerID := fmt.Sprintf("player-%d", i)
		body := map[string]interface{}{
			"gameId":        game.ID,
			"playerId":      playerID,
			"productId":     "com.tfg.sample",
			"timestamp":     time.Now().Unix(),
			"transactionId": uuid.NewV4().String(),
			"id":            offerInstance.ID,
		}
		route := getRoute("/offers/claim")
		res, err := putTo(route, body)
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkImpressionOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	for _, game := range games {
		_, err := createOffers(&db, game, false, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
	}

	var offerInstances []*models.OfferToReturn

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		playerID := fmt.Sprintf("player-%d", i)
		route := getRoute(fmt.Sprintf("/available-offers?game-id=%s&player-id=%s", game.ID, playerID))
		res, err := get(route)
		validateResp(res, err)

		var offersPerPlacement map[string][]*models.OfferToReturn
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = json.Unmarshal(buf.Bytes(), &offersPerPlacement)
		if err != nil {
			panic(err.Error())
		}

		for _, offers := range offersPerPlacement {
			offerInstances = append(offerInstances, offers...)
		}

		res.Body.Close()
	}

	b.ResetTimer()

	for i, offerInstance := range offerInstances {
		gamePosition := i / NumberOfOffersPerGame
		game := games[gamePosition%NumberOfGames]
		playerID := fmt.Sprintf("player-%d", i)
		body := map[string]interface{}{
			"gameId":       game.ID,
			"playerId":     playerID,
			"impressionId": uuid.NewV4().String(),
		}
		route := getRoute(fmt.Sprintf("/offers/%s/impressions", offerInstance.ID))
		res, err := putTo(route, body)
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}
