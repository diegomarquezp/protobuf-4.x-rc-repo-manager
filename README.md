# Protobuf 4.x RC Repo Manager

This command-line tool is designed to manage a collection of Google Cloud Java repositories for the `protobuf-4.x-rc` release candidate branch. It automates cloning, configuration updates, and pushing changes back to GitHub.

## Prerequisites

Before using this tool, please ensure you have the following installed:

*   [Go](https://golang.org/doc/install) (version 1.18 or higher)
*   [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
*   [GitHub CLI (`gh`)](https://cli.github.com/)

You will also need a GitHub Personal Access Token with `repo` scopes. This token should be stored in a file at `~/GITHUB_TOKEN`.

## Installation

You can build the tool from the source:

```bash
go build -o repo-manager main.go
```

## Usage

The `repo-manager` provides several commands to streamline repository management.

```
A CLI tool to manage a list of git repositories

Usage:
  repo-manager [command]

Available Commands:
  add-submodules          Add repositories as submodules
  check-branch          Check the current branch of each repository
  clone                 Clone repositories from a file
  completion            Generate the autocompletion script for the specified shell
  help                  Help about any command
  push                  Commit and push changes for each repository
  update-release-please Update release-please-config.json to set prerelease

Flags:
  -h, --help   help for repo-manager

Use "repo-manager [command] --help" for more information about a command.
```

## Suggested Workflow

Based on our setup, here is the recommended sequence of commands to prepare the repositories for the release candidate.

```mermaid
graph TD
    A[Start] --> B("1. Clone Repositories\n`./repo-manager clone`");
    B --> C("2. Check Branches\n`./repo-manager check-branch`");
    C --> D("3. Update Release Please Config\n`./repo-manager update-release-please`");
    D --> E("4. Add Repositories as Submodules\n`./repo-manager add-submodules`");
    E --> F("5. Commit & Push Changes\n- git add .\n- git commit -m \"feat: configure repos for rc release\"\n- git push");
    F --> G[End];

    style B fill:#d4edda,stroke:#c3e6cb
    style C fill:#d4edda,stroke:#c3e6cb
    style D fill:#d4edda,stroke:#c3e6cb
    style E fill:#d4edda,stroke:#c3e6cb
    style F fill:#d4edda,stroke:#c3e6cb
```

### Step-by-step Guide

1.  **Clone Repositories**: This is the first step. It reads `github_repositories.txt` and clones each repository into your local workspace, checking out the `protobuf-4.x-rc` branch.
    ```bash
    ./repo-manager clone
    ```

2.  **Check Branches**: A verification step to ensure all repositories are on the correct branch before making changes.
    ```bash
    ./repo-manager check-branch --branch "protobuf-4.x-rc"
    ```

3.  **Update Release Please Config**: This command modifies the `release-please-config.json` in each repository to add the `"prerelease": true` flag, which is necessary for creating release candidates.
    ```bash
    ./repo-manager update-release-please
    ```

4.  **Add Repositories as Submodules**: This step converts the cloned repositories into Git submodules, which is a cleaner way to manage project dependencies.
    ```bash
    ./repo-manager add-submodules
    ```

5.  **Commit and Push**: Finally, commit all the changes (including the new `.gitmodules` file) and push them to the central `protobuf-4.x-rc-repo-manager` repository.
    ```bash
    git add .
    git commit -m "feat: configure repositories for rc release and add as submodules"
    git push origin main
    ```

## Miscellaneous

*   The list of target repositories is managed in the `github_repositories.txt` file. You can modify this file to add or remove repositories from the workflow.
*   The `push` command in the tool is designed for pushing changes within the submodules themselves, which may be useful for other automation tasks. For the primary workflow, a manual `git push` from the parent repository is recommended after adding the submodules.
