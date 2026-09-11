package engine

import "fmt"

// Ward applies the ward status. It is a simple "verb the target" effect, so a
// ward that runs beside another status change on the same target folds into one
// phrase in a Sequence (see combinable).

// A ward is a one-shot shield placed on a creature. The first time a warded
// creature would be dealt damage or would leave play, that damage or removal is
// absorbed and the ward is spent instead. Ward intercepts every removal, even the
// controller's own; it covers only damage and leaving play. Warding applies this
// status to each creature the effect targets. When Amount is set, the controller
// chooses that many creatures from the target pool to ward — "ward 2 friendly
// creatures" (Imperium).
type Ward struct {
	Target Target
	Amount int
}

// validate requires an explicit target.
func (e Ward) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Ward")
	}
	return nil
}

func (e Ward) verb() string       { return "ward" }
func (e Ward) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "ward each friendly creature" or, when Amount is
// set, "ward 2 friendly creatures".
func (e Ward) Text() string {
	if e.Amount > 1 {
		noun := singularNoun(e.Target.Text())
		return fmt.Sprintf("%s %d %ss", e.verb(), e.Amount, noun)
	}
	return e.verb() + " " + e.targetText()
}

// applyWard wards one creature. A creature already warded still gets a log line —
// the source still had to choose it — just without a state change.
func (e Ward) applyWard(ctx *EffectContext, id LocalID) {
	if ctx.Resolver.Warded(id) {
		ctx.Resolver.Record(CreatureWarded{Creature: id, By: ctx.Source, AlreadyWarded: true})
		return
	}
	ctx.Resolver.SetWarded(id, true)
	ctx.Resolver.Record(CreatureWarded{Creature: id, By: ctx.Source})
}

// Resolve wards the selected creatures. With Amount set, the controller chooses
// that many distinct creatures from the target pool; otherwise every targeted
// creature is warded.
func (e Ward) Resolve(ctx *EffectContext) {
	if e.Amount > 0 {
		ctx.previewBadge(SelectionBadge{Icon: WardIcon})
		defer ctx.previewBadge(SelectionBadge{})
		chosen := make([]LocalID, 0, e.Amount)
		for i := 0; i < e.Amount; i++ {
			remaining := make([]LocalID, 0)
			for _, id := range e.Target.Select(ctx) {
				picked := false
				for _, c := range chosen {
					if c == id {
						picked = true
						break
					}
				}
				if !picked {
					remaining = append(remaining, id)
				}
			}
			id, ok := ctx.ChooseCard("Choose a creature to ward", remaining)
			if !ok {
				break
			}
			chosen = append(chosen, id)
		}
		for _, id := range chosen {
			e.applyWard(ctx, id)
		}
		return
	}
	for _, id := range e.Target.Select(ctx) {
		e.applyWard(ctx, id)
	}
}
