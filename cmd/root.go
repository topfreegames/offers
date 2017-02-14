// offers
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile is the configuration file used for running a command
var ConfigFile string

// Verbose determines how verbose offers will run under
var Verbose int

// RootCmd is the root command for offers CLI application
var RootCmd = &cobra.Command{
	Use:   "offers",
	Short: "offers-api handles offers",
	Long:  `Use offers-api to handle offers for your game.`,
}

// Execute runs RootCmd to initialize offers CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(
		&Verbose, "verbose", "v", 0,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)

	RootCmd.PersistentFlags().StringVarP(
		&ConfigFile, "config", "c", "./config/local.yaml",
		"config file (default is ./config/local.yaml",
	)
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(ConfigFile)
	}
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("offers")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config file %s failed to load: %s.\n", ConfigFile, err.Error())
		panic("Failed to load config file")
	}
}
