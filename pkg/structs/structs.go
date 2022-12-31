package structs

import "time"

type Build struct {
	ID            int64     `json:"buildID"`
	GithubRepoID  int64     `json:"githubRepoID"`
	Branch        string    `json:"branch"`
	IsSuccessfull bool      `json:"isSuccessfull"`
	CreatedAt     time.Time `json:"createdAt"`
	Stdout        string    `json:"stdout"`
}

type Artifact struct {
	ID      int64  `json:"artifactID"`
	BuildID int64  `json:"buildID"`
	Name    string `json:"name"`
}

type CIBuilderPayload struct {
	RepoName string `json:"repoName"`
	RepoURL  string `json:"repoUrl"`
	Username string `json:"username"`
	Branch   string `json:"branch"`
	Tag      string `json:"tag"`
}

type GithubStatusMessage struct {
	Description string `json:"description"`
	Context     string `json:"context"`
	State       string `json:"state"`
}
