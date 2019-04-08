package rate

import (
	"context"
	"sync"
	"time"

	"github.com/EffDataAly/GithubTraveler/common/util"
	"github.com/google/go-github/github"
)

var locker = new(sync.Mutex)

var remaining int

var reset github.Timestamp

func Initialize(ctx context.Context, client *github.Client) (err error) {
	if util.CheckCtx(ctx) {
		return
	}

	var response *github.Response
	if _, response, err = client.Repositories.List(ctx, "", nil); err != nil {
		return
	}
	remaining = response.Rate.Remaining
	reset = response.Rate.Reset
	return
}

// Get get the visit github locker
func Get() {
	locker.Lock()
	defer locker.Unlock()
	if remaining == 0 {
		ticker := time.NewTicker(500 * time.Millisecond)
		for t := range ticker.C {
			if reset.Before(t) {
				break
			}
		}
	}
}

// Set set github visit locker
func Set(r int) {
	if r > remaining {
		return
	}
	remaining = r
}
