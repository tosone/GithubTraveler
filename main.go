package main

import (
	"fmt"
	"runtime"

	"github.com/EffDataAly/GithubTraveler/cmd"
	"github.com/EffDataAly/GithubTraveler/cmd/version"
	"github.com/EffDataAly/GithubTraveler/common"
	"github.com/tosone/logging"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Version version
var Version = "no provided"

// BuildStamp BuildStamp
var BuildStamp = "no provided"

// GitHash GitHash
var GitHash = "no provided"

func main() {
	if runtime.GOOS == "windows" {
		logging.Panic(fmt.Sprintf("%s not support windows just linux.", common.AppName))
	}

	version.Setting(Version, BuildStamp, GitHash)

	if err := cmd.RootCmd.Execute(); err != nil {
		logging.Panic(err.Error())
	}
}
