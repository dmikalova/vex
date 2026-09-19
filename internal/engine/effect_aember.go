package engine

import (
	"errors"
	"fmt"
)

// This file gathers the Æmber-economy effects. Æmber is the resource players
// spend to forge keys (the way to win). It lives in a player's pool, except when
// it is captured or placed on a creature, where it belongs to no one until that
// creature leaves play.

// To gain Æmber, a player moves that many Æmber from the common supply into
// their pool — the ability's controller by default, or their opponent when the
// card says so. A "for each" clause multiplies the amount by a running count.
// EqualTo instead makes the gain equal a running count ("gain Æmber equal to
// <count>", The Flex gaining half a creature's power); set EqualTo or Amount/Per,
// not both.
type GainAember struct {
	Player  Player
	Amount  int
	Per     Count
	EqualTo Count
}

// validate rejects a GainAember whose player was left unset, or one that sets both
// EqualTo and a fixed Amount/Per (two different ways to say how much to gain).
func (e GainAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("GainAember")
	}
	if e.EqualTo != nil && (e.Amount != 0 || e.Per != nil) {
		return errors.New("GainAember: set EqualTo or Amount/Per, not both")
	}
	return nil
}

// Text renders the effect, e.g. "gain 1 Æmber", "your opponent gains 2 Æmber", or
// "gain Æmber equal to half its power, rounded down". A "for each" count leads the
// sentence (rule 9), e.g. "for each key your opponent has forged, gain 1 Æmber".
func (e GainAember) Text() string {
	if e.EqualTo != nil {
		return e.gainVerb() + " Æmber equal to " + equalToText(e.EqualTo, e.Player)
	}
	phrase := fmt.Sprintf("gain %d Æmber", e.Amount)
	switch e.Player {
	case Opponent:
		phrase = fmt.Sprintf("your opponent gains %d Æmber", e.Amount)
	case EachPlayer:
		phrase = fmt.Sprintf("each player gains %d Æmber", e.Amount)
	case ItsOwner:
		phrase = fmt.Sprintf("its owner gains %d Æmber", e.Amount)
	case ItsController:
		phrase = fmt.Sprintf("its controller gains %d Æmber", e.Amount)
	}
	return forEach(e.Per, phrase)
}

// gainVerb renders the subject and verb for an equal-to gain, whose only player
// forms are the controller, the opponent, and each player.
func (e GainAember) gainVerb() string {
	switch e.Player {
	case Opponent:
		return "your opponent gains"
	case EachPlayer:
		return "each player gains"
	default:
		return "gain"
	}
}

// repeatedText renders the gain as a repetition of an identical one already
// stated — "gain 1 more" — for the later rungs of a threshold ladder. A gain whose
// size is not a plain number, or whose subject gainVerb cannot name, has no such
// form and declines with "".
func (e GainAember) repeatedText() string {
	if e.EqualTo != nil || e.Per != nil {
		return ""
	}
	switch e.Player {
	case Controller, Opponent, EachPlayer:
		return fmt.Sprintf("%s %d more", e.gainVerb(), e.Amount)
	default:
		return ""
	}
}

// Resolve adds the Æmber to the selected player's pool. EachPlayer pays both
// players in turn, each from their own point of view, so a Per count reading a
// "this way" tally scales each player's gain by what they themselves lost
// (Hecatomb, Mating Season).
func (e GainAember) Resolve(ctx *EffectContext) {
	if e.Player == EachPlayer {
		for _, p := range [2]int{ctx.Controller, ctx.Opponent()} {
			fromTheirSide := *ctx
			fromTheirSide.Controller = p
			e.gain(&fromTheirSide, p)
		}
		return
	}
	e.gain(ctx, ctx.PlayerFor(e.Player))
}

// amount is the Æmber this gain is worth as seen from ctx: a running count when
// EqualTo is set, otherwise the fixed Amount scaled by an optional Per.
func (e GainAember) amount(ctx *EffectContext) int {
	if e.EqualTo != nil {
		return e.EqualTo.Value(ctx)
	}
	return scaled(e.Amount, e.Per, ctx)
}

