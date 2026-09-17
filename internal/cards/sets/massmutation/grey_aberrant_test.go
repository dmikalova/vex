package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grey Aberrant
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Monk • Mutant
//
//	Each creature loses each of its traits.
func TestGreyAberrant(t *testing.T) {
	var foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Sanctum,
			InPlay: ct.Cards(GreyAberrant),
		},
		P2: ct.Side{
			House: card.House.Sanctum,
			InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Traits(card.Traits.Mutant, card.Traits.Beast))),
			),
		},
	})

	// While Grey Aberrant is in play, every creature loses its traits.
	if h.Game().HasTrait(foe.ID(), card.Traits.Mutant) {
		t.Error("an enemy creature should lose its traits while Grey Aberrant is in play")
	}
	if got := h.Game().TraitCount(foe.ID()); got != 0 {
		t.Errorf("enemy creature trait count = %d, want 0", got)
	}
}
