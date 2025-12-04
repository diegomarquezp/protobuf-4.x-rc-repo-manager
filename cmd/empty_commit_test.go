package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEmptyCommitCmd(t *testing.T) {
	testCmd := &cobra.Command{Use: "test"}
	testCmd.Flags().String("repo", "", "")
	testCmd.Flags().Bool("all", false, "")
	testCmd.Flags().String("branch", "protobuf-4.x-rc", "")
	testCmd.Flags().String("message", "chore: empty commit", "")

	// Test with --repo
	err := testCmd.Flags().Set("repo", "my-repo")
	assert.NoError(t, err)
	assert.Equal(t, "my-repo", testCmd.Flag("repo").Value.String())

	// Test with --all
	err = testCmd.Flags().Set("all", "true")
	assert.NoError(t, err)
	all, err := testCmd.Flags().GetBool("all")
	assert.NoError(t, err)
	assert.True(t, all)

	// Test with --message
	err = testCmd.Flags().Set("message", "my-message")
	assert.NoError(t, err)
	assert.Equal(t, "my-message", testCmd.Flag("message").Value.String())
}
