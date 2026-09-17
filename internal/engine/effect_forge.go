package engine

import "fmt"

// ForgeKey has the controller forge a key outside the normal start-of-turn step.
// By default they pay the current key cost, if they can afford it; FreeOfCost
// forges without paying. Both paths fire "after you forge a key" abilities and,
// on the final key, win the game. A forge that lands purges the card that made it
// — every forge outside the normal step spends its source (a Vactrol divergence;
// see docs/keyforge-divergences.md).
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
	// Discount switches ReducedBy from trimming the Extra surcharge to discounting
	// the current key cost itself, so the forge can land below the current cost (down
	// to 0) — Desire reaps to forge at current cost reduced by 1 for each friendly Sin
	// creature. It reads with no "+N" surcharge and floors the whole cost at 0.
	Discount bool
	// Or switches Extra to an alternate surcharge when a condition holds, so the card
	// reads "forge a key at +6 Æmber current cost, or +2 if …" instead of a two-armed
	// Otherwise branch (rule 22).
	Or OrAmount
}

// validate rejects a reduction with nothing to reduce.
func (e ForgeKey) validate() error {
	if e.ReducedBy != nil && e.Extra == 0 && !e.Discount {
		return fmt.Errorf("ForgeKey: ReducedBy needs an Extra cost to reduce")
	}
	if e.Discount && e.ReducedBy == nil {
		return fmt.Errorf("ForgeKey: Discount needs a ReducedBy count")
	}
	if e.Discount && e.Extra != 0 {
		return fmt.Errorf(
			"ForgeKey: a Discount forge reduces the current cost, so it carries no Extra",
		)
	}
	if e.FreeOfCost && e.Extra != 0 {
		return fmt.Errorf("ForgeKey: a free forge cannot also cost Extra")
	}
	if e.Or.set() && e.FreeOfCost {
		return fmt.Errorf("ForgeKey: a free forge cannot also carry an Or surcharge")
	}
	return e.Or.validate()
}

// Text renders the effect. The forge gates a self-purge: the card that made it is
// spent only if a key is actually forged.
func (e ForgeKey) Text() string {
	var body string
	switch {
	case e.FreeOfCost:
		body = "forge a key at no cost"
	case e.Discount:
		body = fmt.Sprintf(
			"forge a key at current cost, reduced by 1 Æmber for each %s",
			e.ReducedBy.CountText(),
		)
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
	return body + " -> purge " + SelfName
}

// Resolve forges one key for the controller if affordable, then purges the source
// card when a key was actually forged.
func (e ForgeKey) Resolve(ctx *EffectContext) {
	var forged bool
	if e.FreeOfCost {
		forged = ctx.Resolver.ForgeKeyFree(ctx.Controller)
	} else {
		extra := e.Extra
		if e.Or.set() {
			extra = e.Or.pick(e.Extra, ctx)
		}
		if e.ReducedBy != nil {
			extra -= e.ReducedBy.Value(ctx)
		}
		// A Discount reduces the current cost itself, so its surcharge may go
		// negative; forgeKeyAtExtraCost floors the whole cost at 0. A surcharge
		// reduction (Key Abduction) instead floors at the current cost here.
		if !e.Discount {
			extra = max(extra, 0)
		}
		forged = ctx.Resolver.ForgeKeyAtExtraCost(ctx.Controller, extra)
	}
	if forged {
		PurgeSource{}.Resolve(ctx)
	}
}

// RaiseKeyCost makes a player's keys cost Amount more Æmber for the Duration —
// Lash of Broken Dreams taxes the opponent during their next turn. It mirrors
// LowerKeyCost's windows: OpponentNextTurn stays dormant until the affected
// player's own next turn, whoever plays in between (the window every "during
// your opponent's next turn" surcharge wants); EndOfPlayerNextTurn is instead
// live the moment it resolves and again on the affected player's next turn, so
// a forge forced this turn (Keyfrog) already pays the surcharge; and
// RemainderOfPlayerTurn bites at once and lifts when the current turn ends.
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
	case OpponentNextTurn, RemainderOfPlayerTurn, EndOfPlayerNextTurn:
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
	if e.Duration == RemainderOfPlayerTurn {
		return fmt.Sprintf("%s keys cost +%d Æmber for the remainder of the turn", whose, e.Amount)
	}
	return fmt.Sprintf("keys cost +%d Æmber during %s next turn", e.Amount, whose)
}