// gain hands p the Æmber this effect is worth as seen from ctx, honouring a
// capture that replaces the gain. An equal-to count of zero or less gains nothing.
func (e GainAember) gain(ctx *EffectContext, p int) {
	amount := e.amount(ctx)
	if e.EqualTo != nil && amount <= 0 {
		return
	}
	if capturer, ok := ctx.Resolver.GainAember(p, amount); ok {
		ctx.Resolver.Record(AemberCapturedInsteadOfGain{
			Creature: capturer,
			Player:   p,
			Amount:   amount,
		})
		return
	}
	ctx.Resolver.Record(AemberGained{
		Player: p,
		Amount: amount,
	})
}

// aemberLosers returns the players a lose-Æmber effect drains: both under
// EachPlayer, otherwise the one the relative Player resolves to.
func aemberLosers(ctx *EffectContext, player Player) []int {
	if player == EachPlayer {
		return []int{ctx.Controller, ctx.Opponent()}
	}
	return []int{ctx.PlayerFor(player)}
}

// loseAemberFrom drains amountFor(p) Æmber from each player p, floored at their
// pool so it never goes below zero, tallying each loss (read by a ProducedThisWay
// with TallyAemberLost) and narrating it, and reports whether any Æmber left a
// pool (so a LoseAember can gate a Then). It is the one loss step every LoseAember
// shares, whether it loses a fixed amount, a By share, or a count via EqualTo.
func loseAemberFrom(ctx *EffectContext, players []int, amountFor func(p int) int) bool {
	moved := false
	for _, p := range players {
		lost := min(amountFor(p), ctx.Resolver.Aember(p))
		if lost > 0 {
			moved = true
		}
		ctx.Produced.AemberLost[p] += lost
		ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)-lost)
		ctx.Resolver.Record(AemberLost{
			Player: p,
			Amount: lost,
		})
	}
	return moved
}

// A Loss says how much Æmber to remove from a pool when the amount depends on the
// pool's current size — a Fraction of it (Half/Third, with explicit rounding), or
// all but a fixed remainder. A LoseAember uses one via its By field instead of a
// fixed Amount.
type Loss interface {
	// lose returns how much to remove from a pool of the given size.
	lose(pool int) int
	// object renders the amount as the object of "loses …", using the possessive that
	// fits the loser (e.g. "half of their Æmber, rounded down").
	object(possessive string) string
	// qualifier narrows which players are affected (e.g. "with 6 Æmber or more"), or
	// "" when every player is.
	qualifier() string
}

// AllAember empties a pool entirely, whatever its size — Shatter Storm's "lose
// all your Æmber".
var AllAember Loss = allAember{}

type allAember struct{}

func (allAember) lose(pool int) int { return pool }
func (allAember) object(possessive string) string {
	return "all " + possessive + " Æmber"
}
func (allAember) qualifier() string { return "" }

// AllBut removes everything above keep, leaving a pool of exactly keep; a pool
// already at or below keep is untouched.
func AllBut(keep int) Loss { return allBut{keep: keep} }

type allBut struct{ keep int }

func (a allBut) lose(pool int) int    { return max(0, pool-a.keep) }
func (a allBut) object(string) string { return fmt.Sprintf("all but %d Æmber", a.keep) }
func (a allBut) qualifier() string    { return fmt.Sprintf("with %d Æmber or more", a.keep+1) }

// The economy effects that move Æmber against a pool — StealAember, LoseAember,
// CaptureAember — all express "how much" the same way: a fixed base, optionally
// scaled by a Per count, or a By share of the pool. These three helpers hold that
// shared shape so each effect keeps its own authoring fields but not its own copy
// of the amount logic.

