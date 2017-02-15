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
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func theServerIsUp() error {
	return godog.ErrPending
}

func aGameNamedIsCreatedWithBundleIDOf(arg1, arg2 string) error {
	return godog.ErrPending
}

func theGameExists(arg1 string) error {
	return godog.ErrPending
}

func theGameHasBundleIDOf(arg1, arg2 string) error {
	return godog.ErrPending
}

func aGameNamedIsUpdatedWithBundleIDOf(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastRequestReturnedStatusCode(arg1 int) error {
	return godog.ErrPending
}

func theLastErrorIsWithMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

func theGameDoesNotExist(arg1 string) error {
	return godog.ErrPending
}

func theServerIs(arg1 string) error {
	return godog.ErrPending
}

func theHealthCheckIsDone() error {
	return godog.ErrPending
}

func theLastRequestReturnedStatusCodeAndBody(arg1 int, arg2 string) error {
	return godog.ErrPending
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
	s.Step(`^the server is "([^"]*)"$`, theServerIs)
	s.Step(`^the health check is done$`, theHealthCheckIsDone)
	s.Step(`^the last request returned status code (\d+) and body "([^"]*)"$`, theLastRequestReturnedStatusCodeAndBody)
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
