package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"steamedeo.dev/gitlingo/internal/relay"
)

type Owner struct {
	Login string `json:"login"`
}

type Repository struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CreationTimestamp string `json:"created_at"`
	PrimaryLanguage   string `json:"language"`
	Owner             Owner  `json:"owner"`
	Fork              bool   `json:"fork"`
	Languages         []Language
}

type Language struct {
	Name  string
	Bytes int
}

type GithubStats struct {
	Username      string
	LanguageStats []Language
	TotalBytes    int
}

func ProcessGithubData(ctx context.Context) (*GithubStats, error) {
	languageAggregations := map[string]int{}
	totalBytes := 0

	repos, err := fetchGithubData(ctx)
	if err != nil {
		return nil, err
	}
	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found for the authenticated user")
	}
	for _, repo := range repos {
		for _, lang := range repo.Languages {
			languageAggregations[lang.Name] += lang.Bytes
			totalBytes += lang.Bytes
		}
	}

	// Convert the map to a slice so that it can be sorted
	languageStats := make([]Language, 0, len(languageAggregations))
	for name, bytes := range languageAggregations {
		languageStats = append(languageStats, Language{Name: name, Bytes: bytes})
	}
	sort.Slice(languageStats, func(i, j int) bool {
		return languageStats[i].Bytes > languageStats[j].Bytes
	})
	return &GithubStats{
		LanguageStats: languageStats,
		Username:      repos[0].Owner.Login,
		TotalBytes:    totalBytes,
	}, nil
}

func fetchGithubData(ctx context.Context) ([]Repository, error) {
	page := 1
	repos := []Repository{}
	for {
		reposPage, err := fetchRepositories(ctx, page)
		if err != nil {
			return nil, err
		}
		if len(reposPage) == 0 {
			break
		}
		repos = append(repos, reposPage...)
		page++
	}
	if err := fetchRepoLanguages(ctx, repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func fetchRepositories(ctx context.Context, page int) ([]Repository, error) {
	var repositories []Repository
	url := fmt.Sprintf("https://api.github.com/user/repos?affiliation=owner&per_page=100&page=%d", page)
	relay, err := relay.NewRelay(relay.RelayOptions{
		Context:        ctx,
		URL:            url,
		Method:         "GET",
		RequestPayload: nil,
	})
	if err != nil {
		return nil, err
	}
	res, err := relay.SendHTTPRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var allRepos []Repository
	if err := json.Unmarshal(body, &allRepos); err != nil {
		return nil, err
	}

	// Filter out forks
	for _, repo := range allRepos {
		if !repo.Fork {
			repositories = append(repositories, repo)
		}
	}

	return repositories, nil
}

func fetchRepoLanguages(ctx context.Context, r []Repository) error {
	for i := range r {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", r[i].Owner.Login, r[i].Name)
		relay, err := relay.NewRelay(relay.RelayOptions{
			Context:        ctx,
			URL:            url,
			Method:         "GET",
			RequestPayload: nil,
		})
		if err != nil {
			return err
		}
		res, err := relay.SendHTTPRequest()
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			res.Body.Close()
			return err
		}
		res.Body.Close()
		var langMap map[string]int
		if err := json.Unmarshal(body, &langMap); err != nil {
			return err
		}
		for name, bytes := range langMap {
			r[i].Languages = append(r[i].Languages, Language{Name: name, Bytes: bytes})
		}
	}
	return nil
}
