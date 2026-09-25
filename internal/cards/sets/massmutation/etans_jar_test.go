package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Etan's Jar
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Play: Name a card. Cards with that name cannot be played until Etan's Jar leaves play.
func TestEtansJar(t *testing.T) {
	t.Run("names a card and bars it from being played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(EtansJar, Snarette, RelentlessCreeper),
			},
		})

		h.P1.Play(EtansJar)
		h.P1.ClickOption("Snarette")

		// The named card cannot be played, but a differently named card still can.
		h.P1.ExpectCannotPlay(Snarette)
		h.P1.Play(RelentlessCreeper)
	})
}
