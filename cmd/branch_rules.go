
package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

// GetBranchProtection gets the branch protection for a given repository and branch.
func GetBranchProtection(owner, repo, branch string) (*github.Protection, error) {
	token, err := readToken("~/GITHUB_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	protection, resp, err := client.Repositories.GetBranchProtection(context.Background(), owner, repo, branch)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("branch protection not found for %s/%s branch %s", owner, repo, branch)
		}
		return nil, err
	}
	return protection, nil
}

// ApplyBranchProtection applies branch protection rules to a given repository and branch.
func ApplyBranchProtection(owner, repo, branch string, protection *github.Protection) error {
	token, err := readToken("~/GITHUB_TOKEN")
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	protectionRequest := &github.ProtectionRequest{}
	if protection.RequiredStatusChecks != nil {
		protectionRequest.RequiredStatusChecks = &github.RequiredStatusChecks{
			Strict:   protection.RequiredStatusChecks.Strict,
			Contexts: protection.RequiredStatusChecks.Contexts,
		}
	}
	if protection.RequiredPullRequestReviews != nil {
		protectionRequest.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          protection.RequiredPullRequestReviews.DismissStaleReviews,
			RequireCodeOwnerReviews:      protection.RequiredPullRequestReviews.RequireCodeOwnerReviews,
			RequiredApprovingReviewCount: protection.RequiredPullRequestReviews.RequiredApprovingReviewCount,
		}
	}
	if protection.EnforceAdmins != nil {
		protectionRequest.EnforceAdmins = protection.EnforceAdmins.Enabled
	}

	_, _, err = client.Repositories.UpdateBranchProtection(context.Background(), owner, repo, branch, protectionRequest)
	return err
}
