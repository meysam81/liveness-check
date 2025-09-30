# PR Management Scripts

This directory contains helper scripts for managing pull requests in the repository.

## Available Scripts

### 1. `fix-pr-precommit.sh` - Fix Pre-commit Issues for a Single PR

Fixes pre-commit check failures for a specific PR by running:
- `go mod tidy`
- `go fmt`
- `goimports`
- `golangci-lint --fix`
- `pre-commit run --all-files`

**Usage:**
```bash
./scripts/fix-pr-precommit.sh <pr_number>

# Example:
./scripts/fix-pr-precommit.sh 17
```

**Requirements:**
- Go toolchain
- `goimports` (optional): `go install golang.org/x/tools/cmd/goimports@latest`
- `golangci-lint` (optional): See https://golangci-lint.run/usage/install/
- `pre-commit` (optional): `pip install pre-commit`
- GitHub CLI `gh` (optional, for auto-detecting branch names): https://cli.github.com/

### 2. `approve-and-merge-all.sh` - Bulk Approve and Merge PRs

Processes all open PRs by:
1. Fixing pre-commit issues on each PR
2. Approving each PR
3. Merging PRs (or waiting for auto-merge)

**Usage:**
```bash
./scripts/approve-and-merge-all.sh
```

**Requirements:**
- All requirements from `fix-pr-precommit.sh`
- GitHub CLI `gh` (required) - must be installed and authenticated
- Repository write permissions (for approval and merge)

**Authentication:**
```bash
# Install gh CLI first, then:
gh auth login
```

## Quick Start

### Option 1: Fix All PRs (Recommended)

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run bulk script (requires gh CLI)
./scripts/approve-and-merge-all.sh
```

### Option 2: Fix Individual PRs

```bash
# Fix a specific PR
./scripts/fix-pr-precommit.sh 17

# Then approve and merge manually via GitHub UI or:
gh pr review 17 --approve
gh pr merge 17 --squash
```

### Option 3: Manual Process Without Scripts

```bash
# For each PR, checkout its branch
git fetch origin
git checkout <pr-branch-name>

# Fix issues
go mod tidy
go fmt ./...
pre-commit run --all-files || true

# Commit and push
git add .
git commit -m "fix: pre-commit fixes"
git push origin <pr-branch-name>

# Then approve and merge via GitHub
```

## Current Open PRs

| PR # | Title | Branch | Auto-Merge |
|------|-------|--------|------------|
| #17 | Update actions/setup-go v5→v6 | `renovate/actions-setup-go-6.x` | ❌ |
| #16 | Update golang 1.24→1.25 | `renovate/golang-1.x` | ❌ |
| #15 | Update go dependency 1.24→1.25 | `renovate/go-1.x` | ❌ |
| #14 | Update actions/checkout v4→v5 | `renovate/actions-checkout-5.x` | ❌ |
| #13 | Update urfave/cli v3.3.8→v3.4.1 | `renovate/github.com-urfave-cli-v3-3.x` | ✅ |
| #12 | Update meysam81/x v1.8.2→v1.13.0 | `renovate/github.com-meysam81-x-1.x` | ❌ |
| #11 | Update k8s packages v0.33.1→v0.34.1 | `renovate/kubernetes-go` | ✅ |
| #10 | Pre-commit autoupdate | `pre-commit-ci-update-config` | ❌ |

## Troubleshooting

### "gh: command not found"

Install GitHub CLI:
- macOS: `brew install gh`
- Ubuntu/Debian: `sudo apt install gh`
- Windows: `choco install gh`
- Or download from: https://cli.github.com/

### "gh: not authenticated"

Run: `gh auth login` and follow the prompts.

### "Permission denied" when running scripts

Make scripts executable: `chmod +x scripts/*.sh`

### Pre-commit checks still failing after fixes

Some checks may require specific tools:
- Install `goimports`: `go install golang.org/x/tools/cmd/goimports@latest`
- Install `golangci-lint`: https://golangci-lint.run/usage/install/
- Install `pre-commit`: `pip install pre-commit`

### Can't approve or merge PRs

You need repository maintainer or write permissions. If you don't have these, ask a repository owner to run the scripts or approve/merge manually.

## Notes

- PRs with auto-merge enabled (#11, #13) will merge automatically once checks pass
- The scripts are safe to run multiple times
- Failed scripts will continue processing other PRs
- All scripts include detailed progress output
