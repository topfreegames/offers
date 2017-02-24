// offers api
// https://github.com/topfree/ames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("App", func() {
	var err error
	var l *logrus.Logger
	var config *viper.Viper
	var clock models.Clock

	BeforeEach(func() {
		l = logrus.New()
		l.Level = logrus.FatalLevel
		config, err = GetDefaultConfig()
		Expect(err).NotTo(HaveOccurred())
		clock = MockClock{CurrentTime: 1486678000}
	})

	Describe("NewApp", func() {
		It("should return new app", func() {
			application, err := api.NewApp("0.0.0.0", 8889, config, false, l, clock)
			Expect(err).NotTo(HaveOccurred())
			Expect(application).NotTo(BeNil())
			Expect(application.Address).NotTo(Equal(""))
			Expect(application.Debug).To(BeFalse())
			Expect(application.Router).NotTo(BeNil())
			Expect(application.Server).NotTo(BeNil())
			Expect(application.Config).To(Equal(config))
			Expect(application.DB).NotTo(BeNil())
			Expect(application.Logger).NotTo(BeNil())
			Expect(application.NewRelic).NotTo(BeNil())
			Expect(application.Clock).To(Equal(clock))
		})

		It("should return new app with nil clock", func() {
			application, err := api.NewApp("0.0.0.0", 8889, config, false, l, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(application).NotTo(BeNil())
			var realClock *models.RealClock
			Expect(application.Clock).To(BeAssignableToTypeOf(realClock))
		})

		It("should fail if some error occured", func() {
			config.Set("newrelic.key", 12345)
			application, err := api.NewApp("0.0.0.0", 8889, config, false, l, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("license length is not 40"))
			Expect(application).To(BeNil())
		})

		It("should not fail if no newrelic key is provided", func() {
			config.Set("newrelic.key", "")
			application, err := api.NewApp("0.0.0.0", 8889, config, false, l, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(application).NotTo(BeNil())
		})
	})
})
