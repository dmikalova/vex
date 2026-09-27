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

// iconFallbackAllowed once listed effect types the Iconography pass rendered as
// the abstract "unknown" glyph. It is gone: an unmapped mechanic is a bug, not an
// exemption. The totality tests below fail unconditionally on the fallback, so a
// new mechanic ships only once a glyph is drawn for it — the fix is always to add
// a glyph, never to allowlist the fallback (ADR 0022).

// forEachAbilityEffect walks every triggered-ability effect on every card the
// gallery shows (materialized variants included) and calls fn with it.
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

// TestIconTotality walks every card and fails loud if any triggered ability's
// effect falls back to the abstract glyph instead of a real mapping. A new
// mechanic with no glyph decision fails here: the fix is to add a glyph, never to
// tolerate the fallback — there is no allowlist to exempt it (ADR 0022).
func TestIconTotality(t *testing.T) {
	forEachAbilityEffect(func(e engine.Effect) {
		if _, covered := effectGlyphs(e); !covered {
			t.Errorf("effect %s has no glyph mapping; add one (ADR 0022)", effectTypeName(e))
		}
	})
}

// familyName is a glyph family's function name, so a disjointness failure names
// the families that clashed rather than their positions in the chain.
func familyName(family func(engine.Effect) ([]glyph, bool, bool)) string {
	name := runtime.FuncForPC(reflect.ValueOf(family).Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

// TestGlyphFamiliesAreDisjoint walks every triggered-ability effect on every card
// and fails unless exactly one family in glyphFamilies claims it. The chain is not
// compiler-checked: two families listing the same effect type would make the
// second dead code and draw a wrong glyph rather than an unknown one, which no
// other test catches (ADR 0047). Zero claims is a failure too — an effect the
// chain drops renders the abstract unknown glyph.
func TestGlyphFamiliesAreDisjoint(t *testing.T) {
	forEachAbilityEffect(func(e engine.Effect) {
		var claimed []string
		for _, family := range glyphFamilies {
			if _, _, ok := family(e); ok {
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

// TestNoResidualUnknownGlyph renders every card's full Icon strip and fails if any
// glyph is the abstract unknown — in a composed line's glyphs or in its trigger
// heads. Unlike TestIconTotality — which only reads the top-level covered flag —
// this walks the composed lines, so an unknown buried inside a wrapper (a
// ChooseHouseThen's inner effect, a granted ability) or at a trigger head that a
// masking covered=true would hide is still caught (ADR 0022).
func TestNoResidualUnknownGlyph(t *testing.T) {
	regs := card.Cards()
	for i := range regs {
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			for _, line := range cardGlyphs(&def) {
				for _, tr := range line.triggers {
					if tr == "glyph-unknown" {
						t.Errorf("%s renders an unknown trigger glyph in its Icon strip", def.Name)
					}
				}
				for _, g := range line.glyphs {
					if g.asset == "glyph-unknown" {
						t.Errorf("%s renders an unknown glyph in its Icon strip", def.Name)
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

// TestDistinctMechanicsDrawDistinctGlyphs is the verb axis's gate: it walks every
// triggered ability on every card (the corpus TestIconTotality walks), renders each
// effect, groups the effects by the strip they drew, and fails when two distinct
// leaf mechanics land in one group. Without it the next agent to add an effect node
// could quietly reuse another mechanic's picture, which is how twelve pairs came to
// render identically before ADR 0047.
//
// The composition family is excluded: a Conditional, May, ForEach, Repeat, Then or
// Sequence is *supposed* to render as the effect it wraps. Everything else must
// either draw its own strip or be a declared synonym (glyphSynonyms).
func TestDistinctMechanicsDrawDistinctGlyphs(t *testing.T) {
	byStrip := map[string]map[string]bool{}
	forEachAbilityEffect(func(e engine.Effect) {
		if _, _, ok := composeEffectGlyphs(e); ok {
			return
		}
		gs, _ := effectGlyphs(e)
		strip := fmt.Sprintf("%+v", gs)
		if byStrip[strip] == nil {
			byStrip[strip] = map[string]bool{}
		}
		byStrip[strip][effectTypeName(e)] = true
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
