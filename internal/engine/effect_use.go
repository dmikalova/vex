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
	return fmt.Sprintf("use %d %ss, one at a time", e.Max, noun)
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
		cands := usableCards(ctx, e.Target.Select(ctx))
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

// usableCards keeps the ready cards that can be used by an ability right now: any
// creature, or an artifact with an Action.
func usableCards(ctx *EffectContext, ids []LocalID) []LocalID {
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Exhausted(id) {
			continue
		}
		switch ctx.Resolver.TypeOf(id) {
		case Creature:
			out = append(out, id)
		case Artifact:
			if ctx.Resolver.HasAction(id) {
				out = append(out, id)
			}
		}
	}
	return out
}

// useCard resolves a chosen card's use — an artifact fires its Action, a creature
// is used by choosing how (reap, fight, or Action). Either way the ability
// resolves for the effect's controller, so a card they do not control is used as
// if it were theirs.
func useCard(ctx *EffectContext, id LocalID) {
	if ctx.Resolver.TypeOf(id) == Artifact {
		ctx.Resolver.UseActionOf(ctx.Controller, id)
		return
	}
	UseVerb{}.Apply(ctx, id)
}

// Sentence resolves its child normally but renders it as a complete sentence, so a
// following Sequence child starts a new sentence instead of joining with ", and".
type Sentence struct {
	Effect Effect
}

// Text renders the child as a punctuated sentence.
func (e Sentence) Text() string { return punctuate(e.Effect.Text()) }

// Resolve resolves the child effect.
func (e Sentence) Resolve(ctx *EffectContext) { e.Effect.Resolve(ctx) }

// validate surfaces a configuration error from the child effect.
func (e Sentence) validate() error { return validateEffect(e.Effect) }
