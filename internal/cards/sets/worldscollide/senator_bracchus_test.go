package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Senator Bracchus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	You may spend Æmber on friendly Creatures as if it were in your pool.
//	Fight/Reap: Exalt Senator Bracchus.
func TestSenatorBracchus(t *testing.T) {
	t.Run("reaping exalts Senator Bracchus", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, InPlay: ct.Cards(SenatorBracchus)},
		})

		h.P1.Reap(SenatorBracchus)

		h.Expect(SenatorBracchus).AmberOn(1)
	})

	t.Run("fighting exalts Senator Bracchus", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, InPlay: ct.Cards(SenatorBracchus)},
			// A low-power, well-armored enemy so Bracchus survives to exalt itself.
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.Power(1), ct.Armor(5))),
			)},
		})

		h.P1.Fight(SenatorBracchus, enemy)

		h.Expect(SenatorBracchus).AmberOn(1)
	})

	t.Run(
		"a friendly creature's Æmber pays for a key while Bracchus is in play",
		func(t *testing.T) {
			var bank ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					InPlay: ct.Cards(SenatorBracchus, ct.Bind(&bank, ct.Creature())),
					Amber:  engine.KeyCost - 2,
				},
			})
			h.Game().AddAmberOn(bank.ID(), 2)

			h.P1.EndTurn()
			h.P2.EndTurn()

			// Pool (KeyCost-2) plus the 2 on the friendly creature covers the key.
			h.P1.ExpectKeys(1)
			h.P1.ExpectAmber(0)
			h.Expect(bank).AmberOn(0)
		},
	)

	t.Run("without Bracchus a creature's Æmber cannot pay for a key", func(t *testing.T) {
		var bank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&bank, ct.Creature())),
				Amber:  engine.KeyCost - 2,
			},
		})
		h.Game().AddAmberOn(bank.ID(), 2)

		h.P1.EndTurn()
		h.P2.EndTurn()

		h.P1.ExpectKeys(0)
		h.Expect(bank).AmberOn(2)
	})
}
