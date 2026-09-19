package engine

import (
	"slices"
	"testing"
)

// TestCardsInPlayCountsUpgrades pins that "in play" means every card someone
// controls in play — creatures, artifacts, and the upgrades on either — so an
// effect that does not name a card type reaches an upgrade. A card under another
// card is out of play and is never reached.
func TestCardsInPlayCountsUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(NewCard("Host", Untamed, Creature, Common, WithPower(3)), 0)
	art := g.AddArtifact(NewCard("Relic", Untamed, Artifact, Common), 0)
	onCreature := g.Register(NewCard("Boon", Untamed, Upgrade, Common), 0)
	g.AttachUpgrade(host, onCreature)
	onArtifact := g.Register(NewCard("Mod", Untamed, Upgrade, Common), 0)
	g.AttachUpgrade(art, onArtifact)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	got := resolverCardsInPlay(ctx, 0)
	want := []LocalID{onCreature, host, onArtifact, art}
	if !slices.Equal(got, want) {
		t.Errorf("cards in play = %v, want %v (each host's upgrades ahead of it)", got, want)
	}

	if n := (CardsInPlay{Player: Controller}).Value(ctx); n != 4 {
		t.Errorf("untyped CardsInPlay = %d, want all 4 cards in play", n)
	}
	if n := (CardsInPlay{
		Player: Controller,
		Type:   Upgrade,
	}).Value(ctx); n != 2 {
		t.Errorf("CardsInPlay{Type: Upgrade} = %d, want the 2 attached upgrades", n)
	}
	if n := (CardsInPlay{
		Player: Controller,
		Type:   Creature,
	}).Value(ctx); n != 1 {
		t.Errorf("CardsInPlay{Type: Creature} = %d, want only the creature", n)
	}
}

// TestInPlayCountsUpgrades pins that the two "is this card in play?" predicates
// agree with resolverCardsInPlay about an attached upgrade. They must: a Target
// that hands a removal an upgrade a predicate still calls out of play makes every
// reachability re-check silently skip it.
func TestInPlayCountsUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(NewCard("Host", Untamed, Creature, Common, WithPower(3)), 0)
	up := g.Register(NewCard("Boon", Untamed, Upgrade, Common), 0)
	g.AttachUpgrade(host, up)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if !g.inPlay(up) || !resolverInPlay(ctx, up) {
		t.Error("an attached upgrade should read in play")
	}
	g.PutIntoHand(host)
	if g.inPlay(up) || resolverInPlay(ctx, up) {
		t.Error("an upgrade should leave play with its host")
	}
}

// TestEachCardInPlayReachesUpgrades pins the split between the target kinds that
// say "card in play" and the ones that name the two types. The former reach an
// attached upgrade; "a creature or artifact" does not, because it says what it
// takes.
func TestEachCardInPlayReachesUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(NewCard("Host", Untamed, Creature, Common, WithPower(3)), 0)
	up := g.Register(NewCard("Boon", Untamed, Upgrade, Common), 0)
	g.AttachUpgrade(host, up)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	for _, tc := range []struct {
		kind TargetKind
		want []LocalID
	}{
		{TargetEachCardInPlay, []LocalID{up, host}},
		{TargetEachFriendlyCardInPlay, []LocalID{up, host}},
		{TargetChosenCreatureOrArtifact, []LocalID{host}},
		{TargetChosenFriendlyCreatureOrArtifact, []LocalID{host}},
	} {
		if got := (Target{Kind: tc.kind}).selectBase(ctx); !slices.Equal(got, tc.want) {
			t.Errorf("kind %v selected %v, want %v", tc.kind, got, tc.want)
		}
	}
}

// TestCrossZoneMoverReachesUpgradesInPlay pins that a mover sourcing from InPlay
// gathers the upgrades as cards in their own right, each ahead of its host, and
// that moving one takes it to the destination rather than shedding it to the
// discard pile — removeFromPlay detaches an attached card on the way out, so the
// mover needs no detach of its own.
func TestCrossZoneMoverReachesUpgradesInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(NewCard("Host", Untamed, Creature, Common, WithPower(3)), 0)
	up := g.Register(NewCard("Boon", Untamed, Upgrade, Common), 0)
	g.AttachUpgrade(host, up)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	mover := crossZoneMover{
		Player:  0,
		Dest:    ToHand,
		Sources: []Zone{InPlay, Discard},
	}
	got := mover.gather(ctx, func(LocalID) bool { return true })
	if !slices.Equal(got, []LocalID{up, host}) {
		t.Fatalf("gather = %v, want [%d %d] (upgrade ahead of its host)", got, up, host)
	}
	if z, ok := mover.originOf(ctx, up); z != InPlay || !ok {
		t.Errorf("originOf(upgrade) = (%v, %v), want (InPlay, true)", z, ok)
	}

	mover.move(ctx, up)
	if !g.State.Hand[0].contains(up) {
		t.Errorf("upgrade did not reach the hand")
	}
	if g.State.Discard[0].contains(up) {
		t.Errorf("upgrade was shed to the discard pile instead of moving")
	}
	if n := len(g.Upgrades(host)); n != 0 {
		t.Errorf("host still carries %d upgrades, want the moved one detached", n)
	}
}

// TestCrossZoneMoveSkipsCardThatLeftItsSourceZones pins that a pick which has
// left every source zone since it was picked is skipped, not moved from where it
// no longer is. A power settle inside the batch can destroy one of the picks and
// a granted "Destroyed:" ability can file it somewhere that is not a source zone
// (Biomatrix Backup archives Inka the Spider while Song of Spring shuffles the
// hand, discard pile, and battleline into the deck). Without the skip, originOf
// falls through to the last source and the card is moved out of a zone it left,
// leaving it in two places at once.
func TestCrossZoneMoveSkipsCardThatLeftItsSourceZones(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(NewCard("Inka", Untamed, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	mover := crossZoneMover{
		Player:  0,
		Dest:    ToDeckShuffled,
		Sources: []Zone{Hand, Discard, InPlay},
	}

	// The card leaves play into the archives — no source zone — after being picked.
	g.PutIntoArchives(id)
	if z, ok := mover.originOf(ctx, id); ok {
		t.Fatalf("originOf = (%v, true), want ok false for a card in no source zone", z)
	}

	mover.move(ctx, id)
	if !g.State.Archives[0].contains(id) {
		t.Error("a pick that left every source zone must stay where it landed")
	}
	if g.State.Deck[0].contains(id) {
		t.Error("a pick that left every source zone must not also reach the destination")
	}
}
