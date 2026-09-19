package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Deusillus
//
//	House:  Saurian
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  20
//	Traits: Mutant
//
//	Play: Deusillus captures all your opponent's Æmber. Deal 5 damage to an enemy creature.
//	Fight/Reap: Move 1 Æmber from Deusillus to the common supply. Deal 2 damage to each enemy creature.
func TestDeusillus(t *testing.T) {
	t.Run(
		"play captures all the opponent's Æmber and deals 5 to an enemy creature",
		func(t *testing.T) {
			var deusillus, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					Hand:  ct.Cards(ct.Bind(&deusillus, Deusillus), card.GiganticArt(Deusillus)),
				},
				P2: ct.Side{
					Amber: 3,
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(8))),
					),
				},
			})

			h.P1.Play(deusillus)

			h.P2.ExpectAmber(0)
			h.Expect(deusillus).AmberOn(3)
			h.Expect(enemy).Damage(5)
		},
	)

	t.Run(
		"reap deals 2 damage to each enemy creature",
		func(t *testing.T) {
			var deusillus, enemy1, enemy2 ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					InPlay: ct.Cards(ct.Bind(&deusillus, Deusillus)),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&enemy1, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					ct.Bind(&enemy2, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				)},
			})

			h.P1.Reap(deusillus)

			h.Expect(enemy1).Damage(2)
			h.Expect(enemy2).Damage(2)
		},
	)
}
