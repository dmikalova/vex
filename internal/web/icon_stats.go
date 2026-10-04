package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the stats-and-abilities
// family (ADR 0047) — power/armor changes, keyword and text-box grants, and the
// card-level identity changes (house, creature type) that ride alongside them.

// statsEffectGlyphs claims every stat, keyword and ability-grant effect.
func statsEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.GainStats:
		gs := make([]glyph, 0, 3)
		if v.Power != 0 {
			gs = append(gs, glyph{
				asset: "power",
				qty:   v.Power,
			})
		}
		if v.Armor != 0 {
			gs = append(gs, glyph{
				asset: "shield",
				qty:   v.Armor,
			})
		}
		return append(gs, arrowToTarget(v.Target)...), true
	case engine.OverrideStats:
		gs := make([]glyph, 0, 2)
		if v.HasPower {
			gs = append(gs, glyph{
				asset: "power",
				qty:   v.Power,
			})
		}
		if v.HasArmor {
			gs = append(gs, glyph{
				asset: "shield",
				qty:   v.Armor,
			})
		}
		return gs, true
	case engine.CopyPrintedStats:
		return append(targetGlyphs(v.Source), arrowTo(glyph{
			asset: "type-creature",
			decor: decorThis,
		})), true
	case engine.GainAssault:
		return append([]glyph{{asset: "kw-assault"}}, arrowToTarget(v.Target)...), true
	case engine.GainAssaultUntilNextTurn:
		// The same grant as GainAssault, bounded by a window rather than kept, so the
		// turn marker is the whole difference between the two (ADR 0047). Declaring
		// them synonyms was the alternative — a lasting buff and a permanent one are
		// one act with a different shelf life — but composition costs one slot on one
		// rare effect and says the shelf life out loud, so it is drawn.
		return inWindow(
			append([]glyph{{asset: "kw-assault"}}, arrowToTarget(v.Target)...),
			engine.StartOfPlayerNextTurn,
		), true
	case engine.GainKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return unmapped(k), true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, arrowToTarget(v.Target)...), true
	case engine.LoseKeyword:
		// The board-wide form: it carries no Target because every creature loses the
		// keyword, so the stacked creature noun is what tells it from the targeted
		// LoseKeywords beside it (TestDistinctMechanicsDrawDistinctGlyphs).
		a := keywordIcon(v.Keyword)
		if a == "" {
			return unmapped(v.Keyword), true
		}
		return []glyph{
			{asset: a},
			{asset: "glyph-ban"},
			arrowTo(glyph{asset: "type-creature", decor: decorEach}),
		}, true
	case engine.LoseKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return unmapped(k), true
			}
			gs = append(gs, glyph{asset: a})
		}
		return inWindow(
			append(append(gs, glyph{asset: "glyph-ban"}), arrowToTarget(v.Target)...),
			v.Duration,
		), true
	case engine.GainAbility:
		gs := append(targetGlyphs(v.Target), glyph{asset: triggerIcon(v.Ability.Trigger)})
		return append(gs, composeGlyphs(v.Ability.Effect)...), true
	case engine.GainTextBox:
		return append(targetGlyphs(v.Source), arrowToTarget(v.Target)...), true
	case engine.LendTextBoxFromHand:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorChosen,
			}),
		}, true
	case engine.BlankEnemyText:
		return []glyph{{asset: "type-creature", decor: decorEnemy | decorEach}}, true
	case engine.GainTrait:
		// Which trait is granted stays in the rules text — traitGlyph is one mark for
		// all 70 of them — so the strip says a trait lands on the target.
		return append([]glyph{traitGlyph()}, arrowToTarget(v.Target)...), true
	case engine.BelongToHouse:
		if a := houseIconName(v.House); a != "" {
			return append([]glyph{{asset: a}}, arrowToTarget(v.Target)...), true
		}
		return append([]glyph{{asset: "glyph-swap"}}, arrowToTarget(v.Target)...), true
	case engine.TurnIntoCreature:
		return append(
			targetGlyphs(v.Target),
			arrowTo(glyph{asset: "type-creature"}),
			glyph{asset: "glyph-flank"},
		), true
	case engine.FuseTriggersForTurn:
		return []glyph{
			{asset: "glyph-reap"},
			{asset: "glyph-swap"},
			{asset: "glyph-fight"},
		}, true
	default:
		return nil, false
	}
}
