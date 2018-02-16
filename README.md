# Offers

[![Build Status](https://travis-ci.org/topfreegames/offers.svg?branch=master)](https://travis-ci.org/topfreegames/offers)
[![Coverage Status](https://coveralls.io/repos/github/topfreegames/offers/badge.svg?branch=master)](https://coveralls.io/github/topfreegames/offers?branch=master)
[![Code Climate](https://codeclimate.com/github/topfreegames/offers/badges/gpa.svg)](https://codeclimate.com/github/topfreegames/offers)
[![Go Report Card](https://goreportcard.com/badge/github.com/topfreegames/offers)](https://goreportcard.com/report/github.com/topfreegames/offers)
[![Docs](https://readthedocs.org/projects/offers-api/badge/?version=latest
)](http://offers-api.readthedocs.io/en/latest/)
[![](https://imagelayers.io/badge/tfgco/offers:latest.svg)](https://imagelayers.io/?images=tfgco/offers:latest 'Offers Image Layers')

Offers is a service meant to handle offers/promotions in your games.

### Dependencies
* Go 1.7
* Postgres >= 9.5

### Setup
First, set your $GOPATH ([Go Lang](https://golang.org/doc/install)) env variable and add $GOPATH/bin to your $PATH

```bash
make setup
```

### Building

```bash
make build
```

### Running

```bash
make run-full
```

### Automated tests

Offers has unit, integration and acceptance tests (using cucumber). To run all of them:

```bash
make test
```

### Available Environment variables

Offers uses PostgreSQL to store offers information. The container takes environment variables to specify this connection:

* `OFFERS_POSTGRES_HOST` - PostgreSQL host to connect to;
* `OFFERS_POSTGRES_PORT` - PostgreSQL port to connect to;
* `OFFERS_POSTGRES_DBNAME` - PostgreSQL database to connect to;
* `OFFERS_POSTGRES_PASSWORD` - Password of the PostgreSQL Server to connect to;
* `OFFERS_POSTGRES_USER` - PostgreSQL user;

Offers uses basic auth to restrict access to routes that are not used directly by a client consuming the offers.

* `OFFERS_BASICAUTH_USERNAME` - Basic Auth user;
* `OFFERS_BASICAUTH_PASSWORD` - Basic Auth password;

When a client requests the available offers the API returns a `max-age` header. The cache TTL (in seconds) can be defined using the following variable:

* `OFFERS_CACHE_MAXAGESECONDS` - Max age in seconds;

Other than that, there are a couple more configurations you can pass using environment variables:

* `OFFERS_NEWRELIC_KEY` - If you have a [New Relic](https://newrelic.com/) account, you can use this variable to specify your API Key to populate data with New Relic API;
* `OFFERS_NEWRELIC_APP` - Name of the NewRelic app ;
* `OFFERS_SENTRY_URL` - If you have a [sentry server](https://docs.getsentry.com/hosted/) you can use this variable to specify your project's URL to send errors to.
