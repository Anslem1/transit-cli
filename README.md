<div align="center">

# 🚀 Transit CLI

**The modern, zero-config workflow orchestrator and interactive command runner for developers who hate remembering multi-step terminal commands.**

[![Release](https://img.shields.io/github/v/release/Anslem1/transit-cli?style=flat-square&color=6366f1)](https://github.com/Anslem1/transit-cli/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/Anslem1/transit-cli/release.yml?style=flat-square&logo=github-actions)](https://github.com/Anslem1/transit-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Anslem1/transit?style=flat-square)](https://goreportcard.com/report/github.com/Anslem1/transit)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Homebrew](https://img.shields.io/badge/brew-transit--cli-orange?style=flat-square&logo=homebrew)](https://github.com/Anslem1/homebrew-tap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Anslem1/transit-cli?style=flat-square&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat-square)](#-installation)

<br/>

<p align="center">
  <img src="demo.gif" alt="Transit Interactive Demo" width="780">
</p>

<p align="center">
  <b>Group commands. Run sequentially or concurrently. Share workflows with your team. Never write messy aliases again.</b>
</p>

</div>

---

## 📑 Table of Contents

- [💡 Why Transit?](#-why-transit)
- [⚖️ How Transit Compares](#️-how-transit-compares)
- [✨ Key Features](#-key-features)
- [📦 Installation](#-installation)
  - [Homebrew (macOS & Linux)](#homebrew-macos--linux)
  - [Go Install](#go-install)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Build from Source](#build-from-source)
- [⚡ Quickstart in 60 Seconds](#-quickstart-in-60-seconds)
- [📖 Command Reference](#-command-reference)
  - [Overview Table](#overview-table)
  - [`transit execute`](#transit-execute)
  - [`transit create`](#transit-create)
  - [`transit add`](#transit-add)
  - [`transit list`](#transit-list)
  - [`transit edit`](#transit-edit)
  - [`transit reorder`](#transit-reorder)
  - [`transit remove`](#transit-remove)
  - [`transit delete`](#transit-delete)
  - [`transit search`](#transit-search)
  - [`transit version`](#transit-version)
  - [`transit completion`](#transit-completion)
- [🧩 Dynamic Parameter Interpolation](#-dynamic-parameter-interpolation)
- [🤝 Team Workflows with `.transit.yaml`](#-team-workflows-with-transityaml)
- [🎯 Practical Real-World Cookbooks](#-practical-real-world-cookbooks)
  - [1. Full-Stack Dev Environment (Parallel)](#1-full-stack-dev-environment-parallel)
  - [2. Docker Build, Tag & Push](#2-docker-build-tag--push)
  - [3. Git PR Preparation & Linting](#3-git-pr-preparation--linting)
  - [4. Database Migration & Seeding](#4-database-migration--seeding)
  - [5. CI/CD Pipeline Automation](#5-cicd-pipeline-automation)
- [📁 Storage Specifications & Conventions](#-storage-specifications--conventions)
- [⌨️ Shell Autocompletion Setup](#️-shell-autocompletion-setup)
- [❓ Frequently Asked Questions (FAQ)](#-frequently-asked-questions-faq)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## 💡 Why Transit?

Every developer struggles with command sprawl:
- **Terminal tab overload:** Opening 3 terminal windows every morning just to run `npm run dev`, `air` (or `go run`), and `docker compose up`.
- **Forgotten parameters:** Struggling to recall lengthy Docker build arguments, Kubernetes port-forwarding flags, or AWS S3 synchronization strings.
- **Fragile `.zshrc` / `.bashrc` aliases:** Bloating your dotfiles with dozens of ad-hoc functions and aliases that only exist on one machine.
- **Onboarding friction:** Writing multi-page READMEs instructing new teammates to run command after command in exact sequence.

**Transit** replaces ad-hoc shell scripts, bloated aliases, and messy makefiles with a clean, interactive command orchestration layer:
- **Zero DSL or Complex Syntax:** No Make syntax peculiarities or tab vs. space errors. A transit is simply a named list of standard shell commands.
- **Sequential or Concurrent:** Run sequential pipelines with fail-fast safety, or spin up microservices simultaneously with synchronized, color-coded output tags.
- **Interactive TUI:** Can't remember the exact name of a transit or what commands are inside it? Run `transit execute` without arguments to launch an interactive, arrow-key selection menu.
- **Global & Project-Local:** Maintain your personal global productivity library, or check a `.transit.yaml` directly into git to streamline team onboarding.
- **Dynamic Arguments:** Pass runtime arguments (`$1`, `$2`, `"$@"`) on the fly with dry-run preview capabilities.

---

## ⚖️ How Transit Compares

| Feature | `transit` | `Makefile` | `npm scripts` | `just` / `task` | Shell Scripts (`.sh`) | Shell Aliases (`.zshrc`) |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| **Zero New DSL to Learn** | ✅ Yes | ❌ Tab/syntax rules | ❌ JSON escaping | ❌ Custom DSL | ✅ Yes | ✅ Yes |
| **Interactive TUI Menus** | ✅ Native | ❌ No | ❌ No | ❌ No | ❌ Manual | ❌ No |
| **Parallel / Concurrent Output** | ✅ Clean tagged lines | ⚠️ Basic (`-j`) | ⚠️ Requires `concurrently` | ⚠️ Limited | ❌ Messy output | ❌ No |
| **Personal Global Library** | ✅ Built-in XDG | ❌ Per-directory | ❌ Per-directory | ❌ Per-directory | ⚠️ Loose files | ⚠️ Dotfiles bloat |
| **Project-Local Git Tracking** | ✅ `.transit.yaml` | ✅ `Makefile` | ✅ `package.json` | ✅ `Justfile` | ✅ Repo script | ❌ Machine-only |
| **Dynamic Argument Interpolation** | ✅ `$1`, `$2`, `"$@"` | ⚠️ Cryptic syntax | ⚠️ `--` passing | ✅ Supported | ✅ Supported | ⚠️ Shell functions |
| **Zero Dependencies (Single Binary)** | ✅ Go binary | ⚠️ Needs Make | ❌ Needs Node.js | ✅ Single binary | ⚠️ Shell dependent | ⚠️ Shell dependent |
| **Cross-Platform (macOS/Linux/Win)** | ✅ Native | ❌ Windows issues | ⚠️ Node runtime | ✅ Cross-platform | ❌ Windows issues | ❌ Shell specific |

---

## ✨ Key Features

- ⚡ **Concurrent Execution (`--parallel` / `-p`)**: Run multiple dev servers, watchers, or containers simultaneously. All output lines are synchronized and prepended with distinct, color-coded tags (e.g. `[1: npm run dev]`, `[2: air]`). Pressing `Ctrl+C` gracefully terminates the entire process tree.
- 🛡️ **Fail-Fast by Default**: Sequential execution halts immediately if any command exits with a non-zero code, preventing cascading errors (such as deploying broken builds).
- 🎯 **Interactive Terminal UI**: Built-in interactive menus powered by PromptUI. Omit arguments to browse, select, edit, reorder, or delete transits with your arrow keys.
- 📦 **Project-Local Workflows (`.transit.yaml`)**: Drop a `.transit.yaml` into any repository root. Anyone who clones the repo automatically inherits all project commands without modifying their global config.
- 🧩 **Runtime Parameter Interpolation**: Inject dynamic parameters into commands using `$1`, `$2`, `"$@"`, `*`, or `{{1}}`.
- 🧪 **Dry-Run Inspection (`--dry-run`)**: Preview interpolated commands before executing them to ensure correctness.
- 🔍 **Global Full-Text Search**: Run `transit search <query>` to immediately identify which transit contains a specific tool, URL, or command.
- ⌨️ **Dynamic Shell Autocompletion**: Native completion for transit names across Zsh, Bash, Fish, and PowerShell (`transit execute <TAB>`).
- 🌐 **XDG Compliant & Cross-Platform**: Stores configuration securely adhering to the XDG Base Directory specification across macOS, Linux, and Windows.

---

## 📦 Installation

### Homebrew (macOS & Linux)

The easiest and recommended installation method for macOS and Linux users:

```bash
# Tap the official Transit repository
brew tap Anslem1/homebrew-tap

# Install transit-cli
brew install transit-cli
```

To update to the latest release at any time:
```bash
brew update && brew upgrade transit-cli
```

---

### Go Install

If you have Go (1.22 or newer) installed:

```bash
go install github.com/Anslem1/transit@latest
```

Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is included in your system's `$PATH`.

---

### Pre-built Binaries

Download compiled standalone binaries for your target architecture directly from the [GitHub Releases](https://github.com/Anslem1/transit-cli/releases) page:

| OS | Architecture | Binary Archive |
| :--- | :--- | :--- |
| **macOS** | Apple Silicon (`arm64`) | `transit-cli_*_Darwin_arm64.tar.gz` |
| **macOS** | Intel (`x86_64`) | `transit-cli_*_Darwin_x86_64.tar.gz` |
| **Linux** | 64-bit (`x86_64`) | `transit-cli_*_Linux_x86_64.tar.gz` |
| **Linux** | ARM64 (`aarch64`) | `transit-cli_*_Linux_arm64.tar.gz` |
| **Windows** | 64-bit (`x86_64`) | `transit-cli_*_Windows_x86_64.zip` |

#### Quick Install Script (macOS / Linux)
```bash
curl -sSL https://raw.githubusercontent.com/Anslem1/transit-cli/main/install.sh | bash
```

---

### Build from Source

```bash
git clone https://github.com/Anslem1/transit-cli.git
cd transit-cli
go build -o transit main.go
sudo mv transit /usr/local/bin/
```

Verify your installation:
```bash
transit version
```

---

## ⚡ Quickstart in 60 Seconds

### 1. Create a Transit
```bash
# Create a global transit (accessible from any directory)
transit create dev-servers

# Or create a project-local transit (saved in .transit.yaml in the current directory)
transit create dev-servers --local
```

### 2. Add Commands
```bash
# Add commands directly as arguments:
transit add dev-servers "npm run dev" "air" "docker compose up db"

# Or launch the interactive addition prompt:
transit add dev-servers
```

### 3. Run It!
```bash
# Run sequentially (halts if any command fails):
transit execute dev-servers

# Run concurrently in parallel with synchronized color-coded output:
transit execute dev-servers --parallel
```

### 4. Interactive Mode
Forget the name of your transit? Just run:
```bash
transit execute
```
Use the `↑` and `↓` arrow keys to browse your transits, hit `Enter`, and Transit executes it immediately!

---

## 📖 Command Reference

### Overview Table

| Command | Shorthand | Description |
| :--- | :--- | :--- |
| [`transit execute`](#transit-execute) | `transit exec`, `transit run` | Run commands in a transit (sequential or parallel) |
| [`transit create`](#transit-create) | `transit new` | Initialize a new transit (`--local` for repo-scoped) |
| [`transit add`](#transit-add) | — | Append commands to a transit (inline or interactive) |
| [`transit list`](#transit-list) | `transit ls` | List all transits or inspect commands inside one |
| [`transit edit`](#transit-edit) | — | Interactively modify an existing command in a transit |
| [`transit reorder`](#transit-reorder) | — | Interactively reorder execution sequence of commands |
| [`transit remove`](#transit-remove) | `transit rm` | Remove an individual command from a transit |
| [`transit delete`](#transit-delete) | — | Delete an entire transit |
| [`transit search`](#transit-search) | — | Search across all transits for command substrings |
| [`transit version`](#transit-version) | — | Print current Transit CLI version and build metadata |
| [`transit completion`](#transit-completion) | — | Generate shell autocompletion script |

---

### `transit execute`
Executes all commands defined inside a transit.

```bash
transit execute [transit-name] [arguments...] [flags]
```

#### Flags
| Flag | Shorthand | Description |
| :--- | :---: | :--- |
| `--parallel` | `-p` | Run all commands concurrently with tagged, color-coded line logs |
| `--skip` | `-s` | Skip interactive per-step confirmation prompts |
| `--continue-on-error`| `-c` | Continue executing subsequent commands even if an error occurs |
| `--dry-run` | — | Display interpolated commands without executing them |

#### Examples
```bash
# Execute sequentially (default)
transit execute deploy

# Run concurrently in parallel (ideal for microservices and dev watchers)
transit execute dev --parallel

# Pass dynamic runtime arguments to commands
transit execute release v1.4.0 production

# Preview commands without executing
transit execute release v1.4.0 --dry-run

# Interactive picker (omitting transit name)
transit execute
```

---

### `transit create`
Initializes a new empty transit.

```bash
transit create <transit-name> [flags]
```

#### Flags
| Flag | Shorthand | Description |
| :--- | :---: | :--- |
| `--local` | `-l` | Save transit in `./.transit.yaml` in the current working directory |

#### Examples
```bash
# Create a global transit (stored in your user config directory)
transit create docker-clean

# Create a project-local transit for this repository
transit create test-suite --local
```

---

### `transit add`
Appends one or more commands to an existing transit.

```bash
transit add <transit-name> [commands...]
```

#### Examples
```bash
# Append multiple commands directly from the terminal:
transit add build-all "go build -o api ./cmd/api" "npm run build"

# Interactive input mode (prompts until an empty line or Ctrl+C):
transit add build-all
```

---

### `transit list`
Displays existing transits and inspects their commands.

```bash
transit list [transit-name] [flags]
```

#### Flags
| Flag | Shorthand | Description |
| :--- | :---: | :--- |
| `--all` | `-a` | Print a clean summary table of all transits, showing global vs. local scopes |

#### Examples
```bash
# Summary table of all available transits:
transit list --all

# View the commands configured inside a specific transit:
transit list dev-servers

# Interactive selector to pick a transit to inspect:
transit list
```

---

### `transit edit`
Interactively selects a command from a transit to modify inline. Great for fixing typos without manually opening config files.

```bash
transit edit <transit-name>
```

---

### `transit reorder`
Changes the sequence of commands within a transit using an interactive step-by-step picker.

```bash
transit reorder <transit-name>
```

---

### `transit remove`
Removes a specific command from a transit.

```bash
transit remove <transit-name>
```

---

### `transit delete`
Permanently deletes one or more transits.

```bash
# Delete a specific transit by name:
transit delete <transit-name>

# Interactive menu to select which transit to delete:
transit delete
```

---

### `transit search`
Searches across all global and local transits for commands containing a specific query or keyword.

```bash
transit search <keyword>
```

#### Example
```bash
$ transit search "docker"
Found 2 matching command(s):
  • [global] 'docker-cleanup' (command 1): docker system prune -af --volumes
  • [local]  'dev' (command 3): docker compose up -d postgres redis
```

---

### `transit version`
Prints the installed version, build commit, and build date.

```bash
$ transit version
transit-cli version 1.1.0 (built with go1.24)
```

---

### `transit completion`
Generates shell autocompletion scripts for dynamic `<TAB>` completion.

```bash
transit completion <bash|zsh|fish|powershell>
```

---

## 🧩 Dynamic Parameter Interpolation

Transit supports dynamic runtime parameter substitution. You can write parameterized commands using standard bash-style variables or mustache tokens:

| Syntax | Description | Example |
| :--- | :--- | :--- |
| `$1`, `$2`, `$3` | Positional arguments | `docker build -t myapp:$1 .` |
| `"$@"` or `$*` | All supplied arguments | `git commit -m "$@"` |
| `{{1}}`, `{{2}}` | Mustache positional syntax | `echo "Target: {{1}}"` |

### How It Works

1. Create a transit with placeholders:
   ```bash
   transit create ship
   transit add ship \
     "docker build -t registry.mycompany.com/service:$1 ." \
     "docker push registry.mycompany.com/service:$1" \
     "kubectl set image deployment/service service=registry.mycompany.com/service:$1 -n $2"
   ```

2. Execute with arguments:
   ```bash
   transit execute ship v2.5.1 production
   ```

3. **Validate with `--dry-run` first**:
   ```bash
   $ transit execute ship v2.5.1 production --dry-run
   [dry-run] Would execute: docker build -t registry.mycompany.com/service:v2.5.1 .
   [dry-run] Would execute: docker push registry.mycompany.com/service:v2.5.1
   [dry-run] Would execute: kubectl set image deployment/service service=registry.mycompany.com/service:v2.5.1 -n production
   ```

---

## 🤝 Team Workflows with `.transit.yaml`

Eliminate project onboarding friction. Instead of maintaining disparate `.sh` scripts or complex setup READMEs, check a `.transit.yaml` file into your repository root:

```yaml
# .transit.yaml
transits:
  # Quick onboarding for new engineers
  setup:
    - cp .env.example .env
    - npm install
    - npx prisma db push
    - npx prisma db seed

  # Concurrent local development environment
  dev:
    - npm run dev
    - air
    - docker compose up postgres redis

  # Full CI check before pushing
  test:
    - npm run lint
    - go test -v ./...
    - npm test

  # Database reset routine
  db-reset:
    - docker compose down -v
    - docker compose up -d postgres
    - sleep 2
    - npx prisma migrate reset --force
```

### Team Benefits
- **Zero Configuration for Teammates:** Any team member with `transit` installed automatically sees these transits when running `transit list --all` or `transit execute` from within the repository.
- **Scope Clarity:** Local transits are marked with `[local (.transit.yaml)]`, so there is no ambiguity between your personal commands and project commands.
- **Local Precedence:** If a local transit shares the same name as a global transit, running from that directory cleanly uses the project version.

---

## 🎯 Practical Real-World Cookbooks

### 1. Full-Stack Dev Environment (Parallel)
Start your web frontend, API server, and containerized database simultaneously in a single terminal session:

```bash
transit create dev --local
transit add dev "npm run dev" "air" "docker compose up db"
transit execute dev --parallel
```

**Output:**
```text
[1: npm run dev] > frontend@1.0.0 dev
[1: npm run dev] > ready on http://localhost:3000
[2: air]         > [air] listening on :8080...
[3: docker...]   > PostgreSQL 16 ready for connections
```
*Pressing `Ctrl+C` terminates all three processes cleanly.*

---

### 2. Docker Build, Tag & Push
Package and publish container images across registries with versioning:

```bash
transit create docker-publish
transit add docker-publish \
  "docker build -t org/api:$1 ." \
  "docker tag org/api:$1 org/api:latest" \
  "docker push org/api:$1" \
  "docker push org/api:latest"

# Execute:
transit execute docker-publish v1.2.3
```

---

### 3. Git PR Preparation & Linting
Run pre-flight checks and clean up branches before opening a pull request:

```bash
transit create pr-prep
transit add pr-prep \
  "git fetch origin" \
  "git rebase origin/main" \
  "npm run lint" \
  "npm run typecheck" \
  "npm test"

# Execute with fail-fast safety:
transit execute pr-prep
```
*If linting fails, execution halts immediately—your remote branch stays safe.*

---

### 4. Database Migration & Seeding
Reset local development databases reliably:

```bash
transit create db-refresh --local
transit add db-refresh \
  "docker compose stop postgres" \
  "docker compose rm -f postgres" \
  "docker compose up -d postgres" \
  "sleep 3" \
  "npx prisma migrate deploy" \
  "npx prisma db seed"

transit execute db-refresh
```

---

### 5. CI/CD Pipeline Automation
Use Transit directly in GitHub Actions or GitLab CI to ensure local workflows match CI workflows 1:1:

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      # Install Transit CLI
      - name: Install Transit
        run: |
          curl -sSL https://raw.githubusercontent.com/Anslem1/transit-cli/main/install.sh | bash

      # Run standard project transit
      - name: Run Test Suite
        run: transit execute test
```

---

## 📁 Storage Specifications & Conventions

Transit respects the official **XDG Base Directory Specification** on Linux and standard system conventions on macOS and Windows:

```
macOS:
  ~/Library/Application Support/transit/cmds/
  ├── dev-servers.yaml
  └── docker-cleanup.yaml
  (Note: Legacy path ~/Library/transit/cmds/ is automatically supported as fallback)

Linux:
  $XDG_CONFIG_HOME/transit/cmds/  (Default: ~/.config/transit/cmds/)
  ├── dev-servers.yaml
  └── docker-cleanup.yaml

Windows:
  %APPDATA%\transit\cmds\
  ├── dev-servers.yaml
  └── docker-cleanup.yaml

Project-Local:
  ./.transit.yaml  (In repository root)
```

- **File Permissions:** All config files are written with secure `0644` permissions, and directories with `0755`.
- **Dotfiles Friendly:** To sync global transits across your machines, simply symlink `~/.config/transit/cmds` into your dotfiles git repository!

---

## ⌨️ Shell Autocompletion Setup

Transit generates native autocompletion scripts for dynamic transit name completion.

### Zsh
Add the following line to your `~/.zshrc`:
```zsh
eval "$(transit completion zsh)"
```

### Bash
Ensure `bash-completion` is installed, then add to `~/.bashrc`:
```bash
eval "$(transit completion bash)"
```

### Fish
Add the following to `~/.config/fish/config.fish`:
```fish
transit completion fish | source
```

### PowerShell
Add the following to your PowerShell `$PROFILE`:
```powershell
transit completion powershell | Out-String | Invoke-Expression
```

---

## ❓ Frequently Asked Questions (FAQ)

<details>
<summary><b>Can I use shell pipes (<code>|</code>), redirects (<code>&gt;</code>), and environment variables in commands?</b></summary>
<br/>
<b>Yes!</b> Transit executes all commands using your host shell (<code>/bin/sh -c</code> on Unix/macOS and <code>cmd.exe /c</code> on Windows). This means pipes, boolean operators (<code>&&</code>, <code>||</code>), file redirects, and subshells work natively.
</details>

<details>
<summary><b>How does Transit handle <code>Ctrl+C</code> in parallel execution mode?</b></summary>
<br/>
Transit listens for interrupt signals (<code>SIGINT</code> / <code>SIGTERM</code>) and immediately sends termination signals down the entire process tree to all concurrently running commands, ensuring no orphan background processes continue running.
</details>

<details>
<summary><b>What happens if a project-local transit and a global transit have the same name?</b></summary>
<br/>
When running inside a directory containing a <code>.transit.yaml</code>, the project-local transit takes precedence. Running <code>transit list --all</code> displays both scopes clearly.
</details>

<details>
<summary><b>Are my existing transits from earlier versions migrated?</b></summary>
<br/>
Yes! Transit features transparent backward-compatibility. If you have transits in the legacy macOS directory <code>~/Library/transit/cmds/</code>, Transit detects and reads them automatically.
</details>


---

## 🤝 Contributing

Contributions, bug reports, and feature suggestions are welcome and appreciated!

1. Fork the repository on GitHub.
2. Create a feature branch:
   ```bash
   git checkout -b feat/my-new-feature
   ```
3. Ensure the test suite passes:
   ```bash
   go test -v ./...
   ```
4. Commit your changes with clear messages:
   ```bash
   git commit -m "feat: add support for custom shell environments"
   ```
5. Push to your branch and open a Pull Request.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE) © Anslem1.