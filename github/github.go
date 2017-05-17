package github

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHub struct
type GitHub struct {
	owner  string
	repo   string
	Client *github.Client
}

// Initialize initialize
func Initialize(owner, repo string) *GitHub {
	// Setup Github Client
	// check arguments
	if len(repo) == 0 {
		log.Fatal("repository not specified")
	}
	if len(owner) == 0 {
		log.Fatal("owner not specified")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitHub{
		owner:  owner,
		repo:   repo,
		Client: client,
	}
}

// DeleteReleases delete release tag
func (g *GitHub) DeleteReleases(dry bool, year, months, days int, prefix string) error {
	ctx := context.Background()
	rls := g.ListReleases(year, months, days, prefix)
	if dry {
		return nil
	}
	for _, r := range rls {
		// Delete GitHub Release
		if _, e := g.Client.Repositories.DeleteRelease(ctx, g.owner, g.repo, *r.ID); e != nil {
			log.Fatalf("Deleting release %s failed; error: %s", r.GetName(), e)
		}
		// Delete Git tag
		if _, e := g.Client.Git.DeleteRef(ctx, g.owner, g.repo, "tags/"+*r.TagName); e != nil {
			// Draft release dose not have git tag sometimes.
			log.Printf("Deleting tag %s failed; error: %s", *r.TagName, e)
		}
	}
	return nil
}

// ListReleases get releases
func (g *GitHub) ListReleases(year, months, days int, prefix string) []*github.RepositoryRelease {
	log.Println("ID TagName TargetCommitish CreatedAt")
	var rls []*github.RepositoryRelease
	for page := 1; ; {
		ctx := context.Background()
		rs, res, _ := g.Client.Repositories.ListReleases(ctx, g.owner, g.repo, &github.ListOptions{Page: page})
		for _, r := range rs {
			if isTargetRelease(r, time.Now().AddDate(-1*year, -1*months, -1*days), prefix) {
				log.Printf("%d %s %s %s", *r.ID, r.GetName(), *r.TargetCommitish, *r.CreatedAt)
				rls = append(rls, r)
			}
		}
		// if current page is last page, LastPage value is 0
		if res.LastPage == 0 {
			break
		}
		page = res.NextPage
	}
	return rls
}

func isTargetRelease(r *github.RepositoryRelease, t time.Time, prefix string) bool {
	return r.CreatedAt.Time.Unix() < t.Unix() &&
		strings.HasPrefix(r.GetName(), prefix)
}
