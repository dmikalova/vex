package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// explorerModel browses the whole card database with a fuzzy search box: typing
// filters the list, and the selected card's generated text shows on the right.
type explorerModel struct {
	all      []engine.CardDefinition
	filtered []engine.CardDefinition
	query    string
	cursor   int
	top      int
	width    int
	height   int
}

func newExplorerModel(w, h int) explorerModel {
	all := cards.All()
	return explorerModel{all: all, filtered: all, width: w, height: h}
}

func (m explorerModel) resize(w, h int) explorerModel { m.width, m.height = w, h; return m }

func (m explorerModel) listHeight() int {
	h := m.height - 6 // title + search box + help + padding
	if h < 3 {
		h = 3
	}
	return h
}

func (m explorerModel) clampScroll() explorerModel {
	visible := m.listHeight()
	if m.cursor < m.top {
		m.top = m.cursor
	}
	if m.cursor >= m.top+visible {
		m.top = m.cursor - visible + 1
	}
	return m
}

func (m explorerModel) Update(msg tea.Msg) (explorerModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case tea.KeyRunes, tea.KeySpace:
		m.query += string(key.Runes)
		return m.refilter(), nil
	case tea.KeyBackspace:
		if r := []rune(m.query); len(r) > 0 {
			m.query = string(r[:len(r)-1])
			return m.refilter(), nil
		}
	case tea.KeyEsc:
		if m.query != "" {
			m.query = ""
			return m.refilter(), nil
		}
		return m, gotoScreen(screenMenu)
	}
	return m.clampScroll(), nil
}

// refilter recomputes the visible list from the query (fuzzy, best-first) and
// resets the cursor to the top.
func (m explorerModel) refilter() explorerModel {
	if m.query == "" {
		m.filtered = m.all
	} else {
		type scored struct {
			def   engine.CardDefinition
			score int
		}
		var matches []scored
		for _, d := range m.all {
			if s, ok := fuzzyScore(m.query, d.Name); ok {
				matches = append(matches, scored{d, s})
			}
		}
		sort.SliceStable(
			matches,
			func(i, j int) bool { return matches[i].score > matches[j].score },
		)
		m.filtered = make([]engine.CardDefinition, len(matches))
		for i, s := range matches {
			m.filtered[i] = s.def
		}
	}
	m.cursor, m.top = 0, 0
	return m
}

func (m explorerModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Card explorer") + "\n")
	b.WriteString(faintStyle.Render("search: ") + m.query + "▌\n\n")

	if len(m.filtered) == 0 {
		b.WriteString(faintStyle.Render("  no cards match") + "\n\n")
		b.WriteString(helpStyle.Render("type to search · esc clear/menu"))
		return b.String()
	}

	visible := m.listHeight()
	end := m.top + visible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	var rows []string
	for i := m.top; i < end; i++ {
		name := m.filtered[i].Name
		if i == m.cursor {
			rows = append(rows, selectedStyle.Render(cursor(true)+name))
		} else {
			rows = append(rows, cursor(false)+name)
		}
	}
	const listWidth = 24
	list := lipgloss.NewStyle().Width(listWidth).Render(strings.Join(rows, "\n"))

	def := m.filtered[m.cursor]
	detail := headerStyle.Render(def.Name) + "\n\n" + engine.RenderCardText(&def)
	detailStyle := lipgloss.NewStyle().PaddingLeft(3)
	// Wrap the card text to the space left of the list instead of letting long
	// ability lines overflow the terminal and get truncated.
	if avail := m.width - listWidth; avail > 8 {
		detailStyle = detailStyle.Width(avail)
	}
	detail = detailStyle.Render(detail)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, list, detail))
	b.WriteString(
		"\n\n" + helpStyle.Render(
			fmt.Sprintf(
				"%d/%d · ↑/↓ move · type to search · esc clear/menu",
				m.cursor+1,
				len(m.filtered),
			),
		),
	)
	return b.String()
}

// fuzzyScore reports whether every rune of query appears in order within target
// (case-insensitive, spaces ignored) and returns a score where higher is a better
// match: consecutive runs and a shared prefix score higher.
func fuzzyScore(query, target string) (int, bool) {
	q := []rune(normalizeName(query))
	t := []rune(normalizeName(target))
	if len(q) == 0 {
		return 0, true
	}
	score, ti, streak := 0, 0, 0
	for _, qc := range q {
		matched := false
		for ti < len(t) {
			c := t[ti]
			ti++
			if c == qc {
				score += 1 + streak
				streak++
				matched = true
				break
			}
			streak = 0
		}
		if !matched {
			return 0, false
		}
	}
	if strings.HasPrefix(normalizeName(target), normalizeName(query)) {
		score += 10
	}
	return score, true
}

func normalizeName(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}
