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

package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func theServerIsUp() error {
	return godog.ErrPending
}

func theUserCreatesAGameNamedWithBundleIDOf(arg1, arg2 string) error {
	return godog.ErrPending
}

func theGameExists(arg1 string) error {
	return godog.ErrPending
}

func theGameHasBundleIDOf(arg1, arg2 string) error {
	return godog.ErrPending
}

func theUserUpdatesAGameNamedWithBundleIDOf(arg1, arg2 string) error {
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

func anOfferIsCreatedWith(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func anOfferWithNameExists(arg1 string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the server is up$`, theServerIsUp)
	s.Step(`^the user creates a game named "([^"]*)" with bundle id of "([^"]*)"$`, theUserCreatesAGameNamedWithBundleIDOf)
	s.Step(`^the game "([^"]*)" exists$`, theGameExists)
	s.Step(`^the game "([^"]*)" has bundle id of "([^"]*)"$`, theGameHasBundleIDOf)
	s.Step(`^the user updates a game named "([^"]*)" with bundle id of "([^"]*)"$`, theUserUpdatesAGameNamedWithBundleIDOf)
	s.Step(`^the last request returned status code (\d+)$`, theLastRequestReturnedStatusCode)
	s.Step(`^the last error is "([^"]*)" with message "([^"]*)"$`, theLastErrorIsWithMessage)
	s.Step(`^the game "([^"]*)" does not exist$`, theGameDoesNotExist)
	s.Step(`^the server is "([^"]*)"$`, theServerIs)
	s.Step(`^the health check is done$`, theHealthCheckIsDone)
	s.Step(`^the last request returned status code (\d+) and body "([^"]*)"$`, theLastRequestReturnedStatusCodeAndBody)
	s.Step(`^a game with name "([^"]*)" exists$`, aGameWithNameExists)
	s.Step(`^an offer is created with:$`, anOfferIsCreatedWith)
	s.Step(`^an offer with name "([^"]*)" exists$`, anOfferWithNameExists)
}
