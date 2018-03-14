package crawler

import (
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

func Initialize(tags ...string) {
	var err error
	if err = models.Connect(); err != nil {
		logging.Fatal(err)
	}
	user := new(models.User)
	user.Login = viper.GetString("Crawler.Entrance")
	if err = user.Create(); err != nil {
		logging.Fatal(err)
	}
	userCrawler()
}
