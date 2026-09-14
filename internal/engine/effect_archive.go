package engine

import (
	"fmt"
	"strings"
)

// ArchiveCard sets cards aside into the controller's own archives: they go
// face-down, out of the opponent's reach, and the controller may take them into
// hand after choosing a house on a later turn. A Selection decides how each card
// is picked — the controller chooses one (Chosen, restrictable to type/house), a
// uniformly random card leaves a hidden hand (Random, Eureka!), or a named card is
// pinned (Named, Hyde archiving Velum). Zone names the source, Hand or Discard.
// Amount archives that many; the zero value archives one. An Optional Selection
// (Chosen with Optional) lets the controller archive fewer, down to none, which
// reads as "up to N" (Mobius Scroll). Revealed shows each card before
// archiving it, so a filtered choice is verified (Incubation Chamber). Per scales
// the count by a running tally (Dr. Milli); Or switches Amount to an alternate
// when a condition holds (Velum). It reports whether any card was archived, so it
// can gate a Then (Zyzzix the Many).
type ArchiveCard struct {
	// Zone names the source the cards are archived from: Hand or Discard.
	Zone Zone
	// From names whose zone the cards are drawn from. The zero value (playerUnset)
	// is the controller's own zone; Opponent draws from the opponent's hand into
	// the controller's archives (Hidden Stash), an abduction from hand.
	From Player
	// Selection decides how each card is picked; it must be set.
	Selection Selection
	// Amount is how many cards to archive; the zero value counts as one.
	Amount int
	// Revealed shows each chosen card to the opponent before archiving it, which is
	// how a filtered choice is verified (Incubation Chamber).
	Revealed bool
	// Per scales how many cards are archived by a running count (Dr. Milli). The
	// zero value archives exactly the count.
	Per Count
	// Or switches Amount to an alternate when a condition holds, so the card reads
	// "archive a card, or 2 cards if …" (Velum archives 2 while controlling Hyde).
	Or OrAmount
}

// validate rejects an ArchiveCard whose selection was left unset, a source zone
// other than the hand, discard pile, or deck, a positional selection paired with
// an unordered zone, a non-positional selection on the deck, or an invalid Or
// guard.
func (e ArchiveCard) validate() error {
	if e.Selection == nil {
		return fmt.Errorf("ArchiveCard: selection must be set")
	}
	if e.Zone != Hand && e.Zone != Discard && e.Zone != Deck {
		return fmt.Errorf("ArchiveCard: zone must be Hand, Discard, or Deck")
	}
	if selectionPositional(e.Selection) && !e.Zone.ordered() {
		return fmt.Errorf("ArchiveCard: positional selection needs an ordered zone")
	}
	if e.Zone == Deck && !selectionPositional(e.Selection) {
		return fmt.Errorf("ArchiveCard: a deck archive must be positional")
	}
	if e.From == Opponent && e.Zone != Hand {
		return fmt.Errorf("ArchiveCard: an opponent archive must be from the hand")
	}
	return e.Or.validate()
}

// count is Amount with the zero value treated as one.
func (e ArchiveCard) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// object renders the count-bearing noun phrase the archive acts on, e.g. "a card",
// "2 cards", "up to 2 cards", "the top card", or the pinned "Velum".
func (e ArchiveCard) object() string {
	switch {
	case selectionPositional(e.Selection):
		return positionalObject(e.Selection, e.count())
	case selectionDeclinable(e.Selection) && e.count() > 1:
		return "up to " + countNoun(e.count(), e.Selection.noun())
	case e.count() > 1:
		return countNoun(e.count(), e.Selection.noun())
	default:
		return e.Selection.object()
	}
}

