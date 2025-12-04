package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var updateReleasePleaseCmd = &cobra.Command{
	Use:   "update-release-please",
	Short: "Update release-please-config.json to set prerelease",
	Run: func(cmd *cobra.Command, args []string) {
		prerelease, _ := cmd.Flags().GetBool("prerelease")
		updateReleasePlease(prerelease)
	},
}

func init() {
	rootCmd.AddCommand(updateReleasePleaseCmd)
	updateReleasePleaseCmd.Flags().Bool("prerelease", true, "Set to true for prerelease, false otherwise")
}

func updateReleasePlease(prerelease bool) {
	repos, err := readLines("github_repositories.txt")
	if err != nil {
		fmt.Println("Error reading repositories file:", err)
		return
	}

	for _, repo := range repos {
		repoDir := filepath.Base(repo)
		configPath := filepath.Join(repoDir, "release-please-config.json")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("Skipping '%s': release-please-config.json not found\n", repoDir)
			continue
		}

		updateConfig(configPath, prerelease)
	}
}

func updateConfig(path string, prerelease bool) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", path, err)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Printf("Error unmarshalling %s: %v\n", path, err)
		return
	}

	if _, exists := data["prerelease"]; !exists {
		fmt.Printf("Updating '%s' to set 'prerelease: %v'\n", path, prerelease)

		content := string(file)
		lastBraceIndex := strings.LastIndex(content, "}")
		if lastBraceIndex == -1 {
			fmt.Printf("Could not find closing brace in %s\n", path)
			return
		}

		// Heuristic to check if a comma is needed.
		// This is not a proper JSON parser and might fail on edge cases.
		contentBeforeBrace := strings.TrimSpace(content[:lastBraceIndex])
		var insertion string
		if strings.HasSuffix(contentBeforeBrace, "{") {
			// No comma needed if it's the first element in an object
			insertion = fmt.Sprintf("\n  \"prerelease\": %v\n", prerelease)
		} else {
			insertion = fmt.Sprintf(",\n  \"prerelease\": %v\n", prerelease)
		}

		newContent := content[:lastBraceIndex] + insertion + content[lastBraceIndex:]

		if err := ioutil.WriteFile(path, []byte(newContent), 0644); err != nil {
			fmt.Printf("Error writing to %s: %v\n", path, err)
		}
	} else {
		fmt.Printf("Skipping '%s': 'prerelease' key already exists\n", path)
	}
}
