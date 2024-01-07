package timeline

import "time"

// Commit - a commit struct
type Commit struct {
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
	Author   string    `json:"author"`
	CommitId string    `json:"commitId"`
}

// Timeline - struct of ordered commits representing a timeline of a single file
type Timeline struct {
	Commits []Commit `json:"commits"`
}

// GitRepositoryCommunitor - an interface for communicating with a git repository
type GitRepositoryCommunitor interface {
	CommitFile(string, string, string, []byte) (string, error)
	GetFileContent(string, string) ([]byte, *Commit, error)
	GetFileTimeline(string) (*Timeline, error)
}
