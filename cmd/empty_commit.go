package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var emptyCommitCmd = &cobra.Command{
	Use:   "empty-commit",
	Short: "Push an empty commit to a branch",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		all, _ := cmd.Flags().GetBool("all")
		branch, _ := cmd.Flags().GetString("branch")
		message, _ := cmd.Flags().GetString("message")

		if repo != "" {
			emptyCommit(repo, branch, message)
		} else if all {
			repos, err := readLines("github_repositories.txt")
			if err != nil {
				fmt.Println("Error reading repositories file:", err)
				return
			}
			for _, r := range repos {
				repoDir := filepath.Base(r)
				emptyCommit(repoDir, branch, message)
			}
		} else {
			fmt.Println("Please specify either a single repo with --repo or all repos with --all")
		}
	},
}

func init() {
	rootCmd.AddCommand(emptyCommitCmd)
	emptyCommitCmd.Flags().StringP("repo", "r", "", "The repository to update")
	emptyCommitCmd.Flags().BoolP("all", "a", false, "Update all repositories")
	emptyCommitCmd.Flags().StringP("branch", "b", "protobuf-4.x-rc", "The branch to commit to")
	emptyCommitCmd.Flags().StringP("message", "m", "chore: empty commit", "The commit message")
}

func emptyCommit(repoDir, branch, message string) {
	fmt.Printf("--- Pushing empty commit to '%s' in %s ---\n", branch, repoDir)

	// Fetch
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = repoDir
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error fetching in %s: %s\n%s", repoDir, err, output)
		return
	}

	// Checkout
	checkoutCmd := exec.Command("git", "checkout", branch)
	checkoutCmd.Dir = repoDir
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error checking out branch in %s: %s\n%s", repoDir, err, output)
		return
	}

	// Empty commit
	commitCmd := exec.Command("git", "commit", "--allow-empty", "-m", message)
	commitCmd.Dir = repoDir
	if output, err := commitCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error committing in %s: %s\n%s", repoDir, err, output)
		return
	}

	// Push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Dir = repoDir
	if output, err := pushCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error pushing in %s: %s\n%s", repoDir, err, output)
		return
	}

	fmt.Printf("Successfully pushed empty commit to '%s' in %s\n", branch, repoDir)
}
