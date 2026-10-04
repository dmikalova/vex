package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the creature-state
// family (ADR 0047) — the exhaust/ready cycle, stun, ward, enrage, using a
// creature's action, repeated fights, and the counters that track them.

// creatureEffectGlyphs claims every creature-state effect.
func creatureEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.Stun:
		return append([]glyph{{asset: "stun"}}, arrowToTarget(v.Target)...), true
	case engine.Unstun:
		return append([]glyph{
			{asset: "stun", decor: decorFriendly},
		}, arrowToTarget(v.Target)...), true
	case engine.Enrage:
		return append([]glyph{{asset: "enrage"}}, arrowToTarget(v.Target)...), true
	case engine.Ward:
		return append([]glyph{{asset: "ward"}}, arrowToTarget(v.Target)...), true
	case engine.RemoveWard:
		return append([]glyph{
			{asset: "ward"},
			{asset: "glyph-ban"},
		}, arrowToTarget(v.Target)...), true
	case engine.Exhaust:
		return append([]glyph{{asset: "exhausted"}}, arrowToTarget(v.Target)...), true
	case engine.ExhaustCreatures:
		return append([]glyph{{asset: "exhausted"}}, arrowToTarget(v.Target)...), true
	case engine.Ready:
		return append([]glyph{
			{asset: "exhausted", decor: decorFriendly},
		}, arrowToTarget(v.Target)...), true
	case engine.ReadyCreatures:
		return append([]glyph{
			{asset: "exhausted", decor: decorFriendly},
		}, arrowToTarget(v.Target)...), true
	case engine.Use:
		return append([]glyph{{asset: "glyph-action"}}, arrowToTarget(v.Target)...), true
	case engine.RepeatedFight:
		return append([]glyph{{asset: "glyph-fight"}}, arrowToTarget(v.Target)...), true
	case engine.CancelFight:
		return []glyph{{asset: "glyph-fight"}, {asset: "glyph-ban"}}, true
	case engine.PlaceCounter:
		return append([]glyph{
			{asset: counterAsset(v.Kind), sigil: sigilGain},
		}, arrowToTarget(v.Target)...), true
	case engine.RemoveCounters:
		return append([]glyph{
			{asset: counterAsset(v.Kind), sigil: sigilLose},
		}, arrowToTarget(v.Target)...), true
	case engine.AddPowerCounter:
		if v.Per != nil || v.Equal != nil {
			return append([]glyph{{asset: "power"}}, arrowToTarget(v.Target)...), true
		}
		return append(
			[]glyph{{asset: "power", qty: v.Amount}},
			arrowToTarget(v.Target)...), true
	default:
		return nil, false
	}
}
