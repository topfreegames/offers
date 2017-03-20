// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package bench

import (
	"fmt"
	"github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
	"net/http"
	"testing"
)

var gameResult *http.Response

func BenchmarkListGames(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		route := getRoute("/games")
		res, err := get(route)
		validateResp(res, err)
		res.Body.Close()

		gameResult = res
	}
}

func BenchmarkInsertGame(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		route := getRoute(fmt.Sprintf("/games/%s", uuid.NewV4().String()))
		res, err := putTo(route, map[string]interface{}{
			"name": fmt.Sprintf("game-%d", i),
		})
		validateResp(res, err)
		res.Body.Close()

		gameResult = res
	}
}

func BenchmarkUpdateGame(b *testing.B) {
	db, err := GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	games, err := getGames(&db)
	if err != nil {
		panic(err.Error())
	}
	length := len(games)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		game := games[i%length]
		route := getRoute(fmt.Sprintf("/games/%s", game.ID))
		res, err := putTo(route, map[string]interface{}{
			"name": fmt.Sprintf("%s-%d", game.Name, i),
		})
		validateResp(res, err)
		res.Body.Close()

		gameResult = res
	}
}
