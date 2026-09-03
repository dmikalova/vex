package engine

import (
	"fmt"
	"strings"
)

// LoseKeyword takes a keyword away from every creature in play for the remainder
// of the turn — Sniffer sniffs out each elusive creature. It reaches creatures
// played later in the turn too, because the loss is held on the game rather than
// on the creatures it found.
type LoseKeyword struct {
	Keyword Keyword
}

// validate rejects a keyword that cannot be taken away.
func (e LoseKeyword) validate() error {
	if !e.Keyword.valid() {
		return fmt.Errorf("LoseKeyword: unset keyword")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each creature
// loses elusive".
func (e LoseKeyword) Text() string {
	return "for the remainder of the turn, each creature loses " +
		strings.ToLower(e.Keyword.String())
}

// Resolve records the loss for the rest of the turn.
func (e LoseKeyword) Resolve(ctx *EffectContext) {
	ctx.Resolver.LoseKeyword(e.Keyword)
}
