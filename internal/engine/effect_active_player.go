package engine

import "strings"

// ByActivePlayer resolves its inner effect as the active player — the player
// whose turn it is — rather than as the ability's controller. It carries out the
// KeyForge wordings where a public event names the player who caused it: Giant
// Gnawbill's "After a player chooses an active house, that player destroys an
// artifact of that house" fires on both players' cards, but the chooser — not the
// card's owner — destroys the artifact.
type ByActivePlayer struct {
	Do Effect
}

// Text renders the inner effect in the third person with "that player" as its
// subject, e.g. "that player destroys an artifact of that house".
func (e ByActivePlayer) Text() string {
	verb, rest, found := strings.Cut(e.Do.Text(), " ")
	if !found {
		return "that player " + verb + "s"
	}
	return "that player " + verb + "s " + rest
}

// Resolve carries out the inner effect with the active player as its controller,
// then restores the original controller so a following effect is unaffected.
func (e ByActivePlayer) Resolve(ctx *EffectContext) {
	saved := ctx.Controller
	ctx.Controller = ctx.Resolver.ActivePlayer()
	e.Do.Resolve(ctx)
	ctx.Controller = saved
}

// validate descends into the wrapped effect.
func (e ByActivePlayer) validate() error { return validateEffect(e.Do) }
