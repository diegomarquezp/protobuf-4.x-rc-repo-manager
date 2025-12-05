package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFormatReleasePlease(t *testing.T) {
	t.Run("reorders branch key to be first", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)

		config := `
releaseType: java-yoshi
branches:
  - releaseType: java-backport
    branch: 1.0.x
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		formatReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)

		var root yaml.Node
		err = yaml.Unmarshal(data, &root)
		assert.NoError(t, err)

		branchesNode := root.Content[0].Content[3] // branches
		branchNode := branchesNode.Content[0]
		assert.Equal(t, "branch", branchNode.Content[0].Value)
		assert.Equal(t, "1.0.x", branchNode.Content[1].Value)
	})

	t.Run("does not reorder when branch is already first", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)

		config := `
releaseType: java-yoshi
branches:
  - branch: 1.0.x
    releaseType: java-backport
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		formatReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "branch: 1.0.x")
	})

	t.Run("handles branch with no other keys", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)

		config := `
branches:
  - branch: 1.0.x
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		formatReleasePlease(repoDir)
		// No assertion, just checking for no panic
	})
}
