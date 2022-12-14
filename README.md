# Developer Scripts for Stacked Diff Workflow

These scripts make it easier to build from the command line and to create and update PR's with Github. They facilitates a [stacked diff workflow](https://kastiglione.github.io/git/2020/09/11/git-stacked-commits.html), where you always commit on `main` branch and have can have multiple streams of work all on `main`.

Note: these scripts do *not* facilitate Stacked *Pull Requests*. Github does some things that add friction to using Stacked PR's, even with support from third party software. For example, after merging one of the PR's in the stack, you may require a re-review of the other PR's in the stack. Instead, it's recommended to organize your PR's, as much as reasonably possible, so that they can be all be rebased against main at the same time. When there are dependencies, wait for dependant PR to be merged before putting up the next one. You may find that often you are still working on the next commit while the other is being reviewed/merged.

## TL;DR

Using a stacked diff workflow like this allows you to work on separate streams of work without changing branches.

## Installation

Clone the repository or download the [latest release](releases), and then:

```bash
# Install Github CLI
brew install gh 
# Setup login for Github CLI
gh auth login 
# Add the /bin directory to your PATH. 
# Replace the directory below to wherever you cloned the repository or unzipped the release.
# For example if using zsh and cloned in your home directory:
echo "export PATH=\$PATH:\$HOME/stacked-diff-workflow/bin" >> ~/.zshrc
source ~/.zshrc
```

## Scripts

### For Stacked Diff Workflow

#### git-checkout

`git-checkout <commit hash or pr number>`

Checkout the feature branch associated with a given PR or commit. For when you want to checkout the feature branch to rebase with origin/main, merge with origin/main, or for any other reason. After modifying the feature branch use `replac-commit` or `replace-head` to sync local `main`.

#### gitlog

`gitlog`

Abbreviated git log that only shows what has changed, useful for copying commit hashes.

#### new-pr

Create a new PR with a cherry-pick of the given commit hash

```
new-pr <commitHash>
Usage of new-pr:
  -draft
    	Whether to create the PR as draft (default true)
```

###### Ticket Number

If you prefix the Jira ticket to the git commit summary then `newpr` will populate the `Ticket` section of the PR description.

For example:
`CONV-9999 Add new feature`

###### Templates

The Pull Request Title, Body (aka Description), and Branch Name are created from [golang templates](https://pkg.go.dev/text/template). The defaults are:

- [branch-name.template](cmd/config/branch-name.template)
- [pr-description.template](cmd/config/pr-description.template)
- [pr-title.template](cmd/config/pr-title.template)

The [possible values](config/templates) for the templates are:

- **TicketNumber** - Jira ticket as parsed from the commit summary
- **Username** -  Name as parsed from git config email
- **CommitBody** - Body of the commit message
- **CommitSummary** - Summary line of the commit message
- **CommitSummaryCleaned** - Summary line of the commit message without spaces or special characters
- **CommitSummaryWithoutTicket** - Summary line of the commit message without the prefix of the ticket number

To change a template, copy the default from [cmd/config/](cmd/config/) into `~/.stacked-diff-workflow/` and modify.

#### replace-commit

`replace-commit <commit hash or pr number>`

Reset the main branch with the squashed contents of the given commits associated branch. Sometimes you might want to switch to a feature branch and make changes to it (rebase, amend). With this script you can then ensure that your `main` branch is up to date.

#### replace-head

`replace-head`

Use during rebase of main branch to use the contents of a feature branch that already fixed the merge conflicts.

#### update-pr

`update-pr <commitHash or pullRequestNumber> [fixup commit (defaults to top commit)] [other fixup commit...]`

Add one or more commits to a PR.

### To Help with Github

#### add-reviewers

`add-reviewers <pullRequestNumber or commitHash>`

The `add-reviewers` command will mark your Draft PR as "Ready for Review" and automatically add reviewers that are specified in the PR_REVIEWERS environment variable.
You can specify more than one reviewer using a comma-delimited string.

```bash
export PR_REVIEWERS=first-user,second-user,third-user
```

Add this to your shell rc file (`~/.zshrc` or `~/.bashrc`) and run `source <rc-file>`

```
Usage of add-reviewers:
  -poll-frequency duration
    	Frequency which to poll checks. For valid formats see https://pkg.go.dev/time#ParseDuration (default 5m0s)
  -reviewers string
    	Comma-separated list of Github usernames to add as reviewers
  -silent
    	Whether to use voice output
  -when-checks-pass
    	Poll until all checks pass before adding reviewers (default true)
  <pullRequestNumber>
```

#### git-merge-pr

`git-merge-pr <pull request number>`

Add the given PR to the merge queue

#### git-prs

`git-prs`

Lists all of your open PRs. Useful for copying PR numbers.

## Example Workflow

### Creating and Updating PRs

Use **new-pr** and **update-pr** to create and update PR's while always staying on `main` branch.

### To Update Main

Once a PR has been merged, just rebase main normally. The local PR commit will be replaced by the one that Github created when squasing and merging.

```bash
git fetch && git rebase origin/main
```

### To Fix a Merge Conflict

If you have a merge conflict on your PR, you can use **replace-head** to keep your local `main` up to date.

```bash
# switch to feature branch that has a merge conflict
git-checkout <commit hash or PR number> 
# rebase or merge
git fetch && git rebase origin/main
# ... and address any merge conflicts
# Update your PR
git push origin/xxx 
git switch main
git rebase origin/main
# hit same merge conflicts, use the replace-merge script to copy the fixes you just made
replace-head
# continue with the rebase
git add . && git rebase --continue
# All done... now both the feature branch and your local main are rebased with main, and the merge conflicts only had to be fixed once
```
