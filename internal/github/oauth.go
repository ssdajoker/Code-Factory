package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"
)

const (
	githubDeviceCodeURL = "https://github.com/login/device/code"
	githubTokenURL      = "https://github.com/login/oauth/access_token"
	defaultClientID     = "" // User must provide their own GitHub App client ID
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID string
	Scopes   []string
}

// OAuthFlow handles GitHub device flow authentication
type OAuthFlow struct {
	config OAuthConfig
	client *http.Client
}

// NewOAuthFlow creates a new OAuth flow handler
func NewOAuthFlow(clientID string) *OAuthFlow {
	scopes := []string{"repo", "read:user"}
	return &OAuthFlow{
		config: OAuthConfig{
			ClientID: clientID,
			Scopes:   scopes,
		},
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// InitiateDeviceFlow starts the device authorization flow
func (o *OAuthFlow) InitiateDeviceFlow() (*DeviceCode, error) {
	if o.config.ClientID == "" {
		return nil, fmt.Errorf("GitHub OAuth client ID not configured")
	}

	data := url.Values{}
	data.Set("client_id", o.config.ClientID)
	data.Set("scope", joinScopes(o.config.Scopes))

	req, err := http.NewRequest("POST", githubDeviceCodeURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate device flow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("device flow error: %s", string(body))
	}

	var result struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &DeviceCode{
		DeviceCode:      result.DeviceCode,
		UserCode:        result.UserCode,
		VerificationURI: result.VerificationURI,
		ExpiresIn:       result.ExpiresIn,
		Interval:        result.Interval,
	}, nil
}

// PollForToken polls GitHub for the access token
func (o *OAuthFlow) PollForToken(deviceCode string, interval int) (string, error) {
	if interval < 5 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		data := url.Values{}
		data.Set("client_id", o.config.ClientID)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest("POST", githubTokenURL, bytes.NewBufferString(data.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := o.client.Do(req)
		if err != nil {
			continue
		}

		var result struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			Scope       string `json:"scope"`
			Error       string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		switch result.Error {
		case "":
			if result.AccessToken != "" {
				return result.AccessToken, nil
			}
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5
			continue
		case "expired_token":
			return "", fmt.Errorf("device code expired")
		case "access_denied":
			return "", fmt.Errorf("access denied by user")
		default:
			return "", fmt.Errorf("oauth error: %s", result.Error)
		}
	}
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

func joinScopes(scopes []string) string {
	result := ""
	for i, s := range scopes {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}
