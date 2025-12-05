#!/bin/bash

REPO_DIR=$1

if [ -z "$REPO_DIR" ]; then
  echo "Usage: $0 <repository_directory>"
  exit 1
fi

cd "$REPO_DIR" || exit

BRANCH_NAME="chore/cleanup-release-please"
COMMIT_TITLE="chore: cleanup release-please config"
COMMIT_BODY="- Remove redundant options already declared at the top level.\n- Remove bumpMinorPreMajor for repositories after the first major release."
PR_BODY="This PR cleans up the .github/release-please.yml file by removing redundant options and the bump-minor-pre-major setting for major releases."

echo "--- Processing repository: $(basename "$REPO_DIR") ---"

# 1. Create a new branch
echo "Creating branch $BRANCH_NAME..."
git checkout -b "$BRANCH_NAME"

# 2. Add the modified release-please.yml file
echo "Adding .github/release-please.yml to staging..."
git add .github/release-please.yml

# 3. Commit the changes
echo "Committing changes..."
git commit -m "$COMMIT_TITLE" -m "$COMMIT_BODY"

# 4. Push the new branch to the remote
echo "Pushing branch to remote..."
git push -u origin "$BRANCH_NAME"

# 5. Create a pull request
echo "Creating pull request..."
PR_URL=$(gh pr create --title "$COMMIT_TITLE" --body "$PR_BODY")

if [ -z "$PR_URL" ]; then
  echo "Failed to create pull request for $(basename "$REPO_DIR")"
  exit 1
fi

echo "Pull Request URL: $PR_URL"

# 6. Append the PR URL to prs.txt
echo "$PR_URL" >> ../prs.txt

echo "--- Finished processing $(basename "$REPO_DIR") ---"
