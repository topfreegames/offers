// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package bench

import (
	"fmt"
	. "github.com/topfreegames/offers/testing"
	"gopkg.in/mgutz/dat.v2/dat"
	"net/http"
	"testing"
)

var offerResponse *http.Response

func BenchmarkListOffers(b *testing.B) {
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

	games, err := getGames(&db)
	if err != nil {
		panic(err.Error())
	}

	gamesLength := len(games)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gameID := games[i%gamesLength].ID
		route := getRoute("/offers")
		res, err := postTo(route, map[string]interface{}{
			"gameId":    gameID,
			"name":      fmt.Sprintf("offer-%d", i),
			"period":    dat.JSON([]byte(`{"every": "1s"}`)),
			"frequency": dat.JSON([]byte(`{"every": "1s"}`)),
			"trigger":   dat.JSON([]byte(`{"from": 1487280506000, "to": 1487280507000}`)),
			"placement": "popup",
			"productId": "tfg.com.sample",
			"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
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

	offers, err := getOffers(&db)
	if err != nil {
		panic(err.Error())
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

func BenchmarkDisableOffer(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	offers, err := getOffers(&db)
	if err != nil {
		panic(err.Error())
	}
	offersLength := len(offers)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offer := offers[i%offersLength]
		route := getRoute(fmt.Sprintf("/offers/%s/disable?game-id=%s", offer.ID, offer.GameID))
		res, err := putTo(route, map[string]interface{}{})
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

	offers, err := getOffers(&db)
	if err != nil {
		panic(err.Error())
	}
	offersLength := len(offers)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offer := offers[i%offersLength]
		route := getRoute(fmt.Sprintf("/offers/%s/enable?game-id=%s", offer.ID, offer.GameID))
		res, err := putTo(route, map[string]interface{}{})
		validateResp(res, err)
		res.Body.Close()

		offerResponse = res
	}
}