// Resolve arms the surcharge on the named player for the Duration.
func (e RaiseKeyCost) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	switch e.Duration {
	case OpponentNextTurn:
		ctx.Resolver.RaiseKeyCostNextTurn(p, e.Amount, ctx.Source)
	case RemainderOfPlayerTurn:
		ctx.Resolver.RaiseKeyCostThisTurn(p, e.Amount, ctx.Source)
	case EndOfPlayerNextTurn:
		// Live the moment it resolves and again on the affected player's next turn,
		// so a forge forced this turn (Keyfrog) already pays the surcharge.
		ctx.Resolver.RaiseKeyCostThisTurn(p, e.Amount, ctx.Source)
		ctx.Resolver.RaiseKeyCostNextTurn(p, e.Amount, ctx.Source)
	}
}

// RaiseKeyCostPerHouseCreature raises a player's key cost by Amount for each
// creature of House in play, measured live while the surcharge is active — Waking
// Nightmare taxes the opponent +1 per Dis creature during their next turn. Unlike
// RaiseKeyCost's fixed bump, the surcharge is recomputed at each forge, so a Dis
// creature entering or leaving during the taxed turn changes what a key costs.
//
// It only arms the next-turn window (OpponentNextTurn); a counted surcharge for
// the current turn has no card to want it yet.
type RaiseKeyCostPerHouseCreature struct {
	Player   Player
	Amount   int
	House    House
	Duration Duration
}

// validate requires a player, a positive raise, a house, and the OpponentNextTurn
// window (the only counted surcharge window any card wants).
func (e RaiseKeyCostPerHouseCreature) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("RaiseKeyCostPerHouseCreature")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("RaiseKeyCostPerHouseCreature: Amount must be positive")
	}
	if e.House == HouseNone {
		return fmt.Errorf("RaiseKeyCostPerHouseCreature: House must be set")
	}
	switch e.Duration {
	case OpponentNextTurn:
		return nil
	case durationUnset:
		return errUnsetDuration("RaiseKeyCostPerHouseCreature")
	default:
		return fmt.Errorf(
			"RaiseKeyCostPerHouseCreature: Duration %v is not supported; "+
				"a counted surcharge only arms OpponentNextTurn",
			e.Duration,
		)
	}
}

// Text renders the effect, e.g. "keys cost +1 Æmber for each Dis creature in play
// during your opponent's next turn".
func (e RaiseKeyCostPerHouseCreature) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	return fmt.Sprintf(
		"keys cost +%d Æmber for each %s creature in play during %s next turn",
		e.Amount, e.House, whose,
	)
}

// Resolve arms the counted surcharge on the named player for their next turn.
func (e RaiseKeyCostPerHouseCreature) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	ctx.Resolver.RaiseKeyCostPerHouseNextTurn(p, e.Amount, e.House, ctx.Source)
}

// LowerKeyCost makes a player's keys cost Amount less Æmber for the Duration —
// We Can ALL Win drops each player's keys by 2 until the end of the controller's
// next turn. Amount is written positive and rendered "-N". It shares the key
// surcharge bars with RaiseKeyCost (a lower is a negative bump), so a coexisting
// lower and raise on the same player sum, and the key cost read floors the total
// at 0 before a key is forged.
//
// Player may be EachPlayer, lowering both players' keys at once.
// EndOfPlayerNextTurn is live the moment it resolves, unlike RaiseKeyCost's
// OpponentNextTurn, which waits for the affected player's next turn
// before it bites.
type LowerKeyCost struct {
	Player   Player
	Amount   int
	Duration Duration
}

