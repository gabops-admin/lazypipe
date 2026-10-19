package providers

import "time"

// RunStatus represents the current state of a CI/CD pipeline execution.
type RunStatus string

const (
	StatusSuccess    RunStatus = "SUCCESS"
	StatusFailed     RunStatus = "FAILED"
	StatusRunning    RunStatus = "RUNNING"
	StatusQueued     RunStatus = "QUEUED"
	StatusCancelled  RunStatus = "CANCELLED"
)

// ProviderType identifies the CI/CD platform.
type ProviderType string

const (
	ProviderGitHubActions ProviderType = "GitHub Actions"
	ProviderAzureDevOps   ProviderType = "Azure DevOps"
)

// Step represents a single stage or job step inside a pipeline run.
type Step struct {
	Name     string
	Status   RunStatus
	Duration string
}

// PipelineRun represents a single CI/CD workflow or pipeline execution.
type PipelineRun struct {
	ID         string
	Name       string
	Provider   ProviderType
	Repo       string
	Branch     string
	CommitSHA  string
	CommitMsg  string
	Author     string
	Status     RunStatus
	Duration   string
	StartedAt  time.Time
	Steps      []Step
	Logs       []string
}

// SampleRuns returns realistic multi-provider pipeline runs for immediate demo & TUI preview.
func SampleRuns() []PipelineRun {
	now := time.Now()
	return []PipelineRun{
		{
			ID:        "run-94821",
			Name:      "Enterprise Cloud & Infrastructure CI/CD",
			Provider:  ProviderGitHubActions,
			Repo:      "gabops/cloud-cicd-blueprint-pro",
			Branch:    "main",
			CommitSHA: "3373763",
			CommitMsg: "feat(observability,docs): add OpenTelemetry, Prometheus, Grafana Tempo",
			Author:    "gabriel@gabops.dev",
			Status:    StatusSuccess,
			Duration:  "3m 42s",
			StartedAt: now.Add(-4 * time.Minute),
			Steps: []Step{
				{Name: "Security & IaC Compliance Gate (Trivy)", Status: StatusSuccess, Duration: "28s"},
				{Name: "Deploy Dev (10.10.0.0/16)", Status: StatusSuccess, Duration: "54s"},
				{Name: "Deploy Staging (10.50.0.0/16)", Status: StatusSuccess, Duration: "1m 08s"},
				{Name: "Deploy Production (10.100.0.0/16)", Status: StatusSuccess, Duration: "1m 12s"},
			},
			Logs: []string{
				"[18:50:12] INFO  Starting workflow: Enterprise Cloud & Infrastructure CI/CD",
				"[18:50:14] INFO  Authenticating with Azure via OIDC Workload Identity (passwordless)...",
				"[18:50:16] OK    OIDC Federated Token acquired for subscription e4b1...9a2c",
				"[18:50:19] INFO  Cache hit for key: tf-providers-Linux-production-8fbaa2f",
				"[18:50:22] INFO  Running Aqua Trivy IaC scan on ./terraform (threshold: HIGH,CRITICAL)...",
				"[18:50:40] OK    Trivy SARIF report uploaded: 0 CRITICAL, 0 HIGH misconfigurations found",
				"[18:51:34] INFO  Applying immutable artifact: production.tfplan on gabops-production-aks",
				"[18:53:48] OK    azurerm_kubernetes_cluster.aks: Apply complete! 0 added, 2 changed, 0 destroyed.",
				"[18:53:54] OK    Pipeline completed successfully in 3m 42s.",
			},
		},
		{
			ID:        "ado-2049",
			Name:      "AKS Ephemeral Runners & FinOps Audit",
			Provider:  ProviderAzureDevOps,
			Repo:      "gabops/platform-infra",
			Branch:    "feat/calico-netpol",
			CommitSHA: "8fbaa2f",
			CommitMsg: "fix(aks): tighten egress rules on snet-runners-ephemeral",
			Author:    "gabriel@gabops.dev",
			Status:    StatusFailed,
			Duration:  "1m 19s",
			StartedAt: now.Add(-12 * time.Minute),
			Steps: []Step{
				{Name: "Trivy IaC & Vulnerability Scan", Status: StatusFailed, Duration: "41s"},
				{Name: "Terraform Plan (dev)", Status: StatusCancelled, Duration: "0s"},
				{Name: "Terraform Apply (dev)", Status: StatusCancelled, Duration: "0s"},
			},
			Logs: []string{
				"[18:41:02] INFO  Pool: gabops-ephemeral-runners | Agent: runner-pod-7d9f8b-x4k2p",
				"[18:41:05] INFO  Checking out feat/calico-netpol (commit 8fbaa2f)...",
				"[18:41:10] INFO  Executing Trivy config scan on terraform/modules/aks-eks...",
				"[18:41:38] WARN  AVD-AZU-0041 (HIGH): Network policy allow-all-egress detected in staging override",
				"[18:41:41] ERROR Security Gate failed: 1 HIGH severity misconfiguration detected (exit code 1)",
				"[18:41:42] ERROR Hint: Press [r] in Lazypipe to re-trigger after fixing or [o] to inspect SARIF report.",
			},
		},
		{
			ID:        "run-94825",
			Name:      "BuildKit Multi-Stage Container & CVE Scan",
			Provider:  ProviderGitHubActions,
			Repo:      "gabops/otel-ingest-service",
			Branch:    "main",
			CommitSHA: "eb90c7d",
			CommitMsg: "perf(docker): enable GHA registry layer cache mode=max",
			Author:    "gabriel@gabops.dev",
			Status:    StatusRunning,
			Duration:  "48s (running)",
			StartedAt: now.Add(-48 * time.Second),
			Steps: []Step{
				{Name: "Set up Docker Buildx", Status: StatusSuccess, Duration: "4s"},
				{Name: "Build Image with Registry Layer Cache", Status: StatusSuccess, Duration: "19s"},
				{Name: "Scan Built Container Image (Trivy)", Status: StatusRunning, Duration: "25s"},
				{Name: "Push Verified Image to ghcr.io", Status: StatusQueued, Duration: "-"},
			},
			Logs: []string{
				"[18:56:32] INFO  Initializing Buildx builder instance: gabops-builder",
				"[18:56:36] INFO  Importing cache manifest from ghcr.io/gabops/otel-ingest-service:buildcache",
				"[18:56:39] OK    14/16 layers restored from cache (saved 2m 10s)",
				"[18:56:55] INFO  Image built: ghcr.io/gabops/otel-ingest-service:eb90c7d",
				"[18:56:56] INFO  Scanning OS & Go library dependencies for HIGH/CRITICAL CVEs...",
				"[18:57:18] INFO  Analyzing 84 packages... no critical vulnerabilities found so far...",
			},
		},
		{
			ID:        "ado-2050",
			Name:      "Observability Helm Sync (OTel + Tempo)",
			Provider:  ProviderAzureDevOps,
			Repo:      "gabops/cloud-cicd-blueprint-pro",
			Branch:    "main",
			CommitSHA: "0d5920c",
			CommitMsg: "chore(helm): bump grafana/tempo to 2.4.1 with metricsGenerator",
			Author:    "gabriel@gabops.dev",
			Status:    StatusSuccess,
			Duration:  "2m 05s",
			StartedAt: now.Add(-28 * time.Minute),
			Steps: []Step{
				{Name: "Validate Helm Charts & PrometheusRules", Status: StatusSuccess, Duration: "18s"},
				{Name: "Upgrade kube-prometheus-stack", Status: StatusSuccess, Duration: "52s"},
				{Name: "Upgrade Grafana Tempo", Status: StatusSuccess, Duration: "55s"},
			},
			Logs: []string{
				"[18:29:00] INFO  Connecting to gabops-production-aks (10.100.1.0/24)...",
				"[18:29:15] OK    PrometheusRule gabops-cicd-runner-alerts validated",
				"[18:30:07] OK    Release \"gabops-monitoring\" has been upgraded. Happy Helming!",
				"[18:31:02] OK    Release \"tempo\" has been upgraded. OTLP receivers listening on :4317/:4318",
			},
		},
	}
}
