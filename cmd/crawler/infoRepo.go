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
	"github.com/tosone/logging"
)

func infoRepo(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var body string
	var err error
	var num uint
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		num++
		var repo = new(models.Repo)
		if repo, err = new(models.Repo).FindByID(num); err != nil {
			if err == gorm.ErrRecordNotFound && num == 1 {
				time.Sleep(time.Second * 30)
			}
			num = 0
			continue
		}
		var user = new(models.User)
		if user, err = new(models.User).FindByUserID(repo.UserID); err != nil {
			continue
		}

		if body, _, err = downloader.Get(0, user.Login, repo.Name); err != nil {
			logging.Error(err)
			continue
		}

		var respRepo resp.Repo
		if err = json.Unmarshal([]byte(body), &respRepo); err != nil {
			logging.Error(err)
			continue
		} else {
			var r = new(models.Repo)
			r.UserID = respRepo.Owner.ID
			r.RepoID = respRepo.ID
			r.Name = respRepo.Name
			if respRepo.Description != nil {
				r.Description = *respRepo.Description
			}
			if respRepo.License != nil {
				r.Licence = respRepo.License.Name
			}
			if respRepo.Homepage != nil {
				r.Homepage = *respRepo.Homepage
			}
			if respRepo.Language != nil {
				r.Language = *respRepo.Language
			}
			r.Size = respRepo.Size
			if err = r.Create(); err != nil {
				logging.Error(err)
				continue
			}

			var historyRepoForksNum = new(models.HistoryRepoForksNum)
			historyRepoForksNum.RepoID = respRepo.ID
			historyRepoForksNum.UserID = respRepo.Owner.ID
			historyRepoForksNum.ForksNum = respRepo.ForksCount
			if err = historyRepoForksNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyRepoStarredNum = new(models.HistoryRepoStarredNum)
			historyRepoStarredNum.RepoID = respRepo.ID
			historyRepoStarredNum.UserID = respRepo.Owner.ID
			historyRepoStarredNum.StarredNum = respRepo.StargazersCount
			if err = historyRepoStarredNum.Create(); err != nil {
				logging.Error(err)
			}

			var historyRepoWatchersNum = new(models.HistoryRepoWatchersNum)
			historyRepoWatchersNum.UserID = respRepo.Owner.ID
			historyRepoWatchersNum.RepoID = respRepo.ID
			historyRepoWatchersNum.WatchersNum = respRepo.WatchersCount
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
