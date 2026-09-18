package engine

import "testing"

// TestPurgeFromArchivesGoesToOwnersPile pins that ownership, not the archives
// holder, decides which purge pile receives the card: archives are the one
// out-of-play zone that can hold an enemy card (Hidden Stash abducts one), and
// purging it must not strand a player-1 card in player 0's purge pile.
func TestPurgeFromArchivesGoesToOwnersPile(t *testing.T) {
	g := NewGame("A", "B", 1)
	abducted := g.Register(NewCard("Snaglet", Untamed, Creature, Common, WithPower(1)), 1)
	g.State.Archives[0].add(abducted)

	g.purgeFromArchives(0, abducted)

	if !g.State.Purge[1].contains(abducted) {
		t.Error("the card should sit in its owner's purge pile")
	}
	if g.State.Purge[0].contains(abducted) {
		t.Error("the card should not sit in the archives holder's purge pile")
	}
	if err := g.InvariantError(); err != nil {
		t.Errorf("state should stay sound: %v", err)
	}
}

// TestInvariantOwnershipOutOfPlay checks the invariant catches a card resting in
// an out-of-play zone that is not its owner's — control never follows a card out
// of play, so only archives may hold an enemy card.
func TestInvariantOwnershipOutOfPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	stray := g.Register(NewCard("Snaglet", Untamed, Creature, Common, WithPower(1)), 1)
	g.State.Discard[0].add(stray)
	if g.InvariantError() == nil {
		t.Error("a player-1 card in player 0's discard pile should violate the invariant")
	}

	g.State.Discard[0].remove(stray)
	g.State.Archives[0].add(stray)
	if err := g.InvariantError(); err != nil {
		t.Errorf("archives may hold an abducted enemy card: %v", err)
	}
}

// TestZoneNounNamesEachSource checks each source zone prints the noun card text
// uses, so "from your ..." reads right wherever an effect names a zone.
func TestZoneNounNamesEachSource(t *testing.T) {
	for zone, want := range map[Zone]string{
		Hand:     "hand",
		Archives: "archives",
		Deck:     "deck",
		Purged:   "purge pile",
		InPlay:   "play",
		Discard:  "discard pile",
	} {
		if got := zone.noun(); got != want {
			t.Errorf("noun = %q, want %q", got, want)
		}
	}
}
