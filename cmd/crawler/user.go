package crawler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/tosone/GithubTraveler/common"
	"github.com/tosone/GithubTraveler/common/resp"
	"github.com/tosone/GithubTraveler/models"
	"github.com/tosone/logging"
)

// userCrawler get all of info from username
func userCrawler() {
	var response gorequest.Response
	var body string
	var errs []error
	var err error
	var num uint
	for {
		num++
		var user = new(models.User)
		if user, err = new(models.User).FindByID(num); err != nil {
			num = 0
			continue
		}
		request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
			SetDebug(viper.GetBool("Crawler.Debug")).
			Get(fmt.Sprintf("%s/users/%s/repos", common.GithubApi, user.Login)).
			Query(fmt.Sprintf("clientid=%s", viper.GetString("ClientID"))).
			Query(fmt.Sprintf("clientsecret=%s", viper.GetString("ClientSecret")))
		response, body, errs = request.End()
		log := new(models.Log)
		log.Url = request.Url
		log.Method = request.Method
		log.Response = []byte(body)
		if len(errs) != 0 {
			var errMsg string
			for _, err := range errs {
				errMsg += err.Error()
				logging.Info(err)
			}
			log.ErrMsg = []byte(errMsg)
		}
		if err = log.Create(); err != nil {
			logging.Error(err)
		}
		if response == nil {
			continue
		}
		var repos resp.Repos
		if err = json.Unmarshal([]byte(body), &repos); err != nil {
			for _, repo := range repos {
				var r = new(models.Repo)
				r.Name = repo.Name
				r.RepoID = repo.ID
				r.StargazersCount = repo.StargazersCount
				if err = r.Create(); err != nil {
					logging.Error(err)
					continue
				}
			}
		}
	}
}
