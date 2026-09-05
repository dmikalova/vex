package cards

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/deckgen"
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
		case c.Type == engine.TypeUnset:
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

// TestMaterializedNamesAreUnique extends the duplicate-name check in
// TestAllIsAValidDatabase past materialization: a template (e.g. Master of X)
// never appears in a deck under its own registered name, only under whatever
// name its Materializer gives the concrete variant it produces (Master of 1..5),
// so those variant names must be just as unique as any registered card's name —
// both against every other registered name and against every other template's
// variants. It samples many (House, seed) combinations per template so a
// randomized Materializer (like Master of X rolling its power) exercises its
// whole output range instead of just whatever it happens to produce once.
func TestMaterializedNamesAreUnique(t *testing.T) {
	regs := card.Cards()

	// owner tracks who a name belongs to: a registered card by its own name, or a
	// template by the name it was found under.
	owner := make(map[string]string, len(regs))
	for _, rc := range regs {
		owner[rc.Def.Name] = rc.Def.Name
	}

	const samplesPerHouse = 20
	for _, rc := range regs {
		if rc.Materializer == nil {
			continue
		}
		for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
			for seed := int64(0); seed < samplesPerHouse; seed++ {
				ctx := deckgen.SlotContext{House: h, Rarity: rc.Def.Rarity}
				out := rc.Materializer.Materialize(ctx, rand.New(rand.NewSource(seed)))
				if out.Name == rc.Def.Name {
					// The template face's own name never reaches a deck (materialize
					// always replaces it), so it does not need to out-compete itself.
					continue
				}
				if by, ok := owner[out.Name]; ok && by != rc.Def.Name {
					t.Errorf(
						"%s materializes a variant named %q, which collides with %s",
						rc.Def.Name, out.Name, by,
					)
					continue
				}
				owner[out.Name] = rc.Def.Name
			}
		}
	}
}

// knownDuplicateImplementations are card name pairs (order does not matter)
// confirmed to be intentional identical-behavior reprints under a different
// name — real KeyForge's own equivalent of Fogbank/Foggify — rather than a
// copy-paste mistake. TestNoDuplicateImplementations allows these and flags
// every other collision, so a new hit means a genuinely new pair to triage:
// add it here once confirmed intentional, otherwise fix the card.
var knownDuplicateImplementations = [][2]string{
	// Both "Play: Archive a card." (Logos Common / Shadows Uncommon).
	{"Labwork", "Hidden Stash"},
}

func isKnownDuplicateImplementation(a, b string) bool {
	for _, pair := range knownDuplicateImplementations {
		if (pair[0] == a && pair[1] == b) || (pair[0] == b && pair[1] == a) {
			return true
		}
	}
	return false
}

// TestNoDuplicateImplementations flags two differently named cards whose whole
// behavior — type, stats, keywords, and abilities — is otherwise identical,
// ignoring rarity and the house/name that legitimately differ between a card and
// its reprint. Real KeyForge does this on purpose (Fogbank, Untamed Uncommon, and
// Foggify, Logos Common, both read "Your opponent cannot use creatures to fight
// on their next turn"), so a hit here is not automatically a bug — but outside
// knownDuplicateImplementations it is exactly the shape a copy-pasted card that
// forgot to change its actual effect would take, so it fails until triaged.
func TestNoDuplicateImplementations(t *testing.T) {
	all := All()
	seenBy := make(map[string]string, len(all))
	for _, c := range all {
		sig := c
		sig.Name = ""
		sig.House = engine.HouseNone
		sig.Rarity = ""
		key := fmt.Sprintf("%#v", sig)
		if owner, ok := seenBy[key]; ok {
			if !isKnownDuplicateImplementation(owner, c.Name) {
				t.Errorf(
					"%q and %q have identical implementations (same type, stats, keywords, and abilities); confirm this is an intentional reprint like Fogbank/Foggify, then add it to knownDuplicateImplementations",
					owner,
					c.Name,
				)
			}
			continue
		}
		seenBy[key] = c.Name
	}
}

// TestProvenanceHasNoOverlap checks that no two differently named cards claim the
// same original source printing (set + collector number). Provenance is
// repeatable on a single card (a reprint across sets shares one implementation,
// e.g. Fogbank's four printings), but two different vactrol cards pointing at
// the same original is an authoring mistake — most likely a Provenance tag
// copied from another card's file and never updated.
func TestProvenanceHasNoOverlap(t *testing.T) {
	seenBy := make(map[provenance.Ref]string)
	for _, rc := range card.Cards() {
		for _, ref := range rc.Provenance {
			if owner, ok := seenBy[ref]; ok && owner != rc.Def.Name {
				t.Errorf(
					"%s and %s both claim provenance %s #%d",
					owner, rc.Def.Name, ref.Set.Code, ref.Number,
				)
				continue
			}
			seenBy[ref] = rc.Def.Name
		}
	}
}

// TestEveryReprintResolvesToACard checks that every reprint claim a set makes (its
// 0set.go card.Reprint entries) names a card some set actually implements. A claim
// that resolves to nothing is a stale 0set.go — a card renamed or removed without
// regenerating — which would silently shrink that set's pool; here it fails the
// gate loudly and names the offending claim instead.
func TestEveryReprintResolvesToACard(t *testing.T) {
	byName := make(map[string]bool)
	for _, rc := range card.Cards() {
		byName[normalizeName(rc.Def.Name)] = true
	}
	for _, rp := range card.ReprintRefs() {
		if !byName[normalizeName(rp.Name)] {
			t.Errorf(
				"%s reprint #%d refers to %q, which no set implements "+
					"(stale 0set.go — regenerate with `mage tool:stub %s`)",
				rp.Set.Name, rp.Number, rp.Name, rp.Set.Slug,
			)
		}
	}
}
