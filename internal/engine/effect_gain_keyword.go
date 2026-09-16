package engine

import (
	"fmt"
	"strings"
)

// GainKeyword gives each creature its Target selects a keyword until the start of
// the controller's next turn (Hideaway Hole grants your creatures elusive). The
// grant is held on each creature and lifts at the start of that player's next turn,
// before any start-of-turn ability resolves, so it survives the opponent's turn —
// which is what a defensive keyword like elusive needs. A keyword that only matters on your own turn (Scout's Skirmish) uses the
// GainKeywordVerb building block instead.
type GainKeyword struct {
	Target  Target
	Keyword Keyword
}

// validate requires an explicit target and a real keyword.
func (e GainKeyword) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainKeyword")
	}
	if !e.Keyword.valid() {
		return fmt.Errorf("GainKeyword: unset keyword")
	}
	return nil
}

// Text renders the effect, e.g. "each friendly creature gains elusive until the
// start of your next turn".
func (e GainKeyword) Text() string {
	return fmt.Sprintf("%s gains %s until the start of your next turn",
		e.Target.Text(), strings.ToLower(e.Keyword.String()))
}

// durationSubject names the creature that gains the keyword, folding the "gains"
// verb into the subject so GainUntilNextTurn can share it with a sibling grant.
func (e GainKeyword) durationSubject() string { return e.Target.Text() + " gains" }

// durationPredicate renders the keyword gained, e.g. "elusive".
func (e GainKeyword) durationPredicate() string { return strings.ToLower(e.Keyword.String()) }

// Resolve grants each selected creature the keyword until the controller's next
// turn.
func (e GainKeyword) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GrantKeywordUntilNextTurn(id, e.Keyword)
	}
}
