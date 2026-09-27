package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the stats-and-abilities
// family (ADR 0047) — power/armor changes, keyword and text-box grants, and the
// card-level identity changes (house, creature type) that ride alongside them.

// statsEffectGlyphs claims every stat, keyword and ability-grant effect.
func statsEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
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
		return append(gs, arrowTo(targetGlyph(v.Target))), true, true
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
		return gs, true, true
	case engine.CopyPrintedStats:
		return []glyph{
			targetGlyph(v.Source),
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorThis,
			}),
		}, true, true
	case engine.GainAssault:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainAssaultUntilNextTurn:
		return []glyph{{asset: "kw-assault"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.GainKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}}, true, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, arrowTo(targetGlyph(v.Target))), true, true
	case engine.LoseKeyword:
		if a := keywordIcon(v.Keyword); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-ban"}}, true, true
		}
		return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true, true
	case engine.LoseKeywords:
		gs := make([]glyph, 0, len(v.Keywords)+1)
		for _, k := range v.Keywords {
			a := keywordIcon(k)
			if a == "" {
				return []glyph{{asset: "glyph-unknown"}, {asset: "glyph-ban"}}, true, true
			}
			gs = append(gs, glyph{asset: a})
		}
		return append(gs, glyph{asset: "glyph-ban"}), true, true
	case engine.GainAbility:
		gs := []glyph{targetGlyph(v.Target), {asset: triggerIcon(v.Ability.Trigger)}}
		return append(gs, mustCompose(v.Ability.Effect)...), true, true
	case engine.GainTextBox:
		return []glyph{targetGlyph(v.Source), arrowTo(targetGlyph(v.Target))}, true, true
	case engine.LendTextBoxFromHand:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorChosen,
			}),
		}, true, true
	case engine.BlankEnemyText:
		return []glyph{{asset: "type-creature", decor: decorEnemy | decorEach}}, true, true
	case engine.GainTrait:
		// A trait has no icon in the strip's vocabulary — traits render as the
		// card's text, not glyphs. It only ever folds beside a keyword grant that
		// carries the line's glyph, so it renders nothing yet counts as covered.
		return nil, true, true
	case engine.BelongToHouse:
		if a := houseIconName(v.House); a != "" {
			return []glyph{{asset: a}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{{asset: "glyph-swap"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.TurnIntoCreature:
		return []glyph{
			targetGlyph(v.Target),
			arrowTo(glyph{asset: "type-creature"}),
			{asset: "glyph-flank"},
		}, true, true
	case engine.FuseTriggersForTurn:
		return []glyph{
			{asset: "glyph-reap"},
			{asset: "glyph-swap"},
			{asset: "glyph-fight"},
		}, true, true
	default:
		return nil, false, false
	}
}
