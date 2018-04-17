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
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"gopkg.in/satori/go.uuid.v1"
)

func userFollowing(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var body string
	var err error
	var num uint
	for {
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
		var followingVersion = uuid.NewV4()
		user.Following = followingVersion.String()

		var nextNum = 1

		for nextNum != 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if body, nextNum, err = downloader.Get(nextNum, user.Login); err != nil {
				logging.Error(err)
				continue
			}

			var owners []resp.Owner
			if err = json.Unmarshal([]byte(body), &owners); err != nil {
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
					var reliableFollowing = new(models.UserFollowing)
					reliableFollowing.UserID = user.UserID
					reliableFollowing.Version = followingVersion.String()
					reliableFollowing.FollowingUserID = u.UserID
					if err = reliableFollowing.Create(); err != nil {
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
		}
		if err = user.Update(); err != nil {
			logging.Error(err)
		}
	}
}
