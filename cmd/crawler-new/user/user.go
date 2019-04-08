package user

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler-new/rate"
	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/google/go-github/github"
	"github.com/tosone/logging"
)

// User ..
func User(ctx context.Context, client *github.Client, u string, r string, page ...int) (err error) {
	if !util.CheckCtx(ctx) {
		return
	}
	rate.Get()

	var options = github.ListOptions{PerPage: 60}
	var user []*github.Stargazer
	var response *github.Response
	if user, response, err = client.Activity.ListStargazers(ctx, u, r, &options); err != nil {
		return
	}
	for _, u := range user {
		logging.Infof("%+v", u)
	}
	logging.Infof("%+v", response)
	return
}
