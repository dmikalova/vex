package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kompsos Haruspex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Priest
//
//	Each friendly creature's play effect is a play/reap effect.
func TestKompsosHaruspex(t *testing.T) {
	// playGainer is a friendly creature whose play effect gains 1 Æmber, so the
	// the also-triggers-on rule is visible: with Kompsos in play, reaping it fires that play effect.
	playGainer := card.Build(
		"Play Gainer",
		card.House.Saurian,
		card.Type.Creature,
		card.Rarity.Common,
		card.WithPower(3),
		card.WithAbility(
			card.Trigger.Play,
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		),
	)

	t.Run("with Kompsos in play, reaping fires the play effect", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					KompsosHaruspex,
					ct.Bind(&friend, playGainer),
				),
			},
		})

		h.P1.Reap(friend)
		// +1 from the reap itself, +1 from the play effect the rule fires.
		h.P1.ExpectAmber(2)
	})

	t.Run("without Kompsos, reaping fires no play effect", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&friend, playGainer)),
			},
		})

		h.P1.Reap(friend)
		// Only the reap's own Æmber; the play effect stays dormant.
		h.P1.ExpectAmber(1)
	})
}
