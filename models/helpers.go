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

	_ "github.com/lib/pq" //This is required to use postgres with database/sql
)

//GetDB Connection using the given properties
func GetDB(host string, user string, port int, sslmode string, dbName string, password string, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
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

	return db, nil
}
