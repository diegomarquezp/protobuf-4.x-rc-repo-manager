package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	owner      string
	repo       string
	branch     string
	outputJSON bool
)

var getBranchRulesCmd = &cobra.Command{
	Use:   "get-branch-rules",
	Short: "Get branch protection rules for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		protection, err := GetBranchProtection(owner, repo, branch)
		if err != nil {
			log.Fatalf("Failed to get branch protection: %v", err)
		}

		if outputJSON {
			jsonOutput, err := json.MarshalIndent(protection, "", "  ")
			if err != nil {
				log.Fatalf("Failed to marshal to JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
		} else {
			printBranchProtection(protection)
		}
	},
}

func init() {
	rootCmd.AddCommand(getBranchRulesCmd)
	getBranchRulesCmd.Flags().StringVarP(&owner, "owner", "o", "", "Owner of the repository")
	getBranchRulesCmd.Flags().StringVarP(&repo, "repo", "r", "", "Name of the repository")
	getBranchRulesCmd.Flags().StringVarP(&branch, "branch", "b", "protobuf-4.x-rc", "Branch to get protection rules for")
	getBranchRulesCmd.Flags().BoolVarP(&outputJSON, "output-json", "j", false, "Output in JSON format")
	getBranchRulesCmd.MarkFlagRequired("owner")
	getBranchRulesCmd.MarkFlagRequired("repo")
}

func printBranchProtection(protection *github.Protection) {
	fmt.Println("Branch Protection Rules:")
	if protection.RequiredPullRequestReviews != nil {
		fmt.Printf("  Required Approving Review Count: %d\n", protection.RequiredPullRequestReviews.RequiredApprovingReviewCount)
		fmt.Printf("  Dismiss Stale Reviews: %t\n", protection.RequiredPullRequestReviews.DismissStaleReviews)
		fmt.Printf("  Require Code Owner Reviews: %t\n", protection.RequiredPullRequestReviews.RequireCodeOwnerReviews)
	} else {
		fmt.Println("  No pull request review enforcement")
	}
}
