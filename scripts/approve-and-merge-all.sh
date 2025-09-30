#!/bin/bash
# Script to approve and merge all open PRs
# This requires GitHub CLI (gh) to be installed and authenticated
# Usage: ./approve-and-merge-all.sh

set -e

echo "🚀 Bulk PR Approval and Merge Script"
echo "====================================="
echo ""

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    echo "❌ Error: GitHub CLI (gh) is not installed"
    echo ""
    echo "Please install it from: https://cli.github.com/"
    echo ""
    echo "On macOS: brew install gh"
    echo "On Ubuntu/Debian: sudo apt install gh"
    echo "On Windows: choco install gh"
    exit 1
fi

# Check if gh is authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Error: GitHub CLI is not authenticated"
    echo ""
    echo "Please run: gh auth login"
    exit 1
fi

# List of PRs to process (in reverse order - oldest first)
# This order ensures dependencies are handled correctly
PRS=(10 11 12 13 14 15 16 17)

echo "📋 PRs to process: ${PRS[*]}"
echo ""

# First, run the fix script on all PRs
echo "🔧 Step 1: Fixing pre-commit issues on all PRs"
echo "----------------------------------------------"
for pr in "${PRS[@]}"; do
    echo ""
    echo "Fixing PR #$pr..."
    ./scripts/fix-pr-precommit.sh $pr || echo "⚠️  Warning: Failed to fix PR #$pr, continuing..."
    sleep 2
done

echo ""
echo "✅ All fixes applied!"
echo ""
echo "⏳ Waiting 30 seconds for CI checks to start..."
sleep 30

echo ""
echo "👍 Step 2: Approving and merging PRs"
echo "------------------------------------"

for pr in "${PRS[@]}"; do
    echo ""
    echo "Processing PR #$pr..."
    
    # Get PR status
    PR_STATE=$(gh pr view $pr --json state -q .state)
    
    if [ "$PR_STATE" != "OPEN" ]; then
        echo "⏭️  PR #$pr is $PR_STATE, skipping..."
        continue
    fi
    
    # Approve the PR
    echo "👍 Approving PR #$pr..."
    gh pr review $pr --approve || echo "⚠️  Could not approve (may already be approved)"
    
    # Wait a moment
    sleep 5
    
    # Check if checks are passing
    echo "🔍 Checking CI status..."
    CHECKS_STATUS=$(gh pr checks $pr --json state -q '.[].state' | grep -v "SUCCESS" | wc -l)
    
    if [ "$CHECKS_STATUS" -eq 0 ]; then
        echo "✅ All checks passing!"
        
        # Check if auto-merge is enabled
        AUTO_MERGE=$(gh pr view $pr --json autoMergeRequest -q '.autoMergeRequest')
        
        if [ "$AUTO_MERGE" != "null" ] && [ -n "$AUTO_MERGE" ]; then
            echo "✨ Auto-merge is enabled, PR will merge automatically"
        else
            echo "🔀 Merging PR #$pr..."
            gh pr merge $pr --squash --delete-branch || echo "⚠️  Could not merge PR #$pr"
        fi
        
        echo "✅ PR #$pr processed successfully!"
    else
        echo "⏳ Checks are still running or failing for PR #$pr"
        echo "   View details: gh pr checks $pr"
        echo "   This PR will need manual attention"
    fi
    
    sleep 3
done

echo ""
echo "🎉 Bulk processing complete!"
echo ""
echo "📊 Summary:"
echo "-----------"
gh pr list --state open

echo ""
echo "💡 Tips:"
echo "  - Check failing PRs: gh pr checks <pr_number>"
echo "  - View a PR: gh pr view <pr_number>"
echo "  - Manually merge: gh pr merge <pr_number> --squash"
