#!/bin/bash
set -e

# Pre-commit checks
echo "==> Running pre-commit checks..."
golangci-lint run
go test ./...

# Stage and show changes
echo "==> Staging changes..."
git add -A
git status --short
echo ""
echo "==> Staged diff:"
git diff --staged

# Get commit message from arguments
if [ -z "$1" ]; then
    echo ""
    echo "ERROR: Commit subject required as first argument"
    exit 1
fi

SUBJECT="$1"
DESCRIPTION="${2:-}"

# Build commit message
if [ -n "$DESCRIPTION" ]; then
    COMMIT_MSG="$SUBJECT

$DESCRIPTION

Co-Authored-By: Claude <noreply@anthropic.com>"
else
    COMMIT_MSG="$SUBJECT

Co-Authored-By: Claude <noreply@anthropic.com>"
fi

# Create commit
echo ""
echo "==> Creating commit..."
git commit -m "$COMMIT_MSG"
git status

# Version tagging
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$CURRENT_TAG" ]; then
    echo ""
    echo "==> Current version: $CURRENT_TAG"
    echo "==> Checking if version bump is needed..."

    # If NEW_VERSION is set by caller, create tag and release
    if [ -n "$NEW_VERSION" ]; then
        echo "==> Creating tag $NEW_VERSION..."
        git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

        echo "==> Pushing tag and creating GitHub release..."
        git push origin "$NEW_VERSION"
        gh release create "$NEW_VERSION" --title "$NEW_VERSION" --generate-notes
        echo "==> Created GitHub release $NEW_VERSION"
    else
        echo "==> No NEW_VERSION set, skipping tag"
    fi
else
    echo "==> No existing tags, skipping version bump"
fi

# Sync to ClickUp
echo ""
echo "==> Syncing beans to ClickUp..."
beanup || echo "Warning: beanup failed or not available"

echo ""
echo "==> Done!"
