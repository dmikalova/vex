package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// statsModel is a read-only screen that reports the card database broken down by
// set: how many cards each house has, in total and per card type.
type statsModel struct {
	sets   []cards.Set
	width  int
	height int
}

func newStatsModel(w, h int) statsModel {
	return statsModel{sets: cards.Sets(), width: w, height: h}
}

func (m statsModel) resize(w, h int) statsModel { m.width, m.height = w, h; return m }

func (m statsModel) Update(msg tea.Msg) (statsModel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "esc", "q":
		return m, gotoScreen(screenMenu)
	}
	return m, nil
}

// statsTypes is the fixed column order for the per-type breakdown.
var statsTypes = []engine.CardType{engine.Creature, engine.Tactic, engine.Artifact, engine.Upgrade}

func (m statsModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Card statistics") + "\n\n")

	for _, set := range m.sets {
		b.WriteString(renderSetStats(set) + "\n")
	}

	b.WriteString(helpStyle.Render("esc/q back"))
	return b.String()
}

// renderSetStats renders one set's house/type breakdown as an aligned table.
func renderSetStats(set cards.Set) string {
	byHouse := map[engine.House]map[engine.CardType]int{}
	houseTotal := map[engine.House]int{}
	typeTotal := map[engine.CardType]int{}
	var houses []engine.House
	for _, d := range set.Cards {
		if byHouse[d.House] == nil {
			byHouse[d.House] = map[engine.CardType]int{}
			houses = append(houses, d.House)
		}
		byHouse[d.House][d.Type]++
		houseTotal[d.House]++
		typeTotal[d.Type]++
	}
	sort.Slice(houses, func(i, j int) bool { return houses[i].String() < houses[j].String() })

	var b strings.Builder
	b.WriteString(headerStyle.Render(set.Name) +
		faintStyle.Render(fmt.Sprintf("  (%d cards)", len(set.Cards))) + "\n")

	// Column header.
	header := fmt.Sprintf("  %-8s%7s", "House", "Total")
	for _, ct := range statsTypes {
		header += fmt.Sprintf("%10s", string(ct))
	}
	b.WriteString(faintStyle.Render(header) + "\n")

	// One row per house present in the set.
	for _, h := range houses {
		row := fmt.Sprintf("  %-8s%7d", h.String(), houseTotal[h])
		for _, ct := range statsTypes {
			row += fmt.Sprintf("%10d", byHouse[h][ct])
		}
		b.WriteString(row + "\n")
	}

	// Totals row across all houses.
	total := fmt.Sprintf("  %-8s%7d", "Total", len(set.Cards))
	for _, ct := range statsTypes {
		total += fmt.Sprintf("%10d", typeTotal[ct])
	}
	b.WriteString(selectedStyle.Render(total) + "\n")

	return b.String()
}
