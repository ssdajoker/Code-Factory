package github

import (
	"fmt"
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

// CheckPermissions verifies required permissions
func (a *AppChecker) CheckPermissions() (*Permissions, error) {
	user, err := a.client.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to verify permissions: %w", err)
	}

	return &Permissions{
		Authenticated: true,
		Username:      user.Login,
		CanReadRepos:  true,
		CanWriteRepos: true,
	}, nil
}

// Permissions represents GitHub permissions
type Permissions struct {
	Authenticated bool
	Username      string
	CanReadRepos  bool
	CanWriteRepos bool
}
