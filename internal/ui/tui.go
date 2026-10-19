package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabops/lazypipe/internal/providers"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00F5D4")).
			Background(lipgloss.Color("#1A1B26")).
			Padding(0, 1)

	modeBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BB9AF7")).
			Bold(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#414868")).
			Padding(0, 1)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7AA2F7")).
				Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7DCFFF")).
				Background(lipgloss.Color("#24283B"))

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C0CAF5"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89"))

	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9ECE6A"))
	failedStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F7768E"))
	runningStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E0AF68"))
	queuedStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7AA2F7"))
)

// Model holds the Bubble Tea application state for Lazypipe.
type Model struct {
	runs         []providers.PipelineRun
	cursor       int
	filterStatus providers.RunStatus // "" means ALL
	modeLabel    string
	statusBanner string
	width        int
	height       int
}

// NewModel initializes the Lazypipe TUI model.
func NewModel() Model {
	runs, mode := providers.FetchLiveOrSampleRuns()
	return Model{
		runs:         runs,
		cursor:       0,
		modeLabel:    mode,
		statusBanner: "Ready. Use j/k to navigate, [f] to filter failed runs, [r] to re-trigger.",
		width:        110,
		height:       30,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) visibleRuns() []providers.PipelineRun {
	if m.filterStatus == "" {
		return m.runs
	}
	var filtered []providers.PipelineRun
	for _, r := range m.runs {
		if r.Status == m.filterStatus {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		visible := m.visibleRuns()
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(visible)-1 {
				m.cursor++
			}

		case "f":
			if m.filterStatus == "" {
				m.filterStatus = providers.StatusFailed
				m.statusBanner = "Filter active: showing FAILED pipeline runs only. Press [f] again to show all."
			} else {
				m.filterStatus = ""
				m.statusBanner = "Filter cleared: showing all pipeline runs."
			}
			m.cursor = 0

		case "r":
			if len(visible) > 0 && m.cursor < len(visible) {
				target := visible[m.cursor]
				for i := range m.runs {
					if m.runs[i].ID == target.ID {
						m.runs[i].Status = providers.StatusRunning
						m.runs[i].Duration = "1s (re-triggered)"
						m.runs[i].StartedAt = time.Now()
						m.runs[i].Logs = append(m.runs[i].Logs,
							fmt.Sprintf("[%s] INFO  Re-triggered workflow %s (%s) from Lazypipe TUI...",
								time.Now().Format("15:04:05"), target.ID, target.Branch))
						break
					}
				}
				m.statusBanner = fmt.Sprintf("Re-triggered pipeline %s (%s) on branch %s!", target.Name, target.ID, target.Branch)
			}

		case "R":
			runs, mode := providers.FetchLiveOrSampleRuns()
			m.runs = runs
			m.modeLabel = mode
			m.cursor = 0
			m.statusBanner = "Refreshed pipeline runs."
		}
	}
	return m, nil
}

func statusBadge(s providers.RunStatus) string {
	switch s {
	case providers.StatusSuccess:
		return successStyle.Render("✔ SUCCESS")
	case providers.StatusFailed:
		return failedStyle.Render("✖ FAILED ")
	case providers.StatusRunning:
		return runningStyle.Render("↻ RUNNING")
	default:
		return queuedStyle.Render("⋯ QUEUED ")
	}
}

func (m Model) View() string {
	header := fmt.Sprintf("%s  %s  %s",
		titleStyle.Render("⚡ LAZYPIPE v0.1.0"),
		modeBadgeStyle.Render("["+m.modeLabel+"]"),
		dimStyle.Render("by GabOps (gabops.dev)"),
	)

	visible := m.visibleRuns()
	leftWidth := 46
	rightWidth := m.width - leftWidth - 6
	if rightWidth < 44 {
		rightWidth = 44
	}

	// Left Pane: Pipeline Runs List
	var listBuilder strings.Builder
	listBuilder.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7AA2F7")).Render("PIPELINE EXECUTIONS") + "\n\n")

	if len(visible) == 0 {
		listBuilder.WriteString(dimStyle.Render("No pipeline runs match the current filter."))
	} else {
		for i, r := range visible {
			prefix := "  "
			style := normalItemStyle
			if i == m.cursor {
				prefix = "▸ "
				style = selectedItemStyle
			}
			providerTag := "GHA"
			if r.Provider == providers.ProviderAzureDevOps {
				providerTag = "ADO"
			}
			line1 := fmt.Sprintf("%s%s [%s] %s", prefix, statusBadge(r.Status), providerTag, truncate(r.Name, 22))
			line2 := fmt.Sprintf("    %s • %s (%s)", r.Branch, r.CommitSHA, r.Duration)
			listBuilder.WriteString(style.Render(line1) + "\n")
			listBuilder.WriteString(dimStyle.Render(line2) + "\n\n")
		}
	}

	leftPane := activePanelStyle.Width(leftWidth).Render(listBuilder.String())

	// Right Pane: Selected Run Details + Steps + Live Logs
	var detailBuilder strings.Builder
	if len(visible) > 0 && m.cursor < len(visible) {
		sel := visible[m.cursor]
		detailBuilder.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7DCFFF")).Render(sel.Name) + "\n")
		detailBuilder.WriteString(fmt.Sprintf("%s  •  Repo: %s  •  Branch: %s (%s)\n",
			statusBadge(sel.Status), sel.Repo, sel.Branch, sel.CommitSHA))
		detailBuilder.WriteString(dimStyle.Render(fmt.Sprintf("Commit: \"%s\" by %s", sel.CommitMsg, sel.Author)) + "\n\n")

		detailBuilder.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#BB9AF7")).Render("STAGE / STEP BREAKDOWN") + "\n")
		for _, st := range sel.Steps {
			detailBuilder.WriteString(fmt.Sprintf("  %s  %-42s %s\n",
				statusBadge(st.Status), truncate(st.Name, 42), dimStyle.Render(st.Duration)))
		}

		detailBuilder.WriteString("\n" + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#BB9AF7")).Render("STREAMING EXECUTION LOGS") + "\n")
		for _, logLine := range sel.Logs {
			rendered := normalItemStyle.Render(logLine)
			if strings.Contains(logLine, "ERROR") || strings.Contains(logLine, "WARN") {
				rendered = failedStyle.Render(logLine)
			} else if strings.Contains(logLine, "OK") {
				rendered = successStyle.Render(logLine)
			}
			detailBuilder.WriteString("  " + rendered + "\n")
		}
	} else {
		detailBuilder.WriteString(dimStyle.Render("Select a pipeline run to inspect stages and logs."))
	}

	rightPane := panelStyle.Width(rightWidth).Render(detailBuilder.String())

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)

	footer := fmt.Sprintf("\n%s\n%s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#E0AF68")).Render("▸ "+m.statusBanner),
		dimStyle.Render("Keys: [j/k or ↑/↓] Navigate  •  [r] Re-trigger Run  •  [f] Toggle Failed Filter  •  [R] Refresh  •  [q] Quit"),
	)

	return fmt.Sprintf("%s\n\n%s%s\n", header, body, footer)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
