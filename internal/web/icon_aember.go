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
			return []glyph{
				{asset: "aember", decor: playerDecor(v.Player), sigil: sigilGain},
			}, true, true
		}
		return []glyph{
			{asset: "aember", qty: v.Amount, decor: playerDecor(v.Player), sigil: sigilGain},
		}, true, true
	case engine.LoseAember:
		// The amount is only a numeral in the plain fixed case; a Per/EqualTo loss
		// scales by something that lives in the text, as GainAember's does, and a By
		// loss draws the share of the pool it takes.
		g := glyph{asset: "aember", decor: playerDecor(v.Player), sigil: sigilLose}
		if v.By == nil && v.Per == nil && v.EqualTo == nil {
			g.qty = v.Amount
		}
		return byShare([]glyph{g}, v.By)
	case engine.StealAember:
		return byShare([]glyph{
			{
				asset: "aember",
				qty:   v.Amount,
				decor: decorEnemy | decorChosen,
				sigil: sigilSteal,
			},
		}, v.By)
	case engine.CaptureAember:
		return byShare(
			[]glyph{{asset: "capture", qty: v.Amount, decor: decorEnemy}},
			v.By,
		)
	case engine.CaptureFromAnyPlayer:
		return []glyph{{asset: "capture", qty: v.Amount}}, true, true
	case engine.DistributeCapture:
		// The share is spread over as many friendly creatures as the controller
		// picks, which is what separates it from a single creature's capture.
		return byShare([]glyph{
			{asset: "capture", decor: decorEnemy},
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorFriendly | decorEach,
			}),
		}, v.By)
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
			{asset: "aember", qty: v.Amount, sigil: sigilExalt},
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

// byShare closes an economy verb's run with the share of the pool it moves, and
// claims the effect for the family. A verb with no By moves a plain amount and
// draws nothing extra, so the share costs a horizontal slot only on the verbs
// that have one — the same bargain inWindow makes for a duration.
func byShare(gs []glyph, by engine.Loss) ([]glyph, bool, bool) {
	if by == nil {
		return gs, true, true
	}
	share, covered := lossGlyphs(by)
	return append(gs, share...), covered, true
}

// lossGlyphs renders a Loss as the share of a pool it takes, reporting whether a
// real mapping covered it. The Loss family is gated (TestEveryLossDraws,
// TestIconTotality), so every row here draws.
//
// A loss is not an effect — it is how much the economy verb holding it moves — so
// it draws where that verb draws, the way a Count does. All three shares are the
// Æmber noun under a mark saying how much of the pool goes: the whole pool
// stacked, the remainder it leaves behind the refusal mark, the fraction's
// denominator as a numeral.
func lossGlyphs(l engine.Loss) ([]glyph, bool) {
	switch v := l.(type) {
	case engine.Fraction:
		// A share of the pool rather than a count of it, so the numeral is the
		// denominator the pool is cut into: half draws a 2, a third a 3. Which way it
		// rounds stays in the rules text, as every other amount's fine print does.
		if d := fractionDenominator(v); d > 0 {
			return []glyph{{asset: "aember", qty: d}}, true
		}
		return fallbackGlyphs(), false
	case engine.AllBut:
		// Everything above the remainder, so the pool it leaves is drawn behind the
		// refusal mark — the way every other barred thing is drawn.
		return []glyph{
			{asset: "aember", decor: decorEach},
			{asset: "aember", qty: v.Keep},
			{asset: "glyph-ban"},
		}, true
	}
	if l == engine.AllAember {
		// The whole pool, whatever its size, which is the stack every "each of these"
		// noun wears.
		return []glyph{{asset: "aember", decor: decorEach}}, true
	}
	return fallbackGlyphs(), false
}

// fractionDenominator is the number of parts a Fraction cuts a quantity into, or
// 0 for a fraction the strip cannot number. A Fraction is a comparable value with
// unexported fields built from four package vars, so it is recognised by identity
// against those four — the rounding is what pairs them, and the strip draws the
// share rather than the rounding.
func fractionDenominator(f engine.Fraction) int {
	switch f {
	case engine.HalfRoundedDown, engine.HalfRoundedUp:
		return 2
	case engine.ThirdRoundedDown, engine.ThirdRoundedUp:
		return 3
	}
	return 0
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
