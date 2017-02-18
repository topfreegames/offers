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
	"github.com/mgutz/dat"
	"github.com/topfreegames/offers/errors"

	runner "github.com/mgutz/dat/sqlx-runner"
)

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
