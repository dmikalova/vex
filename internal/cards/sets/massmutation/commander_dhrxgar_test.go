package massmutation

import (
	"testing"

	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commander Dhrxgar
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	After an upgrade enters play, if it is attached to Commander Dhrxgar or one of its neighbors, gain 1 Æmber.
func TestCommanderDhrxgar(t *testing.T) {
	var boost ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			// Brobnar is the default active house, so a vanilla upgrade is playable.
			InPlay: ct.Cards(CommanderDhrxgar),
			Hand:   ct.Cards(ct.Bind(&boost, ct.Upgrade())),
		},
		P2: ct.Side{},
	})

	// Playing the upgrade attaches it to Commander Dhrxgar (its only host), which
	// gains 1 Æmber for the upgrade landing on it.
	h.P1.Play(boost)
	h.P1.ExpectAmber(1)
}
