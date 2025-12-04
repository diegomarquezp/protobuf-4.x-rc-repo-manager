package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addSubmodulesCmd = &cobra.Command{
	Use:   "add-submodules",
	Short: "Add repositories as submodules",
	Run: func(cmd *cobra.Command, args []string) {
		addSubmodules()
	},
}

func init() {
	rootCmd.AddCommand(addSubmodulesCmd)
}

func addSubmodules() {
	repos, err := readLines("github_repositories.txt")
	if err != nil {
		fmt.Println("Error reading repositories file:", err)
		return
	}

	for _, repo := range repos {
		repoDir := filepath.Base(repo)
		url := fmt.Sprintf("https://github.com/%s.git", repo)

		// Remove existing directory
		fmt.Printf("Removing existing directory: %s\n", repoDir)
		if err := os.RemoveAll(repoDir); err != nil {
			fmt.Printf("Error removing directory %s: %v\n", repoDir, err)
			continue
		}

		// Add submodule
		fmt.Printf("Adding submodule for %s\n", repo)
		cmd := exec.Command("git", "submodule", "add", "-b", "protobuf-4.x-rc", url, repoDir)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Error adding submodule %s: %s\n%s", repo, err, output)
			continue
		}
		fmt.Printf("Successfully added submodule for %s\n", repo)
	}

	fmt.Println("All submodules added.")
}
