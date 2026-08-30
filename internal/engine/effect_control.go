package engine

import "fmt"

// TakeControl moves this creature to the ability controller's battleline and makes
// that player its controller for the given Duration. Collar of Subordination takes
// control UntilThisLeavesPlay: the change lasts until the Upgrade granting it leaves
// play. In KeyForge terms, control decides whose battleline the creature is in and
// who may use it; ownership stays fixed and still decides which discard, hand, deck,
// archives, or purge pile the card returns to when it leaves play.
type TakeControl struct {
	Duration Duration
}

// validate requires the (only supported) UntilThisLeavesPlay duration.
func (e TakeControl) validate() error {
	if e.Duration != UntilThisLeavesPlay {
		return fmt.Errorf("TakeControl: duration must be UntilThisLeavesPlay")
	}
	return nil
}

// Text renders Collar of Subordination's control-changing Play ability.
func (e TakeControl) Text() string {
	return "take control of this creature until " + UpgradeName + " leaves play"
}

// Resolve changes control of this creature to the player resolving the ability,
// recording the resolving Upgrade so control reverts when the Upgrade leaves play.
func (e TakeControl) Resolve(ctx *EffectContext) {
	ctx.Resolver.TakeControl(ctx.Source, ctx.Controller, ctx.Upgrade)
}
