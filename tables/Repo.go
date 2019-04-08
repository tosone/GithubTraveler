package tables

import (
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

// Repo repository info struct
type Repo struct {
	gorm.Model
	RepoID      *int64
	NodeID      *string
	Owner       User
	Name        *string
	FullName    *string
	Description *string
	Homepage    *string
	// CodeOfConduct    *CodeOfConduct
	DefaultBranch    *string
	MasterBranch     *string
	CreatedAt        *github.Timestamp
	PushedAt         *github.Timestamp
	UpdatedAt        *github.Timestamp
	HTMLURL          *string
	CloneURL         *string
	GitURL           *string
	MirrorURL        *string
	SSHURL           *string
	SVNURL           *string
	Language         *string
	Fork             *bool
	ForksCount       RepoForksCount
	NetworkCount     RepoNetworkCount
	OpenIssuesCount  RepoOpenIssuesCount
	StargazersCount  RepoStargazersCount
	SubscribersCount RepoSubscribersCount
	WatchersCount    RepoWatchersCount
	Size             *int
	AutoInit         *bool
	// Parent           *Repository
	// Source           *Repository
	// Organization     *Organization
	// Permissions      *map[string]bool
	AllowRebaseMerge *bool
	AllowSquashMerge *bool
	AllowMergeCommit *bool
	Topics           []Topic `gorm:"many2many:repo_topics"`
	Archived         *bool
	License          string
}

// Upsert ..
func (repo *Repo) Upsert() (err error) {
	if err = engine.Model(new(Repo)).Where(Repo{
		RepoID: repo.RepoID,
	}).First(repo).Error; err == gorm.ErrRecordNotFound {
		return engine.Create(repo).Error
	} else if err != nil {
		return
	}

	return engine.Model(new(Repo)).Where(Repo{RepoID: repo.RepoID}).Updates(repo).Error
}
