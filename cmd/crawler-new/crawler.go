package crawler

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler-new/rate"
	"github.com/EffDataAly/GithubTraveler/cmd/crawler-new/repo"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"golang.org/x/oauth2"
)

// Initialize crawler entry
func Initialize() (err error) {
	if err = models.Connect(); err != nil {
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
	if err = repo.Watchers(ctx, client, "tosone", "minimp3"); err != nil {
		logging.Error(err)
		return
	}
	return
}
