package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	sourceOwner      string
	sourceRepo       string
	sourceBranch     string
	destinationOwner string
	destinationRepo  string
	destinationBranch string
)

var applyBranchRulesCmd = &cobra.Command{
	Use:   "apply-branch-rules",
	Short: "Apply branch protection rules from one repository to another",
	Run: func(cmd *cobra.Command, args []string) {
		protection, err := GetBranchProtection(sourceOwner, sourceRepo, sourceBranch)
		if err != nil {
			log.Fatalf("Failed to get branch protection: %v", err)
		}

		err = ApplyBranchProtection(destinationOwner, destinationRepo, destinationBranch, protection)
		if err != nil {
			log.Fatalf("Failed to apply branch protection: %v", err)
		}

		fmt.Println("Successfully applied branch protection rules.")
	},
}

func init() {
	rootCmd.AddCommand(applyBranchRulesCmd)
	applyBranchRulesCmd.Flags().StringVar(&sourceOwner, "source-owner", "", "Owner of the source repository")
	applyBranchRulesCmd.Flags().StringVar(&sourceRepo, "source-repo", "", "Name of the source repository")
	applyBranchRulesCmd.Flags().StringVar(&sourceBranch, "source-branch", "protobuf-4.x-rc", "Branch to get protection rules from")
	applyBranchRulesCmd.Flags().StringVar(&destinationOwner, "destination-owner", "", "Owner of the destination repository")
	applyBranchRulesCmd.Flags().StringVar(&destinationRepo, "destination-repo", "", "Name of the destination repository")
	applyBranchRulesCmd.Flags().StringVar(&destinationBranch, "destination-branch", "protobuf-4.x-rc", "Branch to apply protection rules to")
	applyBranchRulesCmd.MarkFlagRequired("source-owner")
	applyBranchRulesCmd.MarkFlagRequired("source-repo")
	applyBranchRulesCmd.MarkFlagRequired("destination-owner")
	applyBranchRulesCmd.MarkFlagRequired("destination-repo")
}
