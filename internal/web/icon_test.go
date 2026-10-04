package web

import (
	"fmt"
	"maps"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/engine"
)

// The Iconography pass once drew an abstract "unknown" glyph for a node it could
// not map, with an allowlist of the effect types that were allowed to reach it.
// Both are gone: every rule-bearing row of every gated family draws a real,
// distinct glyph, so an unmapped node is a bug rather than an exemption. The
// shipped pass draws nothing for one, and these tests make the hole loud.

// init makes an unmapped node panic for the whole package's tests, which is the
// test half of the split unmappedGlyphsPanic names (icon.go): the shipped wasm
// binary draws nothing, because a panic there is a blank page, while a test run
// stops on the node that has no glyph rather than quietly rendering a gap.
func init() { unmappedGlyphsPanic = true }

// forEachAbilityEffect walks every triggered-ability effect on every card the
// gallery shows (materialized variants included) and calls fn with it.
//
// It is the second population beside the census (icon_catalog.go), not a
// leftover the census replaces. The census holds one bare node per member,
// filled only as far as it needs to render, so it never forms the composed
// shapes a real card does — a ChooseHouseThen wrapping a granted ability, a
// Conditional over a nested sequence — and those compose through the family
// chain differently from the nodes they wrap. Deleting this walk as redundant
// would drop that coverage.
func forEachAbilityEffect(fn func(e engine.Effect)) {
	regs := card.Cards()
	for i := range regs {
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			for _, ab := range def.Abilities {
				fn(ab.Effect)
			}
		}
	}
}

// TestIconTotality fails loud if a rule-bearing mechanic draws nothing. A new
// mechanic with no glyph decision fails here: the fix is to add a glyph, never to
// tolerate the gap — there is no allowlist to exempt it (ADR 0022).
//
// It runs over two populations. The census (icon_catalog.go) is the contract: it
// holds every rule-bearing node internal/engine declares, whether or not a card
// uses one yet, so a node's glyph is owed the moment the node exists. The card
// walk is kept beside it for what the census cannot reach — the continuous rules
// cardGlyphs draws (keywords, a FightRestriction, a DrawModifier, Static bonuses,
// Restricts, Replaces, KeyCostChanges, ConstantAbilities, trigger heads), none of
// which is an Effect, and the composed shapes a real card forms (a ChooseHouseThen
// wrapping a granted ability) that a bare catalogued node, filled only as far as
// it needs to render, never forms.
func TestIconTotality(t *testing.T) {
	forEachGatedIconRow(func(family iconFamily, row iconRow) {
		if len(row.draw()) == 0 {
			t.Errorf(
				"%s node %s owes the term %q and draws no glyph; add one (ADR 0022)",
				family.name,
				row.name,
				row.rules.Term,
			)
		}
	})
	forEachAbilityEffect(func(e engine.Effect) {
		if len(effectGlyphs(e)) == 0 {
			t.Errorf("effect %s draws no glyph; add one (ADR 0022)", effectTypeName(e))
		}
	})
}

// TestDurationGlyphCoversEveryDuration walks engine.Durations() and fails on a
// window that draws nothing. inWindow appends whatever durationGlyph returns, so
// an unmapped window would append an assetless glyph — dropped silently at
// render, which is the card-level counterpart of the empty run the effect level
// already fails on (ADR 0022). The census guards Durations() complete
// against the source, so a new constant cannot hide from this test.
func TestDurationGlyphCoversEveryDuration(t *testing.T) {
	for _, d := range engine.Durations() {
		name := durationGlyph(d)
		if name == "" {
			t.Errorf("duration %v draws no glyph; give it a marker (ADR 0022)", d)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("duration %v: no web/assets/%s.svg", d, name)
		}
	}
}

