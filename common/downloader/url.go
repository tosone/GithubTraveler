package downloader

import (
	"errors"
	"fmt"

	"github.com/EffDataAly/GithubTraveler/common"
)

// urlSwitch get request url
func urlSwitch(t string, params ...string) (url string, err error) {
	switch t {
	case "infoRepo":
		if err = checkParamsNum(2, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/repos/%s/%s", common.GithubAPI, params[0], params[1])
	case "infoUser":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s", common.GithubAPI, params[0])
	case "issueComments":
		if err = checkParamsNum(3, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/repos/%s/%s/issues/%s/comments", common.GithubAPI, params[0], params[1], params[2])
	case "repoIssues":
		if err = checkParamsNum(2, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/repos/%s/%s/issues?state=all", common.GithubAPI, params[0], params[1])
	case "repoStargazers":
		if err = checkParamsNum(2, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/repos/%s/%s/stargazers", common.GithubAPI, params[0], params[1])
	case "repoWatchers":
		if err = checkParamsNum(2, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/repos/%s/%s/watchers", common.GithubAPI, params[0], params[1])
	case "userFollowers":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s/followers", common.GithubAPI, params[0])
	case "userFollowing":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s/following", common.GithubAPI, params[0])
	case "userRepos":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s/repos", common.GithubAPI, params[0])
	case "userStarred":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s/starred", common.GithubAPI, params[0])
	case "userSubscriptions":
		if err = checkParamsNum(1, params...); err != nil {
			return
		}
		url = fmt.Sprintf("%s/users/%s/subscriptions", common.GithubAPI, params[0])
	default:
	}
	return
}

// checkParamsNum check params num
func checkParamsNum(num int, params ...string) error {
	if len(params) != num {
		return errors.New("params is not expected")
	}
	return nil
}
