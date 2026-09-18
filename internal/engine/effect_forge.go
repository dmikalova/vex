package engine

import "fmt"

// ForgeKey has a player forge a key outside the normal start-of-turn step. By
// default the controller forges, paying the current key cost if they can afford it;
// FreeOfCost forges without paying. Player: Opponent instead forces the opponent to
// forge a key at no cost (Turnkey), and the active player — who chooses everything —
// picks its colour; that forced forge spends nothing, so it never purges a source.
// Both controller paths fire "after you forge a key" abilities and, on the final
// key, win the game. A controller forge that lands purges the card that made it —
// every forge outside the normal step spends its source (a Vactrol divergence;
// see docs/keyforge-divergences.md).
type ForgeKey struct {
	// Player is who forges. The zero value forges for the controller; Opponent
	// forces the opponent to forge a key at no cost.
	Player Player
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
	if e.Player == Opponent && !e.FreeOfCost {
		return fmt.Errorf(
			"ForgeKey: an opponent forge is only supported at no cost (FreeOfCost)",
		)
	}
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
// spent only if a key is actually forged. An opponent forge spends nothing, so it
// carries no purge.
func (e ForgeKey) Text() string {
	if e.Player == Opponent {
		return "your opponent forges a key at no cost"
	}
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
// card when a key was actually forged. An opponent forge instead forces the
// opponent to forge for free, spends nothing, and never purges.
func (e ForgeKey) Resolve(ctx *EffectContext) {
	if e.Player == Opponent {
		ctx.Resolver.ForgeKeyFreeForced(ctx.Opponent())
		return
	}
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

// ScheduleOnLeave arms Do to resolve when the source card leaves play, however many
// turns later — Turnkey unforges an opponent's key and, if it does, has the opponent
// forge a key at no cost when Turnkey leaves play. The consequence is held flat as
// an enum-tagged action (ADR 0005), not as a stored effect closure, so Do must be
// one the schedule can carry (scheduledActionOf).
type ScheduleOnLeave struct {
	Do Effect
}

// validate requires a Do the schedule can carry, and a valid Do.
func (e ScheduleOnLeave) validate() error {
	if e.Do == nil {
		return fmt.Errorf("ScheduleOnLeave: Do is required")
	}
	if _, ok := scheduledActionOf(e.Do); !ok {
		return fmt.Errorf("ScheduleOnLeave: %T is not a schedulable effect", e.Do)
	}
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "when <self> leaves play, your opponent forges a
// key at no cost".
func (e ScheduleOnLeave) Text() string {
	return "when " + SelfName + " leaves play, " + e.Do.Text()
}

// Resolve arms the leave-play schedule; the source card's exit resolves it.
func (e ScheduleOnLeave) Resolve(ctx *EffectContext) {
	action, _ := scheduledActionOf(e.Do)
	ctx.Resolver.ScheduleOnLeave(ctx.Source, action)
}

// RaiseKeyCost makes a player's keys cost Amount more Æmber for the Duration —
// Lash of Broken Dreams taxes the opponent during their next turn. Its windows:
// OpponentNextTurn stays dormant until the affected player's own next turn,
// whoever plays in between (the window every "during your opponent's next turn"
// surcharge wants); EndOfPlayerNextTurn is instead live the moment it resolves
// and again on the affected player's next turn, so a forge forced this turn
// (Keyfrog) already pays the surcharge; and RemainderOfPlayerTurn bites at once
// and lifts when the current turn ends.
//
// When House filters, Amount is charged once for each creature it admits in play,
// counted live at each forge rather than frozen (Waking Nightmare, +1 per Dis
// creature). The counted form folds in here rather than as its own node, and only
// the OpponentNextTurn window can carry it. Player may be EachPlayer, taxing both
// players at once.
//
// RaiseKeyCost and LowerKeyCost share this shape, one resolution, and one text:
// a lower is a raise with a negative amount (armKeyCost) and a "-" sign, so the
// two differ only in the sign keyCostText renders.
//
// A surcharge that should last as long as its card is in play is not this
// effect: print it on the card as a KeyCostChange (WithKeyCost), which the key
// cost reads continuously from the cards in play.
type RaiseKeyCost struct {
	Player   Player
	Amount   int
	House    HouseMatcher
	Duration Duration
}

// validate requires a player, a raise, and a duration this bar can express.
func (e RaiseKeyCost) validate() error {
	return validateKeyCost("RaiseKeyCost", e.Player, e.Amount, e.House, e.Duration)
}

// Text renders the effect through keyCostText, e.g. "keys cost +3 Æmber during
// your opponent's next turn". A raise passes the "+" sign; the rest of the
// sentence is shared with LowerKeyCost.
func (e RaiseKeyCost) Text() string {
	return keySurchargeText(e.Player, "+", e.Amount, e.House, e.Duration)
}

// Resolve arms the surcharge on each affected player for the Duration.
func (e RaiseKeyCost) Resolve(ctx *EffectContext) {
	armKeyCost(ctx, e.Player, e.Amount, e.House, e.Duration)
}

// LowerKeyCost makes a player's keys cost Amount less Æmber for the Duration —
// We Can ALL Win drops each player's keys by 2 until the end of the controller's
// next turn. Amount is written positive and rendered "-N". It shares RaiseKeyCost's
// shape and bars (a lower is a negative raise), so a coexisting lower and raise on
// the same player sum, and the key cost read floors the total at 0 before a key is
// forged.
//
// Player may be EachPlayer, lowering both players' keys at once.
// EndOfPlayerNextTurn is live the moment it resolves, unlike RaiseKeyCost's
// OpponentNextTurn, which waits for the affected player's next turn
// before it bites.
type LowerKeyCost struct {
	Player   Player
	Amount   int
	House    HouseMatcher
	Duration Duration
}

// validate requires a player, a positive drop, and a key surcharge window.
func (e LowerKeyCost) validate() error {
	return validateKeyCost("LowerKeyCost", e.Player, e.Amount, e.House, e.Duration)
}

// Text renders the effect through keyCostText, e.g. "each player's keys cost -2
// Æmber until the end of your next turn". A lower passes the "-" sign; the rest of
// the sentence is shared with RaiseKeyCost.
func (e LowerKeyCost) Text() string {
	return keySurchargeText(e.Player, "-", e.Amount, e.House, e.Duration)
}

// Resolve arms the drop on each affected player for the Duration — a negative raise.
func (e LowerKeyCost) Resolve(ctx *EffectContext) {
	armKeyCost(ctx, e.Player, -e.Amount, e.House, e.Duration)
}

// validateKeyCost checks the shared shape of RaiseKeyCost and LowerKeyCost: a set
// player, a positive Amount (each node renders its own sign), and a Duration the
// bars can express. When House filters, the surcharge is counted live per creature
// it admits and only the OpponentNextTurn window can carry it, and the matcher must
// be context-free (a named house or all-but-one) — a surcharge measured live across
// a turn boundary has no resolution context to resolve a chosen or active house.
func validateKeyCost(
	name string,
	player Player,
	amount int,
	house HouseMatcher,
	dur Duration,
) error {
	if !player.valid() {
		return errUnsetPlayer(name)
	}
	if amount <= 0 {
		return fmt.Errorf("%s: Amount must be positive", name)
	}
	if house.filters() {
		switch house.Kind {
		case MatchNamedHouse, MatchExceptHouse:
		default:
			return fmt.Errorf("%s: house matcher %v needs resolution context", name, house.Kind)
		}
		if err := house.validate(); err != nil {
			return err
		}
		if dur == durationUnset {
			return errUnsetDuration(name)
		}
		if dur != OpponentNextTurn {
			return fmt.Errorf("%s: a counted surcharge only arms OpponentNextTurn", name)
		}
		return nil
	}
	switch dur {
	case OpponentNextTurn, RemainderOfPlayerTurn, EndOfPlayerNextTurn:
		return nil
	case durationUnset:
		return errUnsetDuration(name)
	default:
		return fmt.Errorf(
			"%s: Duration %v is not a key surcharge window; for a change that lasts "+
				"while the card is in play use WithKeyCost",
			name, dur,
		)
	}
}

// perCreatureClause renders " for each <house> creature in play" for a counted key
// surcharge (House filters), else "". The clause sits between the amount and the
// window clause, so a caller appends it directly after "Æmber".
func perCreatureClause(house HouseMatcher) string {
	if !house.filters() {
		return ""
	}
	return " for each " + house.qualifyNoun("creature") + " in play"
}

// keySurchargeText renders the sentence RaiseKeyCost and LowerKeyCost share; sign
// is "+" for a raise or "-" for a lower, the only difference between the two. The
// window is framed from the affected player — the opponent's turn reads "your
// opponent's next turn", the controller's or each player's "your next turn". A
// next-turn window already names whose keys are taxed, so the subject stays
// implicit ("keys cost +3 Æmber during your opponent's next turn"); the
// remainder-of-turn window names no player, and EachPlayer taxes both, so those
// name the subject up front ("your keys cost …", "each player's keys cost …"). A
// counted surcharge (House filters) reads "for each <house> creature in play".
func keySurchargeText(
	player Player,
	sign string,
	amount int,
	house HouseMatcher,
	dur Duration,
) string {
	clause := perCreatureClause(house)
	window := windowClause(dur, keySurchargeWindowWhose(player), "")
	if player != EachPlayer && windowNamesAffectedPlayer(dur) {
		return fmt.Sprintf("keys cost %s%d Æmber%s %s", sign, amount, clause, window)
	}
	return fmt.Sprintf(
		"%s keys cost %s%d Æmber%s %s",
		keySurchargeSubject(player), sign, amount, clause, window,
	)
}

// keySurchargeSubject names whose keys a surcharge taxes, for a sentence that
// states its subject up front.
func keySurchargeSubject(player Player) string {
	switch player {
	case Opponent:
		return "your opponent's"
	case EachPlayer:
		return "each player's"
	default:
		return "your"
	}
}

// keySurchargeWindowWhose names the possessive the next-turn window is measured
// against: the opponent's turn for an opponent-only surcharge, the controller's
// otherwise (EachPlayer is framed from the controller, like We Can ALL Win).
func keySurchargeWindowWhose(player Player) string {
	if player == Opponent {
		return "your opponent's"
	}
	return "your"
}

// windowNamesAffectedPlayer reports whether the window clause embeds the affected
// player's possessive — a next-turn window does ("during your opponent's next
// turn"), the remainder-of-turn window does not — so a caller can leave the subject
// implicit when the window already names it.
func windowNamesAffectedPlayer(d Duration) bool {
	switch d {
	case OpponentNextTurn, StartOfPlayerNextTurn, EndOfPlayerNextTurn:
		return true
	default:
		return false
	}
}

// armKeyCost arms the one-turn key surcharge described by a RaiseKeyCost or
// LowerKeyCost on each affected player: signed amount (a lower passes it negative)
// charged flat, or per creature House admits when House filters, over the window
// Duration names. EachPlayer arms both players.
func armKeyCost(ctx *EffectContext, player Player, amount int, house HouseMatcher, dur Duration) {
	if player == EachPlayer {
		for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
			armKeyCostOn(ctx, p, amount, house, dur)
		}
		return
	}
	armKeyCostOn(ctx, ctx.PlayerFor(player), amount, house, dur)
}

// armKeyCostOn arms the surcharge on one player p. A counted surcharge (House
// filters) is stored per house and only arms the next turn; a flat surcharge arms
// the current turn, the next turn, or both — EndOfPlayerNextTurn arms both so it is
// live the moment it resolves and again on p's next turn (Keyfrog).
func armKeyCostOn(ctx *EffectContext, p, amount int, house HouseMatcher, dur Duration) {
	if house.filters() {
		ctx.Resolver.RaiseKeyCostPerHouseNextTurn(p, amount, house, ctx.Source)
		return
	}
	if dur == RemainderOfPlayerTurn || dur == EndOfPlayerNextTurn {
		ctx.Resolver.RaiseKeyCostThisTurn(p, amount, ctx.Source)
	}
	if dur == OpponentNextTurn || dur == EndOfPlayerNextTurn {
		ctx.Resolver.RaiseKeyCostNextTurn(p, amount, ctx.Source)
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
