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
	e "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
	"github.com/topfreegames/offers/testing"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var app *api.App
var logger logrus.Logger
var lastStatus int
var lastBody string
var lastOffers map[string][]*models.OfferToReturn
var clock *testing.MockClock
var cleanDB runner.Connection

func theServerIsUp() error {
	configFile := "../config/acc.yaml"
	config := viper.New()
	config.SetConfigFile(configFile)
	config.SetConfigType("yaml")
	config.SetEnvPrefix("offers")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		return err
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.FatalLevel

	clock := &testing.MockClock{
		CurrentTime: 0,
	}

	var err error
	app, err = api.NewApp("localhost", 9999, config, true, logger, clock)
	if err != nil {
		return err
	}

	cleanDB = app.DB

	return nil
}

func requestGameWithIDAndBundleID(id, bundleID string) error {
	var err error
	lastStatus, lastBody, err = performRequest(app, "PUT", "/games", map[string]interface{}{
		"ID":       replaceString(id),
		"Name":     "Game Awesome Name",
		"BundleID": replaceString(bundleID),
	})

	return err
}

func requestInsertOfferTemplate(name, productID, gameID, contents, period, frequency, trigger, placement string) error {
	code, body, err := performRequest(app, "POST", "/offer-templates", map[string]interface{}{
		"name":      name,
		"productID": productID,
		"gameID":    gameID,
		"contents":  contents,
		"period":    period,
		"frequency": frequency,
		"trigger":   trigger,
		"placement": placement,
	})

	if err != nil {
		return err
	}

	if code != 200 {
		return e.New(body)
	}

	var ot models.OfferTemplate
	json.Unmarshal([]byte(body), &ot)
	if ot.GameID != gameID {
		return e.New("GameID doesn't match")
	}
	return nil
}

//ALIAS to update
func aGameNamedIsCreatedWithBundleIDOf(id, bundleID string) error {
	return requestGameWithIDAndBundleID(id, bundleID)
}

func aGameNamedIsUpdatedWithBundleIDOf(id, bundleID string) error {
	return requestGameWithIDAndBundleID(id, bundleID)
}

func theGameExists(id string) error {
	_, err := models.GetGameByID(app.DB, id, nil)
	return err
}

func theGameHasBundleIDOf(id, bundleID string) error {
	game, err := models.GetGameByID(app.DB, id, nil)
	if err != nil {
		return err
	}
	if game.BundleID != bundleID {
		return fmt.Errorf("Expected game to have bundle ID of %s, but it has %s", bundleID, game.BundleID)
	}
	return nil
}

func theLastRequestReturnedStatusCode(statusCode int) error {
	if lastStatus != statusCode {
		return fmt.Errorf("Expected last request to have status code of %d but it had %d", statusCode, lastStatus)
	}
	return nil
}

func theLastErrorIsWithMessage(code, description string) error {
	var data map[string]string
	err := json.Unmarshal([]byte(lastBody), &data)
	if err != nil {
		return err
	}

	if actualCode, ok := data["code"]; !ok || actualCode != code {
		return fmt.Errorf("Expected status code to be %s, but it was %s", code, actualCode)
	}

	if actualDescription, ok := data["description"]; !ok {
		if !ok {
			return fmt.Errorf("Expected description to be %s, but it was null", description)
		}

		matches := false
		if strings.HasPrefix(description, "*") {
			matches = strings.Contains(actualDescription, description[1:])
		} else {
			matches = actualDescription == description
		}
		if !matches {
			return fmt.Errorf("Expected description to be %s, but it was %s", description, actualDescription)
		}
	}

	return nil
}

func theGameDoesNotExist(id string) error {
	id = replaceString(id)
	if id == "" {
		return nil
	}

	game, err := models.GetGameByID(app.DB, id, nil)
	if err != nil {
		if _, ok := err.(*errors.ModelNotFoundError); ok {
			return nil
		}
		return err
	}
	return fmt.Errorf("The game %s should not exist but it does", game.ID)
}

func aGameWithNameExists(name string) error {
	_, err := newGame(app.DB, name, name)
	return err
}

