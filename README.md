# GitLingo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/steamedeo/gitlingo)](https://goreportcard.com/report/github.com/steamedeo/gitlingo)

![GitLingo Language Stats](img/sample.png)

Checkout your GitHub programming language statistics straight in your terminal.

## Installation

### Prerequisites

- Go 1.25 or higher ([Download here](https://golang.org/dl/))

### Option 1: Local Build (Recommended for beginners)

```bash
# Clone the repository
git clone https://github.com/steamedeo/gitlingo.git
cd gitlingo

# Build the binary (this creates the bin folder)
make build
```

**Setup:** Create a `.env` file in the gitlingo directory with your GitHub token:

```bash
cp .env.example .env
# Edit .env and add: github_token=ghp_yourTokenHere
```

**Usage:** Run from the project directory:

```bash
./bin/gitlingo languages
```

### Option 2: Global Install (For advanced users)

```bash
# Clone and install
git clone https://github.com/steamedeo/gitlingo.git
cd gitlingo
make install
```

**Setup:** Set the `GITHUB_TOKEN` environment variable permanently:

**On Linux/macOS:**
```bash
# Add to ~/.bashrc, ~/.zshrc, or ~/.profile
echo 'export GITHUB_TOKEN=ghp_yourTokenHere' >> ~/.bashrc
source ~/.bashrc
```

**On Windows (PowerShell):**
```powershell
# Set permanently for your user
[System.Environment]::SetEnvironmentVariable('GITHUB_TOKEN', 'ghp_yourTokenHere', 'User')
# Restart your terminal after this
```

**On Windows (CMD):**
```cmd
# Set permanently for your user
setx GITHUB_TOKEN ghp_yourTokenHere
# Restart your terminal after this
```

**Usage:** Run from anywhere:

```bash
gitlingo languages
```

---

### Getting your GitHub Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens > Tokens (classic)](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Give it a name (e.g., "GitLingo")
4. Select the `repo` scope
5. Click "Generate token"
6. Copy the token (starts with `ghp_`)

## Usage

Display your top 15 programming languages by total bytes across all your repositories.

### Options

**Show all languages** (instead of just top 15):

```bash
gitlingo languages --all
```

**Exclude specific languages** (e.g., markup or configuration languages):

```bash
gitlingo languages --exclude CSS,HTML,Dart
```

**Combine flags**:

```bash
gitlingo languages --all --exclude CSS,Python
```

> **Note:** If you used Option 1 (local build), run commands from the project directory with `./bin/gitlingo` instead of just `gitlingo`

## How it works

GitLingo fetches all your GitHub repositories and aggregates language statistics to show you what you code in most.

**Note:** GitLingo analyzes:
- Both **public and private** repositories (requires `repo` scope in your GitHub token)
- Only repositories you **own** (excludes forked repositories)
- Your actual **source code** (GitHub's language detection automatically excludes dependencies like `node_modules`, `vendor`, etc.)

## License

MIT License - see [LICENSE](LICENSE) for details.
