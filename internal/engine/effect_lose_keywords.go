package engine

import (
	"fmt"
	"strings"
)

// LoseKeywords takes one or more keywords away from each creature its Target
// selects for a duration — Niffle Grounds strips a chosen creature of taunt and
// elusive for the remainder of the turn, Reckless Rizzo loses elusive until the
// start of its controller's next turn so the loss survives the opponent's turn. The
// loss is held on the creature; RemainderOfPlayerTurn lifts it at the ready phase,
// StartOfPlayerNextTurn at the start of the controller's next turn. Every keyword
// check honors it meanwhile.
type LoseKeywords struct {
	Target   Target
	Keywords []Keyword
	Duration Duration
}

// validate requires an explicit target, at least one valid keyword, and one of the
// two supported durations.
func (e LoseKeywords) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("LoseKeywords")
	}
	if len(e.Keywords) == 0 {
		return fmt.Errorf("LoseKeywords: no keywords")
	}
	for _, k := range e.Keywords {
		if !k.valid() {
			return fmt.Errorf("LoseKeywords: unset keyword")
		}
	}
	if e.Duration != RemainderOfPlayerTurn && e.Duration != StartOfPlayerNextTurn {
		return fmt.Errorf(
			"LoseKeywords: duration must be RemainderOfPlayerTurn or StartOfPlayerNextTurn",
		)
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, it loses taunt and
// elusive" or "until the start of your next turn, {self} loses elusive".
func (e LoseKeywords) Text() string {
	names := make([]string, 0, len(e.Keywords))
	for _, k := range e.Keywords {
		names = append(names, strings.ToLower(k.String()))
	}
	clause := durationClause(e.Duration, "")
	return clause + ", " + e.Target.Text() + " loses " + oxfordAnd(names)
}

// Resolve takes each keyword away from every selected creature for the chosen
// duration.
func (e LoseKeywords) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		for _, k := range e.Keywords {
			if e.Duration == StartOfPlayerNextTurn {
				ctx.Resolver.LoseKeywordUntilNextTurn(id, k)
			} else {
				ctx.Resolver.LoseKeywordFrom(id, k)
			}
		}
	}
}
