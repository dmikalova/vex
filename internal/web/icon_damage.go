package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the damage-and-destroy
// family (ADR 0047) — dealing, redistributing and healing damage, armor loss, the
// forms of destruction, and the purge zone destruction feeds into.
//
// Two sentence parts a DealDamage carries are drawn here too, beside the verb
// that holds them, the way a Count is drawn where its holder is: the Spread that
// says which related creatures one batch of damage strikes, and the PerTarget
// that says what the amount is measured against on each creature it lands on.
// Both are declared in the engine's effect_damage.go, so both are drawn in the
// icon file that mirrors it.

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
		// draws the shape it strikes; the amounts stay in the text.
		if v.Spread != nil {
			gs, covered := spreadGlyphs(v.Spread)
			return append([]glyph{{asset: "damage"}}, punctuate(gs, punctArrow)...), covered, true
		}
		// A per-target amount is read off each creature hit, so the thing it measures
		// trails the creature it is measured on: "1 damage to each enemy creature for
		// each Æmber on it".
		if v.PerTarget != nil {
			gs, covered := perTargetGlyphs(v.PerTarget)
			return append(
				[]glyph{{asset: "damage"}, arrowTo(targetGlyph(v.Target))},
				gs...,
			), covered, true
		}
		// The remaining per-count variants scale by a board count; the numeral lives
		// in the text, so the glyph drops the qty.
		if v.Per != nil || v.AmountFrom != nil {
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

// spreadGlyphs renders a Spread as the shape one batch of damage strikes,
// reporting whether a real mapping covered it. The Spread family is gated
// (TestEverySpreadDraws, TestIconTotality), so every row here draws.
//
// The five spreads all pick creatures, so what tells them apart is the shape of
// the pick rather than a drawing of its own: a creature and the battleline beside
// it, two separate picks, a bounded handful, an unbounded one, a flank and the
// walk inward. The amounts each creature takes stay in the rules text, as a
// Target's fine filters do — the strip says who is hit, the text says how hard.
func spreadGlyphs(s engine.Spread) ([]glyph, bool) {
	switch v := s.(type) {
	case engine.CreatureAndNeighbors:
		// The chosen creature, then the battleline beside it: the neighbors hang off
		// the pick the way the printed clause hangs them off "it".
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-flank", decor: decorChosen},
			{asset: "type-creature", decor: neighborScopeDecor(v.Scope)},
		}, true
	case engine.DifferentCreatures:
		// Two picks rather than one, which is the whole of what this spread says: the
		// second creature must be a different one.
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "type-creature", decor: decorChosen},
		}, true
	case engine.UpToCreatures:
		// A ceiling the controller may stop short of, so the stacked pick carries the
		// numeral and the may mark — the same pair "up to N cards" draws for a
		// Quantity (quantityGlyphs).
		return []glyph{{
			asset: "type-creature",
			qty:   v.Creatures,
			decor: decorChosen | decorEach,
			sigil: sigilMay,
		}}, true
	case engine.DivideDamage:
		// Any number of creatures, so the stacked pick wears the may mark with no
		// ceiling to print. How the pool is divided among them stays in the text.
		return []glyph{{
			asset: "type-creature",
			decor: decorChosen | decorEach,
			sigil: sigilMay,
		}}, true
	case engine.FlankWalk:
		// The pick is a flank creature and the walk runs inward from it, so the flank
		// mark leads the creatures the walk steps through.
		return []glyph{
			{asset: "glyph-flank", decor: decorChosen},
			{asset: "type-creature", decor: decorEach},
		}, true
	}
	return fallbackGlyphs(), false
}

// neighborScopeDecor is how a spread's neighbors are picked: every neighbor is
// the stack, one chosen neighbor the chosen outline — the same pair every other
// noun in the strip wears, rather than a drawing of its own.
//
// TestNeighborScopesDrawApart holds the two apart, so Mighty Lance's one neighbor
// cannot read as Booby Trap's both.
func neighborScopeDecor(scope engine.NeighborScope) decor {
	if scope == engine.OneNeighbor {
		return decorChosen
	}
	return decorEach
}

// perTargetGlyphs renders a PerTarget as the thing it measures on each creature
// the effect lands on, reporting whether a real mapping covered it. The PerTarget
// family is gated (TestEveryPerTargetDraws, TestIconTotality), so every row here
// draws.
//
// Each axis measures something sitting on the creature, so it draws that thing
// under the chosen outline every other glyph reading the card in context wears
// (countGlyphs): "for each Æmber on it" is the Æmber noun marked as that
// creature's, not a drawing of the "for each" clause.
//
// The members are singletons of unexported types, so they are recognised by
// identity rather than by a type switch; an unlisted axis falls through to the
// unknown glyph, which is what the gate fails on.
func perTargetGlyphs(p engine.PerTarget) ([]glyph, bool) {
	switch p {
	case engine.AemberOnIt:
		return []glyph{{asset: "aember", decor: decorChosen}}, true
	case engine.DamageOnIt:
		return []glyph{{asset: "damage", decor: decorChosen}}, true
	case engine.ArmorLostThisWay:
		// The armor this very effect stripped, so the shield wears the lose verb: it
		// is armor that is gone rather than armor the creature still has.
		return []glyph{{asset: "shield", decor: decorChosen, sigil: sigilLose}}, true
	case engine.UpgradesOnIt:
		return []glyph{{asset: "type-upgrade", decor: decorChosen | decorEach}}, true
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
