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
	// Fraction moves a share of the source's Æmber instead of a fixed Amount —
	// Patronage moves half the Æmber off the chosen creature, rounding up. The zero
	// value is unset; leave it so to use Amount or All.
	Fraction Fraction
	// From selects the eligible source cards; the chosen source must carry Æmber.
	From Target
	// To is the destination pool — Controller's or Opponent's. Leave unset (and set
	// Onto) to move the Æmber onto a card instead.
	To Player
	// Onto selects the card the Æmber moves onto. Leave unset (and set To) to move
	// the Æmber into a pool instead.
	Onto Target
	// Bind selects sources without the "must carry Æmber" pre-filter, so a Refinement
	// on From (MostPowerful) picks from the true candidate set rather than only the
	// Æmber-bearing ones, and leaves the moved-from creature in context (ctx.It) for
	// a following effect — Sic Semper Tyrannosaurus empties the most powerful
	// creature and then destroys that same creature. A source holding no Æmber moves
	// none. When unset, sources are pre-filtered to Æmber-bearers and nothing is
	// bound.
	Bind bool
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
	if err := errAmountOr("MoveAember", "Fraction", e.Amount, e.Fraction.valid()); err != nil {
		return err
	}
	if e.All && e.Fraction.valid() {
		return fmt.Errorf("MoveAember: set All or Fraction, not both")
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
	if e.Fraction.valid() {
		return fmt.Sprintf("move %s the \u00c6mber from %s to %s, %s",
			e.Fraction.word(), e.From.Text(), e.destText(), e.Fraction.roundingPhrase())
	}
	amount := fmt.Sprintf("%d", e.amount())
	if e.All {
		amount = "all"
	}
	return fmt.Sprintf("move %s \u00c6mber from %s to %s", amount, e.From.Text(), e.destText())
}

// Resolve moves Æmber from the chosen source(s) to the destination. Without Bind
// the source choice is restricted to cards carrying Æmber, so it never offers an
// empty card; with Bind the chooser picks from the true candidate set (a source
// holding no Æmber simply moves none) and the moved-from creature is left in
// context (ctx.It). A source holding fewer than Amount moves all it has.
func (e MoveAember) Resolve(ctx *EffectContext) {
	source := e.From
	if !e.Bind {
		source = source.WithAember()
	}
	sources := source.Select(ctx)
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
		if have := ctx.Resolver.AmberOn(from); e.Fraction.valid() {
			moved = e.Fraction.of(have)
		} else if e.All || moved > have {
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
	if e.Bind {
		ctx.It, ctx.HasIt = sources[len(sources)-1], true
	}
}
