package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the damage-and-destroy
// family (ADR 0047) — dealing, redistributing and healing damage, armor loss, the
// forms of destruction, and the purge zone destruction feeds into.

// damageEffectGlyphs claims every damage-dealing and destruction effect.
func damageEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.DealDamage:
		// A follow-up on the damaged creature draws the two clauses joined, not the
		// bare damage glyph.
		if v.Then != nil {
			return claim(damageThenGlyphs(v.Amount, v.Target, v.Then))
		}
		// A Spread carries its own creature targets rather than filling Target, so it
		// draws its own summary noun; the counts and neighbor split stay in the text.
		if v.Spread != nil {
			return []glyph{{asset: "damage"}, arrowTo(spreadTargetGlyph())}, true, true
		}
		// The per-count variants hit several creatures or scale by a board count; the
		// numeral lives in the text, so the glyph drops the qty.
		if v.Per != nil || v.PerTarget != nil || v.AmountFrom != nil {
			return []glyph{{asset: "damage"}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{
			{asset: "damage", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.RedistributeDamage:
		return []glyph{{asset: "damage"}, {asset: "glyph-swap"}}, true, true
	case engine.Heal:
		return []glyph{
			{asset: "glyph-heal", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.LoseArmor:
		return []glyph{
			{asset: "shield"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.TakesExtraDamage:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{
			asset: "damage",
			qty:   v.Amount,
		})}, true, true
	case engine.CannotBeDealtDamage:
		return []glyph{{asset: "shield"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RedirectFightDamage:
		return []glyph{
			{asset: "glyph-fight"},
			{asset: "damage"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.DamageOthersAfterUsingTrait:
		return []glyph{
			traitGlyph(),
			{asset: "glyph-action"},
			{asset: "damage", qty: v.Amount},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.Destroy:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.DestroyChosen:
		return []glyph{{asset: "glyph-destroy"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.BatchDestroy:
		// The gatherer says which creatures the one batch takes, and the shapes differ
		// — two picks from two pools, or every unspared player's battleline — so the
		// strip draws the gatherer rather than a fixed "each creature" that would be
		// true of only one of them.
		gs, covered := gatherGlyphs(v.Gather)
		return append(
			[]glyph{{asset: "glyph-destroy"}},
			punctuate(gs, punctArrow)...,
		), covered, true
	case engine.DestroyEachCreatureAtEndOfTurn:
		// The turn marker is what tells this from BatchDestroy over every creature:
		// the same wipe, scheduled into the end-of-turn window rather than resolved
		// now (ADR 0047).
		return inWindow([]glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, engine.RemainderOfPlayerTurn), true, true
	case engine.PurgeCreature:
		return []glyph{{asset: "zone-purge"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.PurgeSource:
		return []glyph{{asset: "zone-purge", decor: decorThis}}, true, true
	default:
		return nil, false, false
	}
}

// gatherGlyphs renders a Gather as the creatures the one simultaneous batch
// takes, reporting whether a real mapping covered it. The Gather family is gated
// (TestEveryGatherDraws, TestIconTotality), so every row here draws.
//
// A gatherer is not an effect — it is the object of the one verb that carries one
// (BatchDestroy), so it draws where that verb draws, the way a Count does.
func gatherGlyphs(g engine.Gather) ([]glyph, bool) {
	switch v := g.(type) {
	case engine.ChosenFromEach:
		// One pick per pool, drawn as the pools themselves: "a friendly creature and an
		// enemy creature" is the two targets side by side, which is what makes this
		// gatherer different from a wipe.
		gs := make([]glyph, 0, len(v))
		for _, t := range v {
			gs = append(gs, targetGlyph(t))
		}
		return gs, true
	case engine.EachPlayerUnless:
		// Each player's creatures, behind the board the spare condition reads and the
		// refusal mark that says a player fielding it is skipped. Which creatures of an
		// unspared battleline are taken is a Refinement, a family with no glyphs yet,
		// so it stays in the rules text beside the strip.
		return append(
			[]glyph{creaturesCounted(0)},
			append(
				punctuate(cardsInPlayConditionGlyphs(v.Spare), punctGate),
				glyph{asset: "glyph-ban"},
			)...,
		), true
	}
	return fallbackGlyphs(), false
}

// damageThenGlyphs renders a "deal N damage to <target>, then <follow-up>" effect
// as the damage glyph arrowed to its target, followed by the follow-up's glyphs.
func damageThenGlyphs(amount int, target engine.Target, then engine.Effect) ([]glyph, bool) {
	gs := []glyph{{asset: "damage", qty: amount}, arrowTo(targetGlyph(target))}
	more, covered := effectGlyphs(then)
	return append(gs, more...), covered
}
