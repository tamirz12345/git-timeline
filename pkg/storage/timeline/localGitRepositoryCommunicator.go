package timeline

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog"
)

// LocalGitRepository - an implimentation of GitRepositoryCommunitor for local folders
type LocalGitRepositoryCommunicator struct {
	repository *git.Repository
	repoPath   string
	mutex      sync.Mutex
	logger     *zerolog.Logger
}

// NewLocalGitRepositoryCommunicator - creates a new local git repository communicator
func NewLocalGitRepositoryCommunicator(repoPath string, logger *zerolog.Logger) (*LocalGitRepositoryCommunicator, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		if err != nil && strings.Contains(fmt.Sprint(err), "repository does not exist") {
			repo, err := git.PlainInit(repoPath, false)
			if err != nil {
				return nil, fmt.Errorf("failed to init git repository: %w", err)
			}
			return &LocalGitRepositoryCommunicator{
				repository: repo,
				logger:     logger,
				repoPath:   repoPath,
			}, nil
		}
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}
	return &LocalGitRepositoryCommunicator{
		repository: repo,
		logger:     logger,
		repoPath:   repoPath,
	}, nil
}

// CommitFile - adds and commits a file to the repository in case there are changes
func (l *LocalGitRepositoryCommunicator) CommitFile(filename, commitMessage, author string, content []byte) (string, error) {
	l.mutex.Lock() // to not corrupt the repository and to not need to handle conflicts and rebases
	defer l.mutex.Unlock()
	workTree, err := l.repository.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%s", l.repoPath, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(content)
	if err != nil {
		return "", fmt.Errorf("failed to write to output file: %w", err)
	}

	_, err = workTree.Add(filename)
	if err != nil {
		return "", fmt.Errorf("failed to add file: %w", err)
	}

	status, err := workTree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}
	if len(status.String()) == 0 {
		l.logger.Info().Str("filename", filename).Msg("no changes to commit")
		return "", nil
	}
	commitHash, err := workTree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
			Name: author,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}
	l.logger.Info().Str("filename", filename).Str("commitId", commitHash.String()).Msg("commit done")
	return commitHash.String(), nil
}

// GetFileContent - gets the content of a file in a specific commit
func (l *LocalGitRepositoryCommunicator) GetFileContent(filename, commitId string) ([]byte, *Commit, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	commitHash, err := l.repository.ResolveRevision(plumbing.Revision(commitId))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve revision: %w", err)
	}
	commit, err := l.repository.CommitObject(*commitHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get commit: %w", err)
	}
	file, err := commit.File(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %w", err)
	}
	content, err := file.Contents()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file content: %w", err)
	}
	return []byte(content), &Commit{
		Message:  commit.Message,
		Author:   commit.Author.Name,
		Date:     commit.Author.When,
		CommitId: commit.Hash.String(),
	}, nil
}

// GetFileTimeline - gets the timeline of a file
func (l *LocalGitRepositoryCommunicator) GetFileTimeline(filename string) (*Timeline, error) {
	ref, err := l.repository.Head()
	if err != nil {
		l.logger.Error().Err(err).Str("filename", filename).Msg("failed to get repository head")
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}
	historyIterator, err := l.repository.Log(&git.LogOptions{From: ref.Hash(), FileName: &filename})
	if err != nil {
		l.logger.Error().Err(err).Str("filename", filename).Msg("failed to get history iterator")
		return nil, fmt.Errorf("failed to get history iterator: %w", err)
	}
	l.logger.Info().Str("filename", filename).Msg("going over history of changes for file")
	commits := make([]Commit, 0)
	err = historyIterator.ForEach(func(c *object.Commit) error {
		commits = append(commits, Commit{
			Message:  c.Message,
			CommitId: c.Hash.String(),
			Author:   c.Author.Name,
			Date:     c.Author.When,
		})
		return nil
	})
	if err != nil {
		l.logger.Error().Err(err).Str("filename", filename).Msg("error during traversal of history")
		return nil, fmt.Errorf("failed to iterate over history: %w", err)
	}
	return &Timeline{
		Commits: commits,
	}, nil
}
