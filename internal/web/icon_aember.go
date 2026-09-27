package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the Æmber-and-keys
// family (ADR 0047) — gaining, losing, stealing, capturing, exalting and moving
// Æmber, and the key-forging economy that spends it. Key/forge lives here rather
// than in a family of its own: forging is the sink for Æmber's flow, and its six
// cases (ForgeKey, UnforgeKey, RaiseKeyCost, LowerKeyCost, SkipForgePhase,
// CancelForge) do not earn a file to themselves.

// aemberEffectGlyphs claims every Æmber-flow and key-forging effect.
func aemberEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.GainAember:
		if v.Per != nil || v.EqualTo != nil {
			return []glyph{{asset: "aember", decor: playerDecor(v.Player)}}, true, true
		}
		return []glyph{{asset: "aember", qty: v.Amount, decor: playerDecor(v.Player)}}, true, true
	case engine.LoseAember:
		return []glyph{{asset: "aember", text: "−", decor: playerDecor(v.Player)}}, true, true
	case engine.StealAember:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorEnemy | decorChosen},
		}, true, true
	case engine.CaptureAember:
		return []glyph{{asset: "aember", qty: v.Amount, decor: decorEnemy}}, true, true
	case engine.CaptureFromAnyPlayer:
		return []glyph{{asset: "aember", qty: v.Amount}}, true, true
	case engine.DistributeCapture:
		return []glyph{{asset: "aember", decor: decorEnemy}}, true, true
	case engine.GiveAember:
		src := glyph{
			asset: "aember",
			decor: decorEnemy,
		}
		if !v.All {
			src.qty = v.Amount
		}
		return []glyph{src, arrowTo(glyph{asset: "aember"})}, true, true
	case engine.Exalt:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.MoveAember:
		return []glyph{{asset: "aember"}, arrowTo(moveAemberDest(v.Onto, v.To))}, true, true
	case engine.MoveAemberFromPool:
		return []glyph{{asset: "aember", qty: v.Amount}, arrowTo(targetGlyph(v.Target))}, true, true
	case engine.MoveAemberToSupply:
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: decorChosen},
			arrowTo(glyph{asset: "glyph-return"}),
		}, true, true
	case engine.PlaceAemberOnThis:
		return []glyph{
			{asset: "aember", qty: v.Amount},
			arrowTo(glyph{
				asset: "card-back",
				decor: decorThis,
			}),
		}, true, true
	case engine.RedistributeCapturedAember:
		decor := decorFriendly
		if v.Side == engine.Opponent {
			decor = decorEnemy
		}
		return []glyph{{asset: "aember", decor: decor}, {asset: "glyph-swap"}}, true, true
	case engine.ForgeKey:
		return []glyph{{asset: "forge"}}, true, true
	case engine.UnforgeKey:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	case engine.RaiseKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: v.Amount}}, true, true
	case engine.LowerKeyCost:
		return []glyph{{asset: "forge"}, {asset: "aember", qty: -v.Amount}}, true, true
	case engine.SkipForgePhase:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	case engine.CancelForge:
		return []glyph{{asset: "forge"}, {asset: "glyph-ban"}}, true, true
	default:
		return nil, false, false
	}
}

// moveAemberDest is where moved Æmber lands: onto a chosen card, or into a pool.
func moveAemberDest(onto engine.Target, to engine.Player) glyph {
	if onto != (engine.Target{}) {
		return targetGlyph(onto)
	}
	return glyph{
		asset: "aember",
		decor: playerDecor(to),
	}
}
