package agit

import (
	"fmt"
	"os/exec"
	"time"
)

const branchPrefix = "agent-"

func MakeBranchName() string {
	return fmt.Sprintf("%s%d", branchPrefix, time.Now().UnixNano())
}

func GitStatus() (string, error) {
	cmd := exec.Command("git", "status")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
}

func GitSwitchCreate(branch string) (string, error) {
	cmd := exec.Command("git", "switch", "-c", branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
}

func GitAddAll() (string, error) {
	cmd := exec.Command("git", "add", ".")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
}

func GitCommit(commit string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", commit)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
}

func GitPushBranch(branch string) (string, error) {
	cmd := exec.Command("git", "push", "--set-upstream", "origin", branch)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
}
