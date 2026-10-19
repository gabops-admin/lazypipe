package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
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

// FetchLiveGitHubRuns attempts to load recent GitHub Actions runs using the local `gh` CLI if inside a repo.
// Falls back cleanly to SampleRuns() if no runs exist or if offline.
func FetchLiveOrSampleRuns() ([]PipelineRun, string) {
	cmd := exec.Command("gh", "run", "list", "--limit", "8", "--json", "databaseId,displayTitle,workflowName,headBranch,headSha,status,conclusion,createdAt,updatedAt")
	out, err := cmd.Output()
	if err != nil {
		return SampleRuns(), "Demo Mode (GitHub Actions + Azure DevOps)"
	}

	var ghRuns []ghRunJSON
	if err := json.Unmarshal(out, &ghRuns); err != nil || len(ghRuns) == 0 {
		return SampleRuns(), "Demo Mode (GitHub Actions + Azure DevOps)"
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
			Repo:      "local-git-repo",
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
	return runs, "Live GitHub CLI Mode"
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
