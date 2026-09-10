package engine

import (
	"fmt"
	"strings"
)

// Purging a card sets it aside out of the game entirely, in the purge pile, where
// no ability can reach it unless that ability names the purge pile. It is the most
// permanent way a card leaves play: a purged card never enters a discard pile and
// can never be drawn, played, or destroyed again.
// PurgeCard sets cards aside out of the game, taken from a zone the controller
// picks.
// It serves both as a standalone effect (Creeping Oblivion purges up to 2 cards)
// and as the first half of a Then ("purge a creature -> give a +1 power counter"),
// so it reports whether it purged anything.
type PurgeCard struct {
	// Zone is the pile the purge pulls from. It has no default: a Purge must name
	// where it purges from. Only the discard pile is supported today.
	Zone Zone
	// Type restricts the purge to cards of this type; the zero value allows any.
	Type CardType
	// House restricts the purge to cards of this house; HouseNone allows any card.
	House House
	// Amount is how many cards to purge; the zero value counts as one, so a bare
	// Purge reads as "purge a card".
	Amount int
	// UpTo lets the controller purge fewer than Amount, down to none (Creeping
	// Oblivion's "up to 2"). Without it they purge Amount when that many match.
	UpTo bool
}

// validate rejects a Purge that does not name the zone it pulls from.
func (e PurgeCard) validate() error {
	if !e.Zone.valid() {
		return fmt.Errorf("Purge: zone must be set")
	}
	return nil
}

// count is Amount with the zero value treated as one.
func (e PurgeCard) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// noun renders the kind of card purged: the lowercased type when set, else "card",
// house-qualified when House is set (e.g. "Dis card").
func (e PurgeCard) noun() string {
	noun := "card"
	if e.Type != TypeUnset {
		noun = strings.ToLower(e.Type.String())
	}
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return noun
}

// Text renders the effect, e.g. "purge a creature from a discard pile" or "purge
// up to 2 cards from a discard pile".
func (e PurgeCard) Text() string {
	switch {
	case e.UpTo:
		return "purge up to " + countNoun(e.count(), e.noun()) + " from a discard pile"
	case e.count() == 1:
		return "purge " + indefinite(e.noun()) + " from a discard pile"
	default:
		return "purge " + countNoun(e.count(), e.noun()) + " from a discard pile"
	}
}

// Resolve purges the cards, ignoring the report used when Purge gates a Then.
func (e PurgeCard) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate purges up to count matching cards from one discard pile the
// controller picks — the cards one at a time, with a "Done" opt-out when UpTo —
// and reports whether any card was purged.
func (e PurgeCard) resolveGate(ctx *EffectContext) bool {
	matches := func(id LocalID) bool {
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			return false
		}
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			return false
		}
		return true
	}
	// The discard piles holding at least one matching card.
	var piles []int
	for _, p := range []int{ctx.Controller, ctx.Opponent()} {
		if len(discardCardsWhere(ctx, p, matches)) > 0 {
			piles = append(piles, p)
		}
	}
	if len(piles) == 0 {
		return false
	}
	pile := piles[0]
	if len(piles) == 2 {
		pile = piles[ctx.ChooseOption("Choose a discard pile to purge from",
			[]string{"your discard pile", "your opponent's discard pile"})]
	}
	purged := 0
	bonus := 0
	for i := 0; i < e.count(); i++ {
		cands := discardCardsWhere(ctx, pile, matches)
		if len(cands) == 0 {
			break
		}
		var chosen LocalID
		var ok bool
		if e.UpTo {
			chosen, ok = ctx.ChooseCardOptional("Choose a card to purge", cands)
		} else {
			chosen, ok = ctx.ChooseCard("Choose a card to purge", cands)
		}
		if !ok {
			break
		}
		bonus += ctx.Resolver.AemberBonus(chosen)
		purgeFrom(ctx, Discard, pile, chosen)
		purged++
	}
	ctx.Produced.Purged = purged
	ctx.Produced.PurgedAemberBonus = bonus
	return purged > 0
}

// purgeFrom sets one card aside out of the game, dispatching to the resolver
// removal for its current zone — the single move behind every purge (ADR 0031).
func purgeFrom(ctx *EffectContext, from Zone, owner int, id LocalID) {
	switch from {
	case Hand:
		ctx.Resolver.PurgeFromHand(owner, id, ctx.Source)
	case Discard:
		ctx.Resolver.PurgeFromDiscard(owner, id)
	default: // inPlay
		ctx.Resolver.PurgeFromPlay(id)
	}
}

// PurgeFromHand purges cards from a player's hand, with a Selection deciding how
// the cards are picked: chosen by the controller, optionally restricted to a
// house (Imperial Traitor, Greater Oxtet); a random card (Impspector); or each
// matching card (Martians Make Bad Allies, Lesser Oxtet). It reports whether it
// purged anything and records the tally, so it can gate a Then and a following
// effect can scale with it (CardsPurged).
type PurgeFromHand struct {
	// Player whose hand the cards are purged from.
	Player Player
	// Selection decides which cards are purged; it must be set.
	Selection Selection
}

// validate rejects a PurgeFromHand whose player or selection was left unset.
func (e PurgeFromHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("PurgeFromHand")
	}
	if e.Selection == nil {
		return fmt.Errorf("PurgeFromHand: selection must be set")
	}
	return nil
}

// Text renders the effect, e.g. "you may purge a Sanctum card from your
// opponent's hand" or "purge each non-Mars creature from your hand".
func (e PurgeFromHand) Text() string {
	verb := "purge "
	if selectionDeclinable(e.Selection) {
		verb = "you may purge "
	}
	return verb + e.Selection.object() + " from " + whoseHand(e.Player)
}

