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
