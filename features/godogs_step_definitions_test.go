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
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

var app *api.App
var logger logrus.Logger
var lastStatus int
var lastBody string

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

	var err error
	app, err = api.NewApp("localhost", 9999, config, true, logger)
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
		return fmt.Errorf("Expected game to have bundle ID of %s, but it has %s.", bundleID, game.BundleID)
	}
	return nil
}

func theLastRequestReturnedStatusCode(statusCode int) error {
	if lastStatus != statusCode {
		return fmt.Errorf("Expected last request to have status code of %d but it had %d.", statusCode, lastStatus)
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
	return fmt.Errorf("The game %s should not exist but it does.", game.ID)
}

func aGameWithNameExists(arg1 string) error {
	return godog.ErrPending
}

func theFollowingOfferTemplatesExistInTheGame(arg1 string, arg2 *gherkin.DataTable) error {
	return godog.ErrPending
}

func theFollowingPlayersExistInTheGame(arg1 string, arg2 *gherkin.DataTable) error {
	return godog.ErrPending
}

func theCurrentTimeIs(arg1 int) error {
	return godog.ErrPending
}

func playerClaimsOfferInGame(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func theLastRequestReturnedStatusCodeStatusAndBody(arg1 string) error {
	return godog.ErrPending
}

func theCurrentTimeIsD(arg1 int) error {
	return godog.ErrPending
}

func anOfferTemplateIsCreatedInTheGameWith(arg1 string, arg2 *gherkin.DataTable) error {
	return godog.ErrPending
}

func anOfferTemplateWithNameExistsInGame(arg1, arg2 string) error {
	return godog.ErrPending
}

func anOfferTemplateExistsWithNameInGame(arg1, arg2 string) error {
	return godog.ErrPending
}

func anOfferTemplateWithNameDoesNotExistInGame(arg1, arg2 string) error {
	return godog.ErrPending
}

func theGameRequestsOffersForPlayerIn(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func anOfferWithNameIsReturned(arg1 string) error {
	return godog.ErrPending
}

func playerHasSeenOffer(arg1, arg2 string) error {
	return godog.ErrPending
}

func playerHasNotSeenOffer(arg1, arg2 string) error {
	return godog.ErrPending
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
	s.Step(`^the last request returned status code <status> and body "([^"]*)"$`, theLastRequestReturnedStatusCodeStatusAndBody)
	s.Step(`^the current time is (\d+)d$`, theCurrentTimeIsD)
	s.Step(`^an offer template is created in the "([^"]*)" game with:$`, anOfferTemplateIsCreatedInTheGameWith)
	s.Step(`^an offer template with name "([^"]*)" exists in game "([^"]*)"$`, anOfferTemplateWithNameExistsInGame)
	s.Step(`^an offer template exists with name "([^"]*)" in game "([^"]*)"$`, anOfferTemplateExistsWithNameInGame)
	s.Step(`^an offer template with name "([^"]*)" does not exist in game "([^"]*)"$`, anOfferTemplateWithNameDoesNotExistInGame)
	s.Step(`^the game "([^"]*)" requests offers for player "([^"]*)" in "([^"]*)"$`, theGameRequestsOffersForPlayerIn)
	s.Step(`^an offer with name "([^"]*)" is returned$`, anOfferWithNameIsReturned)
	s.Step(`^player "([^"]*)" has seen offer "([^"]*)"$`, playerHasSeenOffer)
	s.Step(`^player "([^"]*)"> has not seen offer "([^"]*)"$`, playerHasNotSeenOffer)
}
