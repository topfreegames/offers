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
	"time"

	_ "github.com/lib/pq" //This is required to use postgres with database/sql
	"github.com/mgutz/dat"
	runner "github.com/mgutz/dat/sqlx-runner"
)

//GetDB Connection using the given properties
func GetDB(
	host string, user string, port int, sslmode string,
	dbName string, password string,
	maxIdleConns, maxOpenConns int,
) (runner.Connection, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s port=%d sslmode=%s dbname=%s",
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

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	return runner.NewDB(db, "postgres"), nil
}

//IsNoRowsInResultSetError returns true if the error is a sqlx error stating that now rows were found
func IsNoRowsInResultSetError(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
