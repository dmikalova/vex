package engine

import "fmt"

// TakeControl moves a card to the ability controller's play area and makes that
// player its controller. With no Target it takes the host creature and reverts
// when the resolving Upgrade leaves play (Collar of Subordination); with a Target
// it takes each chosen creature or artifact. A creature enters the controller's
// battleline and its control reverts when the source card leaves play (Harland
// Mindlock); an artifact enters the controller's artifact row permanently
// (Sneklifter). The Duration says which — UntilThisLeavesPlay pairs with creatures,
// Permanent with artifacts — and drives the rendered "until ... leaves play"
// clause. Ownership stays fixed and still decides which pile the card returns to
// when it leaves play. The last card taken is left in context (ctx.It) so a
// following effect can act on it.
type TakeControl struct {
	// Target picks which cards to take. The zero value means "this creature" — the
	// host of the resolving Upgrade (Collar of Subordination).
	Target   Target
	Duration Duration
}

// validate requires one of the supported durations.
func (e TakeControl) validate() error {
	switch e.Duration {
	case UntilThisLeavesPlay, Permanent:
		return nil
	default:
		return fmt.Errorf("TakeControl: duration must be UntilThisLeavesPlay or Permanent")
	}
}

// Text renders the control change. A reverting form names the card whose leaving
// play reverts it; a Permanent form omits the "until ... leaves play" clause.
func (e TakeControl) Text() string {
	if !e.Target.valid() {
		return "take control of this creature until " + UpgradeName + " leaves play"
	}
	if e.Duration == Permanent {
		return "take control of " + e.Target.Text()
	}
	return "take control of " + e.Target.Text() + " until " + SelfName + " leaves play"
}

// Resolve changes control to the player resolving the ability. A creature moves to
// their battleline anchored to the source (reverting when it leaves play); an
// artifact moves to their artifact row permanently. The host-creature form (no
// Target) is Collar's upgrade-anchored control. The last card taken is recorded as
// "it" for a following effect (Sneklifter's house reassignment).
func (e TakeControl) Resolve(ctx *EffectContext) {
	if !e.Target.valid() {
		ctx.Resolver.TakeControl(ctx.Source, ctx.Controller, ctx.Upgrade)
		ctx.It, ctx.HasIt = ctx.Source, true
		return
	}
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.IsCreature(id) {
			ctx.Resolver.TakeControl(id, ctx.Controller, ctx.Source)
		} else {
			ctx.Resolver.TakeControlOfArtifact(id, ctx.Controller)
		}
		ctx.It, ctx.HasIt = id, true
	}
}
