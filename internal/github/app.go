package github

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AppConfig holds GitHub App configuration
type AppConfig struct {
	AppID          string
	InstallationID string
	ClientID       string
	ClientSecret   string
}

// AppChecker verifies GitHub App installation and permissions
type AppChecker struct {
	client *Client
}

// NewAppChecker creates a new app checker
func NewAppChecker(client *Client) *AppChecker {
	return &AppChecker{client: client}
}

// CheckInstallation verifies the app is installed for a repo
func (a *AppChecker) CheckInstallation(owner, repo string) error {
	_, err := a.client.GetRepo(owner, repo)
	if err != nil {
		return fmt.Errorf("cannot access repository %s/%s: %w", owner, repo, err)
	}
	return nil
}

// CheckPermissions verifies required permissions by inspecting OAuth scopes
func (a *AppChecker) CheckPermissions() (*Permissions, error) {
	user, err := a.client.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to verify permissions: %w", err)
	}

	// Get OAuth scopes from API response headers
	scopes, err := a.getOAuthScopes()
	if err != nil {
		// Fallback: try to determine permissions by testing operations
		return a.checkPermissionsByProbing(user.Login)
	}

	return &Permissions{
		Authenticated: true,
		Username:      user.Login,
		CanReadRepos:  hasScope(scopes, "repo") || hasScope(scopes, "public_repo") || hasScope(scopes, "read:user"),
		CanWriteRepos: hasScope(scopes, "repo") || hasScope(scopes, "public_repo"),
	}, nil
}

// getOAuthScopes retrieves the X-OAuth-Scopes header from GitHub API
func (a *AppChecker) getOAuthScopes() ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	if a.client.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.client.token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Parse X-OAuth-Scopes header
	scopesHeader := resp.Header.Get("X-OAuth-Scopes")
	if scopesHeader == "" {
		return []string{}, nil
	}

	var scopes []string
	for _, s := range strings.Split(scopesHeader, ",") {
		scopes = append(scopes, strings.TrimSpace(s))
	}
	return scopes, nil
}

// checkPermissionsByProbing tests permissions via safe idempotent operations
func (a *AppChecker) checkPermissionsByProbing(username string) (*Permissions, error) {
	perms := &Permissions{
		Authenticated: true,
		Username:      username,
		CanReadRepos:  false,
		CanWriteRepos: false,
	}

	// Test read access by listing user repos
	_, err := a.client.ListRepos()
	if err == nil {
		perms.CanReadRepos = true
	}

	// For write access, we check if we can read repos (write implies read)
	// A more thorough check would require attempting a write operation,
	// but we avoid that to prevent side effects. Instead, we rely on
	// the OAuth scopes check above when available.
	// If scopes aren't available and we can read, assume write is possible
	// if the token was created with repo scope (common case).
	if perms.CanReadRepos {
		// Conservative: only mark write as true if we successfully got scopes
		// Since we're in the fallback path, we can't be certain about write access
		perms.CanWriteRepos = false
	}

	return perms, nil
}

// hasScope checks if a scope is present in the list
func hasScope(scopes []string, target string) bool {
	for _, s := range scopes {
		if s == target {
			return true
		}
	}
	return false
}

// Permissions represents GitHub permissions
type Permissions struct {
	Authenticated bool
	Username      string
	CanReadRepos  bool
	CanWriteRepos bool
}
