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
	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"gopkg.in/satori/go.uuid.v1"
)

func repoStargazers(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "repoStargazers"
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
			if err == gorm.ErrRecordNotFound && num == 1 {
				time.Sleep(time.Second * 30)
			}
			num = 0
			continue
		}
		var stargazersVersion = uuid.NewV4()
		repo.Stargazers = stargazersVersion.String()

		if repo.UserID == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		var user = new(models.User)
		if user, err = new(models.User).FindByUserID(repo.UserID); err != nil {
			continue
		}
		requestURL := fmt.Sprintf("%s/repos/%s/%s/stargazers", common.GithubAPI, user.Login, repo.Name)
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

				var repoStargazers = new(models.RepoStargazers)
				repoStargazers.UserID = owner.ID
				repoStargazers.RepoID = repo.RepoID
				repoStargazers.Version = stargazersVersion.String()
				if err = repoStargazers.Create(); err != nil {
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
