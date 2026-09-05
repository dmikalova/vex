package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Prince Derric, Unifier
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  1
//	Traits: Human • Knight
//
//	Play: If you control creatures from 3 or more houses, gain 3 Æmber.
func TestPrinceDerricUnifier(t *testing.T) {
	t.Run("gains 3 Æmber when your creatures span 3 houses", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(PrinceDerricUnifier),
				InPlay: ct.Cards(
					ct.Creature(ct.Power(4), ct.OfHouse(card.House.Brobnar)),
					ct.Creature(ct.Power(4), ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Play(PrinceDerricUnifier)

		h.P1.ExpectAmber(3)
	})

	t.Run("gains nothing when fewer than 3 houses are represented", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(PrinceDerricUnifier),
				InPlay: ct.Cards(ct.Creature(ct.Power(4), ct.OfHouse(card.House.Brobnar))),
			},
			P2: ct.Side{},
		})

		h.P1.Play(PrinceDerricUnifier)

		h.P1.ExpectAmber(0)
	})
}
