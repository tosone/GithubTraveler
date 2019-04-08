package tables

import (
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

// User ..
type User struct {
	gorm.Model
	Login                   *string
	UserID                  *int64
	NodeID                  *string
	AvatarURL               *string
	HTMLURL                 *string
	GravatarID              *string
	Name                    *string
	Company                 *string
	Blog                    *string
	Location                *string
	Email                   *string
	Hireable                *bool
	Bio                     *string
	PublicRepos             *int
	PublicGists             *int
	Followers               UserFollowersCount
	Following               UserFollowingCount
	CreatedAt               *github.Timestamp
	UpdatedAt               *github.Timestamp
	SuspendedAt             *github.Timestamp
	Type                    *string
	SiteAdmin               *bool
	TotalPrivateRepos       *int
	OwnedPrivateRepos       *int
	PrivateGists            *int
	DiskUsage               *int
	Collaborators           *int
	TwoFactorAuthentication *bool
}

func (user *User) Upsert() (err error) {
	if err = engine.Model(new(User)).Where(User{
		UserID: user.UserID,
	}).First(user).Error; err == gorm.ErrRecordNotFound {
		return engine.Create(user).Error
	} else if err != nil {
		return
	}
	return engine.Model(new(User)).Where(User{UserID: user.UserID}).Updates(user).Error
}
