#!/bin/bash
# Script to fix pre-commit issues for a specific PR
# Usage: ./fix-pr-precommit.sh <pr_number>

set -e

PR_NUMBER=$1
if [ -z "$PR_NUMBER" ]; then
    echo "❌ Error: PR number required"
    echo "Usage: $0 <pr_number>"
    echo "Example: $0 17"
    exit 1
fi

echo "🔧 Fixing pre-commit issues for PR #$PR_NUMBER"
echo "=============================================="

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "⚠️  Warning: GitHub CLI (gh) not found. Will try to determine branch manually."
    
    # List of known PR branches
    declare -A PR_BRANCHES
    PR_BRANCHES[17]="renovate/actions-setup-go-6.x"
    PR_BRANCHES[16]="renovate/golang-1.x"
    PR_BRANCHES[15]="renovate/go-1.x"
    PR_BRANCHES[14]="renovate/actions-checkout-5.x"
    PR_BRANCHES[13]="renovate/github.com-urfave-cli-v3-3.x"
    PR_BRANCHES[12]="renovate/github.com-meysam81-x-1.x"
    PR_BRANCHES[11]="renovate/kubernetes-go"
    PR_BRANCHES[10]="pre-commit-ci-update-config"
    
    BRANCH="${PR_BRANCHES[$PR_NUMBER]}"
    
    if [ -z "$BRANCH" ]; then
        echo "❌ Error: Unknown PR number or branch mapping not found"
        echo "Please install gh CLI or update the PR_BRANCHES mapping in this script"
        exit 1
    fi
else
    # Get PR branch name using gh CLI
    BRANCH=$(gh pr view $PR_NUMBER --json headRefName -q .headRefName 2>/dev/null || true)
    
    if [ -z "$BRANCH" ]; then
        echo "❌ Error: Could not find PR #$PR_NUMBER"
        exit 1
    fi
fi

echo "📍 PR #$PR_NUMBER branch: $BRANCH"

# Fetch latest changes
echo "📥 Fetching latest changes..."
git fetch origin

# Checkout PR branch
echo "🔀 Checking out branch: $BRANCH"
git checkout $BRANCH

# Ensure we have the latest
git pull origin $BRANCH || true

echo ""
echo "🔨 Running fixes..."
echo "-------------------"

# Run go mod tidy
if [ -f "go.mod" ]; then
    echo "▶️  Running: go mod tidy"
    go mod tidy
fi

# Run go fmt
if [ -f "go.mod" ]; then
    echo "▶️  Running: go fmt ./..."
    go fmt ./...
fi

# Run goimports if available
if command -v goimports &> /dev/null; then
    echo "▶️  Running: goimports -w ."
    goimports -w .
fi

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
    echo "▶️  Running: golangci-lint run --fix"
    golangci-lint run --fix || true
fi

# Run pre-commit if available
if command -v pre-commit &> /dev/null; then
    echo "▶️  Running: pre-commit run --all-files"
    pre-commit run --all-files || true
fi

echo ""
echo "📊 Checking for changes..."
echo "-------------------------"

# Check if there are changes
if ! git diff --quiet; then
    echo "✏️  Changes detected, committing..."
    git add .
    git commit -m "fix: apply pre-commit and formatting fixes

- Run go mod tidy
- Apply go fmt
- Apply goimports
- Fix pre-commit hook issues

This commit fixes the pre-commit check failures in PR #$PR_NUMBER"
    
    echo "⬆️  Pushing changes..."
    git push origin $BRANCH
    
    echo ""
    echo "✅ SUCCESS! Fixes have been pushed to PR #$PR_NUMBER"
    echo ""
    echo "🔍 Next steps:"
    echo "   1. Wait for CI checks to complete"
    echo "   2. Review the PR at: https://github.com/meysam81/liveness-check/pull/$PR_NUMBER"
    
    # Check if PR has auto-merge
    if [ $PR_NUMBER -eq 11 ] || [ $PR_NUMBER -eq 13 ]; then
        echo "   3. ✨ This PR has auto-merge enabled - it will merge automatically when checks pass!"
    else
        echo "   3. Approve and merge the PR manually"
    fi
else
    echo "✅ No changes needed - everything is already properly formatted!"
fi

echo ""
echo "Done! 🎉"
