package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the creature-state
// family (ADR 0047) — the exhaust/ready cycle, stun, ward, enrage, using a
// creature's action, repeated fights, and the counters that track them.

// creatureEffectGlyphs claims every creature-state effect.
func creatureEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.Stun:
		return []glyph{{asset: "stun"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Unstun:
		return []glyph{
			{asset: "stun", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Enrage:
		return []glyph{{asset: "enrage"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Ward:
		return []glyph{{asset: "ward"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RemoveWard:
		return []glyph{
			{asset: "ward"},
			{asset: "glyph-ban"},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Exhaust:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.ExhaustCreatures:
		return []glyph{{asset: "exhausted"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.Ready:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.ReadyCreatures:
		return []glyph{
			{asset: "exhausted", decor: decorFriendly},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.Use:
		return []glyph{{asset: "glyph-action"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.RepeatedFight:
		return []glyph{{asset: "glyph-fight"}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.CancelFight:
		return []glyph{{asset: "glyph-fight"}, {asset: "glyph-ban"}}, true, true
	case engine.PlaceCounter:
		return []glyph{
			{asset: counterAsset(v.Kind), sigil: sigilGain},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.RemoveCounters:
		return []glyph{
			{asset: counterAsset(v.Kind), sigil: sigilLose},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.AddPowerCounter:
		if v.Per != nil || v.Equal != nil {
			return []glyph{{asset: "power"}, arrowTo(targetGlyph(v.Target))}, true, true
		}
		return []glyph{{asset: "power", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true, true
	default:
		return nil, false, false
	}
}
