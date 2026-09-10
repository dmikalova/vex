package engine

import (
	"fmt"
	"strings"
)

// GainKeywordForTurn gives each creature its Target selects a keyword for the
// remainder of the turn — Creed of Nature grants a chosen creature skirmish. Unlike
// GainKeyword, which lasts until the controller's next turn so a defensive keyword
// survives the opponent's turn, this clears at the ready phase, for an offensive
// keyword (skirmish) that only matters on the controller's own turn. Because it
// renders a duration body it composes under ForDuration with another per-creature
// grant.
type GainKeywordForTurn struct {
	Target  Target
	Keyword Keyword
}

// validate requires an explicit target and a real keyword.
func (e GainKeywordForTurn) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainKeywordForTurn")
	}
	if !e.Keyword.valid() {
		return fmt.Errorf("GainKeywordForTurn: unset keyword")
	}
	return nil
}

// durationSubject names the creature that gains the keyword, folding the "gains"
// verb into the subject so ForDuration can share it with a sibling grant.
func (e GainKeywordForTurn) durationSubject() string { return e.Target.Text() + " gains" }

// durationPredicate renders the keyword gained, e.g. "skirmish".
func (e GainKeywordForTurn) durationPredicate() string {
	return strings.ToLower(e.Keyword.String())
}

// Text renders the effect, e.g. "for the remainder of the turn, it gains skirmish".
func (e GainKeywordForTurn) Text() string {
	return "for the remainder of the turn, " + e.durationSubject() + " " + e.durationPredicate()
}

// Resolve grants each selected creature the keyword for the remainder of the turn.
func (e GainKeywordForTurn) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GrantKeyword(id, e.Keyword)
	}
}
