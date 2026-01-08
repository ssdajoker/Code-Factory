package github

// GitHubService defines the interface for GitHub operations
type GitHubService interface {
	// OAuth
	InitiateDeviceFlow() (*DeviceCode, error)
	PollForToken(deviceCode string) (string, error)
	
	// Repository operations
	ListRepos() ([]Repository, error)
	GetRepo(owner, name string) (*Repository, error)
	
	// Branch & Commit
	CreateBranch(repo, branch, baseBranch string) error
	Commit(repo, branch string, files []FileChange) error
	
	// Pull Request
	CreatePR(repo string, pr *PullRequest) (*PullRequest, error)
}

// DeviceCode represents OAuth device flow code
type DeviceCode struct {
	DeviceCode      string
	UserCode        string
	VerificationURI string
	ExpiresIn       int
	Interval        int
}

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
	URL   string
}

// FileChange represents a file change for commit
type FileChange struct {
	Path    string
	Content string
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number int
	Title  string
	Body   string
	Head   string
	Base   string
	URL    string
}