// poolAmount is how much Æmber to move against a pool of the given size: the By
// share when set, otherwise base scaled by an optional Per. It does not cap at the
// pool — a share already fits, and a fixed amount is capped by the caller, which
// records the capped figure.
func poolAmount(base int, by Loss, per Count, ctx *EffectContext, pool int) int {
	if by != nil {
		return by.lose(pool)
	}
	return scaled(base, per, ctx)
}

// aemberObject renders the object of an economy verb: a By share against the given
// possessive ("half of their Æmber"), or a plain "N Æmber" count.
func aemberObject(amount int, by Loss, possessive string) string {
	if by != nil {
		return by.object(possessive)
	}
	return fmt.Sprintf("%d Æmber", amount)
}

// To lose Æmber, a player returns that many Æmber from their pool to the common
// supply. A pool can never go below zero, so a player told to lose more Æmber than
// they have simply loses all of it. Player may be EachPlayer, so both players lose.
// The amount lost is a fixed Amount, a By loss of the pool (By: HalfRoundedDown,
// By: AllBut(5)), or a running count via EqualTo ("lose Æmber equal to <count>",
// Power of Fire) — set exactly one.
type LoseAember struct {
	// Player is whose pool loses; the amount is either a fixed Amount or a By loss of
	// the pool (set one, not both).
	Player Player
	Amount int
	By     Loss
	// Per scales a fixed Amount by a running count, the way GainAember does
	// — Phylyx the Disintegrator drains 1 per other friendly Mars creature.
	Per Count
	// EqualTo makes the loss equal a running count instead of a fixed Amount/By,
	// so the sentence reads "lose Æmber equal to <count>" (Power of Fire drains half
	// the sacrificed creature's power). Set EqualTo or Amount/By/Per, not both.
	EqualTo Count
}

// validate rejects a LoseAember with no player, or one that sets both a fixed
// Amount and a By loss (the two are different ways to say how much to lose), or one
// that sets EqualTo alongside a fixed Amount/By/Per.
func (e LoseAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("LoseAember")
	}
	if e.EqualTo != nil {
		if e.Amount != 0 || e.By != nil || e.Per != nil {
			return errors.New("LoseAember: set EqualTo or Amount/By/Per, not both")
		}
		return nil
	}
	return errAmountOr("LoseAember", "By", e.Amount, e.By != nil)
}

// Text renders the effect, e.g. "lose 1 Æmber", "your opponent loses 4 Æmber",
// "each player with 6 Æmber or more loses all but 5 Æmber", or "each player loses
// Æmber equal to half its power, rounded down".
func (e LoseAember) Text() string {
	if e.EqualTo != nil {
		verb := "lose"
		switch e.Player {
		case Opponent:
			verb = "your opponent loses"
		case EachPlayer:
			verb = "each player loses"
		}
		return verb + " Æmber equal to " + e.EqualTo.CountText()
	}
	var subject, verb, possessive string
	switch e.Player {
	case EachPlayer:
		subject, verb, possessive = "each player", "loses", "their"
	case Opponent:
		subject, verb, possessive = "your opponent", "loses", "their"
	case ThatPlayer:
		subject, verb, possessive = "that player", "loses", "their"
	case ItsOwner:
		subject, verb, possessive = "its controller", "loses", "their"
	default:
		subject, verb, possessive = "", "lose", "your"
	}
	object := aemberObject(e.Amount, e.By, possessive)
	qualifier := ""
	if e.By != nil {
		if q := e.By.qualifier(); q != "" {
			qualifier = " " + q
		}
	}
	if subject == "" {
		return forEach(e.Per, verb+" "+object)
	}
	return forEach(e.Per, subject+qualifier+" "+verb+" "+object)
}

// Resolve removes the Æmber from each affected player's pool, never taking a pool
// below zero.
func (e LoseAember) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate removes the Æmber and reports whether any actually left a pool, so a
// LoseAember can gate a Then — Key Charge only forges if Æmber was lost.
func (e LoseAember) resolveGate(ctx *EffectContext) bool {
	return loseAemberFrom(ctx, aemberLosers(ctx, e.Player), func(p int) int {
		return e.amountFor(ctx, p)
	})
}

