// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/lib/pq"
	"github.com/topfreegames/offers/errors"
	"gopkg.in/mgutz/dat.v2/dat"

	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//ClaimOfferPayload has required fields for claiming an offer
type ClaimOfferPayload struct {
	GameID          string `json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	PlayerID        string `json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	ProductID       string `json:"productId" valid:"ascii,stringlength(1|255),required"`
	Timestamp       int64  `json:"timestamp" valid:"int64,required"`
	TransactionID   string `json:"transactionId" valid:"ascii,stringlength(1|1000),required"`
	OfferInstanceID string `json:"id" valid:"uuidv4,optional"`
}

//OfferImpressionPayload has required fields for an offer impression
type OfferImpressionPayload struct {
	GameID       string `json:"gameId" valid:"matches(^[^-][a-zA-Z0-9-_]*$),stringlength(1|255),required"`
	PlayerID     string `json:"playerId" valid:"ascii,stringlength(1|1000),required"`
	ImpressionID string `json:"impressionId" valid:"uuidv4,required"`
}

//GetTransactionsKey returns the key of the player's purchase transactions in redis
func GetTransactionsKey(gameID, playerID string) string {
	return fmt.Sprintf("transactions:%s:%s", gameID, playerID)
}

//GetClaimCounterKey returns the key of the player's claim counter
func GetClaimCounterKey(playerID, offerID string) string {
	return fmt.Sprintf("claim:counter:%s:%s", playerID, offerID)
}

//GetClaimTimestampKey returns the key of the player's last claim timestamp
func GetClaimTimestampKey(playerID, offerID string) string {
	return fmt.Sprintf("claim:timestamp:%s:%s", playerID, offerID)
}

//GetImpressionsKey returns the key of the player's impressions in redis
func GetImpressionsKey(playerID, gameID string) string {
	return fmt.Sprintf("impressions:%s:%s", gameID, playerID)
}

//GetViewCounterKey returns the key of the player's impressions counter
func GetViewCounterKey(playerID, offerID string) string {
	return fmt.Sprintf("view:counter:%s:%s", playerID, offerID)
}

//GetViewTimestampKey returns the key of the player's last impression timestamp
func GetViewTimestampKey(playerID, offerID string) string {
	return fmt.Sprintf("view:timestamp:%s:%s", playerID, offerID)
}

//GetEnabledOffersKey returns the key of the current enabled offers
func GetEnabledOffersKey(gameID string) string {
	return fmt.Sprintf("offers:enabled:%s", gameID)
}

//GetOfferInstancesKey returns the key of the players available offer instances
func GetOfferInstancesKey(offerInstanceID, gameID string) string {
	return fmt.Sprintf("offers:instances:%s:%s", gameID, offerInstanceID)
}

//GetDB Connection using the given properties
func GetDB(
	host string, user string, port int, sslmode string,
	dbName string, password string,
	maxIdleConns, maxOpenConns int,
	connectionTimeoutMS int,
) (runner.Connection, error) {
	if connectionTimeoutMS <= 0 {
		connectionTimeoutMS = 100
	}
	connStr := fmt.Sprintf(
		"host=%s user=%s port=%d sslmode=%s dbname=%s connect_timeout=2",
		host, user, port, sslmode, dbName,
	)
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	ShouldPing(db, time.Duration(connectionTimeoutMS)*time.Millisecond)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 100ms as warnings. (optional)
	runner.LogQueriesThreshold = 100 * time.Millisecond

	return runner.NewDB(db, "postgres"), nil
}

//IsNoRowsInResultSetError returns true if the error is a sqlx error stating that now rows were found
func IsNoRowsInResultSetError(err error) bool {
	return err.Error() == "sql: no rows in result set"
}

//IsForeignKeyViolationError returns true if the error is a pq error stating a foreign key has been violated
func IsForeignKeyViolationError(err error) (*pq.Error, bool) {
	var pqErr *pq.Error
	var ok bool

	if pqErr, ok = err.(*pq.Error); !ok {
		return nil, false
	}

	return pqErr, pqErr.Code == "23503" && strings.Contains(pqErr.Message, "violates foreign key constraint")
}

//ShouldPing the database
func ShouldPing(db *sql.DB, timeout time.Duration) error {
	var err error
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = timeout
	ticker := backoff.NewTicker(b)

	// Ticks will continue to arrive when the previous operation is still running,
	// so operations that take a while to fail could run in quick succession.
	for range ticker.C {
		if err = db.Ping(); err != nil {
			continue
		}

		ticker.Stop()
		return nil
	}

	return fmt.Errorf("could not ping database")
}

//HandleNotFoundError returns the proper error if nothing happens
func HandleNotFoundError(model string, filters map[string]interface{}, err error) error {
	if err != nil {
		if IsNoRowsInResultSetError(err) {
			return errors.NewModelNotFoundError(model, filters)
		}

		return err
	}
	return nil
}

//HandleForeignKeyViolationError returns the proper error if nothing happens
func HandleForeignKeyViolationError(model string, err error) error {
	if err != nil {
		if pqErr, ok := IsForeignKeyViolationError(err); ok {
			return errors.NewInvalidModelError(model, pqErr.Message)
		}
		return err
	}
	return nil
}
