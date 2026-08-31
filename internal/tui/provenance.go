package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

// provStatus is one original card together with the vactrol card(s) that
// reference it (usually zero or one).
type provStatus struct {
	card provenance.Card
	impl []string
}

// provSet is a source set with per-card implementation status.
type provSet struct {
	set         provenance.SourceSet
	cards       []provStatus
	implemented int
}

type provFilter int

const (
	filterAll provFilter = iota
	filterImplemented
	filterMissing
)

func (f provFilter) label() string {
	switch f {
	case filterImplemented:
		return "implemented"
	case filterMissing:
		return "not implemented"
	default:
		return "all"
	}
}

// provenanceModel reports, per source set, how many original cards have a vactrol
// card referencing them (via card.Provenance), and lets you drill into a set to
// see which originals are covered and by what.
type provenanceModel struct {
	sets   []provSet
	cursor int // selected set in the list

	inSet     bool // drilled into a set?
	detailCur int
	filter    provFilter

	width, height int
}

func newProvenanceModel(w, h int) provenanceModel {
	// (source-set slug, collector number) -> implementing vactrol card names.
	impl := map[string]map[int][]string{}
	for _, rc := range card.Cards() {
		for _, ref := range rc.Provenance {
			byNum := impl[ref.Set.Slug]
			if byNum == nil {
				byNum = map[int][]string{}
				impl[ref.Set.Slug] = byNum
			}
			byNum[ref.Number] = append(byNum[ref.Number], rc.Def.Name)
		}
	}

	var sets []provSet
	for _, s := range provenance.Sets() {
		ps := provSet{set: s.SourceSet}
		for _, c := range s.Cards {
			names := impl[s.Slug][c.Number]
			if len(names) > 0 {
				ps.implemented++
			}
			ps.cards = append(ps.cards, provStatus{card: c, impl: names})
		}
		sets = append(sets, ps)
	}
	return provenanceModel{sets: sets, width: w, height: h}
}

func (m provenanceModel) resize(w, h int) provenanceModel { m.width, m.height = w, h; return m }

func (m provenanceModel) Update(msg tea.Msg) (provenanceModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if m.inSet {
		return m.updateDetail(key), nil
	}
	switch key.String() {
	case "esc", "q":
		return m, gotoScreen(screenMenu)
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "l":
		if m.cursor < len(m.sets)-1 {
			m.cursor++
		}
	case "enter", " ":
		if len(m.sets) > 0 {
			m.inSet = true
			m.detailCur = 0
			m.filter = filterAll
		}
	}
	return m, nil
}

func (m provenanceModel) updateDetail(key tea.KeyMsg) provenanceModel {
	rows := m.filtered()
	switch key.String() {
	case "esc", "q":
		m.inSet = false
	case "up", "k":
		if m.detailCur > 0 {
			m.detailCur--
		}
	case "down", "l":
		if m.detailCur < len(rows)-1 {
			m.detailCur++
		}
	case "f":
		m.filter = (m.filter + 1) % 3
		m.detailCur = 0
	}
	return m
}

// filtered returns the selected set's cards under the active filter.
func (m provenanceModel) filtered() []provStatus {
	set := m.sets[m.cursor]
	if m.filter == filterAll {
		return set.cards
	}
	var out []provStatus
	for _, c := range set.cards {
		done := len(c.impl) > 0
		if (m.filter == filterImplemented) == done {
			out = append(out, c)
		}
	}
	return out
}

func (m provenanceModel) View() string {
	if m.inSet {
		return m.detailView()
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("Provenance") + "\n")
	b.WriteString(faintStyle.Render("coverage of original cards, by set") + "\n\n")
	for i, s := range m.sets {
		total := len(s.cards)
		pct := 0
		if total > 0 {
			pct = s.implemented * 100 / total
		}
		line := fmt.Sprintf("%-22s %4d / %-4d  %3d%%", s.set.Name, s.implemented, total, pct)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(cursor(true)+line) + "\n")
		} else {
			b.WriteString(cursor(false) + line + "\n")
		}
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ move · enter open · esc/q back"))
	return b.String()
}

func (m provenanceModel) detailView() string {
	set := m.sets[m.cursor]
	rows := m.filtered()

	var b strings.Builder
	head := fmt.Sprintf("%s — %d/%d implemented", set.set.Name, set.implemented, len(set.cards))
	b.WriteString(titleStyle.Render("Provenance") + "  " + headerStyle.Render(head) + "\n")
	b.WriteString(faintStyle.Render("filter: "+m.filter.label()) + "\n\n")

	// Scroll window keeping the cursor visible.
	visible := m.height - 6
	if visible < 3 {
		visible = 3
	}
	start := 0
	if m.detailCur >= visible {
		start = m.detailCur - visible + 1
	}
	end := start + visible
	if end > len(rows) {
		end = len(rows)
	}

	for i := start; i < end; i++ {
		c := rows[i]
		line := fmt.Sprintf("#%03d %-26s", c.card.Number, c.card.Name)
		if len(c.impl) > 0 {
			line += " → " + strings.Join(c.impl, ", ")
		}
		switch {
		case i == m.detailCur:
			b.WriteString(selectedStyle.Render(cursor(true)+line) + "\n")
		case len(c.impl) > 0:
			b.WriteString(cursor(false) + line + "\n")
		default:
			b.WriteString(cursor(false) + faintStyle.Render(line) + "\n")
		}
	}
	if len(rows) == 0 {
		b.WriteString(faintStyle.Render("  (none)") + "\n")
	}

	b.WriteString("\n" + helpStyle.Render("↑/↓ scroll · f filter · esc back"))
	return b.String()
}
