package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/EffDataAly/GithubTraveler/common"
	"github.com/EffDataAly/GithubTraveler/common/headerlink"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

// userRepos get all of repos from username
func userRepos(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "userRepos"
	wg.Add(1)
	defer wg.Done()

	var response gorequest.Response
	var body string
	var errs []error
	var err error
	var num uint
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		num++
		var user = new(models.User)
		if user, err = new(models.User).FindByID(num); err != nil {
			if err == gorm.ErrRecordNotFound && num == 1 {
				time.Sleep(time.Second * time.Duration(viper.GetInt("Crawler.WaitDataReady")))
			}
			num = 0
			continue
		}
		if user.UserID == 0 {
			time.Sleep(time.Second * time.Duration(viper.GetInt("Crawler.WaitDataReady")))
			continue
		}

		var nextURL = "next"
		var ok bool
		var page = 1
		for nextURL != "" {
			select {
			case <-ctx.Done():
				return
			default:
			}
			requestURL := fmt.Sprintf("%s/users/%s/repos", common.GithubAPI, user.Login)
			if b, _ := ht.Get(requestURL); b {
				continue
			}
			if err = ht.Set(requestURL); err != nil {
				logging.Error(err)
			}
			request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
				SetDebug(viper.GetBool("Crawler.Debug")).
				Get(requestURL).
				Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
				Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret"))).
				Query(fmt.Sprintf("page=%d", page))
			response, body, errs = request.End()
			if response.Header.Get("Link") == "" {
				nextURL = ""
			} else if nextURL, ok = headerlink.Parse(response.Header.Get("Link"))["next"]; ok {
				var u *url.URL
				if u, err = url.Parse(nextURL); err != nil {
					logging.Error(err)
				} else {
					if page, err = strconv.Atoi(u.Query().Get("page")); err != nil {
						logging.Error(err)
					}
				}
			} else {
				nextURL = ""
			}
			log := new(models.Log)
			log.URL = request.Url
			log.Method = request.Method
			log.Response = []byte(body)
			log.Type = crawlerName
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
					if repo.Homepage != nil {
						r.Description = *repo.Homepage
					}
					if repo.Language != nil {
						r.Language = *repo.Language
					}
					r.Size = repo.Size
					if repo.License != nil {
						r.Licence = repo.License.Name
					}
					if repo.Description != nil {
						r.Description = *repo.Description
					}
					if err = r.Create(); err != nil {
						logging.Error(err)
						continue
					}

					var historyRepoForksNum = new(models.HistoryRepoForksNum)
					historyRepoForksNum.UserID = user.UserID
					historyRepoForksNum.RepoID = repo.ID
					historyRepoForksNum.ForksNum = repo.Forks
					if err = historyRepoForksNum.Create(); err != nil {
						logging.Error(err)
					}

					var historyRepoStarredNum = new(models.HistoryRepoStarredNum)
					historyRepoStarredNum.UserID = user.UserID
					historyRepoStarredNum.RepoID = repo.ID
					historyRepoStarredNum.StarredNum = repo.StargazersCount
					if err = historyRepoStarredNum.Create(); err != nil {
						logging.Error(err)
					}

					var historyRepoWatchersNum = new(models.HistoryRepoWatchersNum)
					historyRepoWatchersNum.UserID = user.UserID
					historyRepoWatchersNum.RepoID = repo.ID
					historyRepoWatchersNum.WatchersNum = repo.WatchersCount
					if err = historyRepoWatchersNum.Create(); err != nil {
						logging.Error(err)
					}
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
			}
		}
	}
}
