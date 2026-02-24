package main

import (
	"clidash/internal/ui"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	useTable := flag.Bool("table", false, "Use table view")
	flag.Parse()

	var m tea.Model
	if *useTable {
		m = ui.InitialTableModel()
	} else {
		m = ui.InitialModel()
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
