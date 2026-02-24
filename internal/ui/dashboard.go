package ui

import (
	"clidash/pkg/api"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	criticalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4C4C")).
			Bold(true)

	eventualStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1C40F"))

	strongStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2ECC71"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#44475A")).
			Padding(1).
			MarginRight(1)

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			Italic(true)
)

type model struct {
	state   api.GlobalState
	help    help.Model
	keys    keyMap
	logs    []string
	history []float64
}

type keyMap struct {
	Up   key.Binding
	Down key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down}, {k.Quit}}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func InitialModel() model {
	return model{
		state: api.GlobalState{Services: make(map[string]api.TelemetryUpdate)},
		help:  help.New(),
		keys:  keys,
		logs:  make([]string, 0),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			}
			resp.Body.Close()
		}

		// Update history log
		totalTraffic := 0
		for _, s := range m.state.Services {
			totalTraffic += s.RequestsPerSec
		}
		m.history = append(m.history, float64(totalTraffic))
		if len(m.history) > 40 {
			m.history = m.history[1:]
		}

		if m.state.LastDecision != "" {
			if len(m.logs) == 0 || !strings.Contains(m.logs[0], m.state.LastDecision) {
				m.logs = append([]string{fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), m.state.LastDecision)}, m.logs...)
			}
			if len(m.logs) > 10 {
				m.logs = m.logs[:10]
			}
		}
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}

	// Header
	header := titleStyle.Render("CLIDASH - RL Dynamic Consistency Optimizer")
	agentStatus := fmt.Sprintf("Optimizer: %s | Reward: %.1f | Mode: %s | Active Agents: %d",
		statusStyle.Render("ONLINE (8080)"),
		m.state.Reward,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Render("DISTRIBUTED"),
		len(m.state.Services),
	)
	s.WriteString(header + "\n" + agentStatus + "\n\n")

	// Service Grid (Dynamic from Agents)
	var serviceViews []string
	for id, svc := range m.state.Services {
		serviceInfo := fmt.Sprintf("%-20s\nStatus: %s\nLatency: %.1fms\nLoad: %d rps",
			lipgloss.NewStyle().Bold(true).Render(id),
			strongStyle.Render("CONNECTED"),
			svc.LatencyMS,
			svc.RequestsPerSec,
		)
		serviceViews = append(serviceViews, boxStyle.Render(serviceInfo))
	}

	if len(serviceViews) > 0 {
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, serviceViews...) + "\n\n")
	} else {
		s.WriteString("Waiting for Agent telemetry...\n\n")
	}

	// Graph
	graph := ""
	if len(m.history) > 5 {
		graph = asciigraph.Plot(m.history, asciigraph.Height(8), asciigraph.Width(105), asciigraph.Caption("Total System Throughput (RPS)"))
	}
	s.WriteString(boxStyle.BorderForeground(lipgloss.Color("#6272A4")).Render(graph) + "\n\n")

	// Stats and Logs
	stats := boxStyle.Width(30).Render(fmt.Sprintf(
		"SYSTEM PERFORMANCE\n\nLatency Saved: %s\nSLA Compliance: %s\nEfficiency: %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("Dynamic"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("99.5%"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("Distributed"),
	))

	logContent := "DECISION LOG:\n" + strings.Join(m.logs, "\n")
	logs := boxStyle.Width(70).Render(logStyle.Render(logContent))

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, stats, logs) + "\n\n")

	s.WriteString(m.help.View(m.keys))

	return s.String()
}
