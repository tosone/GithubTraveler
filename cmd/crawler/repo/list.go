package repo

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler/rate"
	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/EffDataAly/GithubTraveler/tables"
	"github.com/google/go-github/github"
)

// List ..
func List(ctx context.Context, client *github.Client, user string, page ...int) (err error) {
	if !util.CheckCtx(ctx) {
		return
	}
	rate.Get()

	var repos []*github.Repository
	var response *github.Response

	var option = github.RepositoryListOptions{Affiliation: "owner"}
	if len(page) == 1 {
		option.Page = page[0]
	}
	if repos, response, err = client.Repositories.List(ctx, user, &option); err != nil {
		return
	}
	for _, repo := range repos {
		var license string
		if repo.License != nil {
			license = *repo.License.Key
		}
		var topics []tables.Topic
		if repo.Topics != nil {
			for _, t := range repo.Topics {
				var topic = &tables.Topic{Content: t}
				if err = topic.Upsert(); err != nil {
					return
				}
				topics = append(topics, *topic)
			}
		}
		var repoForksCount = &tables.RepoForksCount{RepoID: *repo.ID, Num: *repo.ForksCount}
		if repo.ForksCount != nil {
			if err = repoForksCount.Upsert(); err != nil {
				return
			}
		}
		var repoNetworkCountNum int
		if repo.NetworkCount != nil {
			repoNetworkCountNum = *repo.NetworkCount
		}
		var repoNetworkCount = &tables.RepoNetworkCount{RepoID: *repo.ID, Num: repoNetworkCountNum}
		if repo.NetworkCount != nil {
			if err = repoNetworkCount.Upsert(); err != nil {
				return
			}
		}
		var repoOpenIssuesCount = &tables.RepoOpenIssuesCount{RepoID: *repo.ID, Num: *repo.OpenIssuesCount}
		if repo.OpenIssuesCount != nil {
			if err = repoOpenIssuesCount.Upsert(); err != nil {
				return
			}
		}
		var repoStargazersCount = &tables.RepoStargazersCount{RepoID: *repo.ID, Num: *repo.StargazersCount}
		if repo.StargazersCount != nil {
			if err = repoStargazersCount.Upsert(); err != nil {
				return
			}
		}
		var repoSubscribersCountNum int
		if repo.SubscribersCount != nil {
			repoSubscribersCountNum = *repo.SubscribersCount
		}
		var repoSubscribersCount = &tables.RepoSubscribersCount{RepoID: *repo.ID, Num: repoSubscribersCountNum}
		if repo.SubscribersCount != nil {
			if err = repoSubscribersCount.Upsert(); err != nil {
				return
			}
		}
		var repoWatchersCount = &tables.RepoWatchersCount{RepoID: *repo.ID, Num: *repo.WatchersCount}
		if repo.WatchersCount != nil {
			if err = repoWatchersCount.Upsert(); err != nil {
				return
			}
		}
		var r = &tables.Repo{
			RepoID: repo.ID,
			NodeID: repo.NodeID,
			// Owner:            repo.Owner.ID,
			Name:             repo.Name,
			FullName:         repo.FullName,
			Description:      repo.Description,
			Homepage:         repo.Homepage,
			DefaultBranch:    repo.DefaultBranch,
			MasterBranch:     repo.MasterBranch,
			CreatedAt:        repo.CreatedAt,
			PushedAt:         repo.PushedAt,
			UpdatedAt:        repo.UpdatedAt,
			HTMLURL:          repo.HTMLURL,
			CloneURL:         repo.CloneURL,
			GitURL:           repo.GitURL,
			MirrorURL:        repo.MirrorURL,
			SSHURL:           repo.SSHURL,
			SVNURL:           repo.SVNURL,
			Language:         repo.Language,
			Fork:             repo.Fork,
			ForksCount:       *repoForksCount,
			NetworkCount:     *repoNetworkCount,
			OpenIssuesCount:  *repoOpenIssuesCount,
			StargazersCount:  *repoStargazersCount,
			SubscribersCount: *repoSubscribersCount,
			WatchersCount:    *repoWatchersCount,
			Size:             repo.Size,
			AutoInit:         repo.AutoInit,
			AllowRebaseMerge: repo.AllowRebaseMerge,
			AllowSquashMerge: repo.AllowSquashMerge,
			AllowMergeCommit: repo.AllowMergeCommit,
			Topics:           topics,
			Archived:         repo.Archived,
			License:          license,
		}

		if err = r.Upsert(); err != nil {
			return
		}
	}

	rate.Set(response.Rate.Remaining)

	if response.NextPage != 0 {
		return List(ctx, client, user, response.NextPage)
	}
	return
}
