package engine

import "fmt"

// ForgeKey has the controller forge a key outside the normal start-of-turn step.
// By default they pay the current key cost, if they can afford it; FreeOfCost
// forges without paying. Both paths fire "after you forge a key" abilities and,
// on the final key, win the game.
type ForgeKey struct {
	// FreeOfCost forges without paying the key cost.
	FreeOfCost bool
	// Extra raises the cost of this one forge above the current key cost — Key of
	// Darkness forges at +6. It is a surcharge on the forge, not a change to the key
	// cost itself, so it is gone the moment the effect finishes.
	Extra int
	// ReducedBy subtracts a running count from Extra, never below the current key
	// cost — Key Abduction's +9 comes down by 1 for each card in hand.
	ReducedBy Count
	// Or switches Extra to an alternate surcharge when a condition holds, so the card
	// reads "forge a key at +6 Æmber current cost, or +2 if …" instead of a two-armed
	// Otherwise branch (rule 22).
	Or OrAmount
}

// validate rejects a reduction with nothing to reduce.
func (e ForgeKey) validate() error {
	if e.ReducedBy != nil && e.Extra == 0 {
		return fmt.Errorf("ForgeKey: ReducedBy needs an Extra cost to reduce")
	}
	if e.FreeOfCost && e.Extra != 0 {
		return fmt.Errorf("ForgeKey: a free forge cannot also cost Extra")
	}
	if e.Or.set() && e.FreeOfCost {
		return fmt.Errorf("ForgeKey: a free forge cannot also carry an Or surcharge")
	}
	return e.Or.validate()
}

// Text renders the effect.
func (e ForgeKey) Text() string {
	var body string
	switch {
	case e.FreeOfCost:
		return "forge a key at no cost"
	case e.ReducedBy != nil:
		body = fmt.Sprintf(
			"forge a key at +%d Æmber current cost, reduced by 1 Æmber for each %s",
			e.Extra, e.ReducedBy.CountText(),
		)
	case e.Extra != 0:
		body = fmt.Sprintf("forge a key at +%d Æmber current cost", e.Extra)
	default:
		body = "forge a key at current cost"
	}
	if e.Or.set() {
		body += e.Or.tail(fmt.Sprintf("+%d", e.Or.Amount))
	}
	return body
}

// Resolve forges one key for the controller if affordable.
func (e ForgeKey) Resolve(ctx *EffectContext) {
	if e.FreeOfCost {
		ctx.Resolver.ForgeKeyFree(ctx.Controller)
		return
	}
	extra := e.Extra
	if e.Or.set() {
		extra = e.Or.pick(e.Extra, ctx)
	}
	if e.ReducedBy != nil {
		extra -= e.ReducedBy.Value(ctx)
	}
	ctx.Resolver.ForgeKeyAtExtraCost(ctx.Controller, max(extra, 0))
}

// RaiseKeyCost makes a player's keys cost Amount more Æmber for the Duration —
// Lash of Broken Dreams taxes the opponent through their next turn. It mirrors
// the bars (CannotPlay, CannotUse): NextTurn waits for the affected player's own
// next turn whoever plays in between, while EndOfTurn bites at once and lifts
// when the current turn ends.
//
// A surcharge that should last as long as its card is in play is not this
// effect: print it on the card as a KeyCostChange (WithKeyCost), which the key
// cost reads continuously from the cards in play.
type RaiseKeyCost struct {
	Player   Player
	Amount   int
	Duration Duration
}

// validate requires a player, a raise, and a duration this bar can express.
func (e RaiseKeyCost) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("RaiseKeyCost")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("RaiseKeyCost: Amount must be positive")
	}
	switch e.Duration {
	case NextTurn, EndOfTurn:
		return nil
	case durationUnset:
		return errUnsetDuration("RaiseKeyCost")
	default:
		return fmt.Errorf(
			"RaiseKeyCost: Duration %v is not a key surcharge window; "+
				"for a raise that lasts while the card is in play use WithKeyCost",
			e.Duration,
		)
	}
}

// Text renders the effect, e.g. "keys cost +3 Æmber during your opponent's next
// turn".
func (e RaiseKeyCost) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	if e.Duration == EndOfTurn {
		return fmt.Sprintf("%s keys cost +%d Æmber for the remainder of the turn", whose, e.Amount)
	}
	return fmt.Sprintf("keys cost +%d Æmber during %s next turn", e.Amount, whose)
}

// Resolve arms the surcharge on the named player for the Duration.
func (e RaiseKeyCost) Resolve(ctx *EffectContext) {
	switch e.Duration {
	case NextTurn:
		ctx.Resolver.RaiseKeyCostNextTurn(ctx.PlayerFor(e.Player), e.Amount, ctx.Source)
	case EndOfTurn:
		ctx.Resolver.RaiseKeyCostThisTurn(ctx.PlayerFor(e.Player), e.Amount, ctx.Source)
	}
}

// GiveRemainingAemberAfterOpponentForgeKey arms Interdimensional Graft's delayed
// forge penalty: each time the opponent forges a key during their next turn, they
// give their remaining Æmber to the controller. It is durable across this turn's
// end and fires on every forge that turn (a key cheat can forge more than one), then
// expires at the end of that opponent's next turn.
type GiveRemainingAemberAfterOpponentForgeKey struct{}

// Text renders the effect.
func (GiveRemainingAemberAfterOpponentForgeKey) Text() string {
	return "if an opponent forges a key on their next turn, they must give you their remaining Æmber"
}

// Resolve arms the transfer for the opponent's next turn by registering a reaction
// to the opponent's forge, owned by the opponent so it survives this turn, fires on
// each forge during theirs, and clears at the end of their turn.
func (GiveRemainingAemberAfterOpponentForgeKey) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventForgeKey,
		Do:         actGiveRemainingAember,
		Controller: int8(ctx.Opponent()),
	})
}

// UnforgeKey takes a forged key back off a player (Key Hammer). It is the one
// effect that lowers a key count, so it is a node of its own rather than a
// negative ForgeKey: nothing is paid, nothing is refunded, and no "after you
// forge a key" ability fires.
type UnforgeKey struct {
	Player Player
}

// validate rejects an UnforgeKey whose player was left unset.
func (e UnforgeKey) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("UnforgeKey")
	}
	return nil
}

// Text renders the effect, e.g. "unforge one of your opponent's keys".
func (e UnforgeKey) Text() string {
	if e.Player == Opponent {
		return "unforge one of your opponent's keys"
	}
	return "unforge one of your keys"
}

// Resolve takes one key back off the named player.
func (e UnforgeKey) Resolve(ctx *EffectContext) {
	ctx.Resolver.UnforgeKey(ctx.PlayerFor(e.Player))
}
