package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"github.com/unknwon/com"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler"
	"github.com/EffDataAly/GithubTraveler/cmd/version"
	"github.com/EffDataAly/GithubTraveler/common"
)

// RootCmd represents the base command when called without any sub commands
var RootCmd = &cobra.Command{
	Use:   common.AppName,
	Short: "Travel all of the github organizations, users and repositories.",
	Long:  `Travel all of the github organizations, users and repositories.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version",
	Long:  `The version that build detail information.`,
	Run: func(_ *cobra.Command, _ []string) {
		version.Initialize()
	},
}

var crawlerCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Travel all of the github organizations, users and repositories.",
	Long:  `Travel all of the github organizations, users and repositories.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		initConfig()
		if err := crawler.Initialize(); err != nil {
			logging.Fatal(err)
		}
	},
}

// config command line params
var config string

func init() {
	crawlerCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config.yml", "config file")

	RootCmd.AddCommand(crawlerCmd) // crawler commander
	RootCmd.AddCommand(versionCmd) // version commander
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(common.EnvPrefix)
	if com.IsFile(config) {
		viper.SetConfigFile(config)
	} else {
		logging.Fatal("Cannot find config file. Please check.")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logging.Panic("Cannot find the special config file.")
	}
}
