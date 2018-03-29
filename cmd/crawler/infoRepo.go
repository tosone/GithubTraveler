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
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

func infoRepo(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "infoRepo"
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
		var user = new(models.User)
		if user, err = new(models.User).FindByUserID(repo.UserID); err != nil {
			continue
		}

		request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
			SetDebug(viper.GetBool("Crawler.Debug")).
			Get(fmt.Sprintf("%s/repos/%s/%s", common.GithubApi, user.Login, repo.Name)).
			Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
			Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret")))
		response, body, errs = request.End()
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
		var respRepo resp.Repo
		if err = json.Unmarshal([]byte(body), &respRepo); err != nil {
			logging.Error(err)
			continue
		} else {
			var r = new(models.Repo)
			r.UserID = respRepo.Owner.ID
			r.RepoID = respRepo.ID
			r.Name = respRepo.Name
			if respRepo.Description != nil {
				r.Description = *respRepo.Description
			}
			if respRepo.License != nil {
				r.Licence = respRepo.License.Name
			}
			if respRepo.Homepage != nil {
				r.Homepage = *respRepo.Homepage
			}
			if respRepo.Language != nil {
				r.Language = *respRepo.Language
			}
			r.Size = respRepo.Size
			if err = r.Create(); err != nil {
				logging.Error(err)
				continue
			}

			var historyRepoForksNum = new(models.HistoryRepoForksNum)
			historyRepoForksNum.RepoID = respRepo.ID
			historyRepoForksNum.UserID = respRepo.Owner.ID
			historyRepoForksNum.ForksNum = respRepo.ForksCount
			if err = historyRepoForksNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyRepoStarredNum = new(models.HistoryRepoStarredNum)
			historyRepoStarredNum.RepoID = respRepo.ID
			historyRepoStarredNum.UserID = respRepo.Owner.ID
			historyRepoStarredNum.StarredNum = respRepo.StargazersCount
			if err = historyRepoStarredNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyRepoWatchersNum = new(models.HistoryRepoWatchersNum)
			historyRepoWatchersNum.UserID = respRepo.Owner.ID
			historyRepoWatchersNum.RepoID = respRepo.ID
			historyRepoWatchersNum.WatchersNum = respRepo.WatchersCount
			if err = historyRepoWatchersNum.Create(); err != nil {
				logging.Error(err)
			}
		}
	}
}
