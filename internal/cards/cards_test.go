package cards

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// TestAllIsAValidDatabase checks the assembled card database: every set package
// is imported, every card self-registers, and each entry is well-formed with a
// unique name, a real house, and a rarity.
func TestAllIsAValidDatabase(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("All() returned no cards; are the set packages imported?")
	}

	seen := make(map[string]bool, len(all))
	for _, c := range all {
		switch {
		case c.Name == "":
			t.Error("card with empty name")
		case c.House == engine.HouseNone:
			t.Errorf("%q has no house", c.Name)
		case c.Rarity == "":
			t.Errorf("%q has no rarity", c.Name)
		case c.Type == "":
			t.Errorf("%q has no card type", c.Name)
		}
		if seen[c.Name] {
			t.Errorf("duplicate card name %q", c.Name)
		}
		seen[c.Name] = true
	}
}

// TestEveryCreatureAndArtifactHasTrait enforces the card-database policy that
// every creature and artifact carries at least one trait (e.g. Giant, Beast,
// Weapon). Actions and upgrades are exempt.
func TestEveryCreatureAndArtifactHasTrait(t *testing.T) {
	for _, c := range All() {
		if c.Type != engine.Creature && c.Type != engine.Artifact {
			continue
		}
		if len(c.Traits) == 0 {
			t.Errorf(
				"%s (%s) has no trait; every creature and artifact needs at least one",
				c.Name,
				c.Type,
			)
		}
	}
}

// TestAllIsSortedDeterministically verifies the database comes back in a stable
// order (house, then name), independent of package initialization order.
func TestAllIsSortedDeterministically(t *testing.T) {
	all := All()
	for i := 1; i < len(all); i++ {
		prev, cur := all[i-1], all[i]
		if prev.House > cur.House || (prev.House == cur.House && prev.Name > cur.Name) {
			t.Errorf("not sorted at %d: %s (%s) before %s (%s)",
				i, prev.Name, prev.House, cur.Name, cur.House)
		}
	}
}
