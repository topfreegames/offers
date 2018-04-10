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
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	edat "github.com/topfreegames/extensions/dat"
	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
	"github.com/topfreegames/offers/testing"
)

var app *api.App
var logger logrus.Logger
var lastStatus int
var lastBody, lastPlayerID, lastGameID, lastPlacement string
var lastOfferInstances map[string][]*models.OfferToReturn
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

func requestGameWithIDAndName(id, name string) error {
	var err error
	lastStatus, lastBody, err = performRequest(app, "PUT", fmt.Sprintf("/games/%s", replaceString(id)), map[string]interface{}{
		"id":   id,
		"name": replaceString(name),
	})

	return err
}

//ALIAS to update
func aGameWithIDIsCreatedWithAName(id, name string) error {
	return requestGameWithIDAndName(id, name)
}

func aGameWithIDIsUpdatedWithAName(id, name string) error {
	return requestGameWithIDAndName(id, name)
}

func theGameExists(id string) error {
	_, err := models.GetGameByID(app.DB, id, nil)
	return err
}

func theGameHasNameOf(id, name string) error {
	game, err := models.GetGameByID(app.DB, id, nil)
	if err != nil {
		return err
	}
	if game.Name != name {
		return fmt.Errorf("Expected game to have name of %s, but it has %s", name, game.Name)
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

		for j, offerName := range claimedOffers {
			if offerName != "-" {
				unixTime, err := strconv.Atoi(timestamps[j])
				if err != nil {
					return err
				}

				currentTime := time.Unix(int64(unixTime), 0)
				if _, err := models.GetAvailableOffers(app.DB, gameID, playerID, currentTime, nil); err != nil {
					return err
				}
				offer, err := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, playerID, gameID)
				if err != nil {
					return err
				}

				_, _, err = models.ViewOffer(app.DB, gameID, offer.ID, playerID, uuid.NewV4().String(), currentTime, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func aGameWithIDExists(id string) error {
	_, err := newGame(app.DB, id)
	return err
}

func theFollowingOffersExistInTheGame(gameID string, otArgs *gherkin.DataTable) error {
	for i := 1; i < len(otArgs.Rows); i++ {
		offer := &models.Offer{
			Name:      otArgs.Rows[i].Cells[1].Value,
			ProductID: otArgs.Rows[i].Cells[2].Value,
			Contents:  toJSON(otArgs.Rows[i].Cells[3].Value),
			Placement: otArgs.Rows[i].Cells[4].Value,
			Period:    toJSON(otArgs.Rows[i].Cells[5].Value),
			Frequency: toJSON(otArgs.Rows[i].Cells[6].Value),
			Trigger:   toJSON(otArgs.Rows[i].Cells[7].Value),
			GameID:    gameID,
		}

		query := `SELECT id, frequency
							FROM offers
							WHERE name = $1 AND game_id = $2
							LIMIT 1`
		var dbOffer models.Offer
		builder := app.DB.SQL(query, offer.Name, offer.GameID, offer.ProductID)
		builder.Execer = edat.NewExecer(builder.Execer)
		err := builder.QueryStruct(&dbOffer)
		if err != nil {
			if _, err := models.InsertOffer(app.DB, offer, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func theCurrentTimeIs(timestamp string) error {
	var intCurrentTime int64
	var factor int64
	var err error

	unit := timestamp[len(timestamp)-1]

	switch unit {
	case 'd':
		factor = 24 * 60 * 60
		timestamp = timestamp[:len(timestamp)-1]
	default:
		factor = 1
	}

	intCurrentTime, err = strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		return err
	}

	intCurrentTime *= factor

	mockClock := testing.MockClock{
		CurrentTime: intCurrentTime,
	}

	app.Clock = mockClock

	return nil
}

func playerClaimsOfferInstanceInGame(playerID, offerName, gameID string) error {
	var err error
	offerInstance, _ := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, playerID, gameID)

	lastStatus, lastBody, err = performRequest(app, "PUT", "/offers/claim", map[string]interface{}{
		"gameID":        gameID,
		"playerID":      playerID,
		"productId":     "com.tfg.sample",
		"timestamp":     app.Clock.GetTime().Unix(),
		"transactionId": uuid.NewV4().String(),
		"id":            offerInstance.ID,
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

func anOfferIsCreatedInTheGameWithNamePidContentsMetadataPeriodFreqTriggerPlace(gameID, name, pid, contents, metadata, period, freq, trigger, place string) error {
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
	lastStatus, lastBody, err = performRequest(app, "POST", "/offers", payload)

	return err
}

func anOfferWithNameExistsInGame(offerName, gameID string) error {
	var offer models.Offer
	query := "SELECT id FROM offers WHERE name = $1 AND enabled = true"
	builder := app.DB.SQL(query, offerName)
	builder.Execer = edat.NewExecer(builder.Execer)
	err := builder.QueryStruct(&offer)

	if err != nil {
		return err
	}

	return nil
}

func anOfferExistsWithNameInGame(offerName, gameID string) error {
	if err := anOfferWithNameExistsInGame(offerName, gameID); err != nil {
		return insertOffer(app.DB, offerName, gameID)
	}

	return nil
}

func anOfferWithNameDoesNotExistInGame(offerName, gameID string) error {
	if err := anOfferWithNameExistsInGame(offerName, gameID); err == nil {
		return fmt.Errorf("Expected offer %s to not exist in game %s", offerName, gameID)
	}
	return nil
}

func theGameRequestsOfferInstancesForPlayerIn(gameID, playerID, placement string) error {
	var err error

	url := "/available-offers?player-id=" + playerID + "&game-id=" + gameID
	lastStatus, lastBody, err = performRequest(app, "GET", url, nil)

	if err != nil {
		return err
	}

	var newLastOffers map[string][]*models.OfferToReturn
	if err = json.Unmarshal([]byte(lastBody), &newLastOffers); err != nil {
		return err
	}

	lastOfferInstances = newLastOffers
	lastPlayerID = playerID
	lastGameID = gameID
	lastPlacement = placement

	return nil
}

func playerOfGameHasSeenOfferInstance(playerID, gameID, offerName string) error {
	offerInstance, err := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, playerID, gameID)

	if err != nil {
		return err
	}

	offerPlayer, err := models.GetOfferPlayer(app.DB, gameID, playerID, offerInstance.OfferID, nil)
	if err != nil && !models.IsNoRowsInResultSetError(err) {
		return err
	}
	if offerPlayer.ViewCounter == 0 {
		return fmt.Errorf("Expected player %s of game %s to has seen offer %s", playerID, gameID, offerName)
	}

	return nil
}

func playerOfGameHasNotSeenOfferInstance(playerID, gameID, offerName string) error {
	err := playerOfGameHasSeenOfferInstance(playerID, gameID, offerName)

	if err == nil {
		return fmt.Errorf("Expected player %s of game %s not to has seen offer %s", playerID, gameID, offerName)
	}

	return nil
}

func thePlayerOfGameSeesOfferInstanceIn(playerID, gameID, placement string) error {
	var err error
	for _, returnedOffer := range lastOfferInstances[placement] {
		var offerInstance *models.OfferInstance
		offerInstance, err = models.GetOfferInstanceByID(app.DB, gameID, returnedOffer.ID, nil)
		if err != nil {
			return err
		}

		lastStatus, lastBody, err = performRequest(app, "PUT", fmt.Sprintf("/offers/%s/impressions", offerInstance.ID), map[string]interface{}{
			"playerId":     playerID,
			"gameId":       gameID,
			"impressionId": uuid.NewV4().String(),
		})

	}

	return nil
}

func thePlayerOfGameSeesOfferInstanceWithName(playerID, gameID, offerName string) error {
	offerInstance, err := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, playerID, gameID)

	if err != nil {
		return err
	}

	lastStatus, lastBody, err = performRequest(app, "PUT", fmt.Sprintf("/offers/%s/impressions", offerInstance.ID), map[string]interface{}{
		"playerId":     playerID,
		"gameId":       gameID,
		"impressionId": uuid.NewV4().String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func anOfferInstanceWithNameIsReturned(offerName string) error {
	if offerName == "-" {
		if _, ok := lastOfferInstances[lastPlacement]; ok {
			return fmt.Errorf("Expected no offers in placement %s", lastPlacement)
		}

		return nil
	}

	offerInstance, err := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, lastPlayerID, lastGameID)

	if err != nil {
		return err
	}

	for _, returnedOffer := range lastOfferInstances[lastPlacement] {
		if returnedOffer.ID == offerInstance.ID {
			return nil
		}
	}

	return fmt.Errorf("Expected offer %s to be returned", offerName)
}

func theFollowingPlayersClaimedInTheGame(gameID string, players *gherkin.DataTable) error {
	for i := 1; i < len(players.Rows); i++ {
		playerID := players.Rows[i].Cells[0].Value

		claimedOffers := strings.Split(players.Rows[i].Cells[1].Value, ", ")
		timestamps := strings.Split(players.Rows[i].Cells[2].Value, ", ")

		for j, offerName := range claimedOffers {
			if offerName != "-" {
				unixTime, err := strconv.Atoi(timestamps[j])
				if err != nil {
					return err
				}

				currentTime := time.Unix(int64(unixTime), 0)
				if _, err = models.GetAvailableOffers(app.DB, gameID, playerID, currentTime, nil); err != nil {
					return err
				}
				offerInstance, err := selectOfferInstanceByOfferNameAndPlayerAndGame(offerName, playerID, gameID)
				if err != nil {
					return err
				}

				_, _, _, err = models.ClaimOffer(app.DB, gameID, offerInstance.ID, playerID, "", uuid.NewV4().String(), currentTime.Unix(), currentTime, nil)
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
	s.Step(`^a game with id "([^"]*)" is created with a name "([^"]*)"$`, aGameWithIDIsCreatedWithAName)
	s.Step(`^the game "([^"]*)" exists$`, theGameExists)
	s.Step(`^the game "([^"]*)" has name of "([^"]*)"$`, theGameHasNameOf)
	s.Step(`^a game with id "([^"]*)" is updated with a name "([^"]*)"$`, aGameWithIDIsUpdatedWithAName)
	s.Step(`^the last request returned status code (\d+)$`, theLastRequestReturnedStatusCode)
	s.Step(`^the last error is "([^"]*)" with message "([^"]*)"$`, theLastErrorIsWithMessage)
	s.Step(`^the game "([^"]*)" does not exist$`, theGameDoesNotExist)
	s.Step(`^a game with id "([^"]*)" exists$`, aGameWithIDExists)
	s.Step(`^an offer exists with name "([^"]*)" in game "([^"]*)"$`, anOfferExistsWithNameInGame)
	s.Step(`^an offer with name "([^"]*)" does not exist in game "([^"]*)"$`, anOfferWithNameDoesNotExistInGame)
	s.Step(`^the following offers exist in the "([^"]*)" game:$`, theFollowingOffersExistInTheGame)
	s.Step(`^the current time is "([^"]*)"$`, theCurrentTimeIs)
	s.Step(`^the game "([^"]*)" requests offer instances for player "([^"]*)" in "([^"]*)"$`, theGameRequestsOfferInstancesForPlayerIn)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has seen offer instance "([^"]*)"$`, playerOfGameHasSeenOfferInstance)
	s.Step(`^player "([^"]*)" of game "([^"]*)" has not seen offer instance "([^"]*)"$`, playerOfGameHasNotSeenOfferInstance)
	s.Step(`^the current time is (\d+)$`, theCurrentTimeIs)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer instance in "([^"]*)"$`, thePlayerOfGameSeesOfferInstanceIn)
	s.Step(`^the player "([^"]*)" of game "([^"]*)" sees offer instance with name "([^"]*)"$`, thePlayerOfGameSeesOfferInstanceWithName)
	s.Step(`^the following players exist in the "([^"]*)" game:$`, theFollowingPlayersExistInTheGame)
	s.Step(`^player "([^"]*)" claims offer instance "([^"]*)" in game "([^"]*)"$`, playerClaimsOfferInstanceInGame)
	s.Step(`^an offer with name "([^"]*)" exists in game "([^"]*)"$`, anOfferWithNameExistsInGame)
	s.Step(`^the last request returned status code "([^"]*)" and body "([^"]*)"$`, theLastRequestReturnedStatusCodeAndBody)
	s.Step(`^the following players claimed in the "([^"]*)" game:$`, theFollowingPlayersClaimedInTheGame)
	s.Step(`^an offer instance with name "([^"]*)" is returned$`, anOfferInstanceWithNameIsReturned)
	s.Step(`^an offer is created in the "([^"]*)" game with name "([^"]*)" pid "([^"]*)" contents "([^"]*)" metadata "([^"]*)" period "([^"]*)" freq "([^"]*)" trigger "([^"]*)" place "([^"]*)"$`, anOfferIsCreatedInTheGameWithNamePidContentsMetadataPeriodFreqTriggerPlace)
}
