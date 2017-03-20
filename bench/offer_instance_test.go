// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package bench

import (
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

	games, err := getGames(&db)
	if err != nil {
		panic(err.Error())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		playerNumber := (i / NumberOfGames) % NumberOfPlayersPerGame
		playerID := fmt.Sprintf("player-%d", playerNumber)
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

	games, err := getGames(&db)
	if err != nil {
		panic(err.Error())
	}
	gamesLength := len(games)

	var offerInstances []*models.OfferInstance

	for i := 0; len(offerInstances) < b.N; i++ {
		gameID := games[i%gamesLength].ID
		playerID := fmt.Sprintf("player-%d", i)
		newOfferInstances, err := getOfferInstances(gameID, playerID)
		if err != nil {
			panic(err.Error())
		}
		offerInstances = append(offerInstances, newOfferInstances...)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offer := offerInstances[i]
		body := map[string]interface{}{
			"gameId":        offer.GameID,
			"playerId":      offer.PlayerID,
			"productId":     "com.tfg.sample",
			"timestamp":     time.Now().Unix(),
			"transactionId": uuid.NewV4().String(),
			"id":            offer.ID,
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

	games, err := getGames(&db)
	if err != nil {
		panic(err.Error())
	}
	gamesLength := len(games)

	var offerInstances []*models.OfferInstance

	for i := 0; len(offerInstances) < b.N; i++ {
		gameID := games[i%gamesLength].ID
		playerID := fmt.Sprintf("player-%d", i)
		newOfferInstances, err := getOfferInstances(gameID, playerID)
		if err != nil {
			panic(err.Error())
		}
		offerInstances = append(offerInstances, newOfferInstances...)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offerInstance := offerInstances[i]
		body := map[string]interface{}{
			"gameId":       offerInstance.GameID,
			"playerId":     offerInstance.PlayerID,
			"impressionId": uuid.NewV4().String(),
		}
		route := getRoute(fmt.Sprintf("/offers/%s/impressions", offerInstance.ID))
		res, err := putTo(route, body)
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}
