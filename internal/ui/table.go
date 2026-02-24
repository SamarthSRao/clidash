package ui

import (
	"clidash/pkg/api"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type tableModel struct {
	state api.GlobalState
	table table.Model
	help  help.Model
	keys  keyMap
	logs  []string
}

func InitialTableModel() tableModel {
	columns := []table.Column{
		{Title: "Service ID", Width: 25},
		{Title: "Last Op", Width: 15},
		{Title: "Latency (ms)", Width: 15},
		{Title: "Load (RPS)", Width: 15},
		{Title: "Status", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return tableModel{
		state: api.GlobalState{Services: make(map[string]api.TelemetryUpdate)},
		table: t,
		help:  help.New(),
		keys:  keys,
		logs:  make([]string, 0),
	}
}

func (m tableModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case tickMsg:
		// Fetch state from Optimizer
		resp, err := http.Get("http://localhost:8080/state")
		if err == nil {
			var newState api.GlobalState
			if err := json.NewDecoder(resp.Body).Decode(&newState); err == nil {
				m.state = newState

				// Update table rows
				var rows []table.Row
				ids := make([]string, 0, len(m.state.Services))
				for id := range m.state.Services {
					ids = append(ids, id)
				}
				sort.Strings(ids)

				for _, id := range ids {
					svc := m.state.Services[id]
					rows = append(rows, table.Row{
						svc.ServiceID,
						svc.Operation,
						fmt.Sprintf("%.1f", svc.LatencyMS),
						fmt.Sprintf("%d", svc.RequestsPerSec),
						"CONNECTED",
					})
				}
				m.table.SetRows(rows)
			}
			resp.Body.Close()
		}

		if m.state.LastDecision != "" {
			if len(m.logs) == 0 || !strings.Contains(m.logs[0], m.state.LastDecision) {
				m.logs = append([]string{fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), m.state.LastDecision)}, m.logs...)
			}
			if len(m.logs) > 5 {
				m.logs = m.logs[:5]
			}
		}

		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	s := strings.Builder{}

	header := titleStyle.Render("CLIDASH - Tabular Service View")
	agentStatus := fmt.Sprintf("Optimizer: %s | Reward: %.1f | Mode: %s",
		statusStyle.Render("ONLINE"),
		m.state.Reward,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Render("DISTRIBUTED"),
	)

	s.WriteString(header + "\n" + agentStatus + "\n\n")
	s.WriteString(baseStyle.Render(m.table.View()) + "\n\n")

	logContent := "RECENT DECISIONS:\n" + strings.Join(m.logs, "\n")
	s.WriteString(boxStyle.Width(90).Render(logStyle.Render(logContent)) + "\n\n")

	s.WriteString(m.help.View(m.keys))

	return s.String()
}
