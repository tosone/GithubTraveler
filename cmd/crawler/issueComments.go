package crawler

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tosone/GithubTraveler/common/downloader"
	"github.com/tosone/GithubTraveler/common/resp"
	"github.com/tosone/GithubTraveler/models"
	"github.com/tosone/logging"
)

// issueComments get issue's comments
func issueComments(ctx context.Context, wg *sync.WaitGroup) {
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
		var issue = new(models.RepoIssues)
		if issue, err = new(models.RepoIssues).FindByID(num); err != nil {
			if err == gorm.ErrRecordNotFound && num == 1 {
				time.Sleep(time.Second * 30)
			}
			num = 0
			continue
		}

		var user = new(models.User)
		var repo = new(models.Repo)

		if user, err = new(models.User).FindByUserID(issue.UserID); err != nil {
			logging.Error(err)
			continue
		}

		if repo, err = new(models.Repo).FindByRepoID(issue.RepoID); err != nil {
			logging.Error(err)
			continue
		}

		var nextNum = 1
		for nextNum != 0 {
			if body, nextNum, err = downloader.Get(nextNum, user.Login, repo.Name, strconv.FormatUint(issue.Number, 10)); err != nil {
				logging.Error(err)
				continue
			}

			var comments []resp.Comment
			if err = json.Unmarshal([]byte(body), &comments); err != nil {
				logging.Error(body)
				logging.Error(err)
				continue
			} else {
				for _, comment := range comments {
					var c = new(models.IssueComments)
					c.UserID = user.UserID
					c.RepoID = repo.RepoID
					c.Number = issue.Number
					c.Body = comment.Body
					if err = c.Create(); err != nil {
						logging.Error(err)
						continue
					}

					var u = new(models.User)
					u.UserID = comment.User.ID
					u.Login = comment.User.Login
					u.Type = comment.User.Type
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