// amountFor is how much the given player loses: a running count when EqualTo is set
// (floored at zero), the By loss applied to their pool when set, otherwise the
// fixed Amount scaled by Per.
func (e LoseAember) amountFor(ctx *EffectContext, p int) int {
	if e.EqualTo != nil {
		return max(0, e.EqualTo.Value(ctx))
	}
	return poolAmount(e.Amount, e.By, e.Per, ctx, ctx.Resolver.Aember(p))
}

// Moving Æmber from your pool to a card takes that many Æmber out of your pool
// and sets it on the card, where it stays until something moves it off again.
// Parked there it only matters to a card that can spend the Æmber sitting on it —
// Safe Place, Pocket Universe.
type MoveAemberFromPool struct {
	// Amount is how much Æmber to move from the pool onto the card.
	Amount int
	// Target names the card the Æmber moves onto.
	Target Target
	// Source names the pool the Æmber comes from; the zero value moves it from your
	// own pool. ChosenPlayer lets the controller take it from any player's pool
	// (Monument to Shrix while Citizen Shrix is in your discard pile).
	Source Player
}

// Text renders the effect, e.g. "move 1 Æmber from your pool to Safe Place".
func (e MoveAemberFromPool) Text() string {
	from := "your pool"
	if e.Source == ChosenPlayer {
		from = "any player's pool"
	}
	return fmt.Sprintf("move %d Æmber from %s to %s", e.Amount, from, e.Target.Text())
}

// sourcePool names the player whose pool the Æmber is drawn from, prompting the
// controller when the Source is ChosenPlayer.
func (e MoveAemberFromPool) sourcePool(ctx *EffectContext) int {
	if e.Source != ChosenPlayer {
		return ctx.Controller
	}
	if ctx.Resolver.ChooseOption(ctx.Controller, ctx.Source,
		"Move Æmber from which pool?",
		[]string{"your pool", "your opponent's pool"}) == 1 {
		return 1 - ctx.Controller
	}
	return ctx.Controller
}

// Resolve moves as much of the amount as the pool holds onto each target card.
func (e MoveAemberFromPool) Resolve(ctx *EffectContext) {
	from := e.sourcePool(ctx)
	for _, id := range e.Target.Select(ctx) {
		pool := ctx.Resolver.Aember(from)
		moved := min(e.Amount, pool)
		if moved == 0 {
			return
		}
		ctx.Resolver.SetAember(from, pool-moved)
		ctx.Resolver.AddAmberOn(id, moved)
	}
}

// validate rejects a move with no destination or nothing to move.
func (e MoveAemberFromPool) validate() error {
	if e.Amount <= 0 {
		return errors.New("MoveAemberFromPool needs a positive Amount")
	}
	if e.Target == (Target{}) {
		return errors.New("MoveAemberFromPool needs a Target")
	}
	return nil
}

// PlaceAemberOnThis places Amount Æmber from the common supply on the source card,
// where it accrues until the card is destroyed ([REDACTED] hoards Æmber each time
// you choose its house). Unlike MoveAemberFromPool the Æmber comes from the supply,
// not a player's pool.
type PlaceAemberOnThis struct {
	// Amount is how much Æmber to place on the source card.
	Amount int
}

// Text renders the effect, e.g. "place 1 Æmber from the common supply on {self}".
func (e PlaceAemberOnThis) Text() string {
	return fmt.Sprintf("place %d Æmber from the common supply on %s", e.Amount, SelfName)
}

// Resolve places the Æmber on the source card.
func (e PlaceAemberOnThis) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddAmberOn(ctx.Source, e.Amount)
}

// validate rejects a placement of nothing.
func (e PlaceAemberOnThis) validate() error {
	if e.Amount <= 0 {
		return errors.New("PlaceAemberOnThis needs a positive Amount")
	}
	return nil
}
