package engine

import (
	"fmt"
	"strings"
)

// Use uses up to Max cards the controller chooses from Target's pool, one at a
// time — each use fully resolves before the next choice, so a card that exhausts
// itself drops out of later choices. Target is the candidate pool (an "each"
// target); Use offers only its ready, usable members: a creature (choosing whether
// it reaps, fights, or uses its Action) or an artifact with an Action. Combat
// Pheromones uses two other Mars cards.
type Use struct {
	Max    int
	Target Target
	// EvenUnusable offers every artifact in the pool, not only the ready ones with
	// an Action. Poltergeist sets it because its payload is the destroy that
	// follows, not the use: it must be able to choose any artifact — one with no
	// ability (Soul Snatcher), or one already exhausted — and use is incidental,
	// firing only if the chosen artifact can still act.
	EvenUnusable bool
}

// validate requires a target and a positive maximum.
func (e Use) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Use")
	}
	if e.Max <= 0 {
		return fmt.Errorf("Use: Max must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "use an enemy artifact" for a single use, or
// "use 2 other Mars cards, one at a time" for several.
func (e Use) Text() string {
	noun := useNoun(e.Target.Text())
	if e.Max == 1 {
		return "use " + indefinite(noun)
	}
	return "use " + countNoun(e.Max, noun) + ", one at a time"
}

// useNoun turns a Target's collective phrase into the singular noun the "use N ..."
// clause counts, dropping the "each"/"friendly" scaffolding but keeping "other" —
// "each other friendly Mars card" becomes "other Mars card".
func useNoun(phrase string) string {
	phrase = strings.TrimPrefix(phrase, "each ")
	phrase = strings.Replace(phrase, "other friendly ", "other ", 1)
	phrase = strings.Replace(phrase, "friendly ", "", 1)
	return phrase
}

// Resolve chooses and uses up to Max cards from the pool, one at a time. Each use
// fully resolves before the next choice, so an exhausted card drops out. The card
// used most recently is left in context, so a following effect can act on it
// (Poltergeist destroys the artifact it just used).
func (e Use) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Max; i++ {
		cands := usableCards(ctx, e.Target.Select(ctx), e.EvenUnusable)
		if len(cands) == 0 {
			return
		}
		id, ok := ctx.ChooseCard("Choose a card to use", cands)
		if !ok {
			return
		}
		ctx.It, ctx.HasIt = id, true
		useCard(ctx, id)
	}
}

// usableCards keeps the members that this Use can offer. A creature is offered
// while ready; an artifact is offered while ready and holding an Action. When
// evenUnusable is set the artifact filter drops away entirely — every artifact is
// offered, ready or exhausted, ability or none — because the effect's payload is
// what follows the use (Poltergeist destroys the chosen artifact).
func usableCards(ctx *EffectContext, ids []LocalID, evenUnusable bool) []LocalID {
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		switch ctx.Resolver.TypeOf(id) {
		case Creature:
			if !ctx.Resolver.Exhausted(id) {
				out = append(out, id)
			}
		case Artifact:
			if evenUnusable {
				out = append(out, id)
			} else if !ctx.Resolver.Exhausted(id) &&
				ctx.Resolver.HasTrigger(id, TriggerAction) {
				out = append(out, id)
			}
		}
	}
	return out
}

// useCard resolves a chosen card's use — an artifact fires its Action, a creature
// is used by choosing how (reap, fight, or Action). Either way the ability
// resolves for the effect's controller, so a card they do not control is used as
// if it were theirs. An artifact that cannot act — exhausted, or with no Action,
// both reachable only under EvenUnusable — is a no-op here; the effect that
// offered it still acts on it.
func useCard(ctx *EffectContext, id LocalID) {
	if ctx.Resolver.TypeOf(id) == Artifact {
		if !ctx.Resolver.Exhausted(id) && ctx.Resolver.HasTrigger(id, TriggerAction) {
			ctx.Resolver.UseActionOf(ctx.Controller, id)
		}
		return
	}
	UseVerb{}.Apply(ctx, id)
}
