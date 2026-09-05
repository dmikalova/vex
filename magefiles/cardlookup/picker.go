package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// pickSet opens an interactive ↑/↓ list of the source sets and returns the slug
// the user chose, or "" if they cancelled. It is what `mage tool:missing` runs
// when no set is named, so a set can be picked without remembering its slug.
func pickSet() (string, error) {
	sets := provenance.Sets()
	items := make([]setItem, 0, len(sets))
	for _, s := range sets {
		items = append(items, setItem{
			slug:  s.Slug,
			label: fmt.Sprintf("%-5s %s", s.Code, s.Name),
		})
	}
	result, err := tea.NewProgram(setPicker{items: items}).Run()
	if err != nil {
		return "", err
	}
	picked := result.(setPicker)
	if picked.cancelled {
		return "", nil
	}
	return picked.items[picked.cursor].slug, nil
}

// setItem is one selectable set: the slug the command needs and the label shown.
type setItem struct {
	slug  string
	label string
}

// setPicker is the Bubble Tea model backing pickSet: a single-column menu whose
// cursor moves with ↑/↓ (or k/j) and commits on Enter.
type setPicker struct {
	items     []setItem
	cursor    int
	cancelled bool
}

func (m setPicker) Init() tea.Cmd { return nil }

func (m setPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "enter":
		return m, tea.Quit
	case "q", "esc", "ctrl+c":
		m.cancelled = true
		return m, tea.Quit
	}
	return m, nil
}

// pickerCursorStyle highlights the row the cursor is on.
var pickerCursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))

func (m setPicker) View() string {
	var b strings.Builder
	b.WriteString("Select a set (↑/↓ to move, Enter to choose, q to cancel):\n\n")
	for i, it := range m.items {
		if i == m.cursor {
			b.WriteString("▸ " + pickerCursorStyle.Render(it.label) + "\n")
			continue
		}
		b.WriteString("  " + it.label + "\n")
	}
	return b.String()
}
