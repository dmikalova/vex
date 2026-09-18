package engine

import "fmt"

// TakeControl moves a card to the ability controller's play area and makes that
// player its controller. With no Target it takes the host creature and reverts
// when the resolving Upgrade leaves play (Collar of Subordination); with a Target
// it takes each chosen creature or artifact. A creature enters the controller's
// battleline and its control reverts when the source card leaves play (Harland
// Mindlock); an artifact enters the controller's artifact row permanently
// (Sneklifter). The Duration says which — UntilThisLeavesPlay pairs with creatures,
// Forever with artifacts — and drives the rendered "until ... leaves play"
// clause. Ownership stays fixed and still decides which pile the card returns to
// when it leaves play. The last card taken is left in context (ctx.It) so a
// following effect can act on it. ToOpponent gives control away instead of taking
// it — the ability's controller hands the targeted card to their opponent
// (Spangler Box), rather than taking it for themselves.
type TakeControl struct {
	// Target picks which cards to take. The zero value means "this creature" — the
	// host of the resolving Upgrade (Collar of Subordination).
	Target Target
	// Duration says which reversion applies: UntilThisLeavesPlay for creatures,
	// Forever for artifacts. ToOpponent hands the card to the opponent instead of
	// taking it (Spangler Box).
	Duration   Duration
	ToOpponent bool
	// AndExhaust exhausts the seized creature as part of the same effect, rendering
	// "... and exhaust it" (Lord Invidius).
	AndExhaust bool
}

// validate requires one of the supported durations.
func (e TakeControl) validate() error {
	switch e.Duration {
	case UntilThisLeavesPlay, Forever:
		return nil
	default:
		return fmt.Errorf("TakeControl: duration must be UntilThisLeavesPlay or Forever")
	}
}

// Text renders the control change. A reverting form names the card whose leaving
// play reverts it; a Forever form omits the "until ... leaves play" clause.
// ToOpponent renders as the opponent gaining control rather than the resolving
// player taking it.
func (e TakeControl) Text() string {
	if !e.Target.valid() {
		return "take control of this creature " + durationClause(UntilThisLeavesPlay, CardName)
	}
	if e.ToOpponent {
		return "your opponent gains control of " + e.Target.Text()
	}
	text := "take control of " + e.Target.Text()
	if e.Duration != Forever {
		text += " " + durationClause(UntilThisLeavesPlay, SelfName)
	}
	if e.AndExhaust {
		text += " and exhaust it"
	}
	return text
}

// Resolve changes control to the player resolving the ability, or to their
// opponent when ToOpponent is set. A creature moves to the new controller's
// battleline anchored to the source (reverting when it leaves play); an artifact
// moves to their artifact row permanently. The host-creature form (no Target) is
// Collar's upgrade-anchored control. The last card taken is recorded as "it" for a
// following effect (Sneklifter's house reassignment).
func (e TakeControl) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate changes control to the player resolving the ability, or to their
// opponent when ToOpponent is set, and reports whether any card actually moved so
// it can be the First of a Then (Anahita the Trader gives an artifact away -> the
// opponent gives 2 Æmber). A creature moves to the new controller's
// battleline anchored to the source (reverting when it leaves play); an artifact
// moves to their artifact row permanently. The host-creature form (no Target) is
// Collar's upgrade-anchored control. The last card taken is recorded as "it" for a
// following effect (Sneklifter's house reassignment).
func (e TakeControl) resolveGate(ctx *EffectContext) bool {
	if !e.Target.valid() {
		ctx.Resolver.TakeControl(ctx.Source, ctx.Controller, ctx.Upgrade)
		placeSeizedOnFlank(ctx, ctx.Controller, ctx.Source)
		ctx.It, ctx.HasIt = ctx.Source, true
		return true
	}
	newController := ctx.Controller
	if e.ToOpponent {
		newController = ctx.Opponent()
	}
	moved := false
	for _, id := range e.Target.Select(ctx) {
		// A Forever control never reverts to a leaving source, so it anchors to the
		// seized card itself; an UntilThisLeavesPlay control anchors to the resolving
		// card and reverts when that card leaves play.
		source := ctx.Source
		if e.Duration == Forever {
			source = id
		}
		ctx.Resolver.TakeControl(id, newController, source)
		placeSeizedOnFlank(ctx, newController, id)
		if e.AndExhaust {
			ctx.Resolver.SetExhausted(id, true)
		}
		ctx.It, ctx.HasIt = id, true
		moved = true
	}
	return moved
}

// placeSeizedOnFlank places a seized creature on the flank the active player
// chooses (Harland Mindlock, Collar of Subordination). The active player always
// makes the placement call, even for a give-to-opponent seize where the gaining
// controller is the opponent. The engine never assumes a flank: with another
// creature already in the taker's battleline the flank is asked, and with only one
// home there is nothing to ask. A seized artifact has no flank, so it is left where
// control placed it.
func placeSeizedOnFlank(ctx *EffectContext, controller int, id LocalID) {
	if !ctx.Resolver.IsCreature(id) {
		return
	}
	if len(ctx.Resolver.Battleline(controller)) <= 1 {
		return
	}
	right := ctx.Resolver.ChooseOption(ctx.Resolver.ActivePlayer(), ctx.Source,
		"Choose a flank", []string{FlankLeftLabel, FlankRightLabel}) == 1
	ctx.Resolver.MoveToFlank(id, right)
}
