package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type ghRunJSON struct {
	DatabaseID   int64     `json:"databaseId"`
	DisplayTitle string    `json:"displayTitle"`
	WorkflowName string    `json:"workflowName"`
	HeadBranch   string    `json:"headBranch"`
	HeadSha      string    `json:"headSha"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// FetchGitHubRuns queries live GitHub Actions workflow runs via `gh run list`.
func FetchGitHubRuns() ([]PipelineRun, error) {
	cmd := exec.Command("gh", "run", "list", "--limit", "8", "--json", "databaseId,displayTitle,workflowName,headBranch,headSha,status,conclusion,createdAt,UpdatedAt")
	// Note: gh uses updatedAt (camelCase)
	cmd = exec.Command("gh", "run", "list", "--limit", "8", "--json", "databaseId,displayTitle,workflowName,headBranch,headSha,status,conclusion,createdAt,updatedAt")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ghRuns []ghRunJSON
	if err := json.Unmarshal(out, &ghRuns); err != nil {
		return nil, err
	}
	if len(ghRuns) == 0 {
		return nil, fmt.Errorf("no GitHub Actions runs returned")
	}

	runs := make([]PipelineRun, 0, len(ghRuns))
	for _, r := range ghRuns {
		status := mapGHStatus(r.Status, r.Conclusion)
		sha := r.HeadSha
		if len(sha) > 7 {
			sha = sha[:7]
		}
		dur := r.UpdatedAt.Sub(r.CreatedAt).Round(time.Second).String()
		runs = append(runs, PipelineRun{
			ID:        fmt.Sprintf("gh-%d", r.DatabaseID),
			Name:      r.WorkflowName,
			Provider:  ProviderGitHubActions,
			Repo:      "github-repo",
			Branch:    r.HeadBranch,
			CommitSHA: sha,
			CommitMsg: strings.TrimSpace(r.DisplayTitle),
			Author:    "github-actions",
			Status:    status,
			Duration:  dur,
			StartedAt: r.CreatedAt,
			Steps: []Step{
				{Name: r.WorkflowName, Status: status, Duration: dur},
			},
			Logs: []string{
				fmt.Sprintf("Loaded live GitHub Actions run #%d via `gh` CLI", r.DatabaseID),
				fmt.Sprintf("Workflow: %s | Branch: %s | Commit: %s", r.WorkflowName, r.HeadBranch, sha),
				fmt.Sprintf("Status: %s | Duration: %s", status, dur),
			},
		})
	}
	return runs, nil
}

// RetriggerGitHubRun triggers `gh run rerun <run-id>` for a GitHub Actions workflow run.
func RetriggerGitHubRun(runID string) error {
	rawID := strings.TrimPrefix(runID, "gh-")
	cmd := exec.Command("gh", "run", "rerun", rawID)
	return cmd.Run()
}

// FetchLiveOrSampleRuns aggregates live runs from GitHub Actions (`gh`) and Azure DevOps (`az pipelines`),
// falling back cleanly to SampleRuns() when neither provider returns live runs.
func FetchLiveOrSampleRuns() ([]PipelineRun, string) {
	var combined []PipelineRun
	var activeSources []string

	if ghRuns, err := FetchGitHubRuns(); err == nil && len(ghRuns) > 0 {
		combined = append(combined, ghRuns...)
		activeSources = append(activeSources, "GHA Live")
	}

	if adoRuns, err := FetchAzureDevOpsRuns(); err == nil && len(adoRuns) > 0 {
		combined = append(combined, adoRuns...)
		activeSources = append(activeSources, "ADO Live")
	}

	if len(combined) == 0 {
		return SampleRuns(), "Demo Mode (GitHub Actions + Azure DevOps)"
	}

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].StartedAt.After(combined[j].StartedAt)
	})

	return combined, strings.Join(activeSources, " + ")
}

func mapGHStatus(status, conclusion string) RunStatus {
	if status != "completed" {
		if status == "in_progress" {
			return StatusRunning
		}
		return StatusQueued
	}
	switch conclusion {
	case "success":
		return StatusSuccess
	case "failure", "timed_out":
		return StatusFailed
	case "cancelled":
		return StatusCancelled
	default:
		return StatusQueued
	}
}
