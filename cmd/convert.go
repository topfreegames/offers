// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2018 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	"github.com/topfreegames/offers/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

func getDBForConvert() (runner.Connection, error) {
	host := config.GetString("postgres.host")
	user := config.GetString("postgres.user")
	dbName := config.GetString("postgres.dbname")
	password := config.GetString("postgres.password")
	port := config.GetInt("postgres.port")
	sslMode := config.GetString("postgres.sslMode")
	maxIdleConns := config.GetInt("postgres.maxIdleConns")
	maxOpenConns := config.GetInt("postgres.maxOpenConns")
	connectionTimeoutMS := config.GetInt("postgres.connectionTimeoutMS")

	db, err := models.GetDB(
		host, user, port, sslMode, dbName,
		password, maxIdleConns, maxOpenConns,
		connectionTimeoutMS,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getOffers(db runner.Connection, limit, offset int) ([]*models.Offer, error) {
	var offers []*models.Offer
	builder := db.SQL(fmt.Sprintf("SELECT * FROM offers ORDER BY id LIMIT %d OFFSET %d", limit, offset))
	err := builder.QueryStructs(&offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func createOfferVersions(db runner.Connection, offers []*models.Offer) error {
	builder := db.InsertInto("offer_versions").Columns("game_id", "offer_id", "offer_version", "contents", "product_id", "cost")
	for _, offer := range offers {
		offerVersion := &models.OfferVersion{
			GameID:       offer.GameID,
			OfferID:      offer.ID,
			OfferVersion: offer.Version,
			Contents:     offer.Contents,
			ProductID:    offer.ProductID,
			Cost:         offer.Cost,
		}
		builder.Record(offerVersion)
	}
	_, err := builder.Exec()
	return err
}

//RunConversion in selected DB
func RunConversion(writer io.Writer) error {
	database, err := getDBForConvert()

	if err != nil {
		log.Fatal(err)
	}
	cnt := 1
	total := 0
	for cnt > 0 {
		offers, err := getOffers(database, 300, total)
		if err != nil {
			log.Fatal(err)
		}
		cnt = len(offers)
		if cnt <= 0 {
			break
		}
		total += cnt
		err = createOfferVersions(database, offers)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(cnt, total)
	}
	log.Println(total)

	return nil
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "converts the offers",
	Long:  `Converts the offers from offer instances to offer versions`,
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		err := RunConversion(nil)
		if err != nil {
			log.Println(err)
			panic(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)
}