// Text renders the effect, naming the source zone explicitly (rule 17). A
// positional archive reads "of your <zone>" (the top card of the deck), the rest
// "from your <zone>". A Per count leads the sentence; a revealed archive reads as
// the reveal it is; an Or appends its alternate-amount tail.
func (e ArchiveCard) Text() string {
	prep := " from your "
	if selectionPositional(e.Selection) {
		prep = " of your "
	}
	if e.From == Opponent {
		prep = strings.TrimSuffix(prep, "your ") + "your opponent's "
	}
	from := prep + e.Zone.noun()
	if e.Per != nil {
		return forEach(e.Per, "archive "+e.Selection.object()+from)
	}
	if e.Revealed {
		return "reveal " + e.object() + from + " and archive it"
	}
	body := "archive " + e.object() + from
	if e.Or.set() {
		body += e.Or.tail(archiveOrObject(e.Or.Amount, e.Selection.noun()))
	}
	return body
}

// archiveOrObject renders the alternate-amount noun for an Or tail, e.g. "a card"
// or "2 cards".
func archiveOrObject(n int, noun string) string {
	if n == 1 {
		return indefinite(noun)
	}
	return countNoun(n, noun)
}

// Resolve archives the selected cards, stopping early if the source runs out or a
// choice is declined.
func (e ArchiveCard) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the archives and reports whether any card was archived, so
// ArchiveCard can gate a Then. Each archived card is revealed first when Revealed
// is set, then moved into archives from its source zone.
func (e ArchiveCard) resolveGate(ctx *EffectContext) bool {
	archived := false
	amount := scaled(e.Or.pick(e.count(), ctx), e.Per, ctx)
	for i := 0; i < amount; i++ {
		ids := e.Selection.pick(ctx, e.source(ctx))
		if len(ids) == 0 {
			return archived
		}
		e.archive(ctx, ids)
		archived = true
	}
	return archived
}

// declinable reports that a single-card archive with a declinable selection is one
// clickable card, so a "you may" wrapping it (Zyzzix the Many) can be answered by
// clicking that card instead of a separate Yes/No. A multi-card ("up to N") archive
// keeps its own cycle instead.
func (e ArchiveCard) declinable() bool {
	return selectionDeclinable(e.Selection) && e.count() == 1
}

// resolveOptional is resolveGate under a May: the card is asked declinably, with a
// Done to decline.
func (e ArchiveCard) resolveOptional(ctx *EffectContext) bool {
	ids := e.Selection.pick(ctx, e.source(ctx))
	if len(ids) == 0 {
		return false
	}
	e.archive(ctx, ids)
	return true
}

// source returns the cards in the zone the archive draws from, ordered top-first
// for a positional selection so it can read the top at index 0.
func (e ArchiveCard) source(ctx *EffectContext) []LocalID {
	var cards []LocalID
	switch e.Zone {
	case Discard:
		cards = ctx.Resolver.Discard(ctx.Controller)
	case Deck:
		cards = ctx.Resolver.Deck(ctx.Controller)
	default:
		cards = ctx.Resolver.Hand(e.sourcePlayer(ctx))
	}
	if selectionPositional(e.Selection) {
		cards = topFirst(e.Zone, cards)
	}
	return cards
}

// sourcePlayer is the player whose zone the archive draws from: the opponent when
// From is Opponent, otherwise the controller.
func (e ArchiveCard) sourcePlayer(ctx *EffectContext) int {
	if e.From == Opponent {
		return ctx.Opponent()
	}
	return ctx.Controller
}

// archive reveals the picked cards when Revealed is set, then moves each into
// archives from the source zone.
func (e ArchiveCard) archive(ctx *EffectContext, ids []LocalID) {
	for _, id := range ids {
		if e.Revealed {
			ctx.Resolver.Record(CardsRevealedToAll{
				Player: ctx.Controller,
				Cards:  []LocalID{id},
			})
		}
		if e.From == Opponent {
			// Abduction from hand: the enemy card goes into the controller's own
			// archives, not its owner's, and returns home when it leaves them.
			ctx.Resolver.ArchiveEnemyFromHand(ctx.Controller, id)
			continue
		}
		archiveFrom(ctx, e.Zone, ctx.Controller, id)
	}
}

// archiveFrom moves one card into its owner's archives — the archive verb is the
// movement matrix with its destination fixed to archives, so it moves through the
// shared ToArchives destination rather than its own switch (ADR 0031).
func archiveFrom(ctx *EffectContext, from Zone, owner int, id LocalID) {
	ToArchives.moveFrom(ctx, from, owner, id)
}

