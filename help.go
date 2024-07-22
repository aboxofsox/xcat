package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	V    key.Binding // view mode
	N    key.Binding // line numbers
	S    key.Binding // start line
	E    key.Binding // end line
	F    key.Binding // find string
	Quit key.Binding // quit
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.V, k.N, k.S, k.E, k.F}

}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.V, k.N, k.S, k.E, k.F},
		{k.Quit},
	}
}

var keys = keyMap{
	V: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "view mode"),
	),
	N: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "line numbers"),
	),
	S: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start line"),
	),
	E: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "end line"),
	),
	F: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "find string"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type helpModel struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
}

func newHelpModel() helpModel {
	return helpModel{
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m helpModel) Init() tea.Cmd {
	return nil
}

func (m helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m helpModel) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.help.View(m.keys)
}

func ShowHelpModel() {
	if _, err := tea.NewProgram(newHelpModel()).Run(); err != nil {
		fmt.Printf("could not start:\n%v\n", err)
		os.Exit(1)
	}
}
