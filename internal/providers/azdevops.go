package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type adoDefinition struct {
	Name string `json:"name"`
}

type adoProject struct {
	Name string `json:"name"`
}

type adoRequestedFor struct {
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}

type adoRunJSON struct {
	ID            int64           `json:"id"`
	BuildNumber   string          `json:"buildNumber"`
	Status        string          `json:"status"`
	Result        string          `json:"result"`
	SourceBranch  string          `json:"sourceBranch"`
	SourceVersion string          `json:"sourceVersion"`
	StartTime     time.Time       `json:"startTime"`
	FinishTime    time.Time       `json:"finishTime"`
	Definition    adoDefinition   `json:"definition"`
	Project       adoProject      `json:"project"`
	RequestedFor  adoRequestedFor `json:"requestedFor"`
}

// FetchAzureDevOpsRuns queries live pipeline executions via `az pipelines runs list --output json`.
func FetchAzureDevOpsRuns() ([]PipelineRun, error) {
	cmd := exec.Command("az", "pipelines", "runs", "list", "--top", "8", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var adoRuns []adoRunJSON
	if err := json.Unmarshal(out, &adoRuns); err != nil {
		return nil, err
	}
	if len(adoRuns) == 0 {
		return nil, fmt.Errorf("no Azure DevOps runs returned")
	}

	runs := make([]PipelineRun, 0, len(adoRuns))
	for _, r := range adoRuns {
		status := mapADOStatus(r.Status, r.Result)
		branch := strings.TrimPrefix(r.SourceBranch, "refs/heads/")
		sha := r.SourceVersion
		if len(sha) > 7 {
			sha = sha[:7]
		}
		dur := "-"
		if !r.StartTime.IsZero() {
			end := r.FinishTime
			if end.IsZero() {
				end = time.Now()
			}
			dur = end.Sub(r.StartTime).Round(time.Second).String()
		}
		author := r.RequestedFor.UniqueName
		if author == "" {
			author = r.RequestedFor.DisplayName
		}

		runs = append(runs, PipelineRun{
			ID:        fmt.Sprintf("ado-%d", r.ID),
			Name:      r.Definition.Name,
			Provider:  ProviderAzureDevOps,
			Repo:      r.Project.Name,
			Branch:    branch,
			CommitSHA: sha,
			CommitMsg: fmt.Sprintf("Build #%s", r.BuildNumber),
			Author:    author,
			Status:    status,
			Duration:  dur,
			StartedAt: r.StartTime,
			Steps: []Step{
				{Name: r.Definition.Name, Status: status, Duration: dur},
			},
			Logs: []string{
				fmt.Sprintf("Loaded live Azure DevOps run #%d (%s) via `az pipelines` CLI", r.ID, r.BuildNumber),
				fmt.Sprintf("Pipeline: %s | Project: %s | Branch: %s", r.Definition.Name, r.Project.Name, branch),
				fmt.Sprintf("Status: %s (%s) | Duration: %s", r.Status, r.Result, dur),
			},
		})
	}

	return runs, nil
}

// RetriggerAzureDevOpsRun queues a new pipeline run on the same branch using `az pipelines run`.
func RetriggerAzureDevOpsRun(pipelineName, branch string) error {
	cmd := exec.Command("az", "pipelines", "run", "--name", pipelineName, "--branch", branch)
	return cmd.Run()
}

func mapADOStatus(status, result string) RunStatus {
	s := strings.ToLower(status)
	r := strings.ToLower(result)
	if s == "inprogress" {
		return StatusRunning
	}
	if s == "notstarted" || s == "postponed" {
		return StatusQueued
	}
	switch r {
	case "succeeded", "partiallysucceeded":
		return StatusSuccess
	case "failed":
		return StatusFailed
	case "canceled", "cancelled":
		return StatusCancelled
	default:
		return StatusQueued
	}
}
