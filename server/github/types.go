package github

type PushEvent struct {
	Ref          string      `json:"ref"`
	Before       string      `json:"before"`
	After        string      `json:"after"`
	Repository   RepoFull    `json:"repository"`
	Pusher       Pusher      `json:"pusher"`
	Sender       User        `json:"sender"`
	Installation InstallMini `json:"installation"`
	Created      bool        `json:"created"`
	Deleted      bool        `json:"deleted"`
	Forced       bool        `json:"forced"`
	BaseRef      *string     `json:"base_ref"`
	Compare      string      `json:"compare"`
	Commits      []Commit    `json:"commits"`
	HeadCommit   Commit      `json:"head_commit"`
}

type User struct {
	ID                int64  `json:"id"`
	Login             string `json:"login"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	HTMLURL           string `json:"html_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	URL               string `json:"url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	OrganizationsURL  string `json:"organizations_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	GravatarID        string `json:"gravatar_id"`
	UserViewType      string `json:"user_view_type"`
}

type RepoOwner struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type RepoFull struct {
	ID               int64     `json:"id"`
	NodeID           string    `json:"node_id"`
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Private          bool      `json:"private"`
	Owner            RepoOwner `json:"owner"`
	HTMLURL          string    `json:"html_url"`
	Description      *string   `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	CreatedAt        int64     `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	PushedAt         int64     `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SVNURL           string    `json:"svn_url"`
	Homepage         *string   `json:"homepage"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	Language         *string   `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	HasDiscussions   bool      `json:"has_discussions"`
	ForksCount       int       `json:"forks_count"`
	MirrorURL        *string   `json:"mirror_url"`
	Archived         bool      `json:"archived"`
	Disabled         bool      `json:"disabled"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	License          any       `json:"license"`
	AllowForking     bool      `json:"allow_forking"`
	IsTemplate       bool      `json:"is_template"`
	WebCommitSignoff bool      `json:"web_commit_signoff_required"`
	Topics           []string  `json:"topics"`
	Visibility       string    `json:"visibility"`
	Forks            int       `json:"forks"`
	OpenIssues       int       `json:"open_issues"`
	Watchers         int       `json:"watchers"`
	DefaultBranch    string    `json:"default_branch"`
	Stargazers       int       `json:"stargazers"`
	MasterBranch     string    `json:"master_branch"`
}

type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CommitUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	ID        string     `json:"id"`
	TreeID    string     `json:"tree_id"`
	Distinct  bool       `json:"distinct"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    CommitUser `json:"author"`
	Committer CommitUser `json:"committer"`
	Added     []string   `json:"added"`
	Removed   []string   `json:"removed"`
	Modified  []string   `json:"modified"`
}

type InstallMini struct {
	ID     int64  `json:"id"`
	NodeID string `json:"node_id"`
}

type GithubApp struct {
	AppID      int
	PrivateKey string
}

type GithubInstallation struct {
	InstallationID int
	AccessToken    *string
	TokenExpiresAt *string
}

type LatestCommit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	URL     string `json:"html_url"`
	Author  string `json:"author"`
}
