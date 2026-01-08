package github

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	defaultClientID = "" // User must provide their own GitHub App client ID
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID string
	Scopes   []string
}

// OAuthFlow handles GitHub device flow authentication using golang.org/x/oauth2
type OAuthFlow struct {
	config *oauth2.Config
}

// NewOAuthFlow creates a new OAuth flow handler
func NewOAuthFlow(clientID string) *OAuthFlow {
	return &OAuthFlow{
		config: &oauth2.Config{
			ClientID: clientID,
			Scopes:   []string{"repo", "read:user"},
			Endpoint: github.Endpoint,
		},
	}
}

// InitiateDeviceFlow starts the device authorization flow
func (o *OAuthFlow) InitiateDeviceFlow() (*DeviceCode, error) {
	if o.config.ClientID == "" {
		return nil, fmt.Errorf("GitHub OAuth client ID not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	deviceAuth, err := o.config.DeviceAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate device flow: %w", err)
	}

	interval := 5
	if deviceAuth.Interval > 0 {
		interval = int(deviceAuth.Interval)
	}

	return &DeviceCode{
		DeviceCode:      deviceAuth.DeviceCode,
		UserCode:        deviceAuth.UserCode,
		VerificationURI: deviceAuth.VerificationURI,
		ExpiresIn:       int(time.Until(deviceAuth.Expiry).Seconds()),
		Interval:        interval,
	}, nil
}

// PollForToken polls GitHub for the access token using x/oauth2 device flow
func (o *OAuthFlow) PollForToken(deviceCode string, interval int) (string, error) {
	if interval < 5 {
		interval = 5
	}

	// Create a device auth response to use with DeviceAccessToken
	deviceAuth := &oauth2.DeviceAuthResponse{
		DeviceCode: deviceCode,
		Interval:   int64(interval),
		Expiry:     time.Now().Add(15 * time.Minute), // Default expiry
	}

	ctx := context.Background()

	// Use oauth2's built-in polling with custom options
	token, err := o.config.DeviceAccessToken(ctx, deviceAuth)
	if err != nil {
		// Parse common errors
		errStr := err.Error()
		if strings.Contains(errStr, "expired") {
			return "", fmt.Errorf("device code expired")
		}
		if strings.Contains(errStr, "denied") {
			return "", fmt.Errorf("access denied by user")
		}
		return "", fmt.Errorf("oauth error: %w", err)
	}

	return token.AccessToken, nil
}

// OpenBrowser opens the verification URL in the default browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
