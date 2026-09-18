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

	ctx := &EffectContext{Resolver: g, Controller: 0}
	got := resolverCardsInPlay(ctx, 0)
	want := []LocalID{onCreature, host, onArtifact, art}
	if !slices.Equal(got, want) {
		t.Errorf("cards in play = %v, want %v (each host's upgrades ahead of it)", got, want)
	}

	if n := (CardsInPlay{Player: Controller}).Value(ctx); n != 4 {
		t.Errorf("untyped CardsInPlay = %d, want all 4 cards in play", n)
	}
	if n := (CardsInPlay{Player: Controller, Type: Upgrade}).Value(ctx); n != 2 {
		t.Errorf("CardsInPlay{Type: Upgrade} = %d, want the 2 attached upgrades", n)
	}
	if n := (CardsInPlay{Player: Controller, Type: Creature}).Value(ctx); n != 1 {
		t.Errorf("CardsInPlay{Type: Creature} = %d, want only the creature", n)
	}
}