// validate requires a player, a positive drop, and a key surcharge window.
func (e LowerKeyCost) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("LowerKeyCost")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("LowerKeyCost: Amount must be positive")
	}
	switch e.Duration {
	case RemainderOfPlayerTurn, OpponentNextTurn, EndOfPlayerNextTurn:
		return nil
	case durationUnset:
		return errUnsetDuration("LowerKeyCost")
	default:
		return fmt.Errorf(
			"LowerKeyCost: Duration %v is not a key surcharge window", e.Duration,
		)
	}
}

// Text renders the effect, e.g. "each player's keys cost -2 Æmber until the end
// of your next turn".
func (e LowerKeyCost) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	if e.Player == EachPlayer {
		return fmt.Sprintf("each player's keys cost -%d Æmber %s", e.Amount, e.window())
	}
	return fmt.Sprintf("%s keys cost -%d Æmber %s", whose, e.Amount, e.window())
}

// window renders the duration clause. "your" always names the controller, whose
// next turn ends the window even when EachPlayer lowers both sides.
func (e LowerKeyCost) window() string {
	switch e.Duration {
	case RemainderOfPlayerTurn:
		return "for the remainder of the turn"
	case OpponentNextTurn:
		return "during your next turn"
	default: // EndOfPlayerNextTurn
		return "until the end of your next turn"
	}
}

// Resolve arms the drop on each affected player for the Duration.
func (e LowerKeyCost) Resolve(ctx *EffectContext) {
	if e.Player == EachPlayer {
		for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
			e.arm(ctx, p)
		}
		return
	}
	e.arm(ctx, ctx.PlayerFor(e.Player))
}

// arm records the negative bump on player p. EndOfPlayerNextTurn sets both the
// current turn and next turn bars so the drop is live now and again on p's next
// turn.
func (e LowerKeyCost) arm(ctx *EffectContext, p int) {
	if e.Duration == RemainderOfPlayerTurn || e.Duration == EndOfPlayerNextTurn {
		ctx.Resolver.RaiseKeyCostThisTurn(p, -e.Amount, ctx.Source)
	}
	if e.Duration == OpponentNextTurn || e.Duration == EndOfPlayerNextTurn {
		ctx.Resolver.RaiseKeyCostNextTurn(p, -e.Amount, ctx.Source)
	}
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
func (e UnforgeKey) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate takes one key back and reports whether a key was actually removed, so
// a Then can hang off the unforge succeeding (Key Hammer).
func (e UnforgeKey) resolveGate(ctx *EffectContext) bool {
	return ctx.Resolver.UnforgeKey(ctx.PlayerFor(e.Player))
}

// SkipForgePhase makes a player skip their "forge a key" phase at the start of their
// next turn (Miasma).
type SkipForgePhase struct {
	Player Player
}

// validate rejects a SkipForgePhase whose player was left unset.
func (e SkipForgePhase) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("SkipForgePhase")
	}
	return nil
}

// Text renders the effect, e.g. `your opponent skips the "forge a key" phase during
// their next turn`.
func (e SkipForgePhase) Text() string {
	who, whose, verb := "you", "your", "skip"
	if e.Player == Opponent {
		who, whose, verb = "your opponent", "their", "skips"
	}
	return fmt.Sprintf("%s %s the %q phase during %s next turn", who, verb, "forge a key", whose)
}

// Resolve arms the skip on the chosen player's next turn.
func (e SkipForgePhase) Resolve(ctx *EffectContext) {
	ctx.Resolver.SkipForgePhaseNextTurn(ctx.PlayerFor(e.Player), ctx.Source)
}

// CancelForge cancels the opponent's key forge in progress — the forge does not
// happen and no Æmber is spent (Keyforgery). It is the forge counterpart to
// CancelFight: it resolves inside the before-forge window, where beforeForgePrevented
// reads the cancellation and skips the forge. The card whose ability cancels the
// forge is the source, so the log names it.
type CancelForge struct{}

// Text renders the effect.
func (CancelForge) Text() string { return "they do not forge that key" }

// Resolve records the prevented forge and cancels it.
func (CancelForge) Resolve(ctx *EffectContext) {
	ctx.Resolver.Record(KeyForgePrevented{Player: ctx.Opponent(), By: ctx.Source})
	ctx.Resolver.CancelCurrentForge()
}
