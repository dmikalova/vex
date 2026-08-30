//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SoundTheHorns
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Discard cards from the top of your deck until you either discard a Brobnar creature or run out of cards. If you discarded a Brobnar creature this way, put it into your hand.
var SoundTheHorns = card.New(
	"Sound the Horns",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 15),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
