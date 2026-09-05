package engine

import "fmt"

// Stealing Æmber moves it from the opponent's pool into your own. You can only
// steal as much Æmber as the opponent actually has. How much is stolen is either
// a fixed Amount — optionally multiplied by a Per count — or a By share of the
// opponent's pool (By: AllBut(6) leaves them exactly six).
type StealAember struct {
	// Amount is the fixed Æmber to steal.
	Amount int
	// By steals a share of the opponent's pool instead of a fixed Amount.
	By Loss
	// Per multiplies Amount by a running count.
	Per Count
	// Or switches Amount to an alternate when a condition holds, so the card reads
	// "steal 1 Æmber, or 2 if …" instead of a two-armed Otherwise branch (rule 22).
	Or OrAmount
	// Player is who takes the Æmber. Unset (the usual case) means the controller
	// steals from their opponent; Opponent turns the theft around, so the opponent
	// takes from the controller (Magda the Rat as she leaves play).
	Player Player
}

// sides names the two sides of the theft: the controller robs their opponent
// unless Player turns it around.
func (e StealAember) sides(ctx *EffectContext) (player, opponent int) {
	if e.Player == Opponent {
		return ctx.Opponent(), ctx.Controller
	}
	return ctx.Controller, ctx.Opponent()
}

// validate rejects a StealAember that sets both a fixed Amount and a By share
// (the two are different ways to say how much to steal), then checks the Or guard.
func (e StealAember) validate() error {
	if e.Amount != 0 && e.By != nil {
		return fmt.Errorf("StealAember: set Amount or By, not both (got Amount=%d)", e.Amount)
	}
	if e.Or.set() && e.By != nil {
		return fmt.Errorf("StealAember: Or and By both set the amount; use one")
	}
	return e.Or.validate()
}

// Text renders the effect, e.g. "steal 1 Æmber", "steal all but 6 Æmber from your
// opponent", or "for each friendly ready Mars creature, steal 1 Æmber".
func (e StealAember) Text() string {
	object := fmt.Sprintf("%d Æmber", e.Amount)
	if e.By != nil {
		// A share is measured against the opponent's pool, so name whose pool it is:
		// "all but 6 Æmber" on its own does not say.
		object = e.By.object("your opponent's") + " from your opponent"
	}
	verb := "steal "
	if e.Player == Opponent {
		verb = "your opponent steals "
	}
	body := forEach(e.Per, verb+object)
	if e.Or.set() {
		body += e.Or.tail(fmt.Sprintf("%d", e.Or.Amount))
	}
	return body
}

// Resolve moves the Æmber from the opponent's pool to the controller's.
func (e StealAember) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the steal and reports whether any Æmber actually moved, so
// it can gate a following effect (Nerve Blast's "steal 1 -> deal 2 damage").
func (e StealAember) resolveGate(ctx *EffectContext) bool {
	player, opponent := e.sides(ctx)
	if ctx.Resolver.AemberProtected(opponent) {
		return false
	}
	amt := min(e.amount(ctx, opponent), ctx.Resolver.Aember(opponent))
	ctx.Resolver.SetAember(opponent, ctx.Resolver.Aember(opponent)-amt)
	ctx.Resolver.SetAember(player, ctx.Resolver.Aember(player)+amt)
	// Credit the card when the controller is the one stealing, so the line reads
	// from the card's perspective; a turned-around steal (your opponent steals as
	// the card leaves play) has no such agent, so it stays player-attributed.
	ctx.Resolver.Record(AemberStolen{
		Player:    player,
		From:      opponent,
		Amount:    amt,
		Source:    ctx.Source,
		HasSource: e.Player != Opponent,
	})
	return amt > 0
}

// amount is how much to steal: the By share of the opponent's pool when set,
// otherwise the fixed Amount scaled by Per.
func (e StealAember) amount(ctx *EffectContext, opp int) int {
	if e.By != nil {
		return e.By.lose(ctx.Resolver.Aember(opp))
	}
	base := e.Or.pick(e.Amount, ctx)
	if e.Per != nil {
		return base * e.Per.Value(ctx)
	}
	return base
}
