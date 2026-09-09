package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Draco Praeco
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	Reap: You may exalt Draco Praeco, and choose a house - enrage each creature of the chosen house.
func TestDracoPraeco(t *testing.T) {
	t.Run("exalting then choosing a house enrages that house", func(t *testing.T) {
		var draco, brob, unt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&draco, DracoPraeco)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&brob, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				ct.Bind(&unt, ct.Creature(ct.OfHouse(card.House.Untamed))),
			)},
		})
		draco.Ready()

		h.P1.Reap(draco)
		h.P1.ClickCard(draco)
		h.P1.ClickOption("Brobnar")

		h.Expect(draco).AmberOn(1)
		if !h.Game().Enraged(brob.ID()) {
			t.Errorf("%s should be enraged", brob.Name())
		}
		if h.Game().Enraged(unt.ID()) {
			t.Errorf("%s should not be enraged", unt.Name())
		}
	})
}
