package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

const apiBaseURL = "https://api.github.com"

type Repo struct {
	Name string `json:"name"`
}

func GetLastUpdatedRepo(username string) (Repo, error) {
	var repos []Repo
	url := fmt.Sprintf("%s/users/%s/repos?sort=pushed&per_page=1", apiBaseURL, username)

	resp, err := http.Get(url)
	if err != nil {
		return Repo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repo{}, fmt.Errorf("user not found or api error")
	}

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return Repo{}, err
	}

	if len(repos) == 0 {
		return Repo{}, fmt.Errorf("no repositories found")
	}

	return repos[0], nil
}

func GetRepoPrimaryLanguage(username, repoName string) (string, error) {
	var languages map[string]int
	url := fmt.Sprintf("%s/repos/%s/%s/languages", apiBaseURL, username, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("repository not found or api error")
	}

	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return "", err
	}

	if len(languages) == 0 {
		return "N/A", nil
	}

	type langKV struct {
		Key   string
		Value int
	}

	var ss []langKV
	for k, v := range languages {
		ss = append(ss, langKV{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss[0].Key, nil
}
