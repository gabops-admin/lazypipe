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
- **One-Key Re-Trigger:** Press `r` to re-run a failed or flaky workflow directly from the TUI.
- **Zero-Config Live & Demo Modes:** Automatically hooks into your local `gh` CLI when inside a GitHub repository, or launches with a rich multi-provider preview out of the box.

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

## Keybindings

| Key | Action |
| :--- | :--- |
| `j` / `↓` | Move selection down |
| `k` / `↑` | Move selection up |
| `r` | Re-trigger selected pipeline run |
| `f` | Toggle filter for `FAILED` runs only |
| `R` | Refresh live pipeline runs |
| `q` / `Ctrl+C` | Quit Lazypipe |

---

## License

MIT License — Built and maintained by **[GabOps](https://gabops.dev)**.
