package engine

import (
	"fmt"
	"strings"
)

// Purging a card sets it aside out of the game entirely, in the purge pile, where
// no ability can reach it unless that ability names the purge pile. It is the most
// permanent way a card leaves play: a purged card never enters a discard pile and
// can never be drawn, played, or destroyed again.
//
//rulebook:effect Purge

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
	// Count is how many cards to purge; the zero value counts as one, so a bare
	// Purge reads as "purge a card".
	Count int
	// UpTo lets the controller purge fewer than Count, down to none (Creeping
	// Oblivion's "up to 2"). Without it they purge Count when that many match.
	UpTo bool
}

// validate rejects a Purge that does not name the zone it pulls from.
func (e PurgeCard) validate() error {
	if !e.Zone.valid() {
		return fmt.Errorf("Purge: zone must be set")
	}
	return nil
}

// count is Count with the zero value treated as one.
func (e PurgeCard) count() int {
	if e.Count < 1 {
		return 1
	}
	return e.Count
}

// noun renders the kind of card purged: the lowercased type when set, else "card".
func (e PurgeCard) noun() string {
	if e.Type != TypeUnset {
		return strings.ToLower(e.Type.String())
	}
	return "card"
}

// Text renders the effect, e.g. "purge a creature from a discard pile" or "purge
// up to 2 cards from a discard pile".
func (e PurgeCard) Text() string {
	switch {
	case e.UpTo:
		return fmt.Sprintf("purge up to %d %ss from a discard pile", e.count(), e.noun())
	case e.count() == 1:
		return "purge a " + e.noun() + " from a discard pile"
	default:
		return fmt.Sprintf("purge %d %ss from a discard pile", e.count(), e.noun())
	}
}

// Resolve purges the cards, ignoring the report used when Purge gates a Then.
func (e PurgeCard) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate purges up to count matching cards from one discard pile the
// controller picks — the cards one at a time, with a "Done" opt-out when UpTo —
// and reports whether any card was purged.
func (e PurgeCard) resolveGate(ctx *EffectContext) bool {
	matches := func(id LocalID) bool {
		return e.Type == TypeUnset || ctx.Resolver.TypeOf(id) == e.Type
	}
	// The discard piles holding at least one matching card.
	var piles []int
	for _, p := range []int{ctx.Controller, ctx.Opponent()} {
		for _, id := range ctx.Resolver.Discard(p) {
			if matches(id) {
				piles = append(piles, p)
				break
			}
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
	for i := 0; i < e.count(); i++ {
		var cands []LocalID
		for _, id := range ctx.Resolver.Discard(pile) {
			if matches(id) {
				cands = append(cands, id)
			}
		}
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
		ctx.Resolver.PurgeFromDiscard(pile, chosen)
		purged++
	}
	return purged > 0
}

// PurgeFromHand lets the controller choose and purge one card from a player's
// hand, optionally restricted to a house — the "you may choose and purge a
// Sanctum card" of Imperial Traitor. It is a "may": the controller can always
// decline, so nothing is purged when they pass or no card matches.
type PurgeFromHand struct {
	// Player whose hand the card is purged from.
	Player Player
	// House restricts the choice to cards of this house; HouseNone allows any card.
	House House
}

// validate rejects a PurgeFromHand whose player was left unset.
func (e PurgeFromHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("PurgeFromHand")
	}
	return nil
}

// noun renders the kind of card purged: the house-qualified card when House is
// set (e.g. "Sanctum card"), else the generic "card".
func (e PurgeFromHand) noun() string {
	if e.House != HouseNone {
		return e.House.String() + " card"
	}
	return "card"
}

// Text renders the effect, e.g. "you may purge a Sanctum card from your
// opponent's hand".
func (e PurgeFromHand) Text() string {
	whose := "your hand"
	if e.Player == Opponent {
		whose = "your opponent's hand"
	}
	return "you may purge a " + e.noun() + " from " + whose
}

// Resolve offers the matching cards in the player's hand as a declinable choice,
// then purges the chosen one. It does nothing when no card matches or the
// controller declines.
func (e PurgeFromHand) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	var cands []LocalID
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.House == HouseNone || ctx.Resolver.House(id) == e.House {
			cands = append(cands, id)
		}
	}
	chosen, ok := ctx.ChooseCardOptional("Choose a card to purge", cands)
	if !ok {
		return
	}
	ctx.Resolver.PurgeFromHand(owner, chosen)
}

