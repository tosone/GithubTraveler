package crawler

import (
	"time"

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
	user.UserID = uint(viper.GetInt("Crawler.EntranceID"))
	user.Type = viper.GetString("Crawler.EntranceType")
	if err = user.Create(); err != nil {
		logging.Fatal(err)
	}
	go userRepos()
	go userFollowers()
	<-time.After(time.Hour * 3)
}
