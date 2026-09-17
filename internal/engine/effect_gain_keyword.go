package engine

import (
	"fmt"
	"strings"
)

// GainKeywords gives each creature its Target selects one or more keywords for a
// duration — Hideaway Hole grants your creatures elusive until the start of your
// next turn, Creed of Nature grants a chosen creature skirmish for the remainder of
// the turn. StartOfPlayerNextTurn holds the grant on each creature and lifts it at
// the start of that player's next turn, before any start-of-turn ability resolves,
// so it survives the opponent's turn — which a defensive keyword like elusive needs.
// RemainderOfPlayerTurn clears at the ready phase, for an offensive keyword like
// skirmish that only matters on the controller's own turn. It renders a duration
// body, so it composes under GainUntilNextTurn (next turn) or ForDuration
// (remainder) with another per-creature grant. A keyword that only matters on your
// own turn and needs the acting-verb building block (Scout's Skirmish) uses
// GainKeywordVerb instead.
type GainKeywords struct {
	Target   Target
	Keywords []Keyword
	Duration Duration
}

// validate requires an explicit target, at least one real keyword, and one of the
// two supported durations.
func (e GainKeywords) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainKeywords")
	}
	if len(e.Keywords) == 0 {
		return fmt.Errorf("GainKeywords: no keywords")
	}
	for _, k := range e.Keywords {
		if !k.valid() {
			return fmt.Errorf("GainKeywords: unset keyword")
		}
	}
	if e.Duration != RemainderOfPlayerTurn && e.Duration != StartOfPlayerNextTurn {
		return fmt.Errorf(
			"GainKeywords: duration must be RemainderOfPlayerTurn or StartOfPlayerNextTurn",
		)
	}
	return nil
}

// durationSubject names the creature that gains the keywords, folding the "gains"
// verb into the subject so ForDuration and GainUntilNextTurn can share it with a
// sibling grant.
func (e GainKeywords) durationSubject() string { return e.Target.Text() + " gains" }

// durationPredicate renders the keywords gained, e.g. "elusive" or "taunt and
// elusive".
func (e GainKeywords) durationPredicate() string {
	names := make([]string, len(e.Keywords))
	for i, k := range e.Keywords {
		names[i] = strings.ToLower(k.String())
	}
	return oxfordAnd(names)
}

// Text renders the effect standalone, framing the body the way its fold wrapper
// would — a suffix for the next-turn grant ("each friendly creature gains elusive
// until the start of your next turn"), a prefix for the remainder grant ("for the
// remainder of the turn, it gains skirmish").
func (e GainKeywords) Text() string {
	body := e.durationSubject() + " " + e.durationPredicate()
	if e.Duration == StartOfPlayerNextTurn {
		return body + " until the start of your next turn"
	}
	return "for the remainder of the turn, " + body
}

// Resolve grants each selected creature every keyword for the chosen duration.
func (e GainKeywords) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		for _, k := range e.Keywords {
			if e.Duration == StartOfPlayerNextTurn {
				ctx.Resolver.GrantKeywordUntilNextTurn(id, k)
			} else {
				ctx.Resolver.GrantKeyword(id, k)
			}
		}
	}
}