// PurgeEachFromHand purges every card in a player's hand that matches its filters
// — Martians Make Bad Allies purges each non-Mars creature out of the hand it
// just revealed. Unlike PurgeFromHand there is no choice: the filters decide, and
// the tally is recorded so a following effect can scale with it (CardsPurged).
type PurgeEachFromHand struct {
	// Player whose hand the cards are purged from.
	Player Player
	// Type restricts the purge to cards of this type; the zero value allows any.
	Type CardType
	// ExceptHouse spares the cards of that house; HouseNone spares nothing.
	ExceptHouse House
}

// validate rejects a PurgeEachFromHand whose player was left unset.
func (e PurgeEachFromHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("PurgeEachFromHand")
	}
	return nil
}

// noun renders the kind of card purged, e.g. "non-Mars creature".
func (e PurgeEachFromHand) noun() string {
	noun := "card"
	if e.Type != TypeUnset {
		noun = strings.ToLower(e.Type.String())
	}
	if e.ExceptHouse != HouseNone {
		noun = "non-" + e.ExceptHouse.String() + " " + noun
	}
	return noun
}

// Text renders the effect, e.g. "purge each non-Mars creature from your hand".
func (e PurgeEachFromHand) Text() string {
	whose := "your hand"
	if e.Player == Opponent {
		whose = "your opponent's hand"
	}
	return "purge each " + e.noun() + " from " + whose
}

// Resolve purges every matching card from the hand and records the tally.
func (e PurgeEachFromHand) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	var doomed []LocalID
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			continue
		}
		if e.ExceptHouse != HouseNone && ctx.Resolver.House(id) == e.ExceptHouse {
			continue
		}
		doomed = append(doomed, id)
	}
	for _, id := range doomed {
		ctx.Resolver.PurgeFromHand(owner, id)
	}
	ctx.Produced.Purged = len(doomed)
}

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
func (e PurgeCreature) Resolve(ctx *EffectContext) {
	purged := 0
	for _, id := range e.Target.Select(ctx) {
		if resolverInPlay(ctx, id) {
			ctx.Resolver.PurgeFromPlay(id)
			purged++
			continue
		}
		owner := ctx.Resolver.Owner(id)
		for _, d := range ctx.Resolver.Discard(owner) {
			if d == id {
				ctx.Resolver.PurgeFromDiscard(owner, id)
				purged++
				break
			}
		}
	}
	ctx.Produced.Purged = purged
}

// CardsPurged counts the cards the most recent purge in this resolution removed —
// the "for each creature purged this way" tally (One Last Job steals 1 Æmber for
// each creature it purged).
type CardsPurged struct{}

// Value reads the tally the preceding purge recorded.
func (CardsPurged) Value(ctx *EffectContext) int { return ctx.Produced.Purged }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsPurged) CountText() string { return "creature purged this way" }

// PurgeCreatureFromHand purges a creature the controller chooses from their hand
// and puts it in context (ctx.It) for a following effect to act on. It is the
// first half of Custom Virus's "purge a creature from your hand, destroy each
// creature that shares a trait with the purged creature": the purge here, the
// destruction a Destroy whose Target is EachCreature.SharingTrait. The purged
// creature is out of play but its traits still decide what a later effect selects;
// with no creature in hand, or the choice declined, nothing is purged and It stays
// unset (so the following effect finds nothing).
type PurgeCreatureFromHand struct{}

// Text renders the effect.
func (PurgeCreatureFromHand) Text() string {
	return "purge a creature from your hand"
}

// Resolve purges a chosen creature from the controller's hand and sets it as the
// context card.
func (PurgeCreatureFromHand) Resolve(ctx *EffectContext) {
	var inHand []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if ctx.Resolver.IsCreature(id) {
			inHand = append(inHand, id)
		}
	}
	if len(inHand) == 0 {
		return
	}
	purged, ok := ctx.ChooseCreature("Choose a creature to purge from your hand", inHand)
	if !ok {
		return
	}
	ctx.Resolver.PurgeFromHand(ctx.Controller, purged)
	ctx.It, ctx.HasIt = purged, true
}
