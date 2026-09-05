package engine

import "fmt"

// This file holds the log entries that narrate a card changing zones (ADR 0011):
// archiving, purging, discarding, and the several ways a card leaves play. No
// entry decides for itself whether it may name the card that moved — each states
// the zones the move ran between and lets nameMoved apply the one rule.

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

// CardArchivedFromHand narrates a card going from hand to archives.
type CardArchivedFromHand struct {
	Player int
	Card   LocalID
}

// Text renders a card archived out of a hand, which stays hidden throughout.
func (e CardArchivedFromHand) Text(n Namer) string {
	return fmt.Sprintf("%s archives %s",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Hand, Archives))
}

// CardArchivedFromDiscard narrates a card going from the discard pile to
// archives.
type CardArchivedFromDiscard struct {
	Player int
	Card   LocalID
}

// Text renders a card archived out of a discard pile, where it was public.
func (e CardArchivedFromDiscard) Text(n Namer) string {
	return fmt.Sprintf("%s archives %s from their discard pile",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Discard, Archives))
}

// TopOfDeckArchived narrates the top card of a deck going to archives sight
// unseen.
type TopOfDeckArchived struct {
	Player int
	Card   LocalID
}

// Text renders the top of a deck archived sight unseen.
func (e TopOfDeckArchived) Text(n Namer) string {
	return fmt.Sprintf("%s archives %s from the top of their deck",
		n.PlayerName(e.Player), nameMoved(n, e.Card, deck, Archives))
}

// ArchivesDiscarded narrates archives emptying into a discard pile.
type ArchivesDiscarded struct {
	Player int
	Count  int
}

// Text renders how many archived cards were discarded.
func (e ArchivesDiscarded) Text(n Namer) string {
	return fmt.Sprintf("%s discards %s",
		n.PlayerName(e.Player), countNoun(e.Count, "archived card"))
}

// TopOfDeckDiscarded narrates the top card of a deck going to the discard pile.
type TopOfDeckDiscarded struct {
	Player int
	Card   LocalID
}

// Text renders the top of a deck going to the discard pile, which is public, so
// the card lands face up and is named.
func (e TopOfDeckDiscarded) Text(n Namer) string {
	return fmt.Sprintf("%s discards %s from the top of their deck",
		n.PlayerName(e.Player), nameMoved(n, e.Card, deck, Discard))
}

// DeckAndDiscardSwapped narrates a deck and discard pile trading places.
type DeckAndDiscardSwapped struct{ Player int }

// Text renders a deck and discard pile trading places.
func (e DeckAndDiscardSwapped) Text(n Namer) string {
	return fmt.Sprintf("%s swaps their deck and discard pile", n.PlayerName(e.Player))
}

// CardDiscarded narrates a card going from a hand to a discard pile.
type CardDiscarded struct {
	Player int
	Card   LocalID
}

// Text renders the card a player discarded, which the public discard pile names.
func (e CardDiscarded) Text(n Namer) string {
	return fmt.Sprintf("%s discards %s",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Hand, Discard))
}

// CardPurgedFromDiscard narrates a card purged out of a discard pile.
type CardPurgedFromDiscard struct {
	Player int
	Card   LocalID
}

// Text renders the card purged out of a discard pile.
func (e CardPurgedFromDiscard) Text(n Namer) string {
	return fmt.Sprintf("%s purges %s from a discard pile",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Discard, purged))
}

// CardPurgedFromHand narrates a card purged out of a hand. Purging turns it
// face up, so naming it leaks nothing.
type CardPurgedFromHand struct {
	Player int
	Card   LocalID
}

// Text renders the card purged out of a hand.
func (e CardPurgedFromHand) Text(n Namer) string {
	return fmt.Sprintf("%s purges %s from a hand",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Hand, purged))
}

// CardPurged narrates a card in play being purged.
type CardPurged struct{ Card LocalID }

// Text renders the card in play that was purged.
func (e CardPurged) Text(n Namer) string {
	return fmt.Sprintf("%s is purged", nameMoved(n, e.Card, inPlay, purged))
}

// CardPutOnTopOfDeck narrates a card leaving play onto its owner's deck.
type CardPutOnTopOfDeck struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play onto its owner's deck.
func (e CardPutOnTopOfDeck) Text(n Namer) string {
	return fmt.Sprintf("%s is put on top of %s's deck",
		nameMoved(n, e.Card, inPlay, deck), n.PlayerName(e.Owner))
}

// CardReturnedToHand narrates a card leaving play into its owner's hand.
type CardReturnedToHand struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play into its owner's hand.
func (e CardReturnedToHand) Text(n Namer) string {
	return fmt.Sprintf("%s is returned to %s's hand",
		nameMoved(n, e.Card, inPlay, Hand), n.PlayerName(e.Owner))
}

// CardPutIntoArchives narrates a card leaving play into its owner's archives.
type CardPutIntoArchives struct {
	Card  LocalID
	Owner int
}

// Text renders a card leaving play into its owner's archives.
func (e CardPutIntoArchives) Text(n Namer) string {
	return fmt.Sprintf("%s is put into %s's archives",
		nameMoved(n, e.Card, inPlay, Archives), n.PlayerName(e.Owner))
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
		nameMoved(n, e.Card, inPlay, deck), n.PlayerName(e.Owner))
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

// CardReturnedFromDiscardToHand narrates a card recovered out of a discard pile.
type CardReturnedFromDiscardToHand struct {
	Player int
	Card   LocalID
}

// Text renders the card recovered from a discard pile to hand.
func (e CardReturnedFromDiscardToHand) Text(n Namer) string {
	return fmt.Sprintf("%s returns %s from their discard pile to hand",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Discard, Hand))
}

// CardPutFromDeckIntoHand narrates a card searched out of a deck. Both zones are
// hidden, so a search that does not reveal what it found says only that much.
type CardPutFromDeckIntoHand struct {
	Player int
	Card   LocalID
}

// Text renders the card searched out of a deck into hand.
func (e CardPutFromDeckIntoHand) Text(n Namer) string {
	return fmt.Sprintf("%s puts %s from their deck into hand",
		n.PlayerName(e.Player), nameMoved(n, e.Card, deck, Hand))
}

// CardPutFromDiscardOnTopOfDeck narrates a card set back on the deck out of a
// discard pile.
type CardPutFromDiscardOnTopOfDeck struct {
	Player int
	Card   LocalID
}

// Text renders the card set from a discard pile back on the deck.
func (e CardPutFromDiscardOnTopOfDeck) Text(n Namer) string {
	return fmt.Sprintf("%s puts %s from their discard pile on top of their deck",
		n.PlayerName(e.Player), nameMoved(n, e.Card, Discard, deck))
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
		n.PlayerName(e.Player), what, face, n.Name(e.Host))
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
