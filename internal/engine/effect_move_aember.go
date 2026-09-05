package engine

import "fmt"

// MoveAember moves Æmber off a card the controller picks and deposits it
// elsewhere — into a player's pool (Selwyn the Fence moves 1 from a friendly card
// to your pool) or onto another card. From names the eligible sources; only those
// carrying Æmber are offered, so nothing happens when none do. Exactly one
// destination is set: To for a pool, Onto for a card.
type MoveAember struct {
	// Amount of Æmber to move from the chosen source; 0 reads as 1. When the source
	// holds fewer, all of it moves.
	Amount int
	// All moves everything each source carries instead of a fixed Amount — Word of
	// Returning moves all the Æmber off every enemy creature at once.
	All bool
	// From selects the eligible source cards; the chosen source must carry Æmber.
	From Target
	// To is the destination pool — Controller's or Opponent's. Leave unset (and set
	// Onto) to move the Æmber onto a card instead.
	To Player
	// Onto selects the card the Æmber moves onto. Leave unset (and set To) to move
	// the Æmber into a pool instead.
	Onto Target
}

// amount is Amount with the zero value treated as one.
func (e MoveAember) amount() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// toPool reports whether the destination is a pool rather than a card.
func (e MoveAember) toPool() bool { return e.To != playerUnset }

// validate requires a source target and exactly one destination.
func (e MoveAember) validate() error {
	if !e.From.valid() {
		return errUnsetTarget("MoveAember")
	}
	if err := errAmountOr("MoveAember", "All", e.Amount, e.All); err != nil {
		return err
	}
	if e.toPool() == e.Onto.valid() {
		return fmt.Errorf("MoveAember: set exactly one destination (To pool or Onto card)")
	}
	return nil
}

// destText names the destination for the printed phrase.
func (e MoveAember) destText() string {
	if e.toPool() {
		if e.To == Opponent {
			return "your opponent's pool"
		}
		return "your pool"
	}
	return e.Onto.Text()
}

// Text renders the effect, e.g. "move 1 Æmber from a friendly creature or artifact
// to your pool", or "move all Æmber from each enemy creature to your pool".
func (e MoveAember) Text() string {
	amount := fmt.Sprintf("%d", e.amount())
	if e.All {
		amount = "all"
	}
	return fmt.Sprintf("move %s \u00c6mber from %s to %s", amount, e.From.Text(), e.destText())
}

// Resolve moves Æmber from the chosen source(s) to the destination. The source
// choice is restricted to cards carrying Æmber, so it never offers an empty card;
// a source holding fewer than Amount moves all it has.
func (e MoveAember) Resolve(ctx *EffectContext) {
	sources := e.From.WithAember().Select(ctx)
	if len(sources) == 0 {
		return
	}
	var onto LocalID
	if !e.toPool() {
		dest := e.Onto.Select(ctx)
		if len(dest) == 0 {
			return
		}
		onto = dest[0]
	}
	for _, from := range sources {
		moved := e.amount()
		if have := ctx.Resolver.AmberOn(from); e.All || moved > have {
			moved = have
		}
		ctx.Resolver.AddAmberOn(from, -moved)
		if e.toPool() {
			p := ctx.PlayerFor(e.To)
			ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)+moved)
			ctx.Resolver.Record(AemberMovedToPool{
				Player: ctx.Controller,
				From:   from,
				To:     p,
				Amount: moved,
			})
			continue
		}
		ctx.Resolver.AddAmberOn(onto, moved)
		ctx.Resolver.Record(AemberMovedToCard{
			Player: ctx.Controller,
			From:   from,
			To:     onto,
			Amount: moved,
		})
	}
}
