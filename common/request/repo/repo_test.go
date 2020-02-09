package repo

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestRepo(t *testing.T) {
	var err error

	var ctx = context.Background()

	var client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "c77d41f2027f009d0cb088068e86d6465c3ddaf0"})))
	if err = Repo(ctx, client, "", 3); err != nil {
		t.Fatal(err)
	}
}
