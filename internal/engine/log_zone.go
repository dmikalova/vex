package engine

import "fmt"

// This file holds the log entries that narrate a card changing zones (ADR 0011):
// archiving, purging, discarding, and the several ways a card leaves play. No
// entry decides for itself whether it may name the card that moved — each states
// the zones the move ran between and lets nameMoved apply the one rule.

// A zone move is narrated in one of two voices, and each has one composer below
// so the wording cannot drift between the entries that share it. Both take their
// zone words from Zone.noun(), the same source printed card text names zones
// from.

// putIntoZoneText narrates a card leaving play into one of its owner's zones,
// e.g. "Card6 is put into P1's hand". No one is named as doing it: a card leaves
// play into its owner's zone whoever caused it.
func putIntoZoneText(n Namer, id LocalID, owner int, dest Zone) string {
	return fmt.Sprintf("%s is put into %s's %s",
		nameMoved(n, id, inPlay, dest), n.PlayerName(owner), dest.noun())
}

// playerPutsPrefix renders the opening of a "<player> puts <card> from <their>
// <zone>" line, where a player moves one of their own cards between two of their
// zones. It returns the possessive separately for a destination phrase that has
// to name the owner again.
func playerPutsPrefix(n Namer, player int, id LocalID, from, to Zone) (prefix, owner string) {
	who, owner := actorPossessive(n, player)
	return fmt.Sprintf("%s puts %s from %s %s",
		who, nameMoved(n, id, from, to), owner, from.noun()), owner
}

// ArchivesTakenIntoHand narrates a player collecting their archives.
type ArchivesTakenIntoHand struct {
	Player int
	Count  int
}

// Text renders how many archived cards a player took into hand.
func (e ArchivesTakenIntoHand) Text(n Namer) string {
	return fmt.Sprintf("%s takes %s from their archives into hand",
		n.PlayerName(e.Player), countNoun(e.Count, "card"))
}

// CardMoved narrates one specific card moving between zones — archived (To
// Archives), discarded (To Discard), or purged (To purged) from the zone named by
// From. It renders "<player> <verb> <card> <from-phrase>": the verb comes from the
// destination, the card is named only when nameMoved's public/hidden rule allows,
// and the from-phrase varies with the move (purge says "a <zone>", archive and
// discard say "their <zone>", and a move out of a hidden hand names no source).
// The top-of-deck sight-unseen moves stay separate — they name no chosen card.
type CardMoved struct {
	Player int
	Card   LocalID
	From   Zone
	To     Zone
}

// Text renders the move.
func (e CardMoved) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	return fmt.Sprintf("%s %s %s%s",
		who, moveVerb(e.To),
		nameMoved(n, e.Card, e.From, e.To), moveFromPhrase(e.To, e.From, owner))
}

// moveVerb names the action from the destination zone.
func moveVerb(to Zone) string {
	switch to {
	case Archives:
		return "archives"
	case purged:
		return "purges"
	default: // Discard
		return "discards"
	}
}

// moveFromPhrase names the source zone the way each verb prints it — purge says
// "a <zone>", archive and discard name the zone's owner ("their <zone>" or
// "Player 2's <zone>"), and a move out of a hidden hand prints no source.
func moveFromPhrase(to, from Zone, owner string) string {
	if to == purged {
		switch from {
		case Hand:
			return " from a hand"
		case Archives:
			return " from archives"
		case Deck:
			return " from a deck"
		default: // Discard
			return " from a discard pile"
		}
	}
	switch from {
	case Discard:
		return " from " + owner + " discard pile"
	case Archives:
		return " from " + owner + " archives"
	case Deck:
		return " from " + owner + " deck"
	default: // Hand
		return ""
	}
}

// CardDiscarded narrates one card discarded from a hand. A discard forced by a
// card's ability names that card (Old Yurk discards a card) through the record's
// frame; a discard a player makes as their own turn action has no frame and names
// the player.
type CardDiscarded struct {
	// Player is whose hand the card left.
	Player int
	// Card is the discarded card.
	Card LocalID
}

// Text renders the discard, subjected to the forcing card when a frame carries one.
func (e CardDiscarded) Text(n Namer) string {
	return fmt.Sprintf("%s discards %s", subject(n, e.Player), nameMoved(n, e.Card, Hand, Discard))
}

// TopOfDeckArchived narrates the top card of a deck going to archives sight
// unseen.
type TopOfDeckArchived struct {
	Player int
	Card   LocalID
}

// Text renders the top of a deck archived sight unseen.
func (e TopOfDeckArchived) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	return fmt.Sprintf("%s archives %s from the top of %s deck",
		who, nameMoved(n, e.Card, Deck, Archives), owner)
}

// ArchivesDiscarded narrates archives emptying into a discard pile.
type ArchivesDiscarded struct {
	Player int
	Count  int
}

// Text renders how many archived cards were discarded.
func (e ArchivesDiscarded) Text(n Namer) string {
	return fmt.Sprintf("%s discards %s",
		subject(n, e.Player), countNoun(e.Count, "archived card"))
}

