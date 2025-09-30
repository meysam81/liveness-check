# Guide to Approve and Merge Open PRs

## Current Situation

There are 8 open PRs in the repository (as of this analysis):

1. **PR #17**: Update actions/setup-go v5→v6
2. **PR #16**: Update golang Docker tag 1.24→1.25
3. **PR #15**: Update go dependency 1.24→1.25
4. **PR #14**: Update actions/checkout v4→v5
5. **PR #13**: Update urfave/cli/v3 to v3.4.1 (**has auto-merge enabled**)
6. **PR #12**: Update meysam81/x to v1.13.0
7. **PR #11**: Update k8s packages to v0.34.1 (**has auto-merge enabled**)
8. **PR #10**: Pre-commit autoupdate

## Issue

All PRs are currently **failing pre-commit checks**. This prevents them from being merged, even those with auto-merge enabled.

## Why the AI Agent Can't Directly Merge

The GitHub Copilot coding agent has the following limitations:
- ❌ Cannot approve PRs (requires repository permissions)
- ❌ Cannot merge PRs (requires repository permissions)  
- ❌ Cannot use `gh` CLI for approval/merging (no GitHub credentials)
- ✅ CAN fix code issues and prepare PRs for merging
- ✅ CAN create helper scripts and documentation

## Solution Options

### Option 1: Auto-Merge (Recommended for PRs #11 and #13)

These PRs already have auto-merge enabled. Once pre-commit checks pass, they will merge automatically.

**Steps:**
1. Checkout the PR branch
2. Run pre-commit fixes
3. Push the fixes
4. PR will auto-merge when checks pass

### Option 2: Manual Approval and Merge

For other PRs, you'll need to manually approve and merge them.

**Steps:**
1. Review the PR changes
2. Approve the PR through GitHub UI
3. Merge the PR (can use rebase, squash, or merge commit based on your preference)

## How to Fix Pre-Commit Failures

The pre-commit failures are likely due to:
- Formatting issues (go-fmt, go-imports)
- Missing `go mod tidy`
- YAML/JSON formatting

### Quick Fix Script

Save this as `fix-pr-precommit.sh`:

```bash
#!/bin/bash
set -e

PR_NUMBER=$1
if [ -z "$PR_NUMBER" ]; then
    echo "Usage: $0 <pr_number>"
    exit 1
fi

# Get PR branch name
BRANCH=$(gh pr view $PR_NUMBER --json headRefName -q .headRefName)

echo "Fixing PR #$PR_NUMBER (branch: $BRANCH)"

# Checkout PR branch
git fetch origin
git checkout $BRANCH

# Run go mod tidy
go mod tidy

# Run go fmt
go fmt ./...

# Run goimports if available
if command -v goimports &> /dev/null; then
    goimports -w .
fi

# Run pre-commit
pre-commit run --all-files || true

# Commit changes if any
if ! git diff --quiet; then
    git add .
    git commit -m "fix: run pre-commit fixes"
    git push origin $BRANCH
    echo "✅ Fixes pushed to PR #$PR_NUMBER"
else
    echo "✅ No fixes needed for PR #$PR_NUMBER"
fi
```

### Manual Fix for Each PR

For each PR, you can:

```bash
# Example for PR #17
git fetch origin
git checkout renovate/actions-setup-go-6.x
go mod tidy
go fmt ./...
pre-commit run --all-files --show-diff-on-failure || true
git add .
git commit -m "fix: pre-commit fixes" || true
git push origin renovate/actions-setup-go-6.x
```

## Recommended Merge Order

1. **PR #10** (Pre-commit autoupdate) - Should be merged first as it updates pre-commit hooks
2. **PR #11** (K8s packages) - Has auto-merge
3. **PR #13** (urfave/cli) - Has auto-merge  
4. **PR #12** (meysam81/x)
5. **PR #15** (go dependency)
6. **PR #16** (golang Docker tag)
7. **PR #14** (actions/checkout)
8. **PR #17** (actions/setup-go)

## Using GitHub CLI to Approve and Merge

If you have `gh` CLI installed and authenticated:

```bash
# Approve a PR
gh pr review <pr_number> --approve

# Merge a PR (after approval and checks pass)
gh pr merge <pr_number> --squash  # or --rebase or --merge

# Or do both in one command
gh pr review <pr_number> --approve && gh pr merge <pr_number> --squash
```

## Bulk Approve and Merge Script

```bash
#!/bin/bash
# approve-and-merge-all.sh

PRS=(17 16 15 14 13 12 11 10)

for pr in "${PRS[@]}"; do
    echo "Processing PR #$pr..."
    
    # Approve
    gh pr review $pr --approve
    
    # Wait a bit for checks
    sleep 5
    
    # Check if checks are passing
    if gh pr checks $pr | grep -q "All checks have passed"; then
        gh pr merge $pr --squash --auto
        echo "✅ PR #$pr merged"
    else
        echo "⏳ PR #$pr has failing checks, skipping merge"
    fi
done
```

## Important Notes

1. **Auto-merge PRs** (11, 13): These will auto-merge once checks pass
2. **Renovate PRs**: Consider enabling auto-merge in renovate.json for smoother updates
3. **Pre-commit hooks**: Ensure all contributors run pre-commit before pushing
4. **CI/CD**: The repository has comprehensive CI including security scans

## Next Steps for Repository Owner

1. Install required tools: `gh`, `pre-commit`, `go`
2. Run the fix script on failing PRs or fix them manually
3. Approve PRs that look good
4. PRs with auto-merge will merge automatically
5. Manually merge remaining PRs

## Contact

If you need help with any specific PR or want to discuss the merging strategy, please add a comment to this PR.