// Resolve purges the selected cards from the player's hand.
func (e PurgeFromHand) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate purges the selected cards and reports whether any were, recording
// the tally so a Then ("purge a card from your hand -> give two power counters",
// Greater Oxtet) and a following "for each" (CardsPurged) can hang off it. When a
// single card is purged it is put in context (ctx.It) so a following effect can
// name it — Custom Virus destroys each creature sharing a trait with it.
func (e PurgeFromHand) resolveGate(ctx *EffectContext) bool {
	owner := ctx.PlayerFor(e.Player)
	ids := e.Selection.pick(ctx, ctx.Resolver.Hand(owner))
	for _, id := range ids {
		purgeFrom(ctx, Hand, owner, id)
	}
	if len(ids) == 1 {
		ctx.It, ctx.HasIt = ids[0], true
	}
	ctx.Produced.Purged = len(ids)
	return len(ids) > 0
}

// declinable reports whether the purge can be passed — true only for a
// non-mandatory Chosen, whose single card choice is answered by clicking the card
// (or passing) under a May.
func (e PurgeFromHand) declinable() bool { return selectionDeclinable(e.Selection) }

// resolveOptional resolves the purge as its own optional choice under a May.
func (e PurgeFromHand) resolveOptional(ctx *EffectContext) bool { return e.resolveGate(ctx) }

// PurgeCreature purges each creature its Target selects from play into its owner's
// purge pile — the "purge this creature" a card gains (Annihilation Ritual grants
// it to every creature as a Destroyed ability).
type PurgeCreature struct {
	Target Target
}

// validate requires an explicit target.
func (e PurgeCreature) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PurgeCreature")
	}
	return nil
}

// Text renders the effect, e.g. "purge this creature".
func (e PurgeCreature) Text() string { return "purge " + e.Target.Text() }

// Resolve purges each selected creature — from play if it is still there, or from
// its owner's discard pile if it has just been destroyed (Yxilo Bolter purges the
// creature its damage killed). A creature that is in neither zone is left alone.
// The tally is recorded on the context, so a following effect can scale with how
// many were actually purged (see CardsPurged).
func (e PurgeCreature) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate purges and reports whether anything was actually purged, so a Then
// can hang a follow-up off it (Sacrificial Altar only reaches into the discard
// pile if there was a Human to purge).
func (e PurgeCreature) resolveGate(ctx *EffectContext) bool {
	return e.purge(ctx, e.Target.Select(ctx))
}

// declinable reports that the purge is a single clickable creature.
func (e PurgeCreature) declinable() bool { return e.Target.isChosen() }

// vacuous reports that there is nothing here to purge, so a "you may" wrapping it
// need not ask (Buzzle at a flank with no neighbor purges nothing and asks
// nothing).
func (e PurgeCreature) vacuous(ctx *EffectContext) bool { return e.Target.empty(ctx) }

// resolveOptional is resolveGate under a May: the creature is asked declinably, so
// "you may purge a neighboring creature" is answered by clicking that creature
// rather than by a separate Yes/No.
func (e PurgeCreature) resolveOptional(ctx *EffectContext) bool {
	return e.purge(ctx, e.Target.SelectOptional(ctx))
}

// purge carries out the purge of an already-selected set — each creature from play
// if it is still there, or from its owner's discard pile if it has just been
// destroyed (Yxilo Bolter purges the creature its damage killed). It records the
// tally so a following effect can scale with how many were actually purged.
func (e PurgeCreature) purge(ctx *EffectContext, ids []LocalID) bool {
	purged := 0
	for _, id := range ids {
		if resolverInPlay(ctx, id) {
			purgeFrom(ctx, inPlay, 0, id)
			purged++
			continue
		}
		owner := ctx.Resolver.Owner(id)
		for _, d := range ctx.Resolver.Discard(owner) {
			if d == id {
				purgeFrom(ctx, Discard, owner, id)
				purged++
				break
			}
		}
	}
	ctx.Produced.Purged = purged
	return purged > 0
}

// CardsPurged counts the cards the most recent purge in this resolution removed —
// the "for each creature purged this way" tally (One Last Job steals 1 Æmber for
// each creature it purged).
type CardsPurged struct{}

// Value reads the tally the preceding purge recorded.
func (CardsPurged) Value(ctx *EffectContext) int { return ctx.Produced.Purged }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsPurged) CountText() string { return "creature purged this way" }

// PurgedAemberBonus totals the printed Æmber bonus of the cards the most recent
// PurgeCard removed this resolution — Infurnace's opponent loses Æmber equal to the
// total Æmber bonus of the cards it purged. It reads as a whole phrase, not a "for
// each" tally, so a LoseAemberEqualTo names the count directly.
type PurgedAemberBonus struct{}

// Value reads the summed bonus the preceding purge recorded.
func (PurgedAemberBonus) Value(ctx *EffectContext) int { return ctx.Produced.PurgedAemberBonus }

// CountText renders the phrase the loss is measured against.
func (PurgedAemberBonus) CountText() string { return "the total Æmber bonus of the purged cards" }

// PurgeSource purges the card whose ability this is (Library Access purges
// itself). A source still in play — a creature or artifact — is purged from play;
// a resolving action card, not yet in any zone, is instead marked to be set aside
// out of the game when its play completes, rather than going to the discard pile.
type PurgeSource struct{}

// Text renders the effect using the source card's own name.
func (PurgeSource) Text() string { return "purge " + SelfName }

// Resolve purges the source card.
func (PurgeSource) Resolve(ctx *EffectContext) {
	if resolverInPlay(ctx, ctx.Source) {
		ctx.Resolver.PurgeFromPlay(ctx.Source)
		return
	}
	ctx.Resolver.MarkPlayedActionPurged(ctx.Source)
}