func theFollowingOfferTemplatesExistInTheGame(gameID string, otArgs *gherkin.DataTable) error {
	for i := 1; i < len(otArgs.Rows); i++ {
		ot := &models.OfferTemplate{
			Name:      otArgs.Rows[i].Cells[1].Value,
			ProductID: otArgs.Rows[i].Cells[2].Value,
			Contents:  dat.JSON([]byte(otArgs.Rows[i].Cells[3].Value)),
			Placement: otArgs.Rows[i].Cells[4].Value,
			Period:    dat.JSON([]byte(otArgs.Rows[i].Cells[5].Value)),
			Frequency: dat.JSON([]byte(otArgs.Rows[i].Cells[6].Value)),
			Trigger:   dat.JSON([]byte(otArgs.Rows[i].Cells[7].Value)),
			GameID:    gameID,
		}

		if _, err := models.GetOfferTemplateByNameAndGame(app.DB, ot.Name, ot.GameID, nil); err != nil {
			if _, err = models.InsertOfferTemplate(app.DB, ot, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func theCurrentTimeIs(arg1 string) error {
	var intCurrentTime int64
	var err error
	intCurrentTime, err = strconv.ParseInt(arg1, 10, 64)

	if err != nil {
		return err
	}

	mockClock := testing.MockClock{
		CurrentTime: intCurrentTime,
	}

	app.Clock = mockClock

	return nil
}

func playerClaimsOfferInGame(playerID, offerID, gameID string) error {
	var err error
	lastStatus, lastBody, err = performRequest(app, "PUT", "/offer/claim", nil)

	if err != nil {
		return err
	}

	return nil
}

func theLastRequestReturnedStatusCodeStatusAndBody(code int, body string) error {
	codeErr := theLastRequestReturnedStatusCode(code)

	if codeErr != nil {
		return codeErr
	}

	if lastBody != body {
		return fmt.Errorf("Expected last request to have body %s but it had %s", body, lastBody)
	}

	return nil
}

func anOfferTemplateIsCreatedInTheGameWith(gameID string, otArgs *gherkin.DataTable) error {
	payload := map[string]interface{}{
		"gameId":    gameID,
		"name":      otArgs.Rows[1].Cells[0].Value,
		"productId": otArgs.Rows[1].Cells[1].Value,
		"contents":  dat.JSON([]byte(otArgs.Rows[1].Cells[2].Value)),
		"metadata":  dat.JSON([]byte(otArgs.Rows[1].Cells[3].Value)),
		"period":    dat.JSON([]byte(otArgs.Rows[1].Cells[4].Value)),
		"frequency": dat.JSON([]byte(otArgs.Rows[1].Cells[5].Value)),
		"trigger":   dat.JSON([]byte(otArgs.Rows[1].Cells[6].Value)),
		"placement": otArgs.Rows[1].Cells[7].Value,
	}
	var err error
	lastStatus, lastBody, err = performRequest(app, "POST", "/offer-templates", payload)

	return err
}

func anOfferTemplateWithNameExistsInGame(offerTemplateName, gameID string) error {
	if _, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil); err != nil {
		return insertOfferTemplate(app.DB, offerTemplateName, gameID)
	}
	return nil
}

func anOfferTemplateExistsWithNameInGame(offerID, gameID string) error {
	return anOfferTemplateWithNameExistsInGame(offerID, gameID)
}

func anOfferTemplateWithNameDoesNotExistInGame(offerTemplateName, gameID string) error {
	if _, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil); err == nil {
		return fmt.Errorf("Expected offer %s to not exist in game %s", offerTemplateName, gameID)
	}
	return nil
}

func theGameRequestsOffersForPlayerIn(gameID, playerID, placement string) error {
	var err error

	url := "/offers?player-id=" + playerID + "&game-id=" + gameID
	lastStatus, lastBody, err = performRequest(app, "GET", url, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(lastBody), &lastOffers)

	if err != nil {
		return err
	}

	if _, ok := lastOffers[placement]; !ok {
		return fmt.Errorf("Expected to have an offer for placement %s to player %s but the is not", placement, playerID)
	}

	return nil
}

func playerOfGameHasSeenOffer(playerID, gameID, offerTemplateName string) error {
	offerTemplate, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil)

	if err != nil {
		return err
	}

	query := `SELECT id, seen_counter FROM offers 
						WHERE player_id = $1 
							AND game_id = $2 
							AND offer_template_id = $3;`
	var offers []models.Offer

	app.DB.SQL(query, playerID, gameID, offerTemplate.ID).QueryStructs(&offers)

	if len(offers) == 0 || offers[0].SeenCounter == 0 {
		return fmt.Errorf("Expected player %s of game %s to has seen offer %s", playerID, gameID, offerTemplateName)
	}

	return nil
}

func playerOfGameHasNotSeenOffer(playerID, gameID, offerTemplateName string) error {
	err := playerOfGameHasSeenOffer(playerID, gameID, offerTemplateName)

	if err == nil {
		return fmt.Errorf("Expected player %s of game %s not to has seen offer %s", playerID, gameID, offerTemplateName)
	}

	return nil
}

func thePlayerOfGameSeesOfferIn(playerID, gameID, placement string) error {
	var err error
	for _, returnedOffer := range lastOffers[placement] {

		lastStatus, lastBody, err = performRequest(app, "PUT", "/offer/last-seen-at", map[string]interface{}{
			"ID":       returnedOffer.ID,
			"PlayerID": playerID,
			"GameID":   gameID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func thePlayerOfGameSeesOfferWithName(playerID, gameID, offerTemplateName string) error {
	offerTemplate, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil)

	if err != nil {
		return err
	}

	query := `SELECT seen_counter, id FROM offers 
						WHERE player_id = $1 
							AND game_id = $2 
							AND offer_template_id = $3;`
	var offers []models.Offer

	err = app.DB.SQL(query, playerID, gameID, offerTemplate.ID).QueryStructs(&offers)

	if err != nil {
		return err
	}

	for _, offer := range offers {
		if offer.OfferTemplateID == offerTemplate.ID {
			lastStatus, lastBody, err = performRequest(app, "PUT", "/offer/last-seen-at", map[string]interface{}{
				"ID":       offer.ID,
				"PlayerID": playerID,
				"GameID":   gameID,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(interface{}) {
		if cleanDB != nil {
			app.DB = cleanDB
		}
	})

	s.Step(`^the server is up$`, theServerIsUp)
	s.Step(`^a game named "([^"]*)" is created with bundle id of "([^"]*)"$`, aGameNamedIsCreatedWithBundleIDOf)
	s.Step(`^the game "([^"]*)" exists$`, theGameExists)
	s.Step(`^the game "([^"]*)" has bundle id of "([^"]*)"$`, theGameHasBundleIDOf)
	s.Step(`^a game named "([^"]*)" is updated with bundle id of "([^"]*)"$`, aGameNamedIsUpdatedWithBundleIDOf)
	s.Step(`^the last request returned status code (\d+)$`, theLastRequestReturnedStatusCode)
	s.Step(`^the last error is "([^"]*)" with message "([^"]*)"$`, theLastErrorIsWithMessage)
	s.Step(`^the game "([^"]*)" does not exist$`, theGameDoesNotExist)
	s.Step(`^a game with name "([^"]*)" exists$`, aGameWithNameExists)
	s.Step(`^an offer template exists with name "([^"]*)" in game "([^"]*)"$`, anOfferTemplateExistsWithNameInGame)
	s.Step(`^an offer template is created in the "([^"]*)" game with:$`, anOfferTemplateIsCreatedInTheGameWith)
	s.Step(`^an offer template with name "([^"]*)" does not exist in game "([^"]*)"$`, anOfferTemplateWithNameDoesNotExistInGame)
	s.Step(`^the following offer templates exist in the "([^"]*)" game:$`, theFollowingOfferTemplatesExistInTheGame)
	s.Step(`^the current time is "([^"]*)"$`, theCurrentTimeIs)
	s.Step(`^the game "([^"]*)" requests offers for player "([^"]*)" in "([^"]*)"$`, theGameRequestsOffersForPlayerIn)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has seen offer "([^"]*)"$`, playerOfGameHasSeenOffer)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has not seen offer "([^"]*)"$`, playerOfGameHasNotSeenOffer)
	s.Step(`^the current time is (\d+)$`, theCurrentTimeIs)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer in "([^"]*)"$`, thePlayerOfGameSeesOfferIn)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer with name "([^"]*)"$`, thePlayerOfGameSeesOfferWithName)
}
