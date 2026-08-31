package engine

import "fmt"

// This file gathers the Æmber-economy effects. Æmber is the resource players
// spend to forge keys (the way to win). It lives in a player's pool, except when
// it is captured or placed on a creature, where it belongs to no one until that
// creature leaves play.

// To gain Aember, a player moves that many Aember from the common supply into
// their pool — the ability's controller by default, or their opponent when the
// card says so. A "for each" clause multiplies the amount by a running count.
//
//rulebook:effect Gain Aember
type GainAember struct {
	Player Player
	Amount int
	Per    Count
}

// validate rejects a GainAember whose player was left unset.
func (e GainAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("GainAember")
	}
	return nil
}

// Text renders the effect, e.g. "gain 1 Æmber" or "your opponent gains 2 Æmber".
// A "for each" count leads the sentence (rule 9), e.g. "for each key your opponent
// has forged, gain 1 Æmber".
func (e GainAember) Text() string {
	phrase := fmt.Sprintf("gain %d Æmber", e.Amount)
	if e.Player == Opponent {
		phrase = fmt.Sprintf("your opponent gains %d Æmber", e.Amount)
	}
	return forEach(e.Per, phrase)
}

// Resolve adds the Æmber to the selected player's pool.
func (e GainAember) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	amount := e.Amount
	if e.Per != nil {
		amount *= e.Per.Value(ctx)
	}
	if capturer, ok := ctx.Resolver.GainAember(p, amount); ok {
		ctx.Resolver.Logf(
			"%s captures %d Æmber instead of %s gaining it",
			ctx.Resolver.Name(capturer),
			amount,
			ctx.Resolver.PlayerName(p),
		)
		return
	}
	ctx.Resolver.Logf("%s gains %d Æmber", ctx.Resolver.PlayerName(p), amount)
}

// A Loss says how much Æmber to remove from a pool when the amount depends on the
// pool's current size — half of it, or all but a fixed remainder. A LoseAember uses
// one via its By field instead of a fixed Amount.
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

// Half removes half a pool, rounded down.
var Half Loss = half{}

type half struct{}

func (half) lose(pool int) int { return pool / 2 }
func (half) object(possessive string) string {
	return "half of " + possessive + " Æmber, rounded down"
}
func (half) qualifier() string { return "" }

// AllBut removes everything above keep, leaving a pool of exactly keep; a pool
// already at or below keep is untouched.
func AllBut(keep int) Loss { return allBut{keep: keep} }

type allBut struct{ keep int }

func (a allBut) lose(pool int) int    { return max(0, pool-a.keep) }
func (a allBut) object(string) string { return fmt.Sprintf("all but %d Æmber", a.keep) }
func (a allBut) qualifier() string    { return fmt.Sprintf("with %d Æmber or more", a.keep+1) }

// To lose Aember, a player returns that many Aember from their pool to the common
// supply. A pool can never go below zero, so a player told to lose more Aember than
// they have simply loses all of it. Player may be EachPlayer, so both players lose.
// The amount lost is either a fixed Amount or a By loss of the pool (By: Half,
// By: AllBut(5)) — set one, not both.
//
//rulebook:effect Lose Aember
type LoseAember struct {
	Player Player
	Amount int
	By     Loss
}

// validate rejects a LoseAember with no player, or one that sets both a fixed
// Amount and a By loss (the two are different ways to say how much to lose).
func (e LoseAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("LoseAember")
	}
	if e.Amount != 0 && e.By != nil {
		return fmt.Errorf("LoseAember: set Amount or By, not both (got Amount=%d)", e.Amount)
	}
	return nil
}

// Text renders the effect, e.g. "lose 1 Æmber", "your opponent loses 4 Æmber", or
// "each player with 6 Æmber or more loses all but 5 Æmber".
func (e LoseAember) Text() string {
	var subject, verb, possessive string
	switch e.Player {
	case EachPlayer:
		subject, verb, possessive = "each player", "loses", "their"
	case Opponent:
		subject, verb, possessive = "your opponent", "loses", "their"
	default:
		subject, verb, possessive = "", "lose", "your"
	}
	object := fmt.Sprintf("%d Æmber", e.Amount)
	qualifier := ""
	if e.By != nil {
		object = e.By.object(possessive)
		if q := e.By.qualifier(); q != "" {
			qualifier = " " + q
		}
	}
	if subject == "" {
		return verb + " " + object
	}
	return subject + qualifier + " " + verb + " " + object
}

// Resolve removes the Æmber from each affected player's pool, never taking a pool
// below zero.
func (e LoseAember) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate removes the Aember and reports whether any actually left a pool, so a
// LoseAember can gate a Then — Key Charge only forges if Æmber was lost.
func (e LoseAember) resolveGate(ctx *EffectContext) bool {
	moved := false
	for _, p := range e.losers(ctx) {
		lost := min(e.amountFor(ctx, p), ctx.Resolver.Aember(p))
		if lost > 0 {
			moved = true
		}
		ctx.Resolver.SetAember(p, ctx.Resolver.Aember(p)-lost)
		ctx.Resolver.Logf("%s loses %d Æmber", ctx.Resolver.PlayerName(p), lost)
	}
	return moved
}

// losers returns the players who lose Æmber — both for EachPlayer, otherwise the
// one the relative Player resolves to.
func (e LoseAember) losers(ctx *EffectContext) []int {
	if e.Player == EachPlayer {
		return []int{ctx.Controller, ctx.Opponent()}
	}
	return []int{ctx.PlayerFor(e.Player)}
}

// amountFor is how much the given player loses: the By loss applied to their pool
// when set, otherwise the fixed Amount.
func (e LoseAember) amountFor(ctx *EffectContext, p int) int {
	if e.By != nil {
		return e.By.lose(ctx.Resolver.Aember(p))
	}
	return e.Amount
}
