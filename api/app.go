// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/metadata"
	"github.com/topfreegames/offers/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//App is our API application
type App struct {
	Address  string
	Debug    bool
	Router   *mux.Router
	Server   *http.Server
	Config   *viper.Viper
	DB       runner.Connection
	Logger   logrus.FieldLogger
	NewRelic newrelic.Application
	Clock    models.Clock
}

//NewApp ctor
func NewApp(host string, port int, config *viper.Viper, debug bool, logger logrus.FieldLogger, clock models.Clock) (*App, error) {
	if clock == nil {
		clock = &models.RealClock{}
	}
	a := &App{
		Config:  config,
		Address: fmt.Sprintf("%s:%d", host, port),
		Debug:   debug,
		Logger:  logger,
		Clock:   clock,
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
		&MetricsReporterMiddleware{App: a},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("healthcheck")

	r.Handle("/games", Chain(
		&GameHandler{App: a},
		&MetricsReporterMiddleware{App: a},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&ValidationMiddleware{GetPayload: func() interface{} { return &models.Game{} }},
	)).Methods("PUT").Name("game")

	r.Handle("/offer-templates", Chain(
		&OfferTemplateHandler{App: a, Method: "insert"},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&ValidationMiddleware{GetPayload: func() interface{} { return &models.OfferTemplate{} }},
	)).Methods("PUT").Name("offer_templates")

	r.Handle("/offer-templates/set-enabled", Chain(
		&OfferTemplateHandler{App: a, Method: "set-enabled"},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&ValidationMiddleware{GetPayload: func() interface{} { return &models.OfferTemplateToUpdate{} }},
	)).Methods("PUT").Name("offer_templates")

	r.Handle("/offers", Chain(
		&OfferRequestHandler{App: a, Method: "get_offers"},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("offer")

	r.Handle("/offer/claim", Chain(
		&OfferRequestHandler{App: a, Method: "claim_offer"},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&ValidationMiddleware{GetPayload: func() interface{} { return &models.OfferToUpdate{} }},
	)).Methods("PUT").Name("offer")

	r.Handle("/offer/last-seen-at", Chain(
		&OfferRequestHandler{App: a, Method: "update_offer_last_seen_at"},
		&NewRelicMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&ValidationMiddleware{GetPayload: func() interface{} { return &models.OfferToUpdate{} }},
	)).Methods("PUT").Name("offer")

	return r
}

func (a *App) configureApp() error {
	a.configureLogger()

	err := a.configureDatabase()
	if err != nil {
		return err
	}

	err = a.configureNewRelic()
	if err != nil {
		return err
	}

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

	l := a.Logger.WithFields(logrus.Fields{
		"postgres.host":    host,
		"postgres.user":    user,
		"postgres.dbName":  dbName,
		"postgres.port":    port,
		"postgres.sslMode": sslMode,
	})

	l.Debug("Connecting to DB...")
	db, err := models.GetDB(
		host, user, port, sslMode, dbName,
		password, maxIdleConns, maxOpenConns,
		connectionTimeoutMS,
	)
	if err != nil {
		l.WithError(err).Error("Connection to database failed.")
		return nil, err
	}
	l.Debug("Successful connection to database.")
	return db, nil
}

func (a *App) configureLogger() {
	a.Logger = a.Logger.WithFields(logrus.Fields{
		"source":    "offers-api",
		"operation": "initializeApp",
		"version":   metadata.Version,
	})
}

func (a *App) configureNewRelic() error {
	appName := a.Config.GetString("newrelic.app")
	key := a.Config.GetString("newrelic.key")

	l := a.Logger.WithFields(logrus.Fields{
		"appName":   appName,
		"operation": "configureNewRelic",
	})

	if key == "" {
		l.Warning("New Relic key not found. No data will be sent to New Relic.")
		return nil
	}

	l.Debug("Configuring new relic...")
	config := newrelic.NewConfig(appName, key)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		l.WithError(err).Error("Failed to configure new relic.")
		return err
	}
	l.Debug("New Relic configured successfully.")
	a.NewRelic = app
	return nil
}

func (a *App) configureServer() {
	a.Router = a.getRouter()
	a.Server = &http.Server{Addr: a.Address, Handler: a.Router}
}

//HandleError writes an error response with message and status
func (a *App) HandleError(w http.ResponseWriter, status int, msg string, err error) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
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
