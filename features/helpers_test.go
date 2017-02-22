/*
 * Copyright (c) 2017 TFG Co <backend@tfgco.com>
 * Author: TFG Co <backend@tfgco.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package features

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/models"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

func newGame(db runner.Connection, id, bundleID string) (*models.Game, error) {
	game := &models.Game{
		ID:       id,
		Name:     id,
		BundleID: bundleID,
	}
	var c models.RealClock
	err := models.UpsertGame(app.DB, game, c.GetTime(), nil)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func insertOfferTemplate(db runner.Connection, name, gameID string) error {
	offerTemplate := &models.OfferTemplate{
		ID:        uuid.NewV4().String(),
		Name:      name,
		ProductID: "com.tfg.example",
		GameID:    gameID,
		Contents:  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
		Period:    dat.JSON([]byte(`{"type": "once"}`)),
		Frequency: dat.JSON([]byte(`{"every": 24, "unit": "hour"}`)),
		Trigger:   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
		Placement: "popup",
	}
	return models.InsertOfferTemplate(app.DB, offerTemplate, nil)
}

func performRequest(a *api.App, method, url string, payload map[string]interface{}) (int, string, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return 500, "", err
		}
		//fmt.Println(string(data))
		body = strings.NewReader(string(data))
	}
	req := httptest.NewRequest(method, url, body)
	w := httptest.NewRecorder()

	a.Router.ServeHTTP(w, req)

	return w.Code, w.Body.String(), nil
}

func replaceString(val string) string {
	if val == "@VeryBigText@" {
		str := "0123456789"
		for i := 0; i < 10; i++ {
			str += str
		}
		return str
	}
	return val
}
