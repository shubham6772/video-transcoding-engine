#!/bin/bash

set -e

BRANCH_TYPE="$1"
BRANCH_NAME="$2"

if [ -z "$BRANCH_TYPE" ] || [ -z "$BRANCH_NAME" ]; then
    echo "Usage:"
    echo "  ./create-branch.sh <feature|bugfix|hotfix|chore|refactor|docs> <branch-name>"
    exit 1
fi

case "$BRANCH_TYPE" in
    feature|bugfix|chore|refactor|docs)
        BASE_BRANCH="develop"
        ;;
    hotfix)
        BASE_BRANCH="main"
        ;;
    *)
        echo "Invalid branch type: $BRANCH_TYPE"
        echo "Allowed: feature, bugfix, hotfix, chore, refactor, docs"
        exit 1
        ;;
esac

if [[ ! "$BRANCH_NAME" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    echo "Branch name must be lower-kebab-case."
    echo "Example: rabbitmq-consumer"
    exit 1
fi

git rev-parse --is-inside-work-tree >/dev/null 2>&1 || {
    echo "Not inside a git repository"
    exit 1
}

FULL_BRANCH="${BRANCH_TYPE}/${BRANCH_NAME}"

# Check if branch already exists locally
if git show-ref --verify --quiet "refs/heads/$FULL_BRANCH"; then
    echo "Branch already exists locally: $FULL_BRANCH"
    exit 1
fi

echo "Switching to $BASE_BRANCH..."
git checkout "$BASE_BRANCH"

echo "Pulling latest $BASE_BRANCH..."
git pull origin "$BASE_BRANCH"

echo "Creating branch: $FULL_BRANCH"
git checkout -b "$FULL_BRANCH"

echo "Pushing branch to GitHub..."
git push -u origin "$FULL_BRANCH"

echo ""
echo "Successfully created, switched, and pushed:"
echo "  $FULL_BRANCH"
echo ""
echo "Base branch: $BASE_BRANCH"
echo "Remote: origin/$FULL_BRANCH"