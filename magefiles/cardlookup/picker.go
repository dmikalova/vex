package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// pickSet opens an interactive ↑/↓ list of the source sets that still have cards
// to implement and returns the slug the user chose, or "" if they cancelled (or
// every set is already complete). It is what `mage tool:missing` and
// `mage tool:nextCard` run when no set is named, so a set can be picked without
// remembering its slug. Sets with no missing cards are left off the list — there
// is nothing to pick in a complete set.
func pickSet() (string, error) {
	// Touch the aggregator so every set package registers its cards, then read
	// which collector numbers count as implemented.
	_ = cards.All()
	covered := coveredNumbers()

	sets := provenance.Sets()
	items := make([]setItem, 0, len(sets))
	for _, s := range sets {
		missing := 0
		for _, c := range s.Cards {
			if !covered[s.Slug][c.Number] {
				missing++
			}
		}
		if missing == 0 {
			continue // complete set — nothing to implement
		}
		items = append(items, setItem{
			slug:  s.Slug,
			label: fmt.Sprintf("%-5s %-24s %d to implement", s.Code, s.Name, missing),
		})
	}
	if len(items) == 0 {
		return "", nil // every set is complete
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

// pickImportSet opens an interactive ↑/↓ list of every importable set, preceded
// by an "All sets" entry, and returns the chosen slug ("all" for every set) or ""
// if cancelled. It is what `mage tool:importProvenance` runs when no set is named.
// Unlike pickSet it lists sets from the importer registry, not the embedded
// catalogs, so a set with no catalog yet can still be chosen.
func pickImportSet() (string, error) {
	items := []setItem{{slug: "all", label: "All sets"}}
	for _, is := range importSets {
		items = append(items, setItem{
			slug:  is.Set.Slug,
			label: fmt.Sprintf("%-5s %s", is.Set.Code, is.Set.Name),
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
