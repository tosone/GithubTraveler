package repo

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"github.com/tosone/logging"
)

func Repo(ctx context.Context, client *github.Client, user string, page int) (err error) {
	var repos []*github.Repository
	var response *github.Response

	var option = github.RepositoryListOptions{Affiliation: "owner"}
	option.Page = page

	if repos, response, err = client.Repositories.List(ctx, user, &option); err != nil {
		log.Fatalln(err)
	}

	for _, repo := range repos {
		logging.Infof("%+v", repo.License)
		logging.Infof("%+v", repo.Source)
		logging.Infof("%+v", repo.Parent)
		logging.Infof("%+v", repo.CodeOfConduct)
		logging.Infof("%+v", repo)
	}
	logging.Infof("%+v", response.NextPage)
	logging.Infof("%+v", response.Rate)
	return
}