// TopOfDeckDiscarded narrates the top card of a deck going to the discard pile.
// Under a card's ability the source card is the subject and the deck is named by
// its owner ("Card discards X from the top of Player 2's deck"); on a player's own
// action the player is the subject ("Player 2 discards X from the top of their
// deck"). The subject comes from the record's frame, not a field on the entry.
type TopOfDeckDiscarded struct {
	Player int
	Card   LocalID
}

// Text renders the top of a deck going to the discard pile, which is public, so
// the card lands face up and is named.
func (e TopOfDeckDiscarded) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	return fmt.Sprintf("%s discards %s from the top of %s deck",
		who, nameMoved(n, e.Card, Deck, Discard), owner)
}

// DeckAndDiscardSwapped narrates a deck and discard pile trading places.
type DeckAndDiscardSwapped struct{ Player int }

// Text renders a deck and discard pile trading places.
func (e DeckAndDiscardSwapped) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	return fmt.Sprintf("%s swaps %s deck and discard pile", who, owner)
}

// ShuffledIntoDeck narrates one or more zones shuffling into a deck (Screaming
// Cave, Help from Future Self). The discard pile is public, so its cards are
// named; the hand and archives are hidden, so they are only counted (ADR 0011).
type ShuffledIntoDeck struct {
	Player int
	// DiscardCards are the named discard cards that went into the deck.
	DiscardCards []LocalID
	// HandCount and ArchivesCount are the hidden cards, reported as counts only.
	HandCount     int
	ArchivesCount int
}

// Text renders the shuffle, naming the public discard cards and counting the
// hidden ones. With no cards to move it is a bare deck shuffle.
func (e ShuffledIntoDeck) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	var clauses []string
	if len(e.DiscardCards) > 0 {
		clauses = append(clauses, namedCardsAnd(n, e.DiscardCards)+" from "+owner+" discard pile")
	}
	if e.HandCount > 0 {
		clauses = append(clauses, countNoun(e.HandCount, "card")+" from "+owner+" hand")
	}
	if e.ArchivesCount > 0 {
		clauses = append(clauses, countNoun(e.ArchivesCount, "card")+" from "+owner+" archives")
	}
	if len(clauses) == 0 {
		return fmt.Sprintf("%s shuffles %s deck", who, owner)
	}
	return fmt.Sprintf("%s shuffles %s into %s deck", who, oxfordAnd(clauses), owner)
}

// CardPurged narrates a card in play being purged.
type CardPurged struct{ Card LocalID }

// Text renders the card in play that was purged.
func (e CardPurged) Text(n Namer) string {
	return fmt.Sprintf("%s is purged", nameMoved(n, e.Card, inPlay, purged))
}

// CardPurgedFromHand narrates a card purged from a hand, naming the card whose
// ability purged it (from the record's frame) and the player whose hand it left
// (Impspecter purges a card from an opponent's hand).
type CardPurgedFromHand struct {
	Card  LocalID
	Owner int
}

// Text renders the purge, naming the purging card and the hand's owner.
func (e CardPurgedFromHand) Text(n Namer) string {
	if s, ok := framedSource(n); ok {
		return fmt.Sprintf("%s purges %s from %s's hand",
			s, nameMoved(n, e.Card, Hand, purged), n.PlayerName(e.Owner))
	}
	return fmt.Sprintf("%s is purged from %s's hand",
		nameMoved(n, e.Card, Hand, purged), n.PlayerName(e.Owner))
}

// CardPutOnTopOfDeck narrates a card leaving play onto its owner's deck.
type CardPutOnTopOfDeck struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play onto its owner's deck.
func (e CardPutOnTopOfDeck) Text(n Namer) string {
	return fmt.Sprintf("%s is put on top of %s's deck",
		nameMoved(n, e.Card, inPlay, Deck), n.PlayerName(e.Owner))
}

// CardPutIntoHand narrates a card leaving play into its owner's hand.
type CardPutIntoHand struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play into its owner's hand.
func (e CardPutIntoHand) Text(n Namer) string {
	return putIntoZoneText(n, e.Card, e.Owner, Hand)
}

// CardPutIntoArchives narrates a card leaving play into its owner's archives.
type CardPutIntoArchives struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play into its owner's archives.
func (e CardPutIntoArchives) Text(n Namer) string {
	return putIntoZoneText(n, e.Card, e.Owner, Archives)
}

// CardShuffledIntoDeck narrates a card leaving play into its owner's deck, lost
// in the shuffle.
type CardShuffledIntoDeck struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play into its owner's shuffled deck.
func (e CardShuffledIntoDeck) Text(n Namer) string {
	return fmt.Sprintf("%s is shuffled into %s's deck",
		nameMoved(n, e.Card, inPlay, Deck), n.PlayerName(e.Owner))
}

// DeckShuffled narrates a player's deck being shuffled (Borr Nit shuffles the
// cards it revealed but did not purge back into that deck).
type DeckShuffled struct{ Player int }

