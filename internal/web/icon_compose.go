package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the composition family
// (ADR 0047) — the wrapper effects that carry no glyph of their own and instead
// recurse back into effectGlyphs/composeGlyphs over the sub-effect(s) they wrap
// (a repetition, a choice, a conditional, a duration, a scheduled follow-up).

// composeEffectGlyphs claims every effect that composes over sub-effects rather
// than drawing its own glyph.
func composeEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.ForEachHouse:
		return claim(effectGlyphs(v.Do))
	case engine.Repeat:
		return claim(effectGlyphs(v.Do))
	case engine.May:
		return claim(effectGlyphs(v.Do))
	case engine.ByActivePlayer:
		return claim(effectGlyphs(v.Do))
	case engine.Then:
		if first, ok := v.First.(engine.Effect); ok {
			return claim(composeGlyphs(first, v.Result))
		}
		gs, _ := effectGlyphs(v.Result)
		return gs, false, true
	case engine.Sequence:
		return claim(composeGlyphs(v.Effects...))
	case engine.ChooseOne:
		return append([]glyph{{asset: "glyph-choose"}}, dividedOptions(v.Options)...), true, true
	case engine.Conditional:
		return claim(conditionalGlyphs(v))
	case engine.ChooseCreatureThen:
		gs := []glyph{targetGlyph(v.Target)}
		more, covered := effectGlyphs(v.Then)
		return append(gs, more...), covered, true
	case engine.ChooseHouseThen:
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Then)...), true, true
	case engine.OnChooseCreature:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true, true
	case engine.OneAtATime:
		return append([]glyph{targetGlyph(v.Target)}, verbGlyphs(v.Verbs)...), true, true
	case engine.ForEach:
		return claim(effectGlyphs(v.Do))
	case engine.ForEachDiscarded:
		return claim(effectGlyphs(v.Do))
	case engine.TriggerAbility:
		return []glyph{
			{asset: triggerIcon(v.Trigger)},
			arrowTo(targetGlyph(v.Target)),
		}, true, true
	case engine.ScheduleOnLeave:
		return append(
			[]glyph{{asset: "type-creature", decor: decorThis}},
			mustCompose(v.Do)...), true, true
	case engine.Instead:
		return []glyph{{asset: "glyph-swap"}}, true, true
	case engine.ForDuration:
		gs, covered := composeGlyphs(v.Effects...)
		return claim(windowed(gs, covered, v.Duration))
	case engine.ForRemainderOfTurn:
		gs, covered := effectGlyphs(v.Do)
		return claim(windowed(gs, covered, engine.RemainderOfPlayerTurn))
	case engine.ForOpponentNextTurn:
		gs, covered := effectGlyphs(v.Do)
		return claim(windowed(gs, covered, engine.OpponentNextTurn))
	case engine.NextPlayed:
		return append([]glyph{{asset: "glyph-play"}}, mustCompose(v.EntersPlay)...), true, true
	case engine.GainUntilNextTurn:
		gs, covered := composeGlyphs(v.Effects...)
		return claim(windowed(gs, covered, engine.StartOfPlayerNextTurn))
	default:
		return nil, false, false
	}
}

// conditionalGlyphs draws a gated effect as two punctuated clauses: a gate mark
// opens the consequence, and a branch mark opens the alternative when there is
// one (ADR 0047). Without the marks "if x, destroy a creature, otherwise draw a
// card" runs together into one sentence of four glyphs.
//
// The condition itself draws nothing yet — the Condition family has no renderer —
// so the gate mark is currently the only sign that what follows is hypothetical.
// It stays the clause opener once conditions are drawn: the condition's glyphs
// take the mark, and the consequence follows them.
func conditionalGlyphs(v engine.Conditional) ([]glyph, bool) {
	gs, covered := effectGlyphs(v.Then)
	gs = punctuate(gs, punctGate)
	if v.Else == nil {
		return gs, covered
	}
	alt, altCovered := effectGlyphs(v.Else)
	return append(gs, punctuate(alt, punctBranch)...), covered && altCovered
}

// dividedOptions renders a choice's options as one run with a divider mark
// opening every option after the first, so the strip reads as alternatives rather
// than as a sequence of things that all happen (ADR 0047).
func dividedOptions(options []engine.Effect) []glyph {
	out := make([]glyph, 0, len(options)*2)
	for i := range options {
		gs := mustCompose(options[i])
		if i > 0 {
			gs = punctuate(gs, punctDivider)
		}
		out = append(out, gs...)
	}
	return out
}
