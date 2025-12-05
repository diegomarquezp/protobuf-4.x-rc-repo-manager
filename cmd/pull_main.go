package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var pullMainCmd = &cobra.Command{
	Use:   "pull-main",
	Short: "Checkout main and pull the latest changes for all repositories",
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := readLines("github_repositories.txt")
		if err != nil {
			log.Fatalf("Failed to read github_repositories.txt: %v", err)
		}

		for _, repoStr := range repos {
			repoParts := strings.Split(repoStr, "/")
			if len(repoParts) != 2 {
				log.Printf("Skipping invalid repo format: %s", repoStr)
				continue
			}
			repoName := repoParts[1]
			repoDir, err := filepath.Abs(repoName)
			if err != nil {
				log.Printf("Could not get absolute path for %s: %v", repoName, err)
				continue
			}

			fmt.Printf("--- Updating %s ---\n", repoName)

			// Check if main branch exists
			verifyCmd := exec.Command("git", "show-branch", "remotes/origin/main")
			verifyCmd.Dir = repoDir
			if err := verifyCmd.Run(); err != nil {
				fmt.Printf("Repository %s does not have a main branch, skipping.\n", repoName)
				continue
			}
			
			// git checkout main
			checkoutCmd := exec.Command("git", "checkout", "main")
			checkoutCmd.Dir = repoDir
			if output, err := checkoutCmd.CombinedOutput(); err != nil {
				fmt.Printf("Error checking out main in %s: %s\n%s", repoName, err, string(output))
				continue
			}
			fmt.Printf("Checked out main in %s\n", repoName)

			// git pull origin main
			pullCmd := exec.Command("git", "pull", "origin", "main")
			pullCmd.Dir = repoDir
			if output, err := pullCmd.CombinedOutput(); err != nil {
				fmt.Printf("Error pulling main in %s: %s\n%s", repoName, err, string(output))
				continue
			}
			fmt.Printf("Pulled latest changes for main in %s\n", repoName)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullMainCmd)
}
