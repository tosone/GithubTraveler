package cmd

import (
	"github.com/EffDataAly/GithubTraveler/cmd/crawler"
	"github.com/EffDataAly/GithubTraveler/cmd/version"
	"github.com/EffDataAly/GithubTraveler/common"
	"github.com/Unknwon/com"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

var config string

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
		crawler.Initialize()
	},
}

func init() {
	crawlerCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config.yml", "config file")

	RootCmd.AddCommand(crawlerCmd)
	RootCmd.AddCommand(versionCmd)
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
