package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	sourceOwnerRepo string
	sourceBranchAll string
)

var applyToAllCmd = &cobra.Command{
	Use:   "apply-to-all",
	Short: "Apply branch protection rules from one repository to all others in github_repositories.txt",
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := readLines("github_repositories.txt")
		if err != nil {
			log.Fatalf("Failed to read github_repositories.txt: %v", err)
		}

		sourceParts := strings.Split(sourceOwnerRepo, "/")
		if len(sourceParts) != 2 {
			log.Fatalf("Invalid source-repo format. Please use owner/repo.")
		}
		sourceOwner := sourceParts[0]
		sourceRepo := sourceParts[1]

		protection, err := GetBranchProtection(sourceOwner, sourceRepo, sourceBranchAll)
		if err != nil {
			log.Fatalf("Failed to get branch protection from %s/%s: %v", sourceOwner, sourceRepo, err)
		}

		for _, repoStr := range repos {
			repoParts := strings.Split(repoStr, "/")
			if len(repoParts) != 2 {
				log.Printf("Skipping invalid repo format: %s", repoStr)
				continue
			}
			destOwner := repoParts[0]
			destRepo := repoParts[1]

			if destOwner == sourceOwner && destRepo == sourceRepo {
				fmt.Printf("Skipping source repository: %s/%s\n", destOwner, destRepo)
				continue
			}

			fmt.Printf("Applying branch protection to %s/%s...\n", destOwner, destRepo)
			err := ApplyBranchProtection(destOwner, destRepo, sourceBranchAll, protection)
			if err != nil {
				log.Printf("Failed to apply branch protection to %s/%s: %v", destOwner, destRepo, err)
			} else {
				fmt.Printf("Successfully applied branch protection to %s/%s\n", destOwner, destRepo)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyToAllCmd)
	applyToAllCmd.Flags().StringVar(&sourceOwnerRepo, "source-repo", "googleapis/google-auth-library-java", "Source repository in owner/repo format")
	applyToAllCmd.Flags().StringVar(&sourceBranchAll, "source-branch", "protobuf-4.x-rc", "Branch to get protection rules from")
}
