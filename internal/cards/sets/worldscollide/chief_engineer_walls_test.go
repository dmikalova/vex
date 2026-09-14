package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chief Engineer Walls
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human
//
//	Elusive.
//	Play/Fight/Reap: You may put an Upgrade or Robot card from your discard pile into your hand.
func TestChiefEngineerWalls(t *testing.T) {
	// discardSetup seeds the discard with an upgrade, a Robot creature, and a
	// non-matching Human creature, and binds the two matching cards.
	discardSetup := func(upgrade, robot *ct.Card) ct.Side {
		return ct.Side{
			House: card.House.StarAlliance,
			Discard: ct.Cards(
				ct.Bind(upgrade, ct.Upgrade(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(robot, ct.Creature(
					ct.OfHouse(card.House.StarAlliance),
					ct.Traits(card.Traits.Robot),
				)),
				ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Traits(card.Traits.Human)),
			),
		}
	}

	t.Run("reap returns a chosen upgrade or Robot card from the discard pile", func(t *testing.T) {
		var upgrade, robot ct.Card
		p1 := discardSetup(&upgrade, &robot)
		p1.InPlay = ct.Cards(ChiefEngineerWalls)
		h := ct.Play(t, ct.Setup{P1: p1})

		h.P1.Reap(ChiefEngineerWalls)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(upgrade)

		h.Expect(upgrade).At(ct.Hand)
		h.Expect(robot).At(ct.Discard)
	})

	t.Run("play offers the same return", func(t *testing.T) {
		var upgrade, robot ct.Card
		p1 := discardSetup(&upgrade, &robot)
		p1.Hand = ct.Cards(ChiefEngineerWalls)
		h := ct.Play(t, ct.Setup{P1: p1})

		h.P1.Play(ChiefEngineerWalls)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(robot)

		h.Expect(robot).At(ct.Hand)
		h.Expect(upgrade).At(ct.Discard)
	})

	t.Run("fight offers the same return", func(t *testing.T) {
		var upgrade, robot ct.Card
		p1 := discardSetup(&upgrade, &robot)
		p1.InPlay = ct.Cards(ChiefEngineerWalls)
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: p1,
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		h.P1.Fight(ChiefEngineerWalls, foe)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(upgrade)

		h.Expect(upgrade).At(ct.Hand)
	})

	t.Run("declining leaves the discard pile untouched", func(t *testing.T) {
		var upgrade, robot ct.Card
		p1 := discardSetup(&upgrade, &robot)
		p1.InPlay = ct.Cards(ChiefEngineerWalls)
		h := ct.Play(t, ct.Setup{P1: p1})

		h.P1.Reap(ChiefEngineerWalls)
		h.P1.ClickOption("No")

		h.Expect(upgrade).At(ct.Discard)
		h.Expect(robot).At(ct.Discard)
	})
}

// TestUpgradeOrRobot checks the deck-generation predicate that guarantees Chief
// Engineer Walls a couple of Upgrades or Robots to pull back (card.PullsMatching).
func TestUpgradeOrRobot(t *testing.T) {
	cases := []struct {
		name string
		def  card.Definition
		want bool
	}{
		{"upgrade", card.Definition{Type: card.Type.Upgrade}, true},
		{
			"robot creature",
			card.Definition{Type: card.Type.Creature, Traits: []card.Trait{card.Traits.Robot}},
			true,
		},
		{
			"plain creature",
			card.Definition{Type: card.Type.Creature, Traits: []card.Trait{card.Traits.Human}},
			false,
		},
	}
	for _, tc := range cases {
		if got := upgradeOrRobot(tc.def); got != tc.want {
			t.Errorf("upgradeOrRobot(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}
