package repo

import (
	"context"

	"github.com/EffDataAly/GithubTraveler/cmd/crawler-new/rate"
	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/EffDataAly/GithubTraveler/models"
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
		var u = &models.User{
			UserID:     user.ID,
			Login:      user.Login,
			NodeID:     user.NodeID,
			AvatarURL:  user.AvatarURL,
			HTMLURL:    user.HTMLURL,
			GravatarID: user.GravatarID,
			Name:       user.Name,
			Company:    user.Company,
			Blog:       user.Blog,
			Location:   user.Location,
			Email:      user.Email,
			Hireable:   user.Hireable,
			Bio:        user.Bio,
			// PublicRepos             *int
			// PublicGists             *int
			// Followers               UserFollowersCount
			// Following               UserFollowingCount
			CreatedAt:               user.CreatedAt,
			UpdatedAt:               user.UpdatedAt,
			SuspendedAt:             user.SuspendedAt,
			Type:                    user.Type,
			SiteAdmin:               user.SiteAdmin,
			TotalPrivateRepos:       user.TotalPrivateRepos,
			OwnedPrivateRepos:       user.OwnedPrivateRepos,
			PrivateGists:            user.PrivateGists,
			DiskUsage:               user.DiskUsage,
			Collaborators:           user.Collaborators,
			TwoFactorAuthentication: user.TwoFactorAuthentication,
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
