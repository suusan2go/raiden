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
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
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
		rs, res, err := g.Client.Repositories.ListReleases(ctx, g.owner, g.repo, &github.ListOptions{Page: page})
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range rs {
			if isTargetRelease(r, time.Now().AddDate(-1*year, -1*months, -1*days), prefix) {
				log.Printf("%d %s %s %s", *r.ID, releaseName(r), *r.TargetCommitish, *r.CreatedAt)
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

// DeleteTags delete release tag
func (g *GitHub) DeleteTags(dry bool, year, months, days int, prefix string) error {
	// ctx := context.Background()
	tags := g.ListTags(year, months, days, prefix)
	if dry {
		return nil
	}
	for _, t := range tags {
		// Delete Git tag
		ctx := context.Background()
		if _, e := g.Client.Git.DeleteRef(ctx, g.owner, g.repo, "tags/"+t.GetName()); e != nil {
			log.Fatalf("Deleting tag %s failed; error: %s", t.GetName(), e)
		}
	}
	return nil
}

// ListTags get tags from github api
func (g *GitHub) ListTags(year, months, days int, prefix string) []*github.RepositoryTag {
	log.Println("ID TagName TargetCommitish CreatedAt")
	var tags []*github.RepositoryTag
	for page := 1; ; {
		ctx := context.Background()
		rts, res, err := g.Client.Repositories.ListTags(ctx, g.owner, g.repo, &github.ListOptions{Page: page, PerPage: 100})
		if err != nil {
			log.Fatal(err)
		}
		for _, rt := range rts {
			ctx := context.Background()
			c, _, err := g.Client.Git.GetCommit(ctx, g.owner, g.repo, rt.Commit.GetSHA())
			rt.Commit = c
			if err != nil {
				log.Fatal(err)
			}
			if isTargetTag(rt, time.Now().AddDate(-1*year, -1*months, -1*days), prefix) {
				log.Printf("%s %s", rt.GetName(), rt.Commit.Author.GetDate())
				tags = append(tags, rt)
			}
		}
		// if current page is last page, LastPage value is 0
		if res.LastPage == 0 {
			break
		}
		page = res.NextPage
	}
	return tags
}

func releaseName(r *github.RepositoryRelease) string {
	var name string
	if len(r.GetName()) == 0 {
		name = *r.TagName
	} else {
		name = r.GetName()
	}
	return name
}

func isTargetRelease(r *github.RepositoryRelease, t time.Time, prefix string) bool {
	return r.CreatedAt.Time.Unix() < t.Unix() &&
		strings.HasPrefix(releaseName(r), prefix)
}

func isTargetTag(tg *github.RepositoryTag, t time.Time, prefix string) bool {
	if tg.Commit.Author != nil {
		return tg.Commit.Author.Date.Unix() < t.Unix() &&
			strings.HasPrefix(tg.GetName(), prefix)
	}
	return strings.HasPrefix(tg.GetName(), prefix)
}