// familyName is a glyph family's function name, so a disjointness failure names
// the families that clashed rather than their positions in the chain.
func familyName(family func(engine.Effect) ([]glyph, bool)) string {
	name := runtime.FuncForPC(reflect.ValueOf(family).Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

// TestGlyphFamiliesAreDisjoint fails unless exactly one family in glyphFamilies
// claims each effect. The chain is not compiler-checked: two families listing the
// same effect type would make the second dead code and draw a wrong glyph, which
// no other test catches (ADR 0047). Zero claims is a failure too — an effect the
// chain drops is unmapped, and draws nothing.
//
// Like TestIconTotality it runs over both populations: the census, which holds
// every effect node the engine declares, and the card walk, which alone reaches
// the composed shapes a real card forms (a ChooseHouseThen wrapping a granted
// ability) that a bare catalogued node never forms.
func TestGlyphFamiliesAreDisjoint(t *testing.T) {
	forEachCatalogedEffect(func(_ iconRow, e engine.Effect) {
		var claimed []string
		for _, family := range glyphFamilies {
			if _, ok := family(e); ok {
				claimed = append(claimed, familyName(family))
			}
		}
		if len(claimed) != 1 {
			t.Errorf(
				"catalogued effect %s claimed by families %v, want exactly one (ADR 0047)",
				effectTypeName(e),
				claimed,
			)
		}
	})
	forEachAbilityEffect(func(e engine.Effect) {
		var claimed []string
		for _, family := range glyphFamilies {
			if _, ok := family(e); ok {
				claimed = append(claimed, familyName(family))
			}
		}
		if len(claimed) != 1 {
			t.Errorf(
				"effect %s claimed by families %v, want exactly one (ADR 0047)",
				effectTypeName(e),
				claimed,
			)
		}
	})
}

// TestEveryCardLineDraws renders every card's full Icon strip through cardGlyphs
// and fails on a line that draws nothing. A line exists because the card has a
// mechanic, so a line with no glyphs — or a glyph or trigger head naming no
// asset, which the renderer drops silently — is a mechanic the player is shown
// nothing for (ADR 0022).
//
// It is the only population that reaches cardGlyphs, and so the only one that
// reaches a card's continuous rules: its keywords, FightRestriction,
// DrawModifier, Static bonuses, Restricts, Replaces, KeyCostChanges and
// ConstantAbilities, plus the trigger head of every ability. None of those is an
// Effect, so no row of the effect census reaches them. It also reaches the
// composed shapes a real card forms, where a bare catalogued node does not, so a
// gap buried inside a wrapper — a ChooseHouseThen's inner effect, a granted
// ability — is caught here.
func TestEveryCardLineDraws(t *testing.T) {
	regs := card.Cards()
	for i := range regs {
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			for _, line := range cardGlyphs(&def) {
				if len(line.glyphs) == 0 {
					t.Errorf("%s draws an Icon line with no glyphs (ADR 0022)", def.Name)
				}
				for _, tr := range line.triggers {
					if tr == "" {
						t.Errorf(
							"%s draws an Icon line whose trigger head names no asset",
							def.Name,
						)
					}
				}
				for _, g := range line.glyphs {
					if g.asset == "" {
						t.Errorf("%s draws an Icon glyph naming no asset: %+v", def.Name, g)
					}
				}
			}
		}
	}
}

// glyphSynonyms declares the effect nodes that are the same mechanic under two
// node names and so *should* draw the same strip. This is a positive assertion,
// not an allowlist: an undeclared collision is a bug in the vocabulary, the way an
// unmapped mechanic is a bug rather than an exemption (ADR 0047). Nothing here is
// a way to opt out of drawing a distinct glyph — a node whose mechanic differs
// from every node listed beside it does not belong in a group.
var glyphSynonyms = []struct {
	reason string
	types  []string
}{
	{
		reason: "three spellings of destroying creatures; the target does the varying",
		types: []string{
			"engine.Destroy",
			"engine.DestroyChosen",
			"engine.BatchDestroy",
		},
	},
	{
		reason: "readying one creature and readying several are the same verb",
		types: []string{
			"engine.Ready",
			"engine.ReadyCreatures",
		},
	},
	{
		reason: "three spellings of putting a card into play; only the origin differs",
		types: []string{
			"engine.PlayFrom",
			"engine.PlayTopOfDeck",
			"engine.PutIntoPlay",
		},
	},
	{
		reason: "archiving from hand and archiving from play are the same move",
		types: []string{
			"engine.ArchiveCard",
			"engine.ArchiveFromPlay",
		},
	},
	{
		reason: "the general node catalogued at the narrow node's own values: PutCard's " +
			"census row moves a card from the discard pile into the hand, which is " +
			"precisely what PutDiscardedIntoHand does. Forcing them apart would mean " +
			"drawing a difference the two catalogued values do not have",
		types: []string{
			"engine.PutCard",
			"engine.PutDiscardedIntoHand",
		},
	},
	{
		reason: "the Filter census catalogues its type axis as an artifact " +
			"(catalog_filter.go), so that row renders \"each artifact\" — the very noun " +
			"the count of artifacts in play draws. One names those cards to act on " +
			"them and the other to count them, and drawing a difference would be " +
			"drawing one the two catalogued values do not have",
		types: []string{
			"engine.ArtifactsInPlay",
			"engine.Filter",
		},
	},
	{
		reason: "a count of a set of cards and a selection of each of them name the same " +
			"cards on purpose: pickedCardGlyphs narrows in the order " +
			"cardsInPlayCountGlyphs does so a pick and a count of the same cards read " +
			"alike (icon_selection.go)",
		types: []string{
			"engine.CardsInPlay",
			"engine.Each",
		},
	},
}

// synonymReason reports the declared reason a set of effect type names is allowed
// to share one strip, and whether any declared group covers all of them.
func synonymReason(types []string) (string, bool) {
	for _, group := range glyphSynonyms {
		all := true
		for _, t := range types {
			if !slices.Contains(group.types, t) {
				all = false
				break
			}
		}
		if all {
			return group.reason, true
		}
	}
	return "", false
}

// rowTypeName is a catalogued row's Go type name ("engine.Destroy"), which is how
// a row is named to glyphSynonyms and how one node catalogued by two families is
// recognised as one mechanic rather than two: engine.Shuffle is both an Effect and
// a TopAct, and engine.CardsInPlay is one Go value in CountCatalog and
// ConditionCatalog. A row whose family catalogues no value falls back to its
// census name.
func rowTypeName(family iconFamily, row iconRow) string {
	if row.node == nil {
		return family.name + "." + row.name
	}
	return reflect.TypeOf(row.node).String()
}

// TestDistinctMechanicsDrawDistinctGlyphs is the verb axis's gate: it renders
// every rule-bearing row of every gated family, groups the rows by the strip they
// drew, and fails when two distinct mechanics land in one group. Without it the
// next agent to add a node could quietly reuse another mechanic's picture, which
// is how twelve pairs came to render identically before ADR 0047.
//
// It runs over the census rather than the card walk, and over all nine gated
// families at once rather than family by family, because the census holds every
// node the engine declares whether or not a card uses one yet, and because the
// grammar — not the family a node happens to live in — is what carries the
// distinctions. A Condition and an Effect resting on the same noun are therefore
// compared, deliberately: the global assertion is the stronger claim.
//
// The composition family is excluded: a Conditional, May, ForEach, Repeat, Then or
// Sequence is *supposed* to render as the effect it wraps. Everything else must
// either draw its own strip or be a declared synonym (glyphSynonyms).
func TestDistinctMechanicsDrawDistinctGlyphs(t *testing.T) {
	byStrip := map[string]map[string]bool{}
	forEachGatedIconRow(func(family iconFamily, row iconRow) {
		if e, ok := row.node.(engine.Effect); ok {
			if _, claimed := composeEffectGlyphs(e); claimed {
				return
			}
		}
		strip := fmt.Sprintf("%+v", row.draw())
		if byStrip[strip] == nil {
			byStrip[strip] = map[string]bool{}
		}
		byStrip[strip][rowTypeName(family, row)] = true
	})
	for strip, set := range byStrip {
		if len(set) < 2 {
			continue
		}
		types := slices.Sorted(maps.Keys(set))
		if _, ok := synonymReason(types); ok {
			continue
		}
		t.Errorf(
			"%s draw the same strip %s; give one a distinct glyph or sigil, or declare "+
				"them synonyms in glyphSynonyms with the reason (ADR 0047)",
			strings.Join(types, " and "),
			strip,
		)
	}
}

// TestCardGlyphsMergesActionTriggers checks that a Play/Fight/Reap ability — three
// abilities sharing one effect on the three action triggers — renders as one line
// carrying all three trigger glyphs, the icon counterpart of the "Play/Fight/Reap:"
// rules-text shorthand, rather than three repeated lines.
func TestCardGlyphsMergesActionTriggers(t *testing.T) {
	effect := engine.Destroy{Target: engine.Target{Kind: engine.TargetChosenCreature}}
	def := &engine.CardDefinition{Abilities: []engine.Ability{
		{Trigger: engine.TriggerAfterPlay, Effect: effect},
		{Trigger: engine.TriggerAfterFight, Effect: effect},
		{Trigger: engine.TriggerAfterReap, Effect: effect},
	}}
	lines := cardGlyphs(def)
	if len(lines) != 1 {
		t.Fatalf("want one merged line, got %d", len(lines))
	}
	if got := lines[0].triggers; len(got) != 3 ||
		got[0] != "glyph-play" || got[1] != "glyph-fight" || got[2] != "glyph-reap" {
		t.Errorf("triggers = %v, want [glyph-play glyph-fight glyph-reap]", got)
	}
}

// TestCardGlyphsShowsDrawModifier checks a card whose only mechanic is a continuous
// hand-refill change (Mother: draw +1) still draws a glyph strip rather than none.
func TestCardGlyphsShowsDrawModifier(t *testing.T) {
	def := &engine.CardDefinition{
		DrawModifier: engine.DrawModifier{
			Player: engine.Controller,
			Amount: 1,
		},
	}
	lines := cardGlyphs(def)
	if len(lines) != 1 {
		t.Fatalf("want one draw-modifier line, got %d", len(lines))
	}
	if got := lines[0].glyphs; len(got) != 1 ||
		got[0].asset != "zone-hand" || got[0].qty != 1 ||
		got[0].decor&decorFriendly == 0 {
		t.Errorf("draw-modifier glyphs = %+v, want a friendly zone-hand +1", got)
	}
}
