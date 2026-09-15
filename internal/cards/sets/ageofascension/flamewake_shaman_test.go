package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Flamewake Shaman
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Play: Deal 2 damage to a creature.
func TestFlamewakeShaman(t *testing.T) {
	t.Run("deals 2 damage to a chosen creature", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(FlamewakeShaman)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(FlamewakeShaman)
		h.P1.ClickCard(enemy)

		h.Expect(enemy).Damage(2)
	})
}
