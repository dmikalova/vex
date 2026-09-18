package engine

import (
	"fmt"
)

// Purging a card sets it aside out of the game entirely, in the purge pile, where
// no ability can reach it unless that ability names the purge pile. It is the most
// permanent way a card leaves play: a purged card never enters a discard pile and
// can never be drawn, played, or destroyed again.
// PurgeCard sets cards aside out of the game, taken from a hand or a discard
// pile. Zone names which, and Player whose: Controller or Opponent for one named
// side (Imperial Traitor purges from your opponent's hand), ChosenPlayer for a
// side the controller picks among those holding a matching card (Creeping
// Oblivion's "a discard pile"), EachPlayer for both at once (Soldiers to
// Flowers). A Selection decides which cards leave — a chosen card, restricted by
// house or type; a uniformly random one (Impspector); or every match (Lesser
// Oxtet). Amount purges that many per zone; an Optional Selection lets the
// controller purge fewer, which reads "you may purge" for one card and "up to N"
// for several. GainOwnerAember gives each purged card's owner 1 Æmber.
// It serves both as a standalone effect (Creeping Oblivion purges up to 2 cards)
// and as the first half of a Then ("purge a creature -> give a +1 power counter"),
// so it reports whether it purged anything, and it records the tally and the total
// Æmber bonus of the purged cards so a following effect can scale with them
// (Noname's power, Infurnace's Æmber loss).
type PurgeCard struct {
	// Zone names the source the cards are purged from: Hand or Discard.
	Zone Zone
	// Player names whose copy of the zone the purge pulls from: Controller or
	// Opponent for a fixed side, ChosenPlayer for one the controller picks,
	// EachPlayer for both at once.
	Player Player
	// Selection decides which cards are purged; it must be set.
	Selection Selection
	// Amount is how many cards to purge from each zone; the zero value counts as
	// one. It is ignored by an Each Selection, which takes every matching card. An
	// Optional Selection lets the controller purge fewer, down to none, reading as
	// "up to N" (Creeping Oblivion's "up to 2").
	Amount int
	// GainOwnerAember gives each purged card's owner 1 Æmber per purged card
	// (Soldiers to Flowers).
	GainOwnerAember bool
}

