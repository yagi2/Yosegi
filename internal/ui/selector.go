package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagi2/yosegi/internal/git"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Quit   key.Binding
	Delete key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

type SelectorModel struct {
	worktrees    []git.Worktree
	cursor       int
	title        string
	action       string
	allowDelete  bool
	selectedPath string
	quitting     bool
}

type SelectionResult struct {
	Worktree git.Worktree
	Action   string // "select", "delete", "quit"
}

func NewSelector(worktrees []git.Worktree, title, action string, allowDelete bool) SelectorModel {
	return SelectorModel{
		worktrees:   worktrees,
		cursor:      0,
		title:       title,
		action:      action,
		allowDelete: allowDelete,
	}
}

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.worktrees)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Enter):
			if len(m.worktrees) > 0 {
				m.selectedPath = m.worktrees[m.cursor].Path
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Delete):
			if m.allowDelete && len(m.worktrees) > 0 {
				m.selectedPath = m.worktrees[m.cursor].Path
				m.action = "delete"
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m SelectorModel) View() string {
	if m.quitting && m.selectedPath == "" {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(fmt.Sprintf("🌲 %s", m.title)))
	b.WriteString("\n\n")

	if len(m.worktrees) == 0 {
		b.WriteString(ErrorStyle.Render("No worktrees found"))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Press q to quit"))
		return BorderStyle.Render(b.String())
	}

	// Worktree list
	for i, worktree := range m.worktrees {
		var line strings.Builder

		// Status icon
		icon := GetStatusIcon(worktree.IsCurrent)
		if worktree.IsCurrent {
			line.WriteString(CurrentItemStyle.Render(icon + " "))
		} else {
			line.WriteString(NormalItemStyle.Render(icon + " "))
		}

		// Branch and path
		branchInfo := fmt.Sprintf("%s %s", GetBranchIcon(), worktree.Branch)
		pathInfo := fmt.Sprintf("%s %s", GetPathIcon(), shortenPath(worktree.Path))

		content := fmt.Sprintf("%-30s %s", branchInfo, pathInfo)

		if i == m.cursor {
			line.WriteString(SelectedItemStyle.Render(content))
		} else if worktree.IsCurrent {
			line.WriteString(CurrentItemStyle.Render(content))
		} else {
			line.WriteString(NormalItemStyle.Render(content))
		}

		b.WriteString(line.String())
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	helpText := []string{
		"↑/k up", "↓/j down", "enter " + m.action,
	}
	if m.allowDelete {
		helpText = append(helpText, "d delete")
	}
	helpText = append(helpText, "q quit")

	b.WriteString(HelpStyle.Render(strings.Join(helpText, " • ")))

	return BorderStyle.Render(b.String())
}

func (m SelectorModel) GetResult() SelectionResult {
	if m.quitting && m.selectedPath == "" {
		return SelectionResult{Action: "quit"}
	}

	for _, wt := range m.worktrees {
		if wt.Path == m.selectedPath {
			return SelectionResult{
				Worktree: wt,
				Action:   m.action,
			}
		}
	}

	return SelectionResult{Action: "quit"}
}

// shortenPath shortens a path for display
func shortenPath(path string) string {
	if len(path) <= 50 {
		return path
	}
	return "..." + path[len(path)-47:]
}
