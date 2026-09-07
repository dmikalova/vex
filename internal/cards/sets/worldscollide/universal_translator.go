//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UniversalTranslator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Use a friendly non-Star Alliance creature."
var UniversalTranslator = card.New(
	"Universal Translator",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 322),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