// ArchiveFromPlay moves each in-play card its Target selects into its owner's
// archives, shedding damage, armor, upgrades, and other in-play state.
type ArchiveFromPlay struct {
	Target Target
}

// validate requires an explicit target.
func (e ArchiveFromPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ArchiveFromPlay")
	}
	return nil
}

// Text renders the effect, e.g. "archive each friendly Knight creature from
// play".
func (e ArchiveFromPlay) Text() string {
	return fmt.Sprintf("archive %s from play", e.Target.Text())
}

// Resolve archives each selected in-play card.
func (e ArchiveFromPlay) Resolve(ctx *EffectContext) { e.archive(ctx, e.Target.Select(ctx)) }

// declinable reports that the archiving is a single clickable card, so a "you
// may" wrapping it (Vezyma Thinkdrone) can be answered by clicking that card
// instead of a separate Yes/No.
func (e ArchiveFromPlay) declinable() bool { return e.Target.isChosen() }

// resolveOptional is Resolve under a May: the card is asked declinably, with a
// Done to decline.
func (e ArchiveFromPlay) resolveOptional(ctx *EffectContext) bool {
	return e.archive(ctx, e.Target.SelectOptional(ctx))
}

// archive puts an already-selected set of in-play cards into their owners'
// archives, and reports whether anything was. The cards leave play at the same
// time, so archiving one that was buffing another does not destroy that other.
func (e ArchiveFromPlay) archive(ctx *EffectContext, ids []LocalID) bool {
	ctx.Produced.Archived = ids
	ctx.Resolver.PutIntoArchivesEach(ctx.Resolver.ActivePlayer(), ids)
	return len(ids) > 0
}

// ArchiveSource archives the card whose ability this is (Sucker Punch archives
// itself). A source still in play — a creature or artifact — is archived from
// play; a resolving action card, not yet in any zone, is instead marked to go to
// its owner's archives when its play completes, rather than to the discard pile.
type ArchiveSource struct{}

// Text renders the effect using the source card's own name.
func (ArchiveSource) Text() string { return "archive " + SelfName }

// Resolve archives the source card.
func (ArchiveSource) Resolve(ctx *EffectContext) {
	if resolverInPlay(ctx, ctx.Source) {
		ctx.Resolver.PutIntoArchives(ctx.Source)
		return
	}
	ctx.Resolver.MarkPlayedActionArchived(ctx.Source)
}

// ArchiveGrantingUpgrade archives the upgrade whose granted ability this is
// (ctx.Upgrade) — Ghostform grants its host "Fight/Reap: Archive Ghostform",
// which sends Ghostform itself off its host to its owner's archives. The ability
// is printed on the upgrade, so its text names the upgrade through the {card}
// placeholder rather than through {self}, which a granted ability renders as the
// host creature.
type ArchiveGrantingUpgrade struct{}

// Text renders the effect, e.g. "archive Ghostform".
func (ArchiveGrantingUpgrade) Text() string { return "archive " + CardName }

// Resolve archives the granting upgrade off its host.
func (ArchiveGrantingUpgrade) Resolve(ctx *EffectContext) {
	ctx.Resolver.ArchiveUpgrade(ctx.Upgrade)
}

// DiscardArchives moves all of a player's archived cards into their discard pile.
type DiscardArchives struct {
	Player Player
}

// validate rejects a DiscardArchives whose player was left unset.
func (e DiscardArchives) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardArchives")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent discards each of their archived
// cards".
func (e DiscardArchives) Text() string {
	if e.Player == Opponent {
		return "your opponent discards each of their archived cards"
	}
	return "discard each of your archived cards"
}

// Resolve discards the chosen player's archives. The active player performs the
// discard, so discardArchives randomizes the order for an opponent's archives
// (which they cannot see) and lets them order their own.
func (e DiscardArchives) Resolve(ctx *EffectContext) {
	ctx.Resolver.DiscardArchives(ctx.PlayerFor(e.Player))
}
