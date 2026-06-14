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
        exit 1
        ;;
esac

if [[ ! "$BRANCH_NAME" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    echo "Branch name must be lower-kebab-case."
    exit 1
fi

git rev-parse --is-inside-work-tree >/dev/null 2>&1 || {
    echo "Not inside a git repository"
    exit 1
}

echo "Switching to $BASE_BRANCH..."
git checkout "$BASE_BRANCH"

echo "Pulling latest $BASE_BRANCH..."
git pull origin "$BASE_BRANCH"

FULL_BRANCH="${BRANCH_TYPE}/${BRANCH_NAME}"

echo "Creating branch: $FULL_BRANCH"
git checkout -b "$FULL_BRANCH"

echo "Created and switched to $FULL_BRANCH"