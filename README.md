# ⚡ Lazypipe

**A fast, keyboard-driven Terminal UI (TUI) for monitoring, filtering, and re-triggering GitHub Actions & Azure DevOps pipelines without leaving your terminal.**

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Built by GabOps](https://img.shields.io/badge/Built%20by-GabOps-7AA2F7)

---

> 🚀 **Want Production-Ready Cloud & CI/CD Infrastructure out of the box?**
> Tired of slow builds, static secrets, and fragile YAML pipelines? Check out the **[Cloud & CI/CD Production Blueprint (Pro Edition)](https://gabops.dev)** — our battle-tested enterprise starter kit featuring modular Terraform (`dev`/`staging`/`prod`), hardened AKS/EKS clusters, passwordless OIDC pipelines, BuildKit caching, Trivy security gates, and OpenTelemetry observability.

---

## Why Lazypipe?

Context-switching between your terminal and 15 browser tabs just to check if a `terraform apply` or container security scan passed kills engineering flow. **Lazypipe** gives you a single, unified terminal dashboard (`lazygit`-style) for both **GitHub Actions** and **Azure DevOps**:

- **Unified Multi-Provider View:** Monitor GitHub Actions (`GHA`) and Azure DevOps (`ADO`) executions side by side.
- **Instant Failure Triage:** Press `f` to immediately isolate failed runs and highlight `ERROR` / `WARN` lines in Trivy & Terraform logs.
- **One-Key Re-Trigger:** Press `r` to re-run a failed or flaky workflow directly from the TUI (`gh run rerun` / `az pipelines run`).
- **Zero-Config Live & Demo Modes:** Automatically hooks into your local `gh` and `az` CLIs when authenticated, or launches with a rich multi-provider preview out of the box.

---

## Quickstart

### Build & Run from Source

```bash
git clone https://github.com/gabops-admin/lazypipe.git
cd lazypipe
go build -o lazypipe ./cmd/lazypipe

# Launch interactive Terminal UI
./lazypipe

# Non-interactive summary snapshot
./lazypipe --once
```

---

## Authentication & Provider Setup

**Lazypipe** delegates authentication to the official **GitHub CLI (`gh`)** and **Azure CLI (`az devops`)**, so your tokens remain securely stored in your OS keychain rather than plain-text config files.

### 1. GitHub Actions (`GHA`)

Install the GitHub CLI and authenticate:

```bash
# macOS (Homebrew)
brew install gh

# Interactive browser login (stores token in OS Keychain)
gh auth login

# Or authenticate via Personal Access Token (requires 'repo' and 'workflow' scopes)
export GH_TOKEN="ghp_your_personal_access_token"
```

Verify access inside any repository:
```bash
gh run list --limit 5
```

### 2. Azure DevOps (`ADO`)

Install the Azure CLI with the `azure-devops` extension and configure your default organization/project:

```bash
# 1. Install Azure CLI & Azure DevOps extension
brew install azure-cli
az extension add --name azure-devops

# 2. Authenticate using a Personal Access Token (Scopes: Build Read & Execute)
export AZURE_DEVOPS_EXT_PAT="your_azure_devops_pat"
# Or store it persistently via:
az devops login

# 3. Set your target Organization and Project defaults
az devops configure --defaults organization=https://dev.azure.com/YOUR_ORG project=YOUR_PROJECT
```

Verify access:
```bash
az pipelines runs list --top 5 --output table
```

> **Note:** If neither CLI is authenticated (or if you run `lazypipe` outside a configured project), Lazypipe automatically starts in **Demo Mode** so you can explore the interface right away.

---

## Keybindings

| Key | Action |
| :--- | :--- |
| `j` / `↓` | Move selection down |
| `k` / `↑` | Move selection up |
| `r` | Re-trigger selected pipeline run (`gh run rerun` / `az pipelines run`) |
| `f` | Toggle filter for `FAILED` runs only |
| `R` | Refresh live pipeline runs |
| `q` / `Ctrl+C` | Quit Lazypipe |

---

## License

MIT License — Built and maintained by **[GabOps](https://gabops.dev)**.
