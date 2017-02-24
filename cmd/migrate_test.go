// +build integration

// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/topfreegames/offers/cmd"
	"github.com/topfreegames/offers/testing"
)

func dropDB() error {
	cmd := exec.Command("make", "drop-test")
	cmd.Dir = "../"
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func migrateDB() error {
	cmd := exec.Command("make", "migrate-test")
	cmd.Dir = "../"
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

var _ = Describe("Migrate Command", func() {
	BeforeEach(func() {
		err := dropDB()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Migrate Cmd", func() {
		It("Should run migrations", func() {
			ConfigFile = "../config/test.yaml"
			InitConfig()

			w, outC := testing.MockStdout()

			err := RunMigrations(false, w)
			Expect(err).NotTo(HaveOccurred())

			w.Close()
			out := <-outC

			Expect(out).To(ContainSubstring("Migrating database to latest version..."))
			Expect(out).To(ContainSubstring("Current database migration status"))
			Expect(out).To(ContainSubstring("APPLIED"))
			Expect(out).To(ContainSubstring("CreateGamesTable.sql"))
			Expect(out).To(ContainSubstring("Database migrated successfully."))
		})

		It("Should not run migrations if failed to connect to db", func() {
			ConfigFile = "../config/test_bad_db.yaml"
			InitConfig()
			w, _ := testing.MockStdout()
			Expect(func() { RunMigrations(false, w) }).To(Panic())
		})
	})

	Describe("Migrate Info Cmd", func() {
		BeforeEach(func() {
			err := migrateDB()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should show migrations info", func() {
			ConfigFile = "../config/test.yaml"
			InitConfig()

			w, outC := testing.MockStdout()
			err := RunMigrations(true, w)
			Expect(err).NotTo(HaveOccurred())

			w.Close()
			out := <-outC

			Expect(out).To(ContainSubstring("Current database migration status"))
			Expect(out).To(ContainSubstring("APPLIED"))
			Expect(out).To(ContainSubstring("CreateGamesTable.sql"))
		})
	})
})
