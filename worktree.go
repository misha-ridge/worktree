package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CurrentCommitID returns the current commit id
func CurrentCommitID() (string, error) {
	idCmd := exec.Command("git", "rev-parse", "HEAD")
	idCmd.Stderr = os.Stderr

	out, err := idCmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}

func currentCommitTimestamp() (time.Time, error) {
	timestampCmd := exec.Command("git", "show", "--format=%ct", "--no-patch")
	timestampCmd.Stderr = os.Stderr

	out, err := timestampCmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	ts, err := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(ts, 0), nil
}

// This function calculates a hash of the whole checked out worktree, including
// modified files.
//
// This is done by creating a separate git index file, adding the contents of
// the whole worktree to it and then writing it to the repository. ID of this
// object is the hash of the worktree.
//
// This causes some garbage accumulating in the repository, but it is cleaned up
// by a periodic 'git gc', and Git does not have an in-memory "simulate writing
// to repository" mode to avoid it.
func worktreeID() (string, error) {
	tempDir, err := os.MkdirTemp("", "worktree-id-")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.RemoveAll(tempDir) // we don't care about failed cleanup
	}()
	tempIndexFile := tempDir + "/index"

	gitAddCmd := exec.Command("git", "add", "--all")
	gitAddCmd.Stdout = os.Stdout
	gitAddCmd.Stderr = os.Stderr
	gitAddCmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+tempIndexFile)
	if err := gitAddCmd.Run(); err != nil {
		return "", err
	}

	gitWriteTreeCmd := exec.Command("git", "write-tree")
	gitWriteTreeCmd.Stderr = os.Stderr
	gitWriteTreeCmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+tempIndexFile)
	out, err := gitWriteTreeCmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(out), "\n"), nil
}

func semver(commitID string, commitTimestamp time.Time, dirtyWorktreeID string) string {
	ver := fmt.Sprintf("0.0.0-%s-%s", commitTimestamp.UTC().Format("20060102150405"), commitID[:12])
	if dirtyWorktreeID != "" {
		ver += "-dirty-" + dirtyWorktreeID
	}
	return ver
}

// Dirty returns whether the checked out worktree is dirty (contains new or modified files)
func Dirty() (bool, error) {
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Stderr = os.Stderr
	out, err := statusCmd.Output()
	if err != nil {
		return false, err
	}
	return len(out) != 0, nil
}

// CurrentVersion returns the current version of the worktree in SemVer 2.0
// format. Dirty trees get a stable version too (two identically-dirty trees
// will produce an identical version).
func CurrentVersion() (string, error) {
	commitID, err := CurrentCommitID()
	if err != nil {
		return "", err
	}
	commitTimestamp, err := currentCommitTimestamp()
	if err != nil {
		return "", err
	}
	dirty, err := Dirty()
	if err != nil {
		return "", err
	}
	var dirtyWorktreeID string
	if dirty {
		dirtyWorktreeID, err = worktreeID()
		if err != nil {
			return "", err
		}
	}
	return semver(commitID, commitTimestamp, dirtyWorktreeID), nil
}
