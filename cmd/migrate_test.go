// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/topfreegames/offers/cmd"
)

func dropDB() error {
	cmd := exec.Command("make", "drop-test")
	cmd.Dir = "../"
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
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
			err := RunMigrations(false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should show migrations info", func() {
			ConfigFile = "../config/test.yaml"
			InitConfig()
			err := RunMigrations(true)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
