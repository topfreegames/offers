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
)

var app *api.App
var logger logrus.Logger
var lastStatus int
var lastBody string
var lastOffers map[string][]*models.OfferTemplate
var clock *testing.MockClock

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

		_, err := models.InsertOfferTemplate(app.DB, ot, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func theFollowingPlayersExistInTheGame(gameID string, players *gherkin.DataTable) error {
	for i := 1; i < len(players.Rows); i++ {
		if _, err := models.GetAvailableOffers(app.DB, players.Rows[i].Cells[0].Value, gameID, app.Clock.GetTime(), nil); err != nil {
			return err
		}

		//TODO: call update offer last seen at
	}

	return nil
}

func theCurrentTimeIs(currentTime int64) error {
	mockClock := testing.MockClock{
		CurrentTime: currentTime,
	}

	app.Clock = mockClock

	return nil
}

func playerClaimsOfferInGame(playerID, offerID, gameID string) error {
	_, alreadyClaimed, err := models.ClaimOffer(app.DB, offerID, playerID, gameID, app.Clock.GetTime(), nil)

	if alreadyClaimed {
		return fmt.Errorf("Offer %s has already been claimed", offerID)
	}

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

func theCurrentTimeIsD(arg1 int) error {
	return godog.ErrPending
}

func anOfferTemplateIsCreatedInTheGameWith(gameID string, otArgs *gherkin.DataTable) error {
	payload := map[string]interface{}{
		"gameId":    gameID,
		"name":      otArgs.Rows[1].Cells[0].Value,
		"productId": otArgs.Rows[1].Cells[1].Value,
		"contents":  otArgs.Rows[1].Cells[2].Value,
		"metadata":  otArgs.Rows[1].Cells[3].Value,
		"period":    otArgs.Rows[1].Cells[4].Value,
		"frequency": otArgs.Rows[1].Cells[5].Value,
		"trigger":   otArgs.Rows[1].Cells[6].Value,
		"placement": otArgs.Rows[1].Cells[7].Value,
	}
	var err error
	lastStatus, lastBody, err = performRequest(app, "PUT", "/offer-templates", payload)

	fmt.Println(payload)
	fmt.Println(err)

	return err
}

func anOfferTemplateWithNameExistsInGame(offerTemplateName, gameID string) error {
	offerTemplate, err := models.GetOfferTemplateByName(app.DB, offerTemplateName, nil)

	if err != nil {
		return err
	}

	if offerTemplate.GameID != gameID {
		return fmt.Errorf("Offer template %s doesn't exist in game %s", offerTemplateName, gameID)
	}

	return nil
}

func anOfferTemplateExistsWithNameInGame(offerID, gameID string) error {
	return anOfferTemplateWithNameExistsInGame(offerID, gameID)
}

func anOfferTemplateWithNameDoesNotExistInGame(offerID, gameID string) error {
	err := anOfferTemplateWithNameExistsInGame(offerID, gameID)

	if err != nil {
		return nil
	}

	return fmt.Errorf("Expected offer %s to not exist in game %s", offerID, gameID)
}

func theGameRequestsOffersForPlayerIn(gameID, playerID, placement string) error {
	var err error
	lastOffers, err = models.GetAvailableOffers(app.DB, playerID, gameID, app.Clock.GetTime(), nil)

	if err != nil {
		return err
	}

	if _, ok := lastOffers[placement]; !ok {
		return fmt.Errorf("Expected to have an offer for placement %s to player %s but the is not", placement, playerID)
	}

	return nil
}

func anOfferWithNameIsReturned(offerID string) error {
	// To continue, I need the offer ID returned by GetAvailableOffers
	return godog.ErrPending
}

func playerHasSeenOffer(playerID, gameID, offerID string) error {
	// To continue, I need the offer ID returned by GetAvailableOffers
	//	offers, err := models.GetAvailableOffers(app.DB, playerID, gameID, app.Clock.GetTime(), nil)
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, ots : range offers {
	//		for _, offerTemplate := range ots {
	//
	//		}
	//	}
	//
	return godog.ErrPending
}

func playerHasNotSeenOffer(playerID, gameID, offerID string) error {
	err := playerHasSeenOffer(playerID, gameID, offerID)

	if err == nil {
		return fmt.Errorf("Expected player %s of game %s not to has seen offer %s", playerID, gameID, offerID)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the server is up$`, theServerIsUp)
	s.Step(`^a game named "([^"]*)" is created with bundle id of "([^"]*)"$`, aGameNamedIsCreatedWithBundleIDOf)
	s.Step(`^the game "([^"]*)" exists$`, theGameExists)
	s.Step(`^the game "([^"]*)" has bundle id of "([^"]*)"$`, theGameHasBundleIDOf)
	s.Step(`^a game named "([^"]*)" is updated with bundle id of "([^"]*)"$`, aGameNamedIsUpdatedWithBundleIDOf)
	s.Step(`^the last request returned status code (\d+)$`, theLastRequestReturnedStatusCode)
	s.Step(`^the last error is "([^"]*)" with message "([^"]*)"$`, theLastErrorIsWithMessage)
	s.Step(`^the game "([^"]*)" does not exist$`, theGameDoesNotExist)
	s.Step(`^a game with name "([^"]*)" exists$`, aGameWithNameExists)
	s.Step(`^the following offer templates exist in the "([^"]*)" game:$`, theFollowingOfferTemplatesExistInTheGame)
	s.Step(`^the following players exist in the "([^"]*)" game:$`, theFollowingPlayersExistInTheGame)
	s.Step(`^the current time is (\d+)$`, theCurrentTimeIs)
	s.Step(`^player "([^"]*)" claims offer "([^"]*)" in game "([^"]*)"$`, playerClaimsOfferInGame)
	s.Step(`^the last request returned status code "([^"]*)" and body "([^"]*)"$`, theLastRequestReturnedStatusCodeStatusAndBody)
	s.Step(`^an offer template is created in the "([^"]*)" game with:$`, anOfferTemplateIsCreatedInTheGameWith)
	s.Step(`^an offer template with name "([^"]*)" exists in game "([^"]*)"$`, anOfferTemplateWithNameExistsInGame)
	s.Step(`^an offer template exists with name "([^"]*)" in game "([^"]*)"$`, anOfferTemplateExistsWithNameInGame)
	s.Step(`^an offer template with name "([^"]*)" does not exist in game "([^"]*)"$`, anOfferTemplateWithNameDoesNotExistInGame)
	s.Step(`^the game "([^"]*)" requests offers for player "([^"]*)" in "([^"]*)"$`, theGameRequestsOffersForPlayerIn)
	s.Step(`^an offer with name "([^"]*)" is returned$`, anOfferWithNameIsReturned)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has seen offer "([^"]*)"$`, playerHasSeenOffer)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has not seen offer "([^"]*)"$`, playerHasNotSeenOffer)
}
