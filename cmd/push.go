package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Commit and push changes for each repository",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		pushChanges(message)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("message", "m", "feat: update release-please config", "Commit message")
}

func pushChanges(message string) {
	repos, err := readLines("github_repositories.txt")
	if err != nil {
		fmt.Println("Error reading repositories file:", err)
		return
	}

	for _, repo := range repos {
		repoDir := filepath.Base(repo)
		pushChange(repoDir, message)
	}
}

func pushChange(repoDir, message string) {
	fmt.Printf("--- Pushing changes for %s ---\n", repoDir)

	// Add
	addCmd := exec.Command("git", "add", "release-please-config.json")
	addCmd.Dir = repoDir
	if output, err := addCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error adding changes in %s: %s\n%s", repoDir, err, output)
		return
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = repoDir
	if output, err := commitCmd.CombinedOutput(); err != nil {
		// It's possible there are no changes to commit, so we check for that
		if !strings.Contains(string(output), "nothing to commit") {
			fmt.Printf("Error committing in %s: %s\n%s", repoDir, err, output)
			return
		}
		fmt.Printf("No changes to commit in %s\n", repoDir)
		return
	}

	// Push
	pushCmd := exec.Command("git", "push", "origin", "protobuf-4.x-rc")
	pushCmd.Dir = repoDir
	if output, err := pushCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error pushing in %s: %s\n%s", repoDir, err, output)
		return
	}

	fmt.Printf("Successfully pushed changes for %s\n", repoDir)
}
