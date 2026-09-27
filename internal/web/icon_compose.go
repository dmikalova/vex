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
		return append([]glyph{{asset: "glyph-choose"}}, mustCompose(v.Options...)...), true, true
	case engine.Conditional:
		if v.Else != nil {
			return claim(composeGlyphs(v.Then, v.Else))
		}
		return claim(effectGlyphs(v.Then))
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
		return claim(composeGlyphs(v.Effects...))
	case engine.ForRemainderOfTurn:
		return claim(effectGlyphs(v.Do))
	case engine.ForOpponentNextTurn:
		return claim(effectGlyphs(v.Do))
	case engine.NextPlayed:
		return append([]glyph{{asset: "glyph-play"}}, mustCompose(v.EntersPlay)...), true, true
	case engine.GainUntilNextTurn:
		return claim(composeGlyphs(v.Effects...))
	default:
		return nil, false, false
	}
}
