package crawler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/EffDataAly/GithubTraveler/common/downloader"
	"github.com/EffDataAly/GithubTraveler/common/resp"
	"github.com/EffDataAly/GithubTraveler/models"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

// userRepos get all of repos from username
func userRepos(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var body string
	var err error
	var num uint
	for {
		num++
		var user = new(models.User)
		if user, err = new(models.User).FindByID(num); err != nil {
			if err == gorm.ErrRecordNotFound && num == 1 {
				time.Sleep(time.Second * time.Duration(viper.GetInt("Crawler.WaitDataReady")))
			}
			num = 0
			continue
		}
		if user.UserID == 0 {
			time.Sleep(time.Second * time.Duration(viper.GetInt("Crawler.WaitDataReady")))
			continue
		}

		var nextNum = 1

		for nextNum != 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if body, nextNum, err = downloader.Get(nextNum, user.Login); err != nil {
				logging.Error(err)
				continue
			}

			var repos []resp.Repo
			if err = json.Unmarshal([]byte(body), &repos); err != nil {
				logging.Error(err)
				continue
			} else {
				for _, repo := range repos {
					var r = new(models.Repo)
					r.Name = repo.Name
					r.UserID = user.UserID
					r.RepoID = repo.ID
					if repo.Homepage != nil {
						r.Description = *repo.Homepage
					}
					if repo.Language != nil {
						r.Language = *repo.Language
					}
					r.Size = repo.Size
					if repo.License != nil {
						r.Licence = repo.License.Name
					}
					if repo.Description != nil {
						r.Description = *repo.Description
					}
					if err = r.Create(); err != nil {
						logging.Error(err)
						continue
					}

					var historyRepoForksNum = new(models.HistoryRepoForksNum)
					historyRepoForksNum.UserID = user.UserID
					historyRepoForksNum.RepoID = repo.ID
					historyRepoForksNum.ForksNum = repo.Forks
					if err = historyRepoForksNum.Create(); err != nil {
						logging.Error(err)
					}

					var historyRepoStarredNum = new(models.HistoryRepoStarredNum)
					historyRepoStarredNum.UserID = user.UserID
					historyRepoStarredNum.RepoID = repo.ID
					historyRepoStarredNum.StarredNum = repo.StargazersCount
					if err = historyRepoStarredNum.Create(); err != nil {
						logging.Error(err)
					}

					var historyRepoWatchersNum = new(models.HistoryRepoWatchersNum)
					historyRepoWatchersNum.UserID = user.UserID
					historyRepoWatchersNum.RepoID = repo.ID
					historyRepoWatchersNum.WatchersNum = repo.WatchersCount
					if err = historyRepoWatchersNum.Create(); err != nil {
						logging.Error(err)
					}
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
			}
		}
	}
}