// validate rejects a Purge whose player or selection was left unset, or whose
// source zone is one no purge draws from.
func (e PurgeCard) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("PurgeCard")
	}
	if e.Selection == nil {
		return fmt.Errorf("PurgeCard: selection must be set")
	}
	if e.Zone != Hand && e.Zone != Discard {
		return fmt.Errorf("PurgeCard: zone must be Hand or Discard")
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

// object renders the noun phrase purged, e.g. "a creature", "up to 2 cards", or
// "each Untamed creature".
func (e PurgeCard) object() string {
	switch {
	case selectionDeclinable(e.Selection) && e.count() > 1:
		return "up to " + countNoun(e.count(), e.Selection.noun())
	case e.count() > 1:
		return countNoun(e.count(), e.Selection.noun())
	default:
		return e.Selection.object()
	}
}

// declinable reports that a single-card purge with a declinable selection is one
// clickable card, so the controller answers it by clicking that card or passing. A
// multi-card purge renders "up to N" and runs its own cycle instead.
func (e PurgeCard) declinable() bool {
	return selectionDeclinable(e.Selection) && e.count() == 1
}

// Text renders the effect, e.g. "purge a creature from a discard pile", "purge up
// to 2 cards from a discard pile", "you may purge a Sanctum card from your
// opponent's hand", or "purge each Untamed creature from each player's discard
// pile. For each card purged this way, its owner gains 1 Æmber".
func (e PurgeCard) Text() string {
	verb := "purge "
	if e.declinable() {
		verb = "you may purge "
	}
	text := verb + e.object() + " from " + whoseZone(e.Player, e.Zone)
	if e.GainOwnerAember {
		text += ". For each card purged this way, its owner gains 1 Æmber"
	}
	return text
}

// Resolve purges the cards, ignoring the report used when Purge gates a Then.
func (e PurgeCard) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveOptional is resolveGate under a May: the card is asked declinably, with a
// Done to decline.
func (e PurgeCard) resolveOptional(ctx *EffectContext) bool { return e.resolveGate(ctx) }

// cards returns the cards in one player's copy of the source zone.
func (e PurgeCard) cards(ctx *EffectContext, player int) []LocalID {
	if e.Zone == Hand {
		return ctx.Resolver.Hand(player)
	}
	return ctx.Resolver.Discard(player)
}

// sides returns whose copies of the zone the purge acts on: both for EachPlayer,
// the named side for Controller or Opponent, or — for ChosenPlayer — the one the
// controller picks among those holding a card the Selection would take (prompting
// only when both do).
func (e PurgeCard) sides(ctx *EffectContext) []int {
	switch e.Player {
	case EachPlayer:
		return []int{ctx.Controller, ctx.Opponent()}
	case ChosenPlayer:
	default:
		return []int{ctx.PlayerFor(e.Player)}
	}
	var eligible []int
	for _, p := range []int{ctx.Controller, ctx.Opponent()} {
		if len(e.Selection.candidates(ctx, e.cards(ctx, p))) > 0 {
			eligible = append(eligible, p)
		}
	}
	if len(eligible) < 2 {
		return eligible
	}
	chosen := ctx.ChooseOption("Choose a "+e.Zone.noun()+" to purge from",
		[]string{"your " + e.Zone.noun(), "your opponent's " + e.Zone.noun()})
	return eligible[chosen : chosen+1]
}

// resolveGate purges up to count matching cards from each named zone — the cards
// one at a time for a Chosen Selection, with a "Done" opt-out for an Optional
// Selection, and every match at once for an Each Selection — records the tally and
// the total Æmber bonus, pays each owner when GainOwnerAember, and reports whether
// any card was purged (so a Then can gate on it). A purge that took exactly one
// card puts it in context (ctx.It) so a following effect can name it — Custom Virus
// destroys each creature sharing a trait with it.
func (e PurgeCard) resolveGate(ctx *EffectContext) bool {
	purged := 0
	bonus := 0
	var last LocalID
	for _, side := range e.sides(ctx) {
		for i := 0; i < e.count(); i++ {
			ids := e.Selection.pick(ctx, e.cards(ctx, side))
			if len(ids) == 0 {
				break
			}
			for _, id := range ids {
				bonus += ctx.Resolver.AemberBonus(id)
				purgeFrom(ctx, e.Zone, side, id)
				purged++
				last = id
				ctx.Produced.Purged[side]++
				if e.GainOwnerAember {
					ctx.Resolver.GainAember(side, 1)
				}
			}
		}
	}
	if purged == 1 {
		ctx.It, ctx.HasIt = last, true
	}
	ctx.Produced.PurgedAemberBonus = bonus
	return purged > 0
}

// purgeFrom sets one card aside out of the game from its current zone — the purge
// verb is the movement matrix with its destination fixed to out-of-the-game, so it
// moves through the shared toPurged destination rather than its own switch (ADR
// 0031).
func purgeFrom(ctx *EffectContext, from Zone, owner int, id LocalID) {
	toPurged.moveFrom(ctx, from, owner, id)
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
func (e PurgeCreature) vacuous(
	ctx *EffectContext,
) bool {
	return e.Target.empty(ctx)
}

// resolveOptional is resolveGate under a May: the creature is asked declinably, so
// "you may purge a neighboring creature" is answered by clicking that creature
// rather than by a separate Yes/No.
func (e PurgeCreature) resolveOptional(ctx *EffectContext) bool {
	return e.purge(ctx, e.Target.SelectOptional(ctx))
}

// purge carries out the purge of an already-selected set — each creature from play
// if it is still there, or from its owner's discard pile if it has just been
// destroyed (Yxilo Bolter purges the creature its damage killed). It records the
// tally so a following effect can scale with how many were actually purged. When a
// single card is purged it is put in context (ctx.It) so a following effect can
// name it — Reclaimed by Nature resolves the bonus icons on the artifact it purged.
func (e PurgeCreature) purge(ctx *EffectContext, ids []LocalID) bool {
	purged := 0
	for _, id := range ids {
		zone, side, ok := currentZone(ctx, id)
		if !ok {
			continue
		}
		ctx.Produced.Purged[side]++
		if len(ids) == 1 {
			ctx.ItController = side
		}
		purgeFrom(ctx, zone, side, id)
		purged++
	}
	if len(ids) == 1 {
		ctx.It, ctx.HasIt = ids[0], true
	}
	return purged > 0
}

// CardsPurged counts the cards the most recent purge in this resolution removed,
// both players' shares together — the "for each card purged this way" tally. Type
// names the noun the clause repeats: unset it reads "card", Creature it reads
// "creature" (One Last Job steals 1 Æmber for each creature it purged).
type CardsPurged struct {
	Type CardType
}

// Value reads the whole tally the preceding purge recorded, both sides together.
func (CardsPurged) Value(ctx *EffectContext) int {
	return ctx.Produced.Purged[0] + ctx.Produced.Purged[1]
}

// CountText renders the singular noun the "for each" clause repeats.
func (c CardsPurged) CountText() string {
	noun := "card"
	if c.Type == Creature {
		noun = "creature"
	}
	return noun + " purged this way"
}

// PurgedAemberBonus totals the printed Æmber bonus of the cards the most recent
// PurgeCard removed this resolution — Infurnace's opponent loses Æmber equal to the
// total Æmber bonus of the cards it purged. It reads as a whole phrase, not a "for
// each" tally, so a LoseAember with EqualTo names the count directly.
type PurgedAemberBonus struct{}

// Value reads the summed bonus the preceding purge recorded.
func (PurgedAemberBonus) Value(
	ctx *EffectContext,
) int {
	return ctx.Produced.PurgedAemberBonus
}

// CountText renders the phrase the loss is measured against.
func (PurgedAemberBonus) CountText() string { return "the total Æmber bonus of the purged cards" }

// PurgeSource purges the card whose ability this is (Library Access purges
// itself), wherever that card is — in play, or mid-play and in no zone at all.
type PurgeSource struct{}

// Text renders the effect using the source card's own name.
func (PurgeSource) Text() string { return "purge " + SelfName }

// Resolve purges the source card.
func (PurgeSource) Resolve(ctx *EffectContext) { toPurged.moveSource(ctx) }
