// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/mgutz/dat"

	"github.com/jmoiron/sqlx/types"
)

//Game represents a tenant in offers API
type Game struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Metadata  types.JSONText `db:"metadata"`
	CreatedAt dat.NullTime   `db:"created_at"`
	UpdatedAt dat.NullTime   `db:"updated_at"`
}
