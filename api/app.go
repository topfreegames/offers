// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	e "errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/pmylund/go-cache"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/metadata"
	"github.com/topfreegames/offers/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//App is our API application
type App struct {
	Address           string
	Clock             models.Clock
	Config            *viper.Viper
	DB                runner.Connection
	Debug             bool
	Logger            logrus.FieldLogger
	MaxAge            int64
	NewRelic          newrelic.Application
	Router            *mux.Router
	Server            *http.Server
	Cache             *cache.Cache
	OffersCacheMaxAge time.Duration
	Pagination        *Pagination
}

//Pagination holds the page size (limit) and the offset that is the page number
type Pagination struct {
	Limit  uint64
	Offset uint64
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
		&SentryMiddleware{},
		&MetricsReporterMiddleware{App: a},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("healthcheck")

	r.Handle("/games", Chain(
		&GameHandler{App: a, Method: "list"},
		&SentryMiddleware{},
		&MetricsReporterMiddleware{App: a},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("game")

	r.HandleFunc("/games/{id}", Chain(
		&GameHandler{App: a, Method: "upsert"},
		&SentryMiddleware{},
		&MetricsReporterMiddleware{App: a},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewParamKeyMiddleware(a, func(id string) bool {
			return govalidator.Matches(id, "^[^-][a-zA-Z0-9-_]*$") && govalidator.StringLength(id, "1", "255")
		}),
		NewValidationMiddleware(func() interface{} { return &models.Game{} }),
	).ServeHTTP).Methods("PUT").Name("game")

	r.Handle("/offers", Chain(
		&OfferHandler{App: a, Method: "list"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("offers")

	r.Handle("/offers", Chain(
		&OfferHandler{App: a, Method: "insert"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewValidationMiddleware(func() interface{} { return &models.Offer{} }),
	)).Methods("POST").Name("offers")

	r.Handle("/offers/claim", Chain(
		&OfferRequestHandler{App: a, Method: "claim"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewValidationMiddleware(func() interface{} { return &models.ClaimOfferPayload{} }),
	)).Methods("PUT").Name("offer-requests")

	r.Handle("/offers/{id}", Chain(
		&OfferHandler{App: a, Method: "update"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewParamKeyMiddleware(a, govalidator.IsUUIDv4),
		NewValidationMiddleware(func() interface{} { return &models.Offer{} }),
	)).Methods("PUT").Name("offers")

	r.Handle("/offers/{id}/enable", Chain(
		&OfferHandler{App: a, Method: "enable"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewParamKeyMiddleware(a, govalidator.IsUUIDv4),
	)).Methods("PUT").Name("offers")

	r.Handle("/offers/{id}/disable", Chain(
		&OfferHandler{App: a, Method: "disable"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a, useBasicAuth: true},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewParamKeyMiddleware(a, govalidator.IsUUIDv4),
	)).Methods("PUT").Name("offers")

	r.Handle("/available-offers", Chain(
		&OfferRequestHandler{App: a, Method: "get-offers"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("offer-requests")

	r.HandleFunc("/offers/{id}/impressions", Chain(
		&OfferRequestHandler{App: a, Method: "impressions"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		NewParamKeyMiddleware(a, govalidator.IsUUIDv4),
		NewValidationMiddleware(func() interface{} { return &models.OfferImpressionPayload{} }),
	).ServeHTTP).Methods("PUT").Name("offer-requests")

	r.Handle("/offer-info", Chain(
		&OfferRequestHandler{App: a, Method: "offer-info"},
		&SentryMiddleware{},
		&NewRelicMiddleware{App: a},
		&AuthMiddleware{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("offer-requests")

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

	a.configureSentry()

	a.MaxAge = a.Config.GetInt64("cache.maxAgeSeconds")
	a.configurePagination()
	a.configureServer()
	a.configureCache()
	return nil
}

func (a *App) configurePagination() {
	limit := uint64(a.Config.GetInt64("pagination.limit"))
	offset := uint64(a.Config.GetInt64("pagination.offset"))
	a.Pagination = &Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func (a *App) configureCache() {
	maxAge := time.Duration(a.Config.GetInt64("offersCache.maxAgeSeconds")) * time.Second
	cleanupInterval := time.Duration(a.Config.GetInt64("offersCache.cleanupInterval")) * time.Second
	a.Cache = cache.New(maxAge, cleanupInterval)
	a.OffersCacheMaxAge = maxAge
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

func (a *App) configureSentry() {
	l := a.Logger.WithFields(logrus.Fields{
		"operation": "configureSentry",
	})
	sentryURL := a.Config.GetString("sentry.url")
	l.Debug("Configuring sentry...")
	raven.SetDSN(sentryURL)
	raven.SetRelease(metadata.Version)
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
	a.Server = &http.Server{Addr: a.Address, Handler: wrapHandlerWithResponseWriter(a.Router)}
}

//HandleError writes an error response with message and status
func (a *App) HandleError(w http.ResponseWriter, status int, msg string, err interface{}) {
	w.WriteHeader(status)
	var sErr errors.SerializableError
	val, ok := err.(errors.SerializableError)
	if ok {
		sErr = val
	} else {
		sErr = errors.NewGenericError(msg, err.(error))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(sErr.Serialize())

	var rErr error
	errVal, ok := err.(error)
	if ok {
		rErr = errVal
	} else {
		rErr = e.New(msg)
	}
	raven.CaptureError(rErr, nil)
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
