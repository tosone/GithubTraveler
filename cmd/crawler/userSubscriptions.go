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
	"github.com/EffDataAly/GithubTraveler/common/headerLink"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

// userSubscriptions get all of info from username
func userSubscriptions(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "userSubscriptions"
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
				time.Sleep(time.Second * 30)
			}
			num = 0
			continue
		}
		var followersVersion = uuid.NewV4()
		user.Followers = followersVersion.String()

		var nextURL = "next"
		var ok bool
		var page = 1

		for nextURL != "" {
			select {
			case <-ctx.Done():
				return
			default:
			}
			request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
				SetDebug(viper.GetBool("Crawler.Debug")).
				Get(fmt.Sprintf("%s/users/%s/subscriptions", common.GithubApi, user.Login)).
				Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
				Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret"))).
				Query(fmt.Sprintf("page=%d", page))
			response, body, errs = request.End()
			if nextURL, ok = headerLink.Parse(response.Header.Get("Link"))["next"]; ok {
				if u, err := url.Parse(nextURL); err != nil {
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
					var u = new(models.User)
					u.UserID = repo.Owner.ID
					u.Login = repo.Owner.Login
					u.Type = repo.Owner.Type
					if err = u.Create(); err != nil {
						logging.Error(err)
						continue
					}

					var r = new(models.Repo)
					r.UserID = repo.Owner.ID
					r.RepoID = repo.ID
					r.Name = repo.Name
					if repo.Homepage != nil {
						r.Homepage = *repo.Homepage
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

		if err = user.Update(); err != nil {
			logging.Error(err)
		}
	}
}
