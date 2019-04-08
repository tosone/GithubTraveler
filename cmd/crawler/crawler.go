package crawler

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler/rate"
	"github.com/EffDataAly/GithubTraveler/cmd/crawler/user"
	"github.com/EffDataAly/GithubTraveler/database"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"golang.org/x/oauth2"
)

// Initialize crawler entry
func Initialize() (err error) {
	if err = database.Connect(); err != nil {
		return
	}
	var ctx = context.Background()
	var client = github.NewClient(oauth2.NewClient(ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString("AccessToken")}),
	))
	if err = rate.Initialize(ctx, client); err != nil {
		return
	}
	// if err = repo.List(ctx, client, ""); err != nil {
	// 	logging.Error(err)
	// 	return
	// }
	if err = user.Followers(ctx, client, "tosone"); err != nil {
		logging.Error(err)
		return
	}
	return
}
