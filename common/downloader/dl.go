package downloader

import (
	"fmt"
	"net/url"
	"time"

	"strconv"

	"errors"

	"github.com/EffDataAly/GithubTraveler/common/htexpire"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

var ht = htexpire.New()

// Get download the specified url's body
func Get(num int, params ...string) (body string, nextNum int, err error) {
	var response gorequest.Response
	var crawlerName string
	var requestURL string
	var nextURL string
	var errs []error
	var ok bool

	if crawlerName, err = trace(); err != nil {
		return
	}

	if requestURL, err = urlSwitch(crawlerName, params...); err != nil {
		return
	}
	logging.Info(fmt.Sprintf("crawlerName: %s", crawlerName))

	if b, _ := ht.Get(requestURL); b {
		err = errors.New("too frequently request same url")
		return
	}
	if err = ht.Set(requestURL); err != nil {
		logging.Error(err)
	}

	request := gorequest.New().
		Timeout(time.Second * time.Duration(viper.GetInt("Crawler.Timeout"))).
		SetDebug(viper.GetBool("Crawler.Debug")).
		Get(requestURL).
		Query(fmt.Sprintf("client_id=%s", viper.GetString("ClientID"))).
		Query(fmt.Sprintf("client_secret=%s", viper.GetString("ClientSecret")))

	if num != 0 {
		request.Query(fmt.Sprintf("page=%d", num))
	}

	response, body, errs = request.End()

	// write log to database
	log := new(models.Log)
	log.URL = request.Url
	log.Method = request.Method
	log.Response = []byte(body)
	log.Type = crawlerName

	defer func() {
		if err = log.Create(); err != nil {
			logging.Error(err)
		}
	}()

	if len(errs) != 0 {
		var errMsg string
		for _, err = range errs {
			errMsg += err.Error()
		}
		log.ErrMsg = []byte(errMsg)
		return
	}

	if nextURL, ok = headerLink(response.Header.Get("Link"))["next"]; ok {
		var u *url.URL
		if u, err = url.Parse(nextURL); err != nil {
			logging.Error(err)
		} else {
			if nextNum, err = strconv.Atoi(u.Query().Get("page")); err != nil {
				logging.Error(err)
			}
		}
	}

	return
}
