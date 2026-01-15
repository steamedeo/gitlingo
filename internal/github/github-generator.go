package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"steamedeo.dev/gitlingo/internal/config"
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
	Languages         []Language
}

type Language struct {
	Name  string
	Bytes int
}

type GithubStats struct {
	Username      string
	LanguageStats []Language
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func ProcessGithubData(ctx context.Context) (*GithubStats, error) {
	languageAggregations := map[string]int{}

	repos, err := fetchGithubData(ctx)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		for _, lang := range repo.Languages {
			languageAggregations[lang.Name] += lang.Bytes
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+config.AppConfig.GithubToken)
	req.Header.Add("User-Agent", "gitlingo")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, err
	}
	return repositories, nil
}

func fetchRepoLanguages(ctx context.Context, r []Repository) error {
	for i := range r {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", r[i].Owner.Login, r[i].Name)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+config.AppConfig.GithubToken)
		req.Header.Add("User-Agent", "gitlingo")
		res, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
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
