package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosone/GithubTraveler/cmd/crawler"
	"github.com/tosone/GithubTraveler/cmd/version"
	"github.com/tosone/GithubTraveler/common"
	"github.com/tosone/logging"
)

var dir string

// RootCmd represents the base command when called without any sub commands
var RootCmd = &cobra.Command{
	Use:   common.AppName,
	Short: "Travel all of the github organizations, users and repository.",
	Long:  `Travel all of the github organizations, users and repository.`,
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
	Short: "Travel all of the github organizations, users and repository.",
	Long:  `Travel all of the github organizations, users and repository.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		initConfig()
		crawler.Initialize(args...)
	},
}

func init() {
	var err error
	var currPath string
	if currPath, err = os.Getwd(); err != nil {
		logging.Fatal(err)
	}
	crawlerCmd.PersistentFlags().StringVarP(&dir, "dir", "d", currPath, "execute path")

	RootCmd.AddCommand(crawlerCmd)
	RootCmd.AddCommand(versionCmd)
}

func initConfig() {
	viper.SetConfigType("yaml")
	if dir != "" {
		var config = path.Join(dir, common.Config)
		if !com.IsFile(config) {
			logging.Fatal(fmt.Sprintf("Cannot find config file here: %s", config))
		} else {
			viper.SetConfigFile(config)
		}
	} else {
		logging.Fatal("Cannot find config file. Please check.")
	}
	if err := viper.ReadInConfig(); err != nil {
		logging.Panic("Cannot find the special config file.")
	}
}
