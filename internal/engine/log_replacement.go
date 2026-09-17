package engine

import "strings"

// This file holds the one voice every replaced outcome is narrated in (ADR 0011).
// A replacement has three parts the reader needs: the card that caused it, what
// actually happened, and what that happened instead of. Each part is always in
// the same place, so every replacing card reads the same way:
//
//	Dimension Door has Player 1 steal 1 Æmber reaping with Dust Pixie, instead of gaining it
//	Po's Pixies has Player 1 steal 2 Æmber from the common supply, instead of from Player 2's pool
//	Gargantodon has Thief capture 2 Æmber, instead of Player 1 stealing it
//	Ether Spider captures 1 Æmber, instead of Player 1 gaining it
//
// Without a shared seam each replacing card words its own line and the same
// mechanic reads a different way on every card that uses it.

// replacementLine renders one replaced outcome. cause is the card whose
// replacement fired, actor is who carried the outcome out, verb and rest are the
// outcome in its bare form ("capture", "2 Æmber"), and displaced is what the
// outcome happened instead of. When the cause is itself the actor (Ether Spider
// captures the Æmber it intercepted) the causative head is dropped and the verb
// agrees with the cause, since a card does not have itself do something.
func replacementLine(cause, actor, verb, rest, displaced string) string {
	parts := []string{cause}
	if actor == "" || actor == cause {
		parts = append(parts, verb+"s")
	} else {
		parts = append(parts, "has", actor, verb)
	}
	if rest != "" {
		parts = append(parts, rest)
	}
	return strings.Join(parts, " ") + ", instead of " + displaced
}
