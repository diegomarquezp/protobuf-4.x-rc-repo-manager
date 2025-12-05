package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var cleanupReleasePleaseCmd = &cobra.Command{
	Use:   "cleanup-release-please",
	Short: "Cleans up .github/release-please.yml files",
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

			fmt.Printf("---" + " Cleaning up %s ---" + "\n", repoName)
			cleanupReleasePlease(repoDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupReleasePleaseCmd)
}

func isMajorRelease(tag string) bool {
	if tag == "" {
		return false
	}
	versionStr := strings.TrimPrefix(tag, "v")
	parts := strings.Split(versionStr, ".")
	if len(parts) > 0 {
		major, err := strconv.Atoi(parts[0])
		if err == nil && major >= 1 {
			return true
		}
	}
	return false
}

func getLatestTag(repoDir string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "No names found") {
			return "", nil // No tags, not a fatal error
		}
		return "", fmt.Errorf("error getting latest tag: %w, %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

func cleanupReleasePlease(repoDir string) {
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

	// Part 1: Remove redundant options
	topLevelOptions := make(map[string]string)
	var branchesNode *yaml.Node
	for i := 0; i < len(mainMappingNode.Content); i += 2 {
		keyNode := mainMappingNode.Content[i]
		valueNode := mainMappingNode.Content[i+1]
		if keyNode.Value == "branches" {
			branchesNode = valueNode
		} else {
			topLevelOptions[keyNode.Value] = valueNode.Value
		}
	}

	if branchesNode != nil && branchesNode.Kind == yaml.SequenceNode {
		for _, branchNode := range branchesNode.Content {
			if branchNode.Kind == yaml.MappingNode {
				var newContent []*yaml.Node
				for i := 0; i < len(branchNode.Content); i += 2 {
					keyNode := branchNode.Content[i]
					valueNode := branchNode.Content[i+1]
					if topValue, ok := topLevelOptions[keyNode.Value]; ok && topValue == valueNode.Value {
						fmt.Printf("  - Removing redundant option '%s' from branch\n", keyNode.Value)
					} else {
						newContent = append(newContent, keyNode, valueNode)
					}
				}
				branchNode.Content = newContent
			}
		}
	}

	// Part 2: Remove bump-minor-pre-major for major releases
	latestTag, err := getLatestTag(repoDir)
	if err != nil {
		log.Printf("  - Could not get latest tag for %s: %v. Skipping bump-minor-pre-major check.", repoDir, err)
	} else if isMajorRelease(latestTag) {
		fmt.Printf("  - Repo is at major release (%s). Removing 'bump-minor-pre-major'.\n", latestTag)

		var newMainContent []*yaml.Node
		for i := 0; i < len(mainMappingNode.Content); i += 2 {
			keyNode := mainMappingNode.Content[i]
			if keyNode.Value != "bump-minor-pre-major" {
				newMainContent = append(newMainContent, keyNode, mainMappingNode.Content[i+1])
			}
		}
		mainMappingNode.Content = newMainContent

		if branchesNode != nil && branchesNode.Kind == yaml.SequenceNode {
			for _, branchNode := range branchesNode.Content {
				if branchNode.Kind == yaml.MappingNode {
					var newBranchContent []*yaml.Node
					for i := 0; i < len(branchNode.Content); i += 2 {
						keyNode := branchNode.Content[i]
						if keyNode.Value != "bump-minor-pre-major" {
							newBranchContent = append(newBranchContent, keyNode, branchNode.Content[i+1])
						}
					}
					branchNode.Content = newBranchContent
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

	fmt.Printf("Successfully cleaned up release-please config for %s\n", repoDir)
}
