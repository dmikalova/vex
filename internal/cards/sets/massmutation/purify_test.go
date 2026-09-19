package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Purify
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Purge a Mutant creature -> discard cards from the top of its controller's deck until you discard a non-Mutant creature or run out of cards -> put the discarded creature into play under its owner's control.
func TestPurify(t *testing.T) {
	t.Run("purges an enemy Mutant and reanimates the dug non-Mutant", func(t *testing.T) {
		var mutant, skipped, found, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(Purify),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant))),
				),
				Deck: ct.Cards(
					ct.Bind(&skipped, ct.Creature(ct.Traits(card.Traits.Mutant))),
					ct.Bind(&found, ct.Creature()),
					ct.Bind(&buried, ct.Creature()),
				),
			},
		})

		h.P1.Play(Purify)

		h.Expect(mutant).At(ct.Purge)
		h.Expect(skipped).At(ct.Discard)
		h.Expect(found).At(ct.PlayArea)
		h.Expect(buried).At(ct.Deck)
	})

	t.Run("does nothing when there is no Mutant to purge", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(Purify),
			},
			P2: ct.Side{Deck: ct.Cards(ct.Bind(&top, ct.Creature()))},
		})

		h.P1.Play(Purify)

		h.Expect(top).At(ct.Deck)
	})
}
