package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Livia the Elder
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may exalt Livia the Elder -> each friendly creature's fight effects and reap effects are fight/reap effects for the remainder of the turn.
func TestLiviaTheElder(t *testing.T) {
	// fightGainer is a friendly creature whose fight effect gains 1 Æmber, so the
	// fuse is visible: while it is active, reaping the creature fires that fight
	// effect too.
	fightGainer := card.Build(
		"Fight Gainer",
		card.House.Saurian,
		card.Type.Creature,
		card.Rarity.Common,
		card.WithPower(3),
		card.WithAbility(
			card.Trigger.Fight,
			card.GainAember{Player: card.Controller, Amount: 1},
		),
	)

	t.Run("taking the exalt fuses fight and reap for the turn", func(t *testing.T) {
		var livia, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&livia, LiviaTheElder),
					ct.Bind(&friend, fightGainer),
				),
			},
		})

		// Reap Livia and take the optional self-exalt, installing the fuse.
		h.P1.Reap(livia)
		h.P1.ClickCard(livia)
		h.Expect(livia).AmberOn(1)
		h.P1.ExpectAmber(1) // the reap's own Æmber

		// Reaping the friend now also fires its fight effect.
		h.P1.Reap(friend)
		// +1 reap + 1 morphed fight effect on top of Livia's reap Æmber.
		h.P1.ExpectAmber(3)
	})

	t.Run("declining the exalt installs no fuse", func(t *testing.T) {
		var livia, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&livia, LiviaTheElder),
					ct.Bind(&friend, fightGainer),
				),
			},
		})

		h.P1.Reap(livia)
		h.P1.ClickDone()
		h.Expect(livia).AmberOn(0)

		h.P1.Reap(friend)
		// Only the two reaps' Æmber; no fight effect fires.
		h.P1.ExpectAmber(2)
	})

	t.Run("the fuse expires at the end of the turn", func(t *testing.T) {
		var livia, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&livia, LiviaTheElder),
					ct.Bind(&friend, fightGainer),
				),
			},
		})

		h.P1.Reap(livia)
		h.P1.ClickCard(livia)

		// Round the table back to P1: their ready phase drops the fuse.
		h.P1.EndTurn()
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Saurian)

		before := h.P1.Amber()
		h.P1.Reap(friend)
		// The fuse is gone, so only the reap's own Æmber lands.
		h.P1.ExpectAmber(before + 1)
	})
}
