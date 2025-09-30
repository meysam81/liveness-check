# Checklist: Merging All Open PRs

Use this checklist to track progress merging all open PRs.

## Prerequisites

- [ ] Install GitHub CLI (`gh`)
  - macOS: `brew install gh`
  - Ubuntu/Debian: `sudo apt install gh`
  - Windows: `choco install gh`
  - Or download from: https://cli.github.com/

- [ ] Authenticate with GitHub CLI
  ```bash
  gh auth login
  ```

- [ ] Install Go tools (if not already installed)
  - [ ] Go: https://golang.org/dl/
  - [ ] goimports: `go install golang.org/x/tools/cmd/goimports@latest`
  - [ ] golangci-lint: https://golangci-lint.run/usage/install/

- [ ] Install pre-commit
  ```bash
  pip install pre-commit
  ```

## Quick Merge (Automated)

- [ ] Make scripts executable
  ```bash
  chmod +x scripts/*.sh
  ```

- [ ] Run bulk merge script
  ```bash
  ./scripts/approve-and-merge-all.sh
  ```

- [ ] Verify all PRs merged successfully
  ```bash
  gh pr list
  ```

## Manual Merge (If Needed)

If the automated script doesn't work or you prefer manual control:

### PR #10 - Pre-commit autoupdate
- [ ] Fix: `./scripts/fix-pr-precommit.sh 10`
- [ ] Approve: `gh pr review 10 --approve`
- [ ] Merge: `gh pr merge 10 --squash`

### PR #11 - K8s packages update (Auto-merge enabled ✨)
- [ ] Fix: `./scripts/fix-pr-precommit.sh 11`
- [ ] Approve: `gh pr review 11 --approve`
- [ ] Wait for auto-merge (will merge automatically)

### PR #12 - meysam81/x update
- [ ] Fix: `./scripts/fix-pr-precommit.sh 12`
- [ ] Approve: `gh pr review 12 --approve`
- [ ] Merge: `gh pr merge 12 --squash`

### PR #13 - urfave/cli update (Auto-merge enabled ✨)
- [ ] Fix: `./scripts/fix-pr-precommit.sh 13`
- [ ] Approve: `gh pr review 13 --approve`
- [ ] Wait for auto-merge (will merge automatically)

### PR #14 - actions/checkout update
- [ ] Fix: `./scripts/fix-pr-precommit.sh 14`
- [ ] Approve: `gh pr review 14 --approve`
- [ ] Merge: `gh pr merge 14 --squash`

### PR #15 - go dependency update
- [ ] Fix: `./scripts/fix-pr-precommit.sh 15`
- [ ] Approve: `gh pr review 15 --approve`
- [ ] Merge: `gh pr merge 15 --squash`

### PR #16 - golang Docker tag update
- [ ] Fix: `./scripts/fix-pr-precommit.sh 16`
- [ ] Approve: `gh pr review 16 --approve`
- [ ] Merge: `gh pr merge 16 --squash`

### PR #17 - actions/setup-go update
- [ ] Fix: `./scripts/fix-pr-precommit.sh 17`
- [ ] Approve: `gh pr review 17 --approve`
- [ ] Merge: `gh pr merge 17 --squash`

## Verification

- [ ] Verify no open PRs remain
  ```bash
  gh pr list
  ```

- [ ] Check CI is passing on main branch
  ```bash
  gh run list --branch main --limit 1
  ```

- [ ] Clean up local branches (optional)
  ```bash
  git fetch --prune
  git branch -r | grep renovate | xargs -I {} git push origin --delete {}
  ```

## Final Steps

- [ ] Merge or close PR #18 (this PR with the tools)
- [ ] Consider improving Renovate config to avoid future pre-commit failures
- [ ] Consider setting up auto-merge for all Renovate PRs

## Troubleshooting

If you encounter issues:

1. **Check script help**: Each script shows usage when run without arguments
2. **View detailed logs**: Scripts output detailed progress information
3. **Check PR status**: `gh pr view <pr-number>`
4. **Check CI logs**: `gh run view` or visit GitHub Actions page
5. **Manual fallback**: You can always use GitHub UI to approve and merge

## Notes

- ✨ PRs #11 and #13 have auto-merge enabled - they'll merge automatically when checks pass
- All scripts are safe to run multiple times
- Scripts will skip already-merged PRs automatically
- Failed scripts won't stop processing of other PRs

---

**Quick start:** Just run `./scripts/approve-and-merge-all.sh` after installing prerequisites!
