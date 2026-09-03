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
	// Type filters which cards may be chosen; the zero value allows any type.
	Type CardType
	// House filters which cards may be chosen; HouseNone allows any house.
	House House
	// Revealed shows the chosen card to the opponent before archiving it, which
	// is how a filtered choice is verified (Incubation Chamber).
	Revealed bool
	// UpTo lets the controller archive fewer than Count, down to none (Mobius
	// Scroll's "up to 2 cards").
	UpTo bool
}

// Text renders the effect, e.g. "archive a card from your hand" or "archive 2
// cards from your hand". The source zone is named explicitly. A revealed,
// filtered choice reads as the reveal it is: "reveal a Mars creature from your
// hand and archive it".
func (e ArchiveFromHand) Text() string {
	if e.Revealed {
		return "reveal " + indefinite(e.handNoun()) + " from your hand and archive it"
	}
	if e.UpTo {
		return fmt.Sprintf("archive up to %d %ss from your hand", e.Count, e.handNoun())
	}
	if e.Count == 1 {
		return "archive " + indefinite(e.handNoun()) + " from your hand"
	}
	return fmt.Sprintf("archive %d %ss from your hand", e.Count, e.handNoun())
}

// handNoun names the cards the filters admit, e.g. "card" or "Mars creature".
func (e ArchiveFromHand) handNoun() string {
	noun := "card"
	switch e.Type {
	case Creature:
		noun = "creature"
	case Artifact:
		noun = "artifact"
	}
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return noun
}

// candidates returns the cards in hand the filters admit.
func (e ArchiveFromHand) candidates(ctx *EffectContext) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			continue
		}
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			continue
		}
		out = append(out, id)
	}
	return out
}

// Resolve has the controller choose and archive Count cards from their hand,
// stopping early if the hand runs out or the choice is declined.
func (e ArchiveFromHand) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate archives from hand and reports whether anything was archived, so
// "if you do" can follow a reveal (Zyzzix the Many).
func (e ArchiveFromHand) resolveGate(ctx *EffectContext) bool {
	archived := false
	for i := 0; i < e.Count; i++ {
		hand := e.candidates(ctx)
		if len(hand) == 0 {
			return archived
		}
		choose := ctx.ChooseCreature
		if e.UpTo {
			choose = ctx.ChooseCardOptional
		}
		id, ok := choose("Choose a card to archive", hand)
		if !ok {
			return archived
		}
		if e.Revealed {
			ctx.Resolver.Record(CardsRevealedToAll{
				Player: ctx.Controller,
				Cards:  []LocalID{id},
			})
		}
		ctx.Resolver.ArchiveFromHand(id)
		archived = true
	}
	return archived
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
