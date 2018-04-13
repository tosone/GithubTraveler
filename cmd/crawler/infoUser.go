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
)

func infoUser(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "infoUser"
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
		requestURL := fmt.Sprintf("%s/users/%s", common.GithubApi, user.Login)
		if ht.Get(requestURL) {
			continue
		}
		ht.Set(requestURL)
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
		var owner resp.User
		if err = json.Unmarshal([]byte(body), &owner); err != nil {
			logging.Error(err)
			continue
		} else {
			var u = new(models.User)
			u.UserID = owner.UserID
			u.Type = owner.Type
			u.Login = owner.Login
			u.Email = owner.Email
			u.Location = owner.Location
			if err = u.Create(); err != nil {
				logging.Error(err)
				continue
			}

			var historyUserFollowersNum = new(models.HistoryUserFollowersNum)
			historyUserFollowersNum.FollowersNum = owner.Followers
			historyUserFollowersNum.UserID = owner.UserID
			if err = historyUserFollowersNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyUserFollowingNum = new(models.HistoryUserFollowingNum)
			historyUserFollowingNum.FollowingNum = owner.Following
			historyUserFollowingNum.UserID = owner.UserID
			if err = historyUserFollowingNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyUserReposNum = new(models.HistoryUserReposNum)
			historyUserReposNum.UserID = owner.UserID
			historyUserReposNum.ReposNum = owner.PublicRepos
			if err = historyUserReposNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyUserGistNum = new(models.HistoryUserGistNum)
			historyUserGistNum.UserID = owner.UserID
			historyUserGistNum.GistNum = owner.PublicGists
			if err = historyUserGistNum.Create(); err != nil {
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
