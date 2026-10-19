package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabops/lazypipe/internal/providers"
	"github.com/gabops/lazypipe/internal/ui"
)

const version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "Print Lazypipe version and exit")
	dumpOnce := flag.Bool("once", false, "Print pipeline status summary to stdout and exit (non-interactive CI mode)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("lazypipe v%s — Terminal UI for GitHub Actions & Azure DevOps (https://gabops.dev)\n", version)
		return
	}

	if *dumpOnce {
		runs, mode := providers.FetchLiveOrSampleRuns()
		fmt.Printf("⚡ LAZYPIPE v%s [%s]\n", version, mode)
		fmt.Println("--------------------------------------------------------------------------------")
		for _, r := range runs {
			fmt.Printf("%-10s | %-14s | %-38s | %s (%s)\n",
				r.Status, r.Provider, r.Name, r.Branch, r.Duration)
		}
		return
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running lazypipe: %v\n", err)
		os.Exit(1)
	}
}
