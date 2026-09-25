package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Smite
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ready and fight with a friendly creature. Deal 2 damage to each neighbor of the fought creature.
func TestSmite(t *testing.T) {
	var champion, leftFoe, victim, rightFoe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			Hand:  ct.Cards(Smite),
			InPlay: ct.Cards(
				ct.Bind(&champion, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&leftFoe, ct.Creature(ct.Power(5))),
				ct.Bind(&victim, ct.Creature(ct.Power(10))),
				ct.Bind(&rightFoe, ct.Creature(ct.Power(5))),
			),
		},
	})

	h.P1.Play(Smite)
	h.P1.ExpectPrompt("Choose a creature to fight").Source("Smite")
	h.P1.ClickCard(victim)

	// The fought creature survives the trade; its neighbors each take 2.
	h.Expect(victim).At(ct.PlayArea).Damage(3)
	h.Expect(leftFoe).Damage(2)
	h.Expect(rightFoe).Damage(2)
}

// Smite's neighbor damage lands even when the fight destroys the fought creature:
// the neighbors are read from the fought creature's position at the fight, not
// afterward, so a warded neighbor still spends its ward absorbing the 2 damage.
func TestSmiteHitsNeighborsWhenFoughtCreatureDies(t *testing.T) {
	var champion, leftFoe, victim, wardedFoe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			Hand:  ct.Cards(Smite),
			InPlay: ct.Cards(
				ct.Bind(&champion, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(10))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&leftFoe, ct.Creature(ct.Power(4))),
				ct.Bind(&victim, ct.Creature(ct.Power(3))),
				ct.Bind(&wardedFoe, ct.Creature(ct.Power(4))),
			),
		},
	})
	wardedFoe.Ward()

	h.P1.Play(Smite)
	h.P1.ExpectPrompt("Choose a creature to fight").Source("Smite")
	h.P1.ClickCard(victim)

	// The fight destroys the fought creature, but its neighbors still take Smite's
	// 2: the plain neighbor is damaged, and the warded neighbor's ward absorbs the
	// hit and is spent, so it takes no damage but is no longer warded.
	h.Expect(victim).At(ct.Discard)
	h.Expect(leftFoe).Damage(2)
	h.Expect(wardedFoe).Damage(0)
	if wardedFoe.Warded() {
		t.Error("Smite's damage should have spent the neighbor's ward")
	}
}
