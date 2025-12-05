package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func setupGitRepo(t *testing.T, dir string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	err := cmd.Run()
	assert.NoError(t, err)

	// Need to configure user for commit
	configCmd := exec.Command("git", "config", "user.email", "test@example.com")
	configCmd.Dir = dir
	err = configCmd.Run()
	assert.NoError(t, err)

	configCmd = exec.Command("git", "config", "user.name", "Test User")
	configCmd.Dir = dir
	err = configCmd.Run()
	assert.NoError(t, err)
}

func createCommit(t *testing.T, dir string, message string) {
	// Create a file to commit
	err := os.WriteFile(filepath.Join(dir, "test.txt"), []byte(message), 0644)
	assert.NoError(t, err)

	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = dir
	err = addCmd.Run()
	assert.NoError(t, err)

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = dir
	err = commitCmd.Run()
	assert.NoError(t, err)
}

func createTag(t *testing.T, dir string, tagName string) {
	createCommit(t, dir, "commit for "+tagName)
	tagCmd := exec.Command("git", "tag", tagName)
	tagCmd.Dir = dir
	err := tagCmd.Run()
	assert.NoError(t, err)
}

func TestCleanupReleasePlease(t *testing.T) {
	t.Run("removes redundant options", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)

		config := `
release-type: simple
bump-minor-pre-major: true
branches:
  - branch: main
    release-type: simple
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)

		var root yaml.Node
		err = yaml.Unmarshal(data, &root)
		assert.NoError(t, err)

		branchesNode := root.Content[0].Content[5] // branches
		branchNode := branchesNode.Content[0]
		// branch: main should be the only thing left
		assert.Equal(t, 2, len(branchNode.Content))
		assert.Equal(t, "branch", branchNode.Content[0].Value)
		assert.Equal(t, "main", branchNode.Content[1].Value)
	})

	t.Run("removes bump-minor-pre-major for major release", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)
		createTag(t, repoDir, "v1.0.0")

		config := `
release-type: simple
bump-minor-pre-major: true
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.NotContains(t, string(data), "bump-minor-pre-major")
	})

	t.Run("preserves bump-minor-pre-major for pre-major release", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)
		createTag(t, repoDir, "v0.1.0")

		config := `
release-type: simple
bump-minor-pre-major: true
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "bump-minor-pre-major")
	})

	t.Run("handles no config file", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)

		// No config file created
		cleanupReleasePlease(repoDir)
		// No assertion, just checking for no panic
	})

	t.Run("handles no tags", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)

		config := `
release-type: simple
bump-minor-pre-major: true
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "bump-minor-pre-major")
	})

	t.Run("handles .yaml extension", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)
		createTag(t, repoDir, "v1.0.0")

		config := `
release-type: simple
bump-minor-pre-major: true
`
		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yaml")
		err = os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)

		data, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.NotContains(t, string(data), "bump-minor-pre-major")
	})

	t.Run("handles empty config file", func(t *testing.T) {
		repoDir, err := os.MkdirTemp("", "repo")
		assert.NoError(t, err)
		defer os.RemoveAll(repoDir)
		setupGitRepo(t, repoDir)

		githubDir := filepath.Join(repoDir, ".github")
		os.Mkdir(githubDir, 0755)
		configPath := filepath.Join(githubDir, "release-please.yml")
		err = os.WriteFile(configPath, []byte(""), 0644)
		assert.NoError(t, err)

		cleanupReleasePlease(repoDir)
		// No assertion, just checking for no panic
	})
}
