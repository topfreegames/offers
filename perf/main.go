// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package main

import (
	"fmt"
	"time"

	edat "github.com/topfreegames/extensions/dat"
	bench "github.com/topfreegames/offers/bench"
	"github.com/topfreegames/offers/models"
	oTesting "github.com/topfreegames/offers/testing"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

func initParameters() {

}

func populateGames(db *runner.Connection) ([]*models.Game, error) {
	var games []*models.Game
	var err error
	query := `INSERT INTO games(
						id,
						name) 
					SELECT
						uuid_generate_v4(),
						('game-' || generate_series)
					FROM
						generate_series(1, $1)
					RETURNING
						id, name
					`
	builder := (*db).SQL(query, bench.NumberOfGames)
	builder.Execer = edat.NewExecer(builder.Execer)
	err = builder.QueryStructs(&games)

	return games, err
}

//CAST(to_jsonb($2::text) as jsonb),
func populateOffers(db *runner.Connection, games []*models.Game) (map[string][]*models.Offer, error) {
	offersByGame := make(map[string][]*models.Offer)
	var err error

	for _, game := range games {
		query := `INSERT INTO offers(
								game_id,
								name,
								period,
								frequency,
								trigger,
								placement,
								product_id,
								contents)
							SELECT
								$1,
								('offer-' || generate_series),
								'{"every": "1h"}',
								'{"every": "1h"}',
								$2,
								'popup',
								'tfg.com.sample',
								'{"x": 1}'
							FROM
								generate_series(1, $3)
							RETURNING
								*
							`
		to := time.Now().Unix() + 5*60
		trigger := fmt.Sprintf("{\"from\": 1486678000, \"to\": %d}", to)
		var offers []*models.Offer
		builder := (*db).SQL(query, game.ID, trigger, bench.NumberOfOffersPerGame)
		builder.Execer = edat.NewExecer(builder.Execer)
		err = builder.QueryStructs(&offers)
		if err != nil {
			return offersByGame, err
		}

		offersByGame[game.ID] = offers
	}

	return offersByGame, err
}

func populateOfferInstances(db *runner.Connection, offersByGame map[string][]*models.Offer) (map[string][]*models.OfferToReturn, error) {
	offerInstances := make(map[string][]*models.OfferToReturn)
	var err error

	return offerInstances, err
}

func populateTestDB(db *runner.Connection) error {
	games, err := populateGames(db)
	if err != nil {
		return err
	}

	offersByGame, err := populateOffers(db, games)
	if err != nil {
		return err
	}

	_, err = populateOfferInstances(db, offersByGame)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := oTesting.GetPerfDB()
	if err != nil {
		panic(err.Error())
	}

	err = populateTestDB(&db)
	if err != nil {
		panic(err.Error())
	}
}
