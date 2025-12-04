package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var checkBranchCmd = &cobra.Command{
	Use:   "check-branch",
	Short: "Check the current branch of each repository",
	Run: func(cmd *cobra.Command, args []string) {
		branch, _ := cmd.Flags().GetString("branch")
		checkBranches(branch)
	},
}

func init() {
	rootCmd.AddCommand(checkBranchCmd)
	checkBranchCmd.Flags().StringP("branch", "b", "protobuf-4.x-rc", "Branch to check for")
}

func checkBranches(branch string) {
	repos, err := readLines("github_repositories.txt")
	if err != nil {
		fmt.Println("Error reading repositories file:", err)
		return
	}

	for _, repo := range repos {
		repoDir := filepath.Base(repo)
		checkBranch(repoDir, branch)
	}
}

func checkBranch(repoDir, expectedBranch string) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error checking branch for %s: %s\n%s", repoDir, err, output)
		return
	}

	branch := strings.TrimSpace(string(output))
	if branch == expectedBranch {
		fmt.Printf("Repository '%s' is on the correct branch: %s\n", repoDir, branch)
	} else {
		fmt.Printf("Repository '%s' is on an incorrect branch: %s\n", repoDir, branch)
	}
}
