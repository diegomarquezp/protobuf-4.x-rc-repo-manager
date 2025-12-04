package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var updateBranchCmd = &cobra.Command{
	Use:   "update-branch",
	Short: "Update a branch with the latest from main",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		all, _ := cmd.Flags().GetBool("all")
		branch, _ := cmd.Flags().GetString("branch")
		from, _ := cmd.Flags().GetString("from")

		if repo != "" {
			updateBranch(repo, branch, from)
		} else if all {
			repos, err := readLines("github_repositories.txt")
			if err != nil {
				fmt.Println("Error reading repositories file:", err)
				return
			}
			for _, r := range repos {
				repoDir := filepath.Base(r)
				updateBranch(repoDir, branch, from)
			}
		} else {
			fmt.Println("Please specify either a single repo with --repo or all repos with --all")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateBranchCmd)
	updateBranchCmd.Flags().StringP("repo", "r", "", "The repository to update")
	updateBranchCmd.Flags().BoolP("all", "a", false, "Update all repositories")
	updateBranchCmd.Flags().StringP("branch", "b", "protobuf-4.x-rc", "The branch to update")
	updateBranchCmd.Flags().StringP("from", "f", "main", "The branch to merge from")
}

func updateBranch(repoDir, branch, from string) {
	fmt.Printf("---" + " Updating branch '%s' in %s ---" + "\n", branch, repoDir)

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

	// Merge
	mergeCmd := exec.Command("git", "merge", "origin/"+from)
	mergeCmd.Dir = repoDir
	if output, err := mergeCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error merging in %s: %s\n%s", repoDir, err, output)
		return
	}

	// Push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Dir = repoDir
	if output, err := pushCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error pushing in %s: %s\n%s", repoDir, err, output)
		return
	}

	fmt.Printf("Successfully updated branch '%s' in %s\n", branch, repoDir)
}
