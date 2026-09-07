//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Infomancer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  3
//	Traits: Human • Cyborg
//
//	Elusive.
//	Play: Graft an action card from your hand onto Infomancer. (Place it faceup under this card.)
//	Reap: Trigger the play effect of an action card grafted onto Infomancer.
var Infomancer = card.New(
	"Infomancer",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Cyborg),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
