# GitLingo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/steamedeo/gitlingo)](https://goreportcard.com/report/github.com/steamedeo/gitlingo)

![GitLingo Language Stats](img/sample.png)

Checkout your GitHub programming language statistics straight in your terminal.

## Installation

### From Source

```bash
git clone https://github.com/steamedeo/gitlingo.git
cd gitlingo
make build
```

## Setup

1. Create a GitHub Personal Access Token:
   - Go to [GitHub Settings > Developer settings > Personal access tokens > Tokens (classic)](https://github.com/settings/tokens)
   - Click "Generate new token (classic)"
   - Give it a name (e.g., "GitLingo")
   - Select the `repo` scope
   - Click "Generate token"
   - Copy the token (starts with `ghp_`)

2. Create a `.env` file:

   ```bash
   cp .env.example .env
   ```

3. Add your token to the `.env` file:
   ```
   github_token=ghp_yourTokenHere
   ```

## Usage

```bash
./bin/gitlingo languages
```

This will display your top 15 programming languages by total bytes across all your repositories.

To see all languages instead of just the top 15:

```bash
./bin/gitlingo languages --all
```

## How it works

GitLingo fetches all your GitHub repositories and aggregates language statistics to show you what you code in most.

## License

MIT License - see [LICENSE](LICENSE) for details.
