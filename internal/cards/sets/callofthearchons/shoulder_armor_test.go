package callofthearchons_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	cota "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
)

// Shoulder Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	While this creature is on a flank, it gains +2 power and +2 armor.
func TestShoulderArmor(t *testing.T) {
	var flanker, middle ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			InPlay: ct.Cards(
				ct.Bind(&flanker, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				ct.Bind(&middle, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3)),
			),
		},
	})

	flanker.Attach(cota.ShoulderArmor)
	h.Expect(flanker).Power(5)
	h.Expect(flanker).Armor(2)

	// The same upgrade on a creature in the middle of the battleline grants nothing.
	middle.Attach(cota.ShoulderArmor)
	h.Expect(middle).Power(3)
	h.Expect(middle).Armor(0)
}
