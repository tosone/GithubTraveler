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
	"github.com/parnurzeal/gorequest"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

func userFollowing(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "userFollowing"
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
			num = 0
			continue
		}
		var followingVersion = uuid.NewV4()
		user.Following = followingVersion.String()

		var nextUrl = "next"
		var ok bool
		var page = 1
		for nextUrl != "" {
			request := gorequest.New().Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
				SetDebug(viper.GetBool("Crawler.Debug")).
				Get(fmt.Sprintf("%s/users/%s/following", common.GithubApi, user.Login)).
				Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
				Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret"))).
				Query(fmt.Sprintf("page=%d", page))
			response, body, errs = request.End()
			if nextUrl, ok = headerLink.Parse(response.Header.Get("Link"))["next"]; ok {
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