// Text renders a player's deck being shuffled.
func (e DeckShuffled) Text(n Namer) string {
	return fmt.Sprintf("%s's deck is shuffled", n.PlayerName(e.Player))
}

// DiscardRecycledIntoDeck narrates the discard pile becoming a new deck when a
// player draws from an empty one. Routine shuffles are logged too, so that every
// reordering of a deck a player might want to reason about has a line and none
// happens silently.
type DiscardRecycledIntoDeck struct{ Player int }

// Text renders the discard pile recycling into the deck.
func (e DiscardRecycledIntoDeck) Text(n Namer) string {
	return fmt.Sprintf("%s's discard pile is shuffled into their deck",
		n.PlayerName(e.Player))
}

// CardsShuffledIntoDeckBy narrates a card shuffling one owner's creatures into
// that owner's deck as one grouped line ("Lost in the Woods shuffles Murmook and
// Chota Hazri into P2's deck"), so a card that shuffles several creatures at once
// reads as one act rather than a passive line per creature. The source card comes
// from the record's frame.
type CardsShuffledIntoDeckBy struct {
	Owner int
	Cards []LocalID
}

// Text renders the grouped shuffle attributed to its source.
func (e CardsShuffledIntoDeckBy) Text(n Namer) string {
	if s, ok := framedSource(n); ok {
		return fmt.Sprintf("%s shuffles %s into %s's deck",
			s, namedCardsAnd(n, e.Cards), n.PlayerName(e.Owner))
	}
	return fmt.Sprintf("%s shuffles %s into their deck",
		n.PlayerName(e.Owner), namedCardsAnd(n, e.Cards))
}

// CardAbducted narrates a creature taken into the abductor's archives, which is
// the one move that puts a card into a zone belonging to someone other than its
// owner.
type CardAbducted struct {
	Player int
	Card   LocalID
	Owner  int
}

// Text renders the abduction, naming the abducted card's owner.
func (e CardAbducted) Text(n Namer) string {
	return fmt.Sprintf("%s abducts %s (owned by %s) into their archives",
		n.PlayerName(e.Player), nameMoved(n, e.Card, inPlay, Archives), n.PlayerName(e.Owner))
}

// CardPutFromDiscardIntoHand narrates a card recovered out of a discard pile.
type CardPutFromDiscardIntoHand struct {
	Player int
	Card   LocalID
}

// Text renders the card recovered from a discard pile to hand.
func (e CardPutFromDiscardIntoHand) Text(n Namer) string {
	prefix, _ := playerPutsPrefix(n, e.Player, e.Card, Discard, Hand)
	return prefix + " into " + Hand.noun()
}

// CardPutFromDeckIntoHand narrates a card searched out of a deck. Both zones are
// hidden, so a search that does not reveal what it found says only that much.
type CardPutFromDeckIntoHand struct {
	Player int
	Card   LocalID
}

// Text renders the card searched out of a deck into hand.
func (e CardPutFromDeckIntoHand) Text(n Namer) string {
	prefix, _ := playerPutsPrefix(n, e.Player, e.Card, Deck, Hand)
	return prefix + " into " + Hand.noun()
}

// CardPutFromDiscardOnTopOfDeck narrates a card set back on the deck out of a
// discard pile.
type CardPutFromDiscardOnTopOfDeck struct {
	Player int
	Card   LocalID
}

// Text renders the card set from a discard pile back on the deck.
func (e CardPutFromDiscardOnTopOfDeck) Text(n Namer) string {
	prefix, owner := playerPutsPrefix(n, e.Player, e.Card, Discard, Deck)
	return prefix + " on top of " + owner + " " + Deck.noun()
}

// CardPutUnder narrates a card placed under a host (Masterplan, Jargogle,
// Graft). Unlike a move between two of Zone's fixed zones, whether the moved
// card may be named here depends on this particular card, not on a zone
// identity nameMoved can look up — a facedown card is exactly as hidden as one
// in a hand, so it is named only when placed faceup.
type CardPutUnder struct {
	Player   int
	Card     LocalID
	Host     LocalID
	FaceDown bool
}

// Text renders the card placed under its host, naming it only if it went
// faceup.
func (e CardPutUnder) Text(n Namer) string {
	what, face := n.Name(e.Card), "faceup"
	if e.FaceDown {
		what, face = "a card", "facedown"
	}
	return fmt.Sprintf("%s puts %s %s under %s",
		subject(n, e.Player), what, face, n.Name(e.Host))
}

// CardGrafted narrates a card in play being grafted faceup under a host, out of
// play (Spangler Box). Graft always places its card faceup, so the card is
// always named.
type CardGrafted struct {
	Card LocalID
	Host LocalID
}

// Text renders the grafted card and the host it now sits under.
func (e CardGrafted) Text(n Namer) string {
	return fmt.Sprintf("%s is grafted onto %s", n.Name(e.Card), n.Name(e.Host))
}
