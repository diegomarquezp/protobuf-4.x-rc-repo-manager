#!/bin/bash

set -e

REPO_DIR=$1
COMMIT_MESSAGE=$2

if [ -z "$REPO_DIR" ] || [ -z "$COMMIT_MESSAGE" ]; then
  echo "Usage: $0 <repository_directory> \"<commit_message>\""
  exit 1
fi

if [ ! -d "$REPO_DIR" ]; then
  echo "Error: Directory '$REPO_DIR' not found."
  exit 1
fi

cd "$REPO_DIR"

BRANCH_NAME="chore/cleanup-release-please"

echo "--- Updating PR for repository: $(basename "$PWD") ---"

# 1. Checkout the branch
echo "Checking out branch $BRANCH_NAME..."
git checkout "$BRANCH_NAME"

# 2. Fetch the latest changes from origin
echo "Fetching latest changes from origin..."
git fetch origin

# 3. Merge the latest changes from main or master
echo "Merging latest changes..."
if git show-ref --verify --quiet refs/remotes/origin/main; then
    git merge origin/main
elif git show-ref --verify --quiet refs/remotes/origin/master; then
    git merge origin/master
else
    echo "Could not find main or master branch in $(basename "$PWD"). Skipping merge."
fi

# 4. Stage the changes (if any)
echo "Adding changes to staging..."
git add .

# 5. Commit the changes
echo "Committing changes..."
# Check if there are any changes to commit after the merge
if git diff-index --quiet HEAD --; then
    echo "No changes to commit after merge. Pushing existing commits."
else
    git commit -m "$COMMIT_MESSAGE"
fi


# 6. Push the changes to the remote
echo "Pushing changes to remote..."
git push

echo "--- Successfully updated PR for $(basename "$PWD") ---"