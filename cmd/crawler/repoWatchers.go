package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/EffDataAly/GithubTraveler/common"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/parnurzeal/gorequest"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

func repoWatchers(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "repoWatchers"
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
		var repo = new(models.Repo)
		if repo, err = new(models.Repo).FindByID(num); err != nil {
			num = 0
			continue
		}
		var watchersVersion = uuid.NewV4()
		repo.Stargazers = watchersVersion.String()

		if repo.UserID == 0 {
			continue
		}

		var user = new(models.User)

		if user, err = new(models.User).FindByUserID(repo.UserID); err != nil {
			continue
		}
		request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
			SetDebug(viper.GetBool("Crawler.Debug")).
			Get(fmt.Sprintf("%s/repos/%s/%s/watchers", common.GithubApi, user.Login, repo.Name)).
			Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
			Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret")))
		response, body, errs = request.End()
		log := new(models.Log)
		log.Url = request.Url
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
		var owners []resp.Owner
		if err = json.Unmarshal([]byte(body), &owners); err != nil {
			logging.Error(body)
			logging.Error(err)
			continue
		} else {
			for _, owner := range owners {
				var u = new(models.User)
				u.UserID = owner.ID
				u.Login = owner.Login
				u.Type = owner.Type
				if err = u.Create(); err != nil {
					logging.Error(err)
					continue
				}

				var repoWatchers = new(models.RepoWatchers)
				repoWatchers.UserID = owner.ID
				repoWatchers.RepoID = repo.RepoID
				repoWatchers.Version = watchersVersion.String()
				if err = repoWatchers.Create(); err != nil {
					logging.Error(err)
					continue
				}
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}
		if err = repo.Update(); err != nil {
			logging.Error(err)
		}
	}
}
