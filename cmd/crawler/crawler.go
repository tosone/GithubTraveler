package crawler

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/EffDataAly/GithubTraveler/common/htexpire"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

var ht = htexpire.New()

// Initialize initialize
func Initialize() {
	var err error
	if err = models.Connect(); err != nil {
		logging.Fatal(err)
	}
	user := new(models.User)
	user.Login = viper.GetString("Crawler.Entrance")
	if err = user.Create(); err != nil {
		logging.Fatal(err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	var wgAll = new(sync.WaitGroup)

	go infoRepo(ctx, wgAll)
	go infoUser(ctx, wgAll)

	go userFollowers(ctx, wgAll)
	go userFollowing(ctx, wgAll)
	go userRepos(ctx, wgAll)
	go userStarred(ctx, wgAll)
	go userSubscriptions(ctx, wgAll)

	go repoStargazers(ctx, wgAll)
	go repoWatchers(ctx, wgAll)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	ctxCancel()
	wgAll.Wait()
	logging.Info("Exit correctly already.")
}
