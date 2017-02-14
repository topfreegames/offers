// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/GuiaBolso/darwin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	migrations "github.com/topfreegames/offers/migrations"
	"github.com/topfreegames/offers/models"
)

var migrationInfo bool

// MigrationError identified rigrations running error
type MigrationError struct {
	Message string
}

func (err *MigrationError) Error() string {
	return fmt.Sprintf("Could not run migrations: %s", err.Message)
}

func getVersion(migName string) float64 {
	parts := strings.Split(filepath.Base(migName), "-")
	migNumber, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Panicf("Failed to parse migration name: %s (error: %s)", migName, err.Error())
	}
	return migNumber
}

func getDescription(migName string) string {
	parts := strings.Split(filepath.Base(migName), "-")
	return parts[1]
}

func getMigrations() []darwin.Migration {
	migNames := migrations.AssetNames()
	sort.Sort(sort.StringSlice(migNames))
	migs := make([]darwin.Migration, len(migNames))

	for i, migName := range migNames {
		contents, err := migrations.Asset(migName)
		if err != nil {
			panic(fmt.Sprintf("Could not read migration %s!", migName))
		}
		migs[i] = darwin.Migration{
			Version:     getVersion(migName),
			Description: getDescription(migName),
			Script:      string(contents),
		}
	}

	return migs
}

func getDB() (*sql.DB, error) {
	host := viper.GetString("postgres.host")
	user := viper.GetString("postgres.user")
	dbName := viper.GetString("postgres.dbname")
	password := viper.GetString("postgres.password")
	port := viper.GetInt("postgres.port")
	sslMode := viper.GetString("postgres.sslMode")
	maxIdleConns := viper.GetInt("postgres.maxIdleConns")
	maxOpenConns := viper.GetInt("postgres.maxOpenConns")

	return models.GetDB(host, user, port, sslMode, dbName, password, maxIdleConns, maxOpenConns)
}

func printStatus(d darwin.Darwin) error {
	infos, err := d.Info()
	if err != nil {
		return err
	}
	fmt.Println("")
	fmt.Println("Current database migrations status")
	fmt.Println("==================================")
	fmt.Println("")
	fmt.Println("Version  | Status          | Name")
	for _, info := range infos {
		status := info.Status.String()
		for i := 0; i < 15-len(info.Status.String()); i++ {
			status += " "
		}
		fmt.Printf("%.1f      | %s | %s\n", info.Migration.Version, status, info.Migration.Description)
	}
	fmt.Println("")

	return nil
}

//RunMigrations in selected DB
func RunMigrations(info bool) error {
	migrations := getMigrations()

	database, err := getDB()

	if err != nil {
		log.Fatal(err)
	}

	driver := darwin.NewGenericDriver(database, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	if info {
		err := printStatus(d)
		if err != nil {
			return err
		}
	} else {
		log.Println("Migrating database to latest version...")
		err = d.Migrate()

		if err != nil {
			return err
		}

		printStatus(d)
		log.Println("Database migrated successfully.")
	}
	return nil
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrates the database up or down",
	Long:  `Migrate the database specified in the configuration file to the given version (or latest if none provided)`,
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		err := RunMigrations(migrationInfo)
		if err != nil {
			log.Println(err)
			panic(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVarP(&migrationInfo, "info", "i", false, "Get database info")
}
