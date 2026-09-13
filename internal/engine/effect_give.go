package engine

import "fmt"

// GiveAember moves Æmber from the controller's opponent into the controller's pool
// — "your opponent gives you Æmber". Unlike a steal it is not theft, and unlike a
// capture the Æmber lands in a pool rather than on a creature; it simply changes
// hands. A fixed Amount is a toll's price (Tentacus, Customs Office charge 1 to act
// with an artifact); All hands over the opponent's whole pool (Interdimensional
// Graft, after the opponent forges). The opponent can give only what their pool
// holds, so the amount is capped at their current Æmber.
type GiveAember struct {
	// Amount is the fixed Æmber the opponent gives; All gives their whole pool instead.
	Amount int
	// All gives the opponent's entire pool rather than a fixed Amount.
	All bool
}

// validate requires exactly one of a fixed Amount or All.
func (e GiveAember) validate() error {
	if err := errAmountOr("GiveAember", "All", e.Amount, e.All); err != nil {
		return err
	}
	if e.Amount == 0 && !e.All {
		return fmt.Errorf("GiveAember: set Amount or All")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent gives you 1 Æmber" or "your opponent
// gives you all their Æmber".
func (e GiveAember) Text() string {
	if e.All {
		return "your opponent gives you all their Æmber"
	}
	return fmt.Sprintf("your opponent gives you %d Æmber", e.Amount)
}

// Resolve moves Æmber from the opponent's pool into the controller's, capped at what
// the opponent holds.
func (e GiveAember) Resolve(ctx *EffectContext) {
	giver := ctx.Opponent()
	receiver := ctx.Controller
	amount := e.Amount
	if e.All {
		amount = ctx.Resolver.Aember(giver)
	}
	if amount > ctx.Resolver.Aember(giver) {
		amount = ctx.Resolver.Aember(giver)
	}
	if amount <= 0 {
		return
	}
	ctx.Resolver.SetAember(giver, ctx.Resolver.Aember(giver)-amount)
	ctx.Resolver.SetAember(receiver, ctx.Resolver.Aember(receiver)+amount)
	ctx.Resolver.Record(AemberGiven{Giver: giver, Receiver: receiver, Amount: amount})
}
