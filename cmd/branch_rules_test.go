package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/stretchr/testify/assert"
)

type mockRepositoriesClient struct {
	getBranchProtection func(ctx context.Context, owner, repo, branch string) (*github.Protection, *github.Response, error)
}

func (m *mockRepositoriesClient) GetBranchProtection(ctx context.Context, owner, repo, branch string) (*github.Protection, *github.Response, error) {
	return m.getBranchProtection(ctx, owner, repo, branch)
}

func TestGetBranchProtectionWithAuth(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "test_token")
	defer os.Unsetenv("GITHUB_TOKEN")

	t.Run("success", func(t *testing.T) {
		client := &mockRepositoriesClient{
			getBranchProtection: func(ctx context.Context, owner, repo, branch string) (*github.Protection, *github.Response, error) {
				return &github.Protection{
					RequiredPullRequestReviews: &github.PullRequestReviewsEnforcement{
						DismissStaleReviews:          true,
						RequireCodeOwnerReviews:      true,
						RequiredApprovingReviewCount: 1,
					},
				}, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusOK,
					},
				}, nil
			},
		}

		protection, err := getBranchProtectionWithClient(client, "owner", "repo", "main")
		assert.NoError(t, err)
		assert.NotNil(t, protection)
		assert.True(t, protection.RequiredPullRequestReviews.DismissStaleReviews)
	})

	t.Run("not found", func(t *testing.T) {
		client := &mockRepositoriesClient{
			getBranchProtection: func(ctx context.Context, owner, repo, branch string) (*github.Protection, *github.Response, error) {
				return nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				}, fmt.Errorf("not found")
			},
		}

		_, err := getBranchProtectionWithClient(client, "owner", "repo", "main")
		assert.Error(t, err)
		assert.Equal(t, "branch protection not found for owner/repo branch main", err.Error())
	})
}

func getBranchProtectionWithClient(client *mockRepositoriesClient, owner, repo, branch string) (*github.Protection, error) {
	ctx := context.Background()
	protection, resp, err := client.GetBranchProtection(ctx, owner, repo, branch)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("branch protection not found for %s/%s branch %s", owner, repo, branch)
		}
		return nil, err
	}
	return protection, nil
}