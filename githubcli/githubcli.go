package githubcli

import (
	"context"
	"fmt"
	"github.com/google/go-github/v42/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sort"
)

type GitHubRepo struct {
	Name     string
	Stars    int
	Url      string
	PushedAt github.Timestamp
}

type GitHubInterface struct {
	orgaName string
}

func (gi *GitHubInterface) SetOrgaName(name string) {
	gi.orgaName = name
}
func (gi *GitHubInterface) GetReposByStars() []GitHubRepo {
	var repos = gi.getRepos()
	fmt.Print("Get repos")
	for _, repo := range repos {
		fmt.Print(repo.Name)
	}
	return gi.sortReposByStars(repos)
}

func (gi *GitHubInterface) sortReposByStars(repos []*github.Repository) []GitHubRepo {
	var repoContainer []GitHubRepo
	for _, repo := range repos {
		repoContainer = append(repoContainer, GitHubRepo{*repo.Name, *repo.StargazersCount, *repo.HTMLURL, *repo.PushedAt})
	}

	sort.Slice(repoContainer, func(i, j int) bool {
		return repoContainer[i].PushedAt.Before(repoContainer[j].PushedAt.Time)
	})
	return repoContainer
}

func (gi *GitHubInterface) getRepos() []*github.Repository {
	client := github.NewClient(nil)

	var repos, _, err = client.Repositories.List(context.Background(), viper.GetString("github.organame"), nil)
	if err != nil {
		log.Errorf(err.Error())
	}
	return repos
}
