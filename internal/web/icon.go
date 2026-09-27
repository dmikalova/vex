package web

import (
	"reflect"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is the Iconography pass (ADR 0022): the Visitor that turns a card's
// mechanics into the composed Glyphs of its Icon strip, the visual counterpart of
// rules-text generation. It type-switches over the engine's exported Effect AST
// and Targets. It lives here in the client, not in package engine, because the
// engine is held at 100% coverage and an in-engine Visitor would force an icon
// test for every one of the ~150 effect nodes; internal/web is ungated. The
// totality test (icon_test.go) still walks every card and fails loud on any
// effect that falls back to the abstract unknown glyph. There is no allowlist to
// exempt one: an unmapped mechanic is a bug, and the fix is always to draw a
// glyph, so a new mechanic cannot ship without a glyph decision.

// decor is a set of edge/overlay treatments applied to a noun glyph to carry a
// Target's shape without spending a horizontal slot: an enemy tint, a friendly
// tint, an "each" stack, a "chosen" outline, or a "this/self" marker.
type decor uint8

const (
	decorEnemy decor = 1 << iota
	decorFriendly
	decorEach
	decorChosen
	decorThis
)

// glyph is one composed icon in the strip. A glyph with an asset renders that SVG
// (tinted and decorated); a glyph with no asset renders its text as a small chip,
// which is how a not-yet-drawn mechanic still shows something readable. Qty, when
// positive, prints a numeral badge on the glyph (the "3" on 3 damage).
type glyph struct {
	asset string
	text  string
	qty   int
	decor decor
	arrow bool // render a leading result-gate arrow (→) before this glyph
}

// glyphLine is one ability's transcription: its trigger glyph(s) and the glyphs
// its effect composes into. Everything is an icon — no words — so triggers render
// as their own glyph too (ADR 0022). Adjacent abilities that share one effect and
// fire on distinct action triggers (Play/Fight/Reap) merge into one line carrying
// all their trigger glyphs, mirroring the "Play/Fight/Reap:" rules-text shorthand.
type glyphLine struct {
	triggers []string // trigger glyph asset stems, in canonical order
	glyphs   []glyph
	covered  bool // false when the effect fell back to the abstract unknown glyph
}

// cardGlyphs is the Iconography pass: it renders a card's triggered abilities and
// its continuous rules — keywords, static bonuses, restrictions, tolls, key-cost
// and Æmber-flow replacements, and card-level flags — as glyph lines, so a card
// whose only mechanic is a continuous rule still draws a strip.
func cardGlyphs(def *engine.CardDefinition) []glyphLine {
	lines := make([]glyphLine, 0, len(def.Abilities)+len(def.ConstantAbilities)+2)
	if kw := keywordGlyphs(def); len(kw) > 0 {
		lines = append(lines, glyphLine{
			glyphs:  kw,
			covered: true,
		})
	}
	if def.FightRestriction != (engine.Target{}) {
		lines = append(lines, glyphLine{
			glyphs:  fightRestrictionGlyphs(def.FightRestriction),
			covered: true,
		})
	}
	if def.DrawModifier.Amount != 0 {
		lines = append(lines, glyphLine{
			glyphs:  drawModifierGlyphs(def.DrawModifier),
			covered: true,
		})
	}
	lines = append(lines, staticLines(def.Static)...)
	lines = append(lines, restrictionLines(def.Restricts)...)
	if def.Replaces != (engine.Instead{}) {
		lines = append(lines, glyphLine{
			glyphs:  replaceGlyphs(def.Replaces),
			covered: true,
		})
	}
	if len(def.KeyCostChanges) > 0 {
		lines = append(lines, glyphLine{
			glyphs:  keyCostChangeGlyphs(),
			covered: true,
		})
	}
	lines = append(lines, cardFeatureLines(def)...)
	for i := 0; i < len(def.Abilities); {
		ab := def.Abilities[i]
		gs, covered := effectGlyphs(ab.Effect)
		triggers := []string{triggerIcon(ab.Trigger)}
		j := i + 1
		if isActionTrigger(ab.Trigger) {
			text := ab.Effect.Text()
			for j < len(def.Abilities) &&
				isActionTrigger(def.Abilities[j].Trigger) &&
				def.Abilities[j].Effect.Text() == text {
				triggers = append(triggers, triggerIcon(def.Abilities[j].Trigger))
				j++
			}
		}
		lines = append(lines, glyphLine{
			triggers: triggers,
			glyphs:   gs,
			covered:  covered,
		})
		i = j
	}
	for _, ca := range def.ConstantAbilities {
		lines = append(lines, constantLines(ca)...)
	}
	return lines
}

// glyphFamilies is the dispatch chain the Iconography pass runs an effect down:
// each family claims the effect types of one mechanic domain, coarsely mirroring
// the engine's effect_<mechanic>.go split, and the first family to claim an effect
// wins (ADR 0047). A family returns (glyphs, covered, ok).
//
// The third value is load-bearing and must not be collapsed into covered, because
// covered = false already means something else: it says "this drew the abstract
// unknown glyph", and four cases return real glyphs with covered = false — a Then
// whose First is not an Effect, and PutFromPlay / PutChosen / PutCard with a
// destination that has no glyph. A family therefore cannot signal "not mine" by
// returning covered = false; ok is how it declines.
//
// The compiler cannot see across the chain, so two families claiming the same
// effect type would silently make the second dead code — drawing a wrong glyph
// rather than an unknown one. TestGlyphFamiliesAreDisjoint polices that.
//
// The chain is filled in init rather than at its declaration because every family
// recurses through effectGlyphs, which reads the slice, and Go rejects that as an
// initialization cycle.
var glyphFamilies []func(engine.Effect) (gs []glyph, covered, ok bool)

func init() {
	glyphFamilies = []func(engine.Effect) (gs []glyph, covered, ok bool){
		aemberEffectGlyphs,
		damageEffectGlyphs,
		creatureEffectGlyphs,
		statsEffectGlyphs,
		zoneEffectGlyphs,
		revealEffectGlyphs,
		boardEffectGlyphs,
		turnEffectGlyphs,
		composeEffectGlyphs,
	}
}

// effectGlyphs renders one effect to its glyphs, reporting whether the effect was
// covered by a real mapping (false means it fell back to the abstract unknown
// glyph). It runs the effect down glyphFamilies and takes the first claim.
func effectGlyphs(e engine.Effect) ([]glyph, bool) {
	for _, family := range glyphFamilies {
		if gs, covered, ok := family(e); ok {
			return gs, covered
		}
	}
	return fallbackGlyphs(e), false
}

// claim adapts a two-valued glyph rendering — a recursive effectGlyphs, a
// composeGlyphs over sub-effects — into a family's three-valued return, saying
// the family claimed the effect.
func claim(gs []glyph, covered bool) ([]glyph, bool, bool) {
	return gs, covered, true
}

// mustCompose renders sub-effects like composeGlyphs but discards the covered
// flag: the caller has already prefixed a real glyph, so the line reads even when
// a nested effect falls back to the abstract glyph.
func mustCompose(effects ...engine.Effect) []glyph {
	gs, _ := composeGlyphs(effects...)
	return gs
}

// composeGlyphs renders a run of sub-effects one after another, reporting covered
// only when every sub-effect maps to a real glyph.
func composeGlyphs(effects ...engine.Effect) ([]glyph, bool) {
	out := make([]glyph, 0, len(effects)*2)
	covered := true
	for _, e := range effects {
		gs, c := effectGlyphs(e)
		out = append(out, gs...)
		covered = covered && c
	}
	return out, covered
}

// fallbackGlyphs renders an effect the pass does not yet map as the single
// abstract "unknown" glyph. The strip is pure icons — no words — so an unmapped
// mechanic shows a placeholder sigil rather than its printed text (ADR 0022).
func fallbackGlyphs(engine.Effect) []glyph {
	return []glyph{{asset: "glyph-unknown"}}
}

// effectTypeName is the effect's concrete Go type name (e.g. "engine.DealDamage"),
// used by the totality test to name an uncovered effect.
func effectTypeName(e engine.Effect) string {
	return reflect.TypeOf(e).String()
}
