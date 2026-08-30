package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Customs Office
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Your opponent must give you 1 Æmber in order to play an artifact.
func TestCustomsOffice(t *testing.T) {
	t.Run("opponent pays the controller 1 Æmber to play an artifact", func(t *testing.T) {
		var toll ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(CustomsOffice)},
			P2: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(ct.Bind(&toll, ct.Artifact())), Amber: 2},
		})

		h.P1.EndTurn() // pass to the opponent
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Play(toll)

		h.P2.ExpectAmber(1) // paid 1 of its 2
		h.P1.ExpectAmber(1) // received the toll
		h.Expect(toll).At(ct.PlayArea)
	})

	t.Run("opponent that cannot pay cannot play the artifact", func(t *testing.T) {
		var toll ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(CustomsOffice)},
			P2: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(ct.Bind(&toll, ct.Artifact())), Amber: 0},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		if _, err := h.Game().PlayArtifact(1, 0); err != engine.ErrCannotPayToll {
			t.Fatalf("PlayArtifact = %v, want ErrCannotPayToll", err)
		}
		h.Expect(toll).At(ct.Hand)
	})
}
