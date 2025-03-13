package github

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HtmlUrl     string `json:"html_url"`
	Languages   []string
}

// RepoList defines repos to fetch
var RepoList = []string{
	"go-inactivity-ping",
	"flock_manager",
	"db_shuffler",
	"ark-realtors",
	"arnoderry-movers",
}

var (
	repoCache     = make(map[string]*Repo)
	cacheMutex    sync.RWMutex
	lastCacheTime time.Time
	cacheDuration = 30 * time.Minute
)

// GetRepos fetches repos with caching
func GetRepos() ([]*Repo, error) {

	cacheMutex.RLock()
	cacheValid := !lastCacheTime.IsZero() && time.Since(lastCacheTime) < cacheDuration
	cacheMutex.RUnlock()

	if cacheValid {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()

		result := make([]*Repo, 0, len(repoCache))
		for _, repo := range repoCache {
			result = append(result, repo)
		}
		return result, nil
	}

	// Cache invalid, fetch fresh data
	repos, err := fetchAllRepos()
	if err != nil {
		return nil, err
	}

	// Update cache
	cacheMutex.Lock()
	lastCacheTime = time.Now()
	for _, repo := range repos {
		repoCache[repo.Name] = repo
	}
	cacheMutex.Unlock()

	return repos, nil
}

// fetchAllRepos fetches all repos in parallel
func fetchAllRepos() ([]*Repo, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	var wg sync.WaitGroup
	resultChan := make(chan *Repo, len(RepoList))
	errorChan := make(chan error, len(RepoList))

	// Start goroutines to fetch repos
	for _, repoName := range RepoList {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			repo, err := fetchRepoInfo(name)
			if err != nil {
				errorChan <- fmt.Errorf("error fetching %s: %v", name, err)
				return
			}
			resultChan <- repo
		}(repoName)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Check for errors
	if err := <-errorChan; err != nil {
		return nil, err
	}

	// Collect results
	var repos []*Repo
	for repo := range resultChan {
		repos = append(repos, repo)
	}

	return repos, nil
}

// fetchRepoInfo fetches info for a single repo
func fetchRepoInfo(repoName string) (*Repo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/ngenohkevin/%s", repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	}
	req.Header.Set("User-Agent", "Go-Github-Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var repo Repo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	// Fetch languages
	languages, err := fetchRepoLanguages(repoName)
	if err != nil {
		log.Printf("Warning: Error fetching languages for %s: %v", repoName, err)
	} else {
		repo.Languages = languages
	}

	return &repo, nil
}

// fetchRepoLanguages fetches languages for a repo
func fetchRepoLanguages(repoName string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/ngenohkevin/%s/languages", repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("User-Agent", "Go-GitHub-Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var langMap map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&langMap); err != nil {
		return nil, err
	}

	// Extract language names
	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}

	return languages, nil
}
