package tui

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// fields finds the table row whose first column equals key and returns it as
// space-separated tokens, so assertions don't depend on the exact padding.
func fields(lines []string, key string) string {
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) > 0 && f[0] == key {
			return strings.Join(f, " ")
		}
	}
	return ""
}

func TestRenderSetStats(t *testing.T) {
	set := cards.Set{
		Name: "Test Set",
		Cards: []engine.CardDefinition{
			engine.NewCard("A", engine.Brobnar, engine.Creature, engine.Common),
			engine.NewCard("B", engine.Brobnar, engine.Tactic, engine.Common),
			engine.NewCard("C", engine.Brobnar, engine.Artifact, engine.Common),
			engine.NewCard("D", engine.Logos, engine.Creature, engine.Common),
		},
	}
	out := renderSetStats(set)
	lines := strings.Split(out, "\n")

	if !strings.Contains(out, "Test Set") || !strings.Contains(out, "(4 cards)") {
		t.Errorf("missing set name/count:\n%s", out)
	}
	// Columns are: House Total Creature Action Artifact Upgrade.
	if got := fields(lines, "Brobnar"); got != "Brobnar 3 1 1 1 0" {
		t.Errorf("Brobnar row = %q", got)
	}
	if got := fields(lines, "Logos"); got != "Logos 1 1 0 0 0" {
		t.Errorf("Logos row = %q", got)
	}
	if got := fields(lines, "Total"); got != "Total 4 2 1 1 0" {
		t.Errorf("Total row = %q", got)
	}
}
