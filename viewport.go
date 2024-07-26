package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = true

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "|"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Left = "|"
		return titleStyle.BorderStyle(b)
	}()

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	contentStyle = lipgloss.NewStyle().Padding(2, 4).Margin(2, 1)
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	title    string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" || k == "enter" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {

			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(contentStyle.Render(m.content))
			m.ready = true

			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {

			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render(m.title)
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}

func (m model) footerView() string {
	help := helpStyle.Render("↑↓/jk/scroll, q/esc/enter to quit")
	return lipgloss.JoinHorizontal(lipgloss.Center, help)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clearScreen() {
	fmt.Println("\033[H\033[2J")
}

func Show(title, content string, useLineNumbers bool) {
	if useLineNumbers {
		content = applyLineNumbers(content)
	}
	p := tea.NewProgram(model{content: content, title: title})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}

func applyLineNumbers(content string) string {
	b := strings.Builder{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for i := 1; scanner.Scan(); i++ {
		b.WriteString(fmt.Sprintf("%s %s\n", lineNumberStyle.Render(strconv.Itoa(i)), scanner.Text()))
	}
	return b.String()
}

func Quick(content string) {
	fmt.Println(content)
}
