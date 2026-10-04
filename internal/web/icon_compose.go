package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the composition family
// (ADR 0047) — the wrapper effects that carry no glyph of their own and instead
// recurse back into effectGlyphs/composeGlyphs over the sub-effect(s) they wrap
// (a repetition, a choice, a conditional, a duration, a scheduled follow-up).

// composeEffectGlyphs claims every effect that composes over sub-effects rather
// than drawing its own glyph.
func composeEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.ForEachHouse:
		return effectGlyphs(v.Do), true
	case engine.Repeat:
		return append(effectGlyphs(v.Do), repeatGateGlyphs(v.Gate)...), true
	case engine.May:
		return effectGlyphs(v.Do), true
	case engine.ByActivePlayer:
		return effectGlyphs(v.Do), true
	case engine.Then:
		// A GatingEffect embeds Effect, so the gate draws like any other sub-effect.
		// A Then with no gate at all is a half-built value no card can carry, and it
		// draws its Result rather than nothing.
		if v.First == nil {
			return effectGlyphs(v.Result), true
		}
		return composeGlyphs(v.First, v.Result), true
	case engine.Sequence:
		return composeGlyphs(v.Effects...), true
	case engine.ChooseOne:
		return append([]glyph{{asset: "glyph-choose"}}, dividedOptions(v.Options)...), true
	case engine.Conditional:
		return conditionalGlyphs(v), true
	case engine.ChooseCreatureThen:
		return append(targetGlyphs(v.Target), effectGlyphs(v.Then)...), true
	case engine.ChooseHouseThen:
		return append([]glyph{{asset: "glyph-choose"}}, composeGlyphs(v.Then)...), true
	case engine.OnChooseCreature:
		return append(targetGlyphs(v.Target), verbGlyphs(v.Verbs)...), true
	case engine.OneAtATime:
		return append(targetGlyphs(v.Target), verbGlyphs(v.Verbs)...), true
	case engine.ForEach:
		return effectGlyphs(v.Do), true
	case engine.ForEachDiscarded:
		return effectGlyphs(v.Do), true
	case engine.TriggerAbility:
		return append([]glyph{
			{asset: triggerIcon(v.Trigger)},
		}, arrowToTarget(v.Target)...), true
	case engine.ScheduleOnLeave:
		return append(
			[]glyph{{asset: "type-creature", decor: decorThis}},
			composeGlyphs(v.Do)...), true
	case engine.Instead:
		return []glyph{{asset: "glyph-swap"}}, true
	case engine.ForDuration:
		return inWindow(composeGlyphs(v.Effects...), v.Duration), true
	case engine.ForRemainderOfTurn:
		return inWindow(effectGlyphs(v.Do), engine.RemainderOfPlayerTurn), true
	case engine.ForOpponentNextTurn:
		return inWindow(effectGlyphs(v.Do), engine.OpponentNextTurn), true
	case engine.NextPlayed:
		return append([]glyph{{asset: "glyph-play"}}, composeGlyphs(v.EntersPlay)...), true
	case engine.GainUntilNextTurn:
		return inWindow(composeGlyphs(v.Effects...), engine.StartOfPlayerNextTurn), true
	default:
		return nil, false
	}
}

// repeatGateGlyphs renders the gate a Repeat consults before running its effect
// again. The RepeatGate family is gated (TestEveryRepeatGateDraws,
// TestIconTotality), so every row here draws.
//
// Three of the four are plumbing: a While, a WhileYouDo and a MayWhileYouDo only
// loop, and everything a player reads comes from the effect being repeated and
// the condition's own glyphs, so they fold into those the way every other
// plumbing node folds into its children — the condition behind the gate mark
// that already opens a hypothetical clause, wearing the may mark when the
// controller chooses each further pass. While and WhileYouDo draw alike: the
// result-gate arrow their text wears for "only while the effect keeps happening"
// is a wording difference, not a second picture. ByExalting is not plumbing:
// exalting is a cost it charges to repeat, and a cost is a rule of its own.
func repeatGateGlyphs(g engine.RepeatGate) []glyph {
	switch v := g.(type) {
	case engine.While:
		return conditionGlyphs(v.Cond)
	case engine.WhileYouDo:
		return conditionGlyphs(v.Cond)
	case engine.MayWhileYouDo:
		return markVerb(conditionGlyphs(v.Cond), sigilMay)
	case engine.ByExalting:
		// The Æmber the repeat costs, exalted onto the creature that pays it, which
		// the controller may decline — so the may mark rides that creature, as it does
		// on every other declinable pick (optionalVerb).
		paid := markVerb(targetGlyphs(v.Creature), sigilMay)
		return append(
			[]glyph{{asset: "aember", sigil: sigilExalt}},
			punctuate(paid, punctArrow)...,
		)
	}
	return unmapped(g)
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
func conditionalGlyphs(v engine.Conditional) []glyph {
	gs := punctuate(effectGlyphs(v.Then), punctGate)
	if v.Else == nil {
		return gs
	}
	return append(gs, punctuate(effectGlyphs(v.Else), punctBranch)...)
}

// dividedOptions renders a choice's options as one run with a divider mark
// opening every option after the first, so the strip reads as alternatives rather
// than as a sequence of things that all happen (ADR 0047).
func dividedOptions(options []engine.Effect) []glyph {
	out := make([]glyph, 0, len(options)*2)
	for i := range options {
		gs := composeGlyphs(options[i])
		if i > 0 {
			gs = punctuate(gs, punctDivider)
		}
		out = append(out, gs...)
	}
	return out
}
