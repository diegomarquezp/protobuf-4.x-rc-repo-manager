package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var formatReleasePleaseCmd = &cobra.Command{
	Use:   "format-release-please",
	Short: "Formats .github/release-please.yml files to have 'branch' as the first key.",
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

			fmt.Printf("--- Formatting %s ---\n", repoName)
			formatReleasePlease(repoDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(formatReleasePleaseCmd)
}

func formatReleasePlease(repoDir string) {
	configPath := filepath.Join(repoDir, ".github", "release-please.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(repoDir, ".github", "release-please.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("No release-please config found for %s, skipping.\n", repoDir)
			return
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Failed to read %s: %v", configPath, err)
		return
	}

	var root yaml.Node
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Printf("Failed to unmarshal YAML from %s: %v", configPath, err)
		return
	}

	if len(root.Content) == 0 {
		return // empty file
	}
	mainMappingNode := root.Content[0]

	var branchesNode *yaml.Node
	for i := 0; i < len(mainMappingNode.Content); i += 2 {
		keyNode := mainMappingNode.Content[i]
		if keyNode.Value == "branches" {
			branchesNode = mainMappingNode.Content[i+1]
			break
		}
	}

	if branchesNode != nil && branchesNode.Kind == yaml.SequenceNode {
		for _, branchNode := range branchesNode.Content {
			if branchNode.Kind == yaml.MappingNode {
				// Find the branch key
				var branchKeyNode, branchValueNode *yaml.Node
				var branchIndex = -1
				for i := 0; i < len(branchNode.Content); i += 2 {
					if branchNode.Content[i].Value == "branch" {
						branchKeyNode = branchNode.Content[i]
						branchValueNode = branchNode.Content[i+1]
						branchIndex = i
						break
					}
				}

				// If branch key exists and is not the first key, move it to the front
				if branchIndex > 0 {
					// Remove from current position
					branchNode.Content = append(branchNode.Content[:branchIndex], branchNode.Content[branchIndex+2:]...)
					// Add to the front
					branchNode.Content = append([]*yaml.Node{branchKeyNode, branchValueNode}, branchNode.Content...)
					fmt.Println("  - Reordered 'branch' key to be first.")
				}
			}
		}
	}

	marshaledData, err := yaml.Marshal(&root)
	if err != nil {
		log.Printf("Failed to marshal YAML for %s: %v", configPath, err)
		return
	}

	err = os.WriteFile(configPath, marshaledData, 0644)
	if err != nil {
		log.Printf("Failed to write %s: %v", configPath, err)
		return
	}
	fmt.Printf("Successfully formatted release-please config for %s\n", repoDir)
}
