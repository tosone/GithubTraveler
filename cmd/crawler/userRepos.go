package crawler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/EffDataAly/GithubTraveler/common"
	"github.com/EffDataAly/GithubTraveler/common/headerLink"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

// userRepos get all of repos from username
func userRepos() {
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
		var nextUrl = "next"
		var ok bool
		var page = 1
		for nextUrl != "" {
			request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
				SetDebug(viper.GetBool("Crawler.Debug")).
				Get(fmt.Sprintf("%s/users/%s/repos", common.GithubApi, user.Login)).
				Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
				Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret"))).
				Query(fmt.Sprintf("page=%d", page))
			response, body, errs = request.End()
			if response.Header.Get("Link") == "" {
				nextUrl = ""
			} else if nextUrl, ok = headerLink.Parse(response.Header.Get("Link"))["next"]; ok {
				if u, err := url.Parse(nextUrl); err != nil {
					logging.Error(err)
				} else {
					if page, err = strconv.Atoi(u.Query().Get("page")); err != nil {
						logging.Error(err)
					}
				}
			} else {
				nextUrl = ""
			}
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
			var repos []resp.Repo
			if err = json.Unmarshal([]byte(body), &repos); err != nil {
				logging.Error(err)
				continue
			} else {
				for _, repo := range repos {
					var r = new(models.Repo)
					r.Name = repo.Name
					r.UserID = user.UserID
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
}
