package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneRepo(repoURL, branch, workdir string) (string, error) {
	if workdir == "" {
		return "", errors.New("WORKDIR environment variable is not set")
	}

	if branch == "" {
		branch = "main"
	}

	// Create a unique folder name with the branch or tag suffix
	folderName := filepath.Join(workdir, fmt.Sprintf("%s-branch-%s", filepath.Base(repoURL), branch))
	repoPath, err := filepath.Abs(folderName)
	if err != nil {
		return "", fmt.Errorf("failed to create repository path: %v", err)
	}

	// Check if the repository folder exists
	_, err = os.Stat(repoPath)
	if os.IsNotExist(err) {
		// Clone the repository if the folder does not exist
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			return "", fmt.Errorf("failed to clone repository: %v", err)
		}
	} else if err == nil {
		// Update the repository if the folder already exists
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return "", fmt.Errorf("failed to open repository: %v", err)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return "", fmt.Errorf("failed to get worktree: %v", err)
		}

		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			fmt.Print(err)
			return "", fmt.Errorf("failed to pull repository '%s' (branch: '%s'): %v", repoURL, branch, err)
		}

	} else {
		return "", fmt.Errorf("failed to access repository path: %v", err)
	}

	return repoPath, nil
}
