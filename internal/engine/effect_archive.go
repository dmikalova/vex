package engine

import "fmt"

// Archiving moves cards into your archives: they are set aside face-down, out of
// the opponent's reach, and you may take them into your hand after picking a
// house on a later turn. Archiving from hand lets you choose which cards to set
// aside.
//
//rulebook:effect Archive
type ArchiveFromHand struct {
	Count int
}

// Text renders the effect, e.g. "archive a card from your hand" or "archive 2
// cards from your hand". The source zone is named explicitly.
func (e ArchiveFromHand) Text() string {
	if e.Count == 1 {
		return "archive a card from your hand"
	}
	return fmt.Sprintf("archive %d cards from your hand", e.Count)
}

// Resolve has the controller choose and archive Count cards from their hand,
// stopping early if the hand runs out or the choice is declined.
func (e ArchiveFromHand) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Count; i++ {
		hand := ctx.Resolver.Hand(ctx.Controller)
		if len(hand) == 0 {
			return
		}
		id, ok := ctx.ChooseCreature("Choose a card to archive", hand)
		if !ok {
			return
		}
		ctx.Resolver.ArchiveFromHand(id)
	}
}

// ArchiveTopOfDeck moves the top Count cards of the controller's deck into their
// archives, with no choice involved.
type ArchiveTopOfDeck struct {
	Count int
}

// Text renders the effect, e.g. "archive the top card of your deck".
func (e ArchiveTopOfDeck) Text() string {
	if e.Count == 1 {
		return "archive the top card of your deck"
	}
	return fmt.Sprintf("archive the top %d cards of your deck", e.Count)
}

// Resolve archives the top cards in order, stopping early if the deck runs out.
func (e ArchiveTopOfDeck) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Count; i++ {
		if !ctx.Resolver.ArchiveTopOfDeck(ctx.Controller) {
			return
		}
	}
}

// ArchiveFromDiscard moves a card the controller chooses from their discard pile
// into their archives.
type ArchiveFromDiscard struct{}

// Text renders the effect, naming the source zone explicitly.
func (e ArchiveFromDiscard) Text() string {
	return "archive a card from your discard pile"
}

// Resolve has the controller choose one card from their discard pile and archive
// it, doing nothing when the discard pile is empty or the choice is declined.
func (e ArchiveFromDiscard) Resolve(ctx *EffectContext) {
	discard := ctx.Resolver.Discard(ctx.Controller)
	if len(discard) == 0 {
		return
	}
	id, ok := ctx.ChooseCreature("Choose a card to archive", discard)
	if !ok {
		return
	}
	ctx.Resolver.ArchiveFromDiscard(ctx.Controller, id)
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

// Text renders the effect, e.g. "archive each friendly Knight trait creature from
// play".
func (e ArchiveFromPlay) Text() string {
	return fmt.Sprintf("archive %s from play", e.Target.Text())
}

// Resolve archives each selected in-play card.
func (e ArchiveFromPlay) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.PutIntoArchives(id)
	}
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
