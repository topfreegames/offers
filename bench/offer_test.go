// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package bench

import (
	"fmt"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
	"net/http"
	"testing"
)

var offerResponse *http.Response

func BenchmarkListOffers(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < NumberOfGames; i++ {
		_, err = createOffers(&db, games[i], true, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		route := getRoute(fmt.Sprintf("/offers?game-id=%s", game.ID))
		res, err := get(route)
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkInsertOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	var offers []*models.Offer
	for _, game := range games {
		newOffers := getOffers(game, true, NumberOfOffersPerGame)
		offers = append(offers, newOffers...)
	}

	offersLength := len(offers)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offer := offers[i%offersLength]
		route := getRoute("/offers")
		res, err := postTo(route, map[string]interface{}{
			"gameId":    offer.GameID,
			"name":      offer.Name,
			"period":    offer.Period,
			"frequency": offer.Frequency,
			"trigger":   offer.Trigger,
			"placement": offer.Placement,
			"productId": offer.ProductID,
			"contents":  offer.Contents,
		})
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkUpdateOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	var offers []*models.Offer
	for _, game := range games {
		newOffers, err := createOffers(&db, game, true, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
		offers = append(offers, newOffers...)
	}

	offersLength := len(offers)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offer := offers[i%offersLength]
		route := getRoute(fmt.Sprintf("/offers/%s", offer.ID))
		res, err := putTo(route, map[string]interface{}{
			"gameId":    offer.GameID,
			"name":      fmt.Sprintf("%s-new", offer.Name),
			"period":    offer.Period,
			"frequency": offer.Frequency,
			"trigger":   offer.Trigger,
			"placement": offer.Placement,
			"productId": offer.ProductID,
			"contents":  offer.Contents,
		})
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkEnableOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	offersByGame := make(map[string][]*models.Offer)
	for _, game := range games {
		offers, err := createOffers(&db, game, false, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
		offersByGame[game.ID] = offers
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		offer := offersByGame[game.ID][i%NumberOfOffersPerGame]
		route := getRoute(fmt.Sprintf("/offers/%s/enable?game-id=%s", offer.ID, game.ID))
		res, err := putTo(route, map[string]interface{}{})
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}

func BenchmarkDisableOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := createGames(&db, NumberOfGames)
	if err != nil {
		panic(err.Error())
	}

	offersByGame := make(map[string][]*models.Offer)
	for _, game := range games {
		offers, err := createOffers(&db, game, false, NumberOfOffersPerGame)
		if err != nil {
			panic(err.Error())
		}
		offersByGame[game.ID] = offers
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%NumberOfGames]
		offer := offersByGame[game.ID][i%NumberOfOffersPerGame]
		route := getRoute(fmt.Sprintf("/offers/%s/disable?game-id=%s", offer.ID, game.ID))
		res, err := putTo(route, map[string]interface{}{})
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}
