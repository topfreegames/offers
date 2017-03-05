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
	"regexp"
	"strconv"
	"strings"
	"time"

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
var lastBody, lastPlayerID, lastGameID, lastPlacement string
var lastOffers map[string][]*models.OfferToReturn
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

func SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID string) (*models.Offer, error) {
	query := `SELECT 
							offers.id, offers.game_id, offers.player_id, offers.seen_counter
						FROM 
							offers
						INNER JOIN offer_templates ON offer_templates.id = offers.offer_template_id
						WHERE offers.player_id = $1
							AND offers.game_id = $2
							AND offer_templates.name = $3`
	var offer models.Offer
	err := app.DB.SQL(query, playerID, gameID, offerTemplateName).QueryStruct(&offer)

	return &offer, err
}

func requestGameWithIDAndBundleID(id, bundleID string) error {
	var err error
	lastStatus, lastBody, err = performRequest(app, "PUT", fmt.Sprintf("/games/%s", replaceString(id)), map[string]interface{}{
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
		if strings.TrimSpace(description) == strings.TrimSpace(lastBody) {
			return nil
		}
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

func theFollowingPlayersExistInTheGame(gameID string, players *gherkin.DataTable) error {
	for i := 1; i < len(players.Rows); i++ {
		playerID := players.Rows[i].Cells[0].Value

		claimedOffers := strings.Split(players.Rows[i].Cells[1].Value, ", ")
		timestamps := strings.Split(players.Rows[i].Cells[2].Value, ", ")

		for j, offerTemplateName := range claimedOffers {
			if offerTemplateName != "-" {
				unixTime, err := strconv.Atoi(timestamps[j])
				if err != nil {
					return err
				}

				currentTime := time.Unix(int64(unixTime), 0)
				if _, err := models.GetAvailableOffers(app.DB, playerID, gameID, currentTime, nil); err != nil {
					return err
				}

				offer, err := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID)
				if err != nil {
					return err
				}

				err = models.UpdateOfferLastSeenAt(app.DB, offer.ID, playerID, gameID, currentTime, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
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
			Contents:  toJSON(otArgs.Rows[i].Cells[3].Value),
			Placement: otArgs.Rows[i].Cells[4].Value,
			Period:    toJSON(otArgs.Rows[i].Cells[5].Value),
			Frequency: toJSON(otArgs.Rows[i].Cells[6].Value),
			Trigger:   toJSON(otArgs.Rows[i].Cells[7].Value),
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
	var factor int
	var err error

	ok, regexErr := regexp.MatchString(`^\d+d$`, arg1)

	if regexErr != nil {
		return regexErr
	} else if ok {
		arg1 = arg1[:len(arg1)-1]
		factor = 24 * 60 * 60
	} else {
		factor = 1
	}

	intCurrentTime, err = strconv.ParseInt(arg1, 10, 64)

	if err != nil {
		return err
	}

	intCurrentTime *= int64(factor)

	mockClock := testing.MockClock{
		CurrentTime: intCurrentTime,
	}

	app.Clock = mockClock

	return nil
}

func playerClaimsOfferInGame(playerID, offerTemplateName, gameID string) error {
	var err error
	offer, _ := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID)

	lastStatus, lastBody, err = performRequest(app, "PUT", fmt.Sprintf("/offers/%s/claim", offer.ID), map[string]interface{}{
		"gameID":   gameID,
		"playerID": playerID,
	})

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

	body = strings.Replace(body, "'", "\"", -1)
	body = strings.Replace(body, " ", "", -1)
	lastTrimmedBody := strings.Replace(lastBody, " ", "", -1)

	if lastTrimmedBody != body {
		return fmt.Errorf("Expected last request to have body %s but it had %s", body, lastBody)
	}

	return nil
}

func theLastRequestReturnedStatusCodeAndBody(code, body string) error {
	intCode, err := strconv.Atoi(code)

	if err != nil {
		return err
	}

	return theLastRequestReturnedStatusCodeStatusAndBody(intCode, body)
}

func toJSON(str string) dat.JSON {
	if len(str) == 0 {
		return nil
	}

	return dat.JSON([]byte(strings.Replace(str, "'", "\"", -1)))
}

func anOfferTemplateIsCreatedInTheGameWithNamePidContentsMetadataPeriodFreqTriggerPlace(gameID, name, pid, contents, metadata, period, freq, trigger, place string) error {
	payload := map[string]interface{}{
		"gameId":    gameID,
		"name":      name,
		"productId": pid,
		"contents":  toJSON(contents),
		"metadata":  toJSON(metadata),
		"period":    toJSON(period),
		"frequency": toJSON(freq),
		"trigger":   toJSON(trigger),
		"placement": place,
	}
	var err error
	lastStatus, lastBody, err = performRequest(app, "POST", "/templates", payload)

	return err
}

func anOfferTemplateWithNameExistsInGame(offerTemplateName, gameID string) error {
	_, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil)

	if err != nil {
		return err
	}

	return nil
}

func anOfferTemplateExistsWithNameInGame(offerTemplateName, gameID string) error {
	if _, err := models.GetOfferTemplateByNameAndGame(app.DB, offerTemplateName, gameID, nil); err != nil {
		return insertOfferTemplate(app.DB, offerTemplateName, gameID)
	}

	return nil
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

	var newLastOffers map[string][]*models.OfferToReturn
	if err = json.Unmarshal([]byte(lastBody), &newLastOffers); err != nil {
		return err
	}

	lastOffers = newLastOffers
	lastPlayerID = playerID
	lastGameID = gameID
	lastPlacement = placement

	return nil
}

func playerOfGameHasSeenOffer(playerID, gameID, offerTemplateName string) error {
	offer, err := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID)

	if err != nil {
		return err
	}

	if offer.SeenCounter == 0 {
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

		lastStatus, lastBody, err = performRequest(app, "POST", fmt.Sprintf("/offers/%s/impressions", returnedOffer.ID), map[string]interface{}{
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
	offer, err := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID)

	if err != nil {
		return err
	}

	lastStatus, lastBody, err = performRequest(app, "POST", fmt.Sprintf("/offers/%s/impressions", offer.ID), map[string]interface{}{
		"ID":       offer.ID,
		"PlayerID": playerID,
		"GameID":   gameID,
	})

	if err != nil {
		return err
	}

	return nil
}

func anOfferWithNameIsReturned(offerTemplateName string) error {
	if offerTemplateName == "-" {
		if _, ok := lastOffers[lastPlacement]; ok {
			return fmt.Errorf("Expected no offers in placement %s", lastPlacement)
		}

		return nil
	}

	offer, err := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, lastPlayerID, lastGameID)

	if err != nil {
		return err
	}

	for _, returnedOffer := range lastOffers[lastPlacement] {
		if returnedOffer.ID == offer.ID {
			return nil
		}
	}

	return fmt.Errorf("Expected offer %s to be returned", offerTemplateName)
}

func theFollowingPlayersClaimedInTheGame(gameID string, players *gherkin.DataTable) error {
	for i := 1; i < len(players.Rows); i++ {
		playerID := players.Rows[i].Cells[0].Value

		claimedOffers := strings.Split(players.Rows[i].Cells[1].Value, ", ")
		timestamps := strings.Split(players.Rows[i].Cells[2].Value, ", ")

		for j, offerTemplateName := range claimedOffers {
			if offerTemplateName != "-" {
				unixTime, err := strconv.Atoi(timestamps[j])
				if err != nil {
					return err
				}

				currentTime := time.Unix(int64(unixTime), 0)
				if _, err := models.GetAvailableOffers(app.DB, playerID, gameID, currentTime, nil); err != nil {
					return err
				}

				offer, err := SelectOfferByOfferTemplateNameAndPlayerAndGame(offerTemplateName, playerID, gameID)
				if err != nil {
					return err
				}

				_, _, err = models.ClaimOffer(app.DB, offer.ID, playerID, gameID, currentTime, nil)
				if err != nil {
					return err
				}
			}
		}
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
	s.Step(`^an offer template exists with name "([^"]*)" in game "([^"]*)"$`, anOfferTemplateExistsWithNameInGame)
	s.Step(`^an offer template with name "([^"]*)" does not exist in game "([^"]*)"$`, anOfferTemplateWithNameDoesNotExistInGame)
	s.Step(`^the following offer templates exist in the "([^"]*)" game:$`, theFollowingOfferTemplatesExistInTheGame)
	s.Step(`^the current time is "([^"]*)"$`, theCurrentTimeIs)
	s.Step(`^the game "([^"]*)" requests offers for player "([^"]*)" in "([^"]*)"$`, theGameRequestsOffersForPlayerIn)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has seen offer "([^"]*)"$`, playerOfGameHasSeenOffer)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has not seen offer "([^"]*)"$`, playerOfGameHasNotSeenOffer)
	s.Step(`^the current time is (\d+)$`, theCurrentTimeIs)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer in "([^"]*)"$`, thePlayerOfGameSeesOfferIn)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer with name "([^"]*)"$`, thePlayerOfGameSeesOfferWithName)
	s.Step(`^the following players exist in the "([^"]*)" game:$`, theFollowingPlayersExistInTheGame)
	s.Step(`^player "([^"]*)" claims offer "([^"]*)" in game "([^"]*)"$`, playerClaimsOfferInGame)
	s.Step(`^an offer template with name "([^"]*)" exists in game "([^"]*)"$`, anOfferTemplateWithNameExistsInGame)
	s.Step(`^an offer template is created in the "([^"]*)" game with name "([^"]*)" pid "([^"]*)" contents "([^"]*)" metadata "([^"]*)" period "([^"]*)" freq "([^"]*)" trigger "([^"]*)" place "([^"]*)"$`, anOfferTemplateIsCreatedInTheGameWithNamePidContentsMetadataPeriodFreqTriggerPlace)
	s.Step(`^the last request returned status code "([^"]*)" and body "([^"]*)"$`, theLastRequestReturnedStatusCodeAndBody)
	s.Step(`^an offer with name "([^"]*)" is returned$`, anOfferWithNameIsReturned)
	s.Step(`^the following players claimed in the "([^"]*)" game:$`, theFollowingPlayersClaimedInTheGame)
}
