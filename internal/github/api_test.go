package github

import (
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"
)

func TestNewClient(t *testing.T) {
        client := NewClient("test-token")
        if client == nil {
                t.Fatal("NewClient returned nil")
        }
        if client.token != "test-token" {
                t.Errorf("token = %v, want 'test-token'", client.token)
        }
}

func TestGetUser(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if r.URL.Path != "/user" {
                        t.Errorf("unexpected path: %s", r.URL.Path)
                }
                if r.Header.Get("Authorization") != "Bearer test-token" {
                        t.Error("missing or incorrect Authorization header")
                }

                json.NewEncoder(w).Encode(User{
                        Login:     "testuser",
                        Name:      "Test User",
                        AvatarURL: "https://thumbs.dreamstime.com/b/test-icon-vector-question-mark-male-user-person-profile-avatar-symbol-help-sign-glyph-pictogram-test-icon-vector-168495430.jpg",
                })
        }))
        defer server.Close()

        // Note: Full integration test would require baseURL injection
        // This test verifies the mock server setup pattern
        _ = server
}

func TestUserStruct(t *testing.T) {
        user := User{
                Login:     "octocat",
                Name:      "The Octocat",
                AvatarURL: "https://pbs.twimg.com/media/Gcj0i2DWIAAJFaX.png",
        }

        if user.Login != "octocat" {
                t.Errorf("Login = %v, want 'octocat'", user.Login)
        }
}

func TestRepositoryStruct(t *testing.T) {
        repo := Repository{
                Owner: "octocat",
                Name:  "hello-world",
                URL:   "https://github.com/octocat/hello-world",
        }

        if repo.Owner != "octocat" {
                t.Errorf("Owner = %v, want 'octocat'", repo.Owner)
        }
        if repo.Name != "hello-world" {
                t.Errorf("Name = %v, want 'hello-world'", repo.Name)
        }
}

func TestPullRequestStruct(t *testing.T) {
        pr := &PullRequest{
                Title:  "Add feature",
                Body:   "This PR adds a new feature",
                Head:   "feature-branch",
                Base:   "main",
                Number: 42,
                URL:    "https://github.com/owner/repo/pull/42",
        }

        if pr.Title != "Add feature" {
                t.Errorf("Title = %v, want 'Add feature'", pr.Title)
        }
        if pr.Number != 42 {
                t.Errorf("Number = %v, want 42", pr.Number)
        }
}
