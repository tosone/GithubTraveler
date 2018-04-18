package crawler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/EffDataAly/GithubTraveler/common/downloader"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/jinzhu/gorm"
	"github.com/tosone/logging"
)

// infoUser get user's detail info
func infoUser(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var body string
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

		if body, _, err = downloader.Get(0, user.Login); err != nil {
			logging.Error(err)
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
