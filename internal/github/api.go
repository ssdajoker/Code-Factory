package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const githubAPIURL = "https://api.github.com"

// Client is a GitHub API client
type Client struct {
	token   string
	httpCli *http.Client
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUser returns the authenticated user
func (c *Client) GetUser() (*User, error) {
	var user User
	if err := c.get("/user", &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// ListRepos lists repositories for the authenticated user
func (c *Client) ListRepos() ([]Repository, error) {
	var repos []repoResponse
	if err := c.get("/user/repos?per_page=100", &repos); err != nil {
		return nil, err
	}

	result := make([]Repository, len(repos))
	for i, r := range repos {
		result[i] = Repository{
			Owner: r.Owner.Login,
			Name:  r.Name,
			URL:   r.HTMLURL,
		}
	}
	return result, nil
}

type repoResponse struct {
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
	Owner   struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// GetRepo gets a specific repository
func (c *Client) GetRepo(owner, name string) (*Repository, error) {
	var repo repoResponse
	if err := c.get(fmt.Sprintf("/repos/%s/%s", owner, name), &repo); err != nil {
		return nil, err
	}
	return &Repository{
		Owner: repo.Owner.Login,
		Name:  repo.Name,
		URL:   repo.HTMLURL,
	}, nil
}

// CreateBranch creates a new branch from a base branch
func (c *Client) CreateBranch(owner, repo, branch, baseBranch string) error {
	// Get base branch SHA
	var ref struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := c.get(fmt.Sprintf("/repos/%s/%s/git/ref/heads/%s", owner, repo, baseBranch), &ref); err != nil {
		return fmt.Errorf("failed to get base branch: %w", err)
	}

	// Create new branch
	body := map[string]string{
		"ref": "refs/heads/" + branch,
		"sha": ref.Object.SHA,
	}
	return c.post(fmt.Sprintf("/repos/%s/%s/git/refs", owner, repo), body, nil)
}

// GetFileContent gets the content of a file
func (c *Client) GetFileContent(owner, repo, path, branch string) (string, string, error) {
	var file struct {
		Content string `json:"content"`
		SHA     string `json:"sha"`
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/contents/%s?ref=%s", owner, repo, path, branch)
	if err := c.get(endpoint, &file); err != nil {
		return "", "", err
	}
	return file.Content, file.SHA, nil
}

// CreateOrUpdateFile creates or updates a file in the repository
func (c *Client) CreateOrUpdateFile(owner, repo, path, branch, message, content, sha string) error {
	body := map[string]string{
		"message": message,
		"content": content,
		"branch":  branch,
	}
	if sha != "" {
		body["sha"] = sha
	}
	return c.put(fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), body, nil)
}

// CreatePR creates a pull request
func (c *Client) CreatePR(owner, repo string, pr *PullRequest) (*PullRequest, error) {
	body := map[string]string{
		"title": pr.Title,
		"body":  pr.Body,
		"head":  pr.Head,
		"base":  pr.Base,
	}

	var result struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
	}
	if err := c.post(fmt.Sprintf("/repos/%s/%s/pulls", owner, repo), body, &result); err != nil {
		return nil, err
	}

	pr.Number = result.Number
	pr.URL = result.HTMLURL
	return pr, nil
}

func (c *Client) get(endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", githubAPIURL+endpoint, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, result)
}

func (c *Client) post(endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", githubAPIURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	return c.doRequest(req, result)
}

func (c *Client) put(endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", githubAPIURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	return c.doRequest(req, result)
}

func (c *Client) doRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github API error %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}
