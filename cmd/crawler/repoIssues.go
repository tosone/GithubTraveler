package crawler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tosone/GithubTraveler/common/downloader"
	"github.com/tosone/GithubTraveler/common/resp"
	"github.com/tosone/GithubTraveler/models"
	"github.com/tosone/logging"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// repoIssues get repo's issues
func repoIssues(ctx context.Context, wg *sync.WaitGroup) {
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
		var watchersVersion = uuid.NewV4()
		repo.Stargazers = watchersVersion.String()

		if repo.UserID == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		var user = new(models.User)

		if user, err = new(models.User).FindByUserID(repo.UserID); err != nil {
			continue
		}
		var nextNum = 1
		for nextNum != 0 {
			if body, nextNum, err = downloader.Get(nextNum, user.Login, repo.Name); err != nil {
				logging.Error(err)
				continue
			}

			var issues []resp.Issue
			if err = json.Unmarshal([]byte(body), &issues); err != nil {
				logging.Error(body)
				logging.Error(err)
				continue
			} else {
				for _, issue := range issues {
					var i = new(models.RepoIssues)
					i.UserID = issue.User.ID
					i.RepoID = repo.RepoID
					i.Number = issue.Number
					i.Comments = issue.Commits
					i.Title = issue.Title
					i.Body = issue.Body
					if err = i.Create(); err != nil {
						logging.Error(err)
						continue
					}

					var u = new(models.User)
					u.UserID = issue.User.ID
					u.Login = issue.User.Login
					u.Type = issue.User.Type
					if err = u.Create(); err != nil {
						logging.Error(err)
						continue
					}

					select {
					case <-ctx.Done():
						return
					default:
					}
				}
			}
			if err = repo.Update(); err != nil {
				logging.Error(err)
			}
		}
	}
}
