package repo

import (
	"context"

	"github.com/google/go-github/github"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler/rate"
	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/EffDataAly/GithubTraveler/tables"
)

// Stargazers ..
func Stargazers(ctx context.Context, client *github.Client, u string, r string, page ...int) (err error) {
	if !util.CheckCtx(ctx) {
		return
	}
	rate.Get()

	var options = github.ListOptions{}
	var stargazers []*github.Stargazer
	var response *github.Response
	if stargazers, response, err = client.Activity.ListStargazers(ctx, u, r, &options); err != nil {
		return
	}
	for _, stargazer := range stargazers {
		var followersCount int
		if stargazer.User.Followers != nil {
			followersCount = *stargazer.User.Followers
		}
		var userFollowersCount = &tables.UserFollowersCount{
			UserID: *stargazer.User.ID,
			Num:    followersCount,
		}
		if followersCount != 0 {
			if stargazer.User.Followers != nil {
				if err = userFollowersCount.Upsert(); err != nil {
					return
				}
			}
		} else {
			userFollowersCount = &tables.UserFollowersCount{}
		}

		var followingCount int
		if stargazer.User.Followers != nil {
			followingCount = *stargazer.User.Following
		}
		var userFollowingCount = &tables.UserFollowingCount{
			UserID: *stargazer.User.ID,
			Num:    followingCount,
		}
		if followingCount != 0 {
			if stargazer.User.Following != nil {
				if err = userFollowingCount.Upsert(); err != nil {
					return
				}
			}
		} else {
			userFollowingCount = &tables.UserFollowingCount{}
		}

		var u = &tables.User{
			UserID:            stargazer.User.ID,
			Login:             stargazer.User.Login,
			NodeID:            stargazer.User.NodeID,
			AvatarURL:         stargazer.User.AvatarURL,
			HTMLURL:           stargazer.User.HTMLURL,
			GravatarID:        stargazer.User.GravatarID,
			Name:              stargazer.User.Name,
			Company:           stargazer.User.Company,
			Blog:              stargazer.User.Blog,
			Location:          stargazer.User.Location,
			Email:             stargazer.User.Email,
			Hireable:          stargazer.User.Hireable,
			Bio:               stargazer.User.Bio,
			PublicRepos:       stargazer.User.PublicRepos,
			PublicGists:       stargazer.User.PublicGists,
			Followers:         *userFollowersCount,
			Following:         *userFollowingCount,
			CreatedAt:         stargazer.User.CreatedAt,
			UpdatedAt:         stargazer.User.UpdatedAt,
			SuspendedAt:       stargazer.User.SuspendedAt,
			Type:              stargazer.User.Type,
			SiteAdmin:         stargazer.User.SiteAdmin,
			TotalPrivateRepos: stargazer.User.TotalPrivateRepos,
			OwnedPrivateRepos: stargazer.User.OwnedPrivateRepos,
			PrivateGists:      stargazer.User.PrivateGists,
			DiskUsage:         stargazer.User.DiskUsage,
			Collaborators:     stargazer.User.Collaborators,
		}
		if err = u.Upsert(); err != nil {
			return
		}
	}
	rate.Set(response.Rate.Remaining)

	if response.NextPage != 0 {
		return Stargazers(ctx, client, u, r, response.NextPage)
	}
	return
}
