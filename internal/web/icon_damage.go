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
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.DestroyEachCreatureAtEndOfTurn:
		return []glyph{
			{asset: "glyph-destroy"},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorEach,
			}),
		}, true, true
	case engine.PurgeCreature:
		return []glyph{{asset: "zone-purge"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.PurgeSource:
		return []glyph{{asset: "zone-purge", decor: decorThis}}, true, true
	default:
		return nil, false, false
	}
}

// damageThenGlyphs renders a "deal N damage to <target>, then <follow-up>" effect
// as the damage glyph arrowed to its target, followed by the follow-up's glyphs.
func damageThenGlyphs(amount int, target engine.Target, then engine.Effect) ([]glyph, bool) {
	gs := []glyph{{asset: "damage", qty: amount}, arrowTo(targetGlyph(target))}
	more, covered := effectGlyphs(then)
	return append(gs, more...), covered
}
