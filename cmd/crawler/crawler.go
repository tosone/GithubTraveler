package crawler

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/EffDataAly/GithubTraveler/common/request/repo"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
	"golang.org/x/oauth2"
)

// Initialize crawler entry
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

	// go infoRepo(ctx, wgAll)
	// go infoUser(ctx, wgAll)

	// go userFollowers(ctx, wgAll)
	// go userFollowing(ctx, wgAll)
	// go userRepos(ctx, wgAll)
	// go userStarred(ctx, wgAll)
	// go userSubscriptions(ctx, wgAll)

	// go repoStargazers(ctx, wgAll)
	// go repoWatchers(ctx, wgAll)
	// go repoIssues(ctx, wgAll)

	// go issueComments(ctx, wgAll)

	var tc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "c77d41f2027f009d0cb088068e86d6465c3ddaf0"}))
	var client = github.NewClient(tc)
	go repo.Repo(ctx, client, "", 0)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel // catch the ctrl-c
	ctxCancel()     // stop all of the crawlers
	wgAll.Wait()    // wait the crawler stopped
	logging.Info("Exit correctly already.")
}
