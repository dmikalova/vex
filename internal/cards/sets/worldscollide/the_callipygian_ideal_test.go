package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// The Callipygian Ideal
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "You may spend Æmber on this creature as if it were in your pool."
//	Play: Exalt this creature.
func TestTheCallipygianIdeal(t *testing.T) {
	t.Run("playing it exalts the creature it upgrades", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&host, ct.Creature())),
				Hand:   ct.Cards(TheCallipygianIdeal),
			},
		})

		// One creature in play, so it is the sole host and is chosen automatically.
		h.P1.Play(TheCallipygianIdeal)

		h.Expect(TheCallipygianIdeal).At(ct.Attached)
		h.Expect(host).AmberOn(1)
	})

	t.Run("the host's Æmber pays for a key while upgraded", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature()), TheCallipygianIdeal),
				),
				Amber: engine.KeyCost - 3,
			},
		})
		h.Game().AddAmberOn(host.ID(), 3)

		h.P1.EndTurn()
		h.P2.EndTurn()

		// Pool (KeyCost-3) plus the 3 on the upgraded creature covers the key.
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0)
		h.Expect(host).AmberOn(0)
	})

	t.Run("without the upgrade the creature's Æmber cannot pay for a key", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&host, ct.Creature())),
				Amber:  engine.KeyCost - 3,
			},
		})
		h.Game().AddAmberOn(host.ID(), 3)

		h.P1.EndTurn()
		h.P2.EndTurn()

		h.P1.ExpectKeys(0)
		h.Expect(host).AmberOn(3)
	})
}
