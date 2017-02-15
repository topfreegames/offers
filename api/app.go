// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"io"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	runner "github.com/mgutz/dat/sqlx-runner"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/metadata"
	"github.com/topfreegames/offers/models"
)

//App is our API application
type App struct {
	Address string
	Router  *mux.Router
	Server  *http.Server
	Config  *viper.Viper
	DB      runner.Connection
	Logger  logrus.FieldLogger
}

//NewApp ctor
func NewApp(config *viper.Viper) (*App, error) {
	a := &App{
		Config: config,
	}
	err := a.configureApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) getRouter() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/healthcheck", Chain(
		&HealthcheckHandler{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		AllowedMethods("GET"),
	)).Name("healthcheck")

	r.Handle("/game/upsert", Chain(
		&GameHandler{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Name("game")

	return r
}

func (a *App) configureApp() error {
	err := a.configureDatabase()
	if err != nil {
		return err
	}
	a.configureLogger()
	a.configureServer()
	return nil
}

func (a *App) configureDatabase() error {
	db, err := a.getDB()
	if err != nil {
		return err
	}

	a.DB = db
	return nil
}

func (a *App) getDB() (runner.Connection, error) {
	host := a.Config.GetString("postgres.host")
	user := a.Config.GetString("postgres.user")
	dbName := a.Config.GetString("postgres.dbname")
	password := a.Config.GetString("postgres.password")
	port := a.Config.GetInt("postgres.port")
	sslMode := a.Config.GetString("postgres.sslMode")
	maxIdleConns := a.Config.GetInt("postgres.maxIdleConns")
	maxOpenConns := a.Config.GetInt("postgres.maxOpenConns")
	connectionTimeoutMS := viper.GetInt("postgres.connectionTimeoutMS")

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

func (a *App) configureLogger() {
	a.Logger = logrus.New().WithFields(logrus.Fields{
		"app":     "offers-api",
		"version": metadata.Version,
	})
}

func (a *App) configureServer() {
	a.Router = a.getRouter()
	a.Server = &http.Server{Addr: a.Address, Handler: a.Router}
}

//ListenAndServe requests
func (a *App) ListenAndServe() (io.Closer, error) {
	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		return nil, err
	}

	err = a.Server.Serve(listener)
	if err != nil {
		listener.Close()
		return nil, err
	}

	return listener, nil
}
