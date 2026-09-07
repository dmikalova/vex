package engine

import (
	"fmt"
	"strings"
)

// LoseKeywords takes one or more keywords away from each creature its Target
// selects for the remainder of the turn (Niffle Grounds strips a chosen creature
// of taunt and elusive). The loss is held on the creature, so the ready phase
// lifts it; every keyword check honors it.
type LoseKeywords struct {
	Target   Target
	Keywords []Keyword
}

// validate requires an explicit target and at least one valid keyword.
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
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, it loses taunt
// and elusive".
func (e LoseKeywords) Text() string {
	names := make([]string, 0, len(e.Keywords))
	for _, k := range e.Keywords {
		names = append(names, strings.ToLower(k.String()))
	}
	return "for the remainder of the turn, " + e.Target.Text() +
		" loses " + oxfordAnd(names)
}

// Resolve takes each keyword away from every selected creature for the turn.
func (e LoseKeywords) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		for _, k := range e.Keywords {
			ctx.Resolver.LoseKeywordFrom(id, k)
		}
	}
}
