package repo

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler/rate"
	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/EffDataAly/GithubTraveler/tables"
	"github.com/google/go-github/github"
)

// Watchers ..
func Watchers(ctx context.Context, client *github.Client, u string, r string, page ...int) (err error) {
	if !util.CheckCtx(ctx) {
		return
	}
	rate.Get()

	var options = github.ListOptions{}
	var users []*github.User
	var response *github.Response
	if users, response, err = client.Activity.ListWatchers(ctx, u, r, &options); err != nil {
		return
	}
	for _, user := range users {
		var followersCount int
		if user.Followers != nil {
			followersCount = *user.Followers
		}
		var userFollowersCount = &tables.UserFollowersCount{
			UserID: *user.ID,
			Num:    followersCount,
		}
		if followersCount != 0 {
			if user.Followers != nil {
				if err = userFollowersCount.Upsert(); err != nil {
					return
				}
			}
		} else {
			userFollowersCount = &tables.UserFollowersCount{}
		}

		var followingCount int
		if user.Followers != nil {
			followingCount = *user.Following
		}
		var userFollowingCount = &tables.UserFollowingCount{
			UserID: *user.ID,
			Num:    followingCount,
		}
		if followingCount != 0 {
			if user.Following != nil {
				if err = userFollowingCount.Upsert(); err != nil {
					return
				}
			}
		} else {
			userFollowingCount = &tables.UserFollowingCount{}
		}
		var u = &tables.User{
			UserID:            user.ID,
			Login:             user.Login,
			NodeID:            user.NodeID,
			AvatarURL:         user.AvatarURL,
			HTMLURL:           user.HTMLURL,
			GravatarID:        user.GravatarID,
			Name:              user.Name,
			Company:           user.Company,
			Blog:              user.Blog,
			Location:          user.Location,
			Email:             user.Email,
			Hireable:          user.Hireable,
			Bio:               user.Bio,
			PublicRepos:       user.PublicRepos,
			PublicGists:       user.PublicGists,
			Followers:         *userFollowersCount,
			Following:         *userFollowingCount,
			CreatedAt:         user.CreatedAt,
			UpdatedAt:         user.UpdatedAt,
			SuspendedAt:       user.SuspendedAt,
			Type:              user.Type,
			SiteAdmin:         user.SiteAdmin,
			TotalPrivateRepos: user.TotalPrivateRepos,
			OwnedPrivateRepos: user.OwnedPrivateRepos,
			PrivateGists:      user.PrivateGists,
			DiskUsage:         user.DiskUsage,
			Collaborators:     user.Collaborators,
		}
		if err = u.Upsert(); err != nil {
			return
		}
	}
	rate.Set(response.Rate.Remaining)

	if response.NextPage != 0 {
		return Watchers(ctx, client, u, r, response.NextPage)
	}
	return
}
