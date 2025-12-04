package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories from a file",
	Run: func(cmd *cobra.Command, args []string) {
		branch, _ := cmd.Flags().GetString("branch")
		cloneRepos(branch)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringP("branch", "b", "protobuf-4.x-rc", "Branch to clone")
}

func cloneRepos(branch string) {
	repos, err := readLines("github_repositories.txt")
	if err != nil {
		fmt.Println("Error reading repositories file:", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	token, err := readToken(filepath.Join(homeDir, "GITHUB_TOKEN"))
	if err != nil {
		fmt.Println("Error reading GITHUB_TOKEN:", err)
		return
	}

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			cloneRepo(repo, token, branch)
		}(repo)
	}

	wg.Wait()
	fmt.Println("All repositories cloned.")
}

func cloneRepo(repo, token, branch string) {
	url := fmt.Sprintf("https://%s@github.com/%s.git", token, repo)
	cmd := exec.Command("git", "clone", "--branch", branch, "--depth", "1", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error cloning %s: %s\n%s", repo, err, output)
		return
	}
	fmt.Printf("Successfully cloned %s\n%s", repo, output)
}
