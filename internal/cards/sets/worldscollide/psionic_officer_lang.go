//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PsionicOfficerLang
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After an enemy creature reaps, archive the top card of your deck.
var PsionicOfficerLang = card.New(
	"Psionic Officer Lang",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 337),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
