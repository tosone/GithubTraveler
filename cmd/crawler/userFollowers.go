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
	"gopkg.in/satori/go.uuid.v1"
)

// userFollowers get all of info from username
func userFollowers(ctx context.Context, wg *sync.WaitGroup) {
	const crawlerName = "userFollowers"
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
			requestURL := fmt.Sprintf("%s/users/%s/followers", common.GithubAPI, user.Login)
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
			if nextURL, ok = headerlink.Parse(response.Header.Get("Link"))["next"]; ok {
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
					var reliableFollowers = new(models.UserFollowers)
					reliableFollowers.UserID = user.UserID
					reliableFollowers.Version = followersVersion.String()
					reliableFollowers.FollowerUserID = u.UserID
					if err = reliableFollowers.Create(); err != nil {
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